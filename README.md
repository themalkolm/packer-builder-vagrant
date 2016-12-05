# packer-builder-vagrant [![Build Status](https://travis-ci.org/themalkolm/packer-builder-vagrant.svg?branch=master)](https://travis-ci.org/themalkolm/packer-builder-vagrant)

Builder proxy. It finds the corresponding vagrant box in the local vagrant cache and
forwards the actual building the to corresponding builder.

Installation
------------

### Pre-built binaries

The easiest way to install this post-processor is to download a pre-built binary from release page.
Download the correct binary for your platform and place it in one of the following places:

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

Add the builder to your packer template:

```json
{
  "builders": [
    {
      "type": "vagrant",

      "box_name": "bento/centos-7.2",
      "box_provider": "virtualbox",
      "box_file": ".ovf",

      "config": {
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

Note that `config` key contains full configuration for `virtualbox-ovf` builder except source_path.
`vagrant` builder will find a locally cached file for this specified box and automatically add `source_path`
key to the `virtualbox-ovf` configuration.
