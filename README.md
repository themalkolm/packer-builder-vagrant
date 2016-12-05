# packer-builder-vagrant [![Build Status](https://travis-ci.org/themalkolm/packer-builder-vagrant.svg?branch=master)](https://travis-ci.org/themalkolm/packer-builder-vagrant)

Builder proxy. It finds the corresponding vagrant box in the local vagrant cache and
forwards the actual building the to corresponding builder.

**WARNING** Currently works only with `virtualbox-ovf` due me being lazy.

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
underlaying actuall builder with the updated configuration.

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
$ cp packer-builder-vagrant-linux-amd64 ~/.packer.d/plugins/packer-builder-vagrant
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
