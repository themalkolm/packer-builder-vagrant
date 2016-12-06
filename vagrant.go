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
)

const (
	defaultVersion = "0"
)

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

// Find a file from the box's directory matching the provided pattern
func findBoxFile(name, version, provider string, pattern string) (string, error) {
	box, err := findCachedBox(name, version, provider)
	if err != nil {
		return "", err
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

// Find a box with the provided name, provide and version.
func findCachedBox(name, version, provider string) (*vagrantutil.Box, error) {
	v, err := vagrantutil.NewVagrant(".")
	if err != nil {
		return nil, err
	}

	boxes, err := v.BoxList()
	if err != nil {
		return nil, err
	}

	found := make([]*vagrantutil.Box, 0)
	for _, b := range boxes {
		if b.Name == name && b.Provider == provider {
			found = append(found, b)
		}
	}

	if len(found) == 0 {
		return nil, fmt.Errorf("Can't find box: %s (%s, %s)", name, provider, version)
	}

	s, err := newBoxSorter(found)
	if err != nil {
		return nil, err
	}

	sort.Sort(s)

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

	if box == nil {
		return nil, fmt.Errorf("Can't find box: %s (%s, %s)", name, provider, version)
	}
	return box, nil
}

// Return directory where all box files are stored
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
