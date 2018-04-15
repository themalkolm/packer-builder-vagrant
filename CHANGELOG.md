## 2018.04.15

* Migrated to packer 1.2.2
* Migrated from semver to hashicorp/go-version (issue #11)

## 2018.03.01

* Migrated to golang 1.10
* Migrated to packer 1.2.0

## 2018.01.31

* Migrated to packer 1.1.3

## 2017.10.17

* Migrated to golang 1.9
* Migrated to packer 1.0.4

## 2017.07.14

* Migrated to packer 1.0.2
* Migrated mitchellh/packer -> hashicorp/packer
* Copying packer code as is, removed glide

## 2017.05.10-1

* Do not render builder's config (issues/#5)
* Pass packer_on_error & packer_user_variables to wrapped builder
* Bump test VM to bento/centos-7.3

## 2017.05.10

* Migrated to golang 1.8
* Migrated to packer 1.0.0
* Strip binaries (~50M -> ~30M)
