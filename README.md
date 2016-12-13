# packer-builder-vagrant [![Build Status](https://travis-ci.org/themalkolm/packer-builder-vagrant.svg?branch=master)](https://travis-ci.org/themalkolm/packer-builder-vagrant)

Builder proxy. It finds the corresponding vagrant box in the local vagrant cache and
forwards the actual building the to corresponding builder.

**WARNING** It essentially compiles **whole** packer code to allow us to configure any
builder plugin.

Brief
-----

Packer allows to package existing boxes and save results as vagrant boxes. The problem is that all
[builders](https://www.packer.io/docs/templates/builders.html) expect you to provide absolute path to box file
under `source_path` e.g.:

```
{
  "type": "virtualbox-ovf",
  "source_path": "path/to/source.ovf",
  "ssh_username": "packer",
  "ssh_password": "packer",
  "shutdown_command": "echo 'packer' | sudo -S shutdown -P now"
}
```

This is not convenient if you think about it. Won't it be nice for packer to automagically find existing vagrant box,
provision it and build a new version right away?

`packer-builder-vagrant` does exactly that. It works as a proxy and finds existing vagrant boxes before starting the
underlying actuall builder with the updated configuration.

Installation
------------

### Pre-built binaries

The easiest way to install this post-processor is to download a pre-built binary from release page. Download the
correct binary for your platform and place it in one of the following places:

* The directory where packer is, or the executable directory.
* `~/.packer.d/plugins` on Unix systems or `%APPDATA%/packer.d/plugins` on Windows.
* The current working directory.

Don't forget to strip off os and arch information from the executable i.e:

```
$ mkdir -p ~/.packer.d/plugins
$ cp packer-0.12.0_packer-builder-vagrant_linux_amd64 ~/.packer.d/plugins/packer-builder-vagrant
```

See [docs](https://www.packer.io/docs/extend/plugins.html) for more info.

Usage
-----

You can convert your existing configuraiton in the following way to enable automatic vagrant boxes discovery:

```json
{
  "builders": [
    {
      "type": "virtualbox-ovf",
      "source_path": "/home/john.doe/.vagrant.d/boxes/bento-VAGRANTSLASH-centos-7.2/2.3.1/virtualbox/box.ovf",
      "guest_additions_mode": "disable",
      "headless": true,
      "ssh_username": "vagrant",
      "ssh_password": "vagrant",
      "ssh_pty": true,
      "ssh_private_key_file": "key/vagrant",
      "shutdown_command": "echo '/sbin/halt -h -p' > /tmp/shutdown.sh; echo 'vagrant'|sudo -S sh '/tmp/shutdown.sh'"
    }
  ],
  ...
}
```

Wrap `virtualbox-ovf` in `vagrant` builder the following way:

```json
{
  "builders": [
    {
      "type": "vagrant",

      "box_name": "bento/centos-7.2",
      "box_provider": "virtualbox",
      "box_file": ".ovf",

      "builder": {
        "type": "virtualbox-ovf",
        "guest_additions_mode": "disable",
        "headless": true,
        "ssh_username": "vagrant",
        "ssh_password": "vagrant",
        "ssh_pty": true,
        "ssh_private_key_file": "key/vagrant",
        "shutdown_command": "echo '/sbin/halt -h -p' > /tmp/shutdown.sh; echo 'vagrant'|sudo -S sh '/tmp/shutdown.sh'"
      }
    }
  ],
  ...
}
```

Note that `config` key contains full configuration for `virtualbox-ovf` builder except `source_path`.
`vagrant` builder will find a locally cached file for this specified box and automatically add `source_path`
key to the `virtualbox-ovf` configuration.

Configuration
-------------

All configuration properties are **required**, except where noted.

### box_name

The name of the box to find.

### box_version (optional)

Box version to find. It uses most recent version available by default.

### box_provider

Provider to find box for e.g.:

* `virtualbox`
* `vmware_desktop`
* [...](https://www.vagrantup.com/docs/providers/)

### box_file

Box file to look for in the box directory. This file will provided in `source_path` to the underlying builder. Different boxes have different file names so it is a *regexp pattern* that must mach **only one** file in the directory.

### builder

Here you should provide configuration for the actual builder to be used. `source_path` will be automatically populated by box discovery. Note that you can always explicitly define `source_path` for test purposes. In this case `source_path` won't be changed.

