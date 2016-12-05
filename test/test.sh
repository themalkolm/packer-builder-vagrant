#!/bin/bash -xe

BOX_NAME="bento/centos-7.2"
BOX_NAME_42="bento42/centos-7.2"
BOX_FILE="bento-VAGRANTSLASH-centos-7.2-virtualbox.box"

function say_green {
    #
    # Using \x1B as OSX has old bash
    # http://stackoverflow.com/questions/28782394/how-to-get-osx-shell-script-to-show-colors-in-echo
    #
    echo -e "\x1B[32m${@}\x1B0"
}

if vagrant box list | grep ${BOX_NAME_42} ; then
    vagrant box remove -f ${BOX_NAME_42} --all
fi
vagrant destroy -f || :
rm -f Vagrantfile

######################################################################
# SETUP
######################################################################

#
# Maje sure to skip vbguest installations
#
if ! vagrant plugin list | grep vagrant-vbguest ; then
    vagrant plugin install vagrant-vbguest
fi

#
# Make sure to have test box
#
if ! vagrant box list | grep ${BOX_NAME} ; then
    vagrant box add --force --provider virtualbox ${BOX_NAME}
fi

#
# Verify initial state
#
cat > Vagrantfile << HERE
Vagrant.configure("2") do |config|
  config.vm.provider "virtualbox"

  config.vm.box = "${BOX_NAME}"
  config.vbguest.auto_update = false
end
HERE
vagrant up
vagrant ssh -- "test ! -f /home/vagrant/foobar"
vagrant destroy -f
rm Vagrantfile

say_green "${BOX_NAME} OK!"

######################################################################
# TEST
######################################################################

#
# Package
#
packer build \
    -var "box_name=${BOX_NAME}" \
    -var "box_file=${BOX_FILE}" \
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

######################################################################
# TEARDOWN
######################################################################

if vagrant box list | grep ${BOX_NAME_42} ; then
    vagrant box remove -f ${BOX_NAME_42} --all
fi
vagrant destroy -f
rm -f Vagrantfile

######################################################################
# REPORT
######################################################################

#
# Success
#
say_green "ALL OK!"
