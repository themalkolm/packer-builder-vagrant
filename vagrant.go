package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"regexp"
	"sort"
	"strings"
	"path/filepath"

	"github.com/blang/semver"
	"github.com/koding/vagrantutil"
	"github.com/mitchellh/packer/packer"
)

const (
	defaultVersion = "0"
)

type Vagrant struct {
	ui      packer.Ui
	vagrant *vagrantutil.Vagrant
}

func NewVagrant(ui packer.Ui) (*Vagrant, error) {
	v, err := vagrantutil.NewVagrant(".")
	if err != nil {
		return nil, err
	}
	return &Vagrant{
		ui: ui,
		vagrant: v,
	}, nil
}

type boxSorter struct {
	boxes    []*vagrantutil.Box
	versions []*semver.Version
}

func newBoxSorter(boxes []*vagrantutil.Box) (*boxSorter, error) {
	versions := make([]*semver.Version, len(boxes))
	for i, b := range boxes {
		version := b.Version
		if b.Version == defaultVersion {
			version = semver.Version{}.String() // 0.0.0
		}

		v, err := semver.ParseTolerant(version)
		if err != nil {
			return nil, err
		}
		versions[i] = &v
	}
	return &boxSorter{
		boxes:    boxes,
		versions: versions,
	}, nil
}

func (s *boxSorter) Len() int {
	return len(s.boxes)
}

func (s *boxSorter) Swap(i, j int) {
	s.boxes[i], s.boxes[j] = s.boxes[j], s.boxes[i]
	s.versions[i], s.versions[j] = s.versions[j], s.versions[i]
}

func (s *boxSorter) Less(i, j int) bool {
	return (*s.versions[i]).LT(*s.versions[j])
}

func (v *Vagrant) fetchBoxFile(url, name, version, provider, pattern string) (string, error) {
	box, err := v.fetchBox(url, name, version, provider)
	if err != nil {
		return "", err
	}

	if box == nil {
		return "", fmt.Errorf("Can't find box: %s (%s, %s)", name, provider, version)
	}

	root, err := boxDir(box)
	if err != nil {
		return "", err
	}

	fs, err := ioutil.ReadDir(root)
	if err != nil {
		return "", err
	}

	p, err := regexp.Compile(pattern)
	if err != nil {
		return "", err
	}

	found := make([]string, 0)
	for _, f := range fs {
		if f.IsDir() {
			continue
		}

		b := []byte(f.Name())
		if p.Match([]byte(b)) {
			found = append(found, filepath.Join(root, f.Name()))
		}
	}

	if len(found) == 0 {
		return "", fmt.Errorf("Can't find a file for box: %s (%s, %s)", box.Name, box.Provider, box.Version)
	}

	if len(found) > 1 {
		return "", fmt.Errorf("More than one (%d) file matched pattern (%s) for box %s (%s, %s): %s",
			len(found), pattern, box.Name, box.Provider, box.Version, strings.Join(found, ", "))
	}

	return found[0], nil
}

func (v *Vagrant) downloadBox(nameOrUrl, version, provider string) (bool, error) {
	box := vagrantutil.Box{
		Name: nameOrUrl,
		Provider: provider,
		Version: version,
	}

	output, err := v.vagrant.BoxAdd(&box)
	if err != nil {
		return false, err
	}

	var outputErr error = nil
	for res := range output {
		if res.Error != nil {
			outputErr = packer.MultiErrorAppend(outputErr, res.Error)
			continue
		}

		line := res.Line
		if strings.HasPrefix(line, "==> ") {
			line = strings.TrimPrefix(line, "==> ")
		}
		v.ui.Message(fmt.Sprintf("(vagrant) %s", line))
	}
	if outputErr != nil {
		return false, outputErr
	}
	return true, nil
}

func (v *Vagrant) fetchBox(url, name, version, provider string) (*vagrantutil.Box, error) {
	box, err := v.findBox(name, version, provider)
	if err != nil {
		return nil, err
	}

	if box != nil {
		return box, nil
	}

	nameOrUrl := name
	if url != "" {
		nameOrUrl = url
	}

	v.ui.Message(fmt.Sprintf("(vagrant) Downloading box: %s (%s, %s)", nameOrUrl, provider, version))
	ok, err := v.downloadBox(nameOrUrl, version, provider)
	if err != nil {
		return nil, fmt.Errorf("Error while downloading box: %s", err)
	}
	v.ui.Message(fmt.Sprintf("(vagrant) Downloaded box: %s (%s, %s)", nameOrUrl, provider, version))

	if !ok {
		return nil, fmt.Errorf("Failed to cache box: %s (%s, %s)", name, provider, version)
	}

	return v.findBox(name, version, provider)
}

func (v *Vagrant) findBox(name, version, provider string) (*vagrantutil.Box, error) {
	boxes, err := v.vagrant.BoxList()
	if err != nil {
		return nil, err
	}

	found := make([]*vagrantutil.Box, 0)
	for _, b := range boxes {
		if b.Name == name && b.Provider == provider {
			found = append(found, b)
		}
	}

	s, err := newBoxSorter(found)
	if err != nil {
		return nil, err
	}

	sort.Sort(s)

	if len(found) == 0 {
		return nil, nil
	}

	var box *vagrantutil.Box = nil
	if version != "" {
		for _, b := range found {
			if b.Version == version {
				box = b
				break
			}
		}
	} else {
		box = found[len(found) - 1]
	}

	return box, nil
}

func boxDir(b *vagrantutil.Box) (string, error) {
	name := b.Name
	version := b.Version
	provider := b.Provider

	if strings.Contains(name, "/") {
		name = strings.Replace(name, "/", "-VAGRANTSLASH-", -1)
	}

	if version == "" {
		version = defaultVersion
	}

	u, err := user.Current()
	if err != nil {
		return "", nil
	}

	home := filepath.Join(u.HomeDir, ".vagrant.d")
	if v, ok := os.LookupEnv("VAGRANT_HOME"); ok {
		home = v
	}
	return fmt.Sprintf("%s/boxes/%s/%s/%s", home, name, version, provider), nil
}
