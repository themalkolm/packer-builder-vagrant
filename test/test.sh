#!/bin/bash -xe

BOX_NAME="bento/centos-7.2"
BOX_NAME_42="bento42/centos-7.2"
BOX_FILE="bento-VAGRANTSLASH-centos-7.2-virtualbox.box"

if vagrant box list | grep ${BOX_NAME_42} ; then
    vagrant box remove ${BOX_NAME_42} --all
fi
vagrant destroy -f || :
rm -f Vagrantfile

######################################################################
# SETUP
######################################################################

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

######################################################################
# TEARDOWN
######################################################################

if vagrant box list | grep ${BOX_NAME_42} ; then
    vagrant box remove ${BOX_NAME_42} --all
fi
vagrant destroy -f
rm -f Vagrantfile

######################################################################
# REPORT
######################################################################

#
# Success
#
echo "OK"