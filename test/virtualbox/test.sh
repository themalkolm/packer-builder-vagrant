#!/bin/bash

set -xeu

BOX_NAME="bento/centos-7.2"
BOX_VERSION="2.3.1"
BOX_NAME_42="bento42/centos-7.2"
BOX_FILE="bento-VAGRANTSLASH-centos-7.2-virtualbox.box"

function say_green {
    #
    # Using \x1B as OSX has old bash
    # http://stackoverflow.com/questions/28782394/how-to-get-osx-shell-script-to-show-colors-in-echo
    #
    echo -e "\x1B[32m${@}\x1B0"
}

function setup {
    if ! vagrant plugin list | grep vagrant-vbguest ; then
        vagrant plugin install vagrant-vbguest
    fi
}

function teardown {
    if vagrant box list | grep ${BOX_NAME_42} ; then
        vagrant box remove -f ${BOX_NAME_42} --all
    fi
    vagrant destroy -f || :
    rm -f Vagrantfile
    rm -f *.box || :
}

teardown

######################################################################
# SETUP
######################################################################

setup

######################################################################
# TEST
######################################################################

#
# Package
#
packer build \
    -var "box_name=${BOX_NAME}" \
    -var "box_version=${BOX_VERSION}" \
    -var "box_file=${BOX_FILE}" \
    -var "vm_cpus=1" \
    template.json
vagrant box add --force --name ${BOX_NAME_42} ${BOX_FILE}

say_green "${BOX_NAME} REPACKAGED!"

######################################################################
# VERIFY
######################################################################

#
# Verify changed state
#
cat > Vagrantfile << HERE
Vagrant.configure("2") do |config|
  config.vm.provider "virtualbox"

  config.vm.box = "${BOX_NAME_42}"
  config.vbguest.auto_update = false
end
HERE
vagrant up
vagrant ssh -- "test -f /home/vagrant/foobar"
vagrant ssh -- "cat /home/vagrant/foobar"

say_green "${BOX_NAME} VERIFIED!"

teardown

say_green "ALL OK!"
