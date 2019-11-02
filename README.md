Terraform Linux Provider [![Build Status](https://travis-ci.org/mavidser/terraform-provider-linux.svg?branch=master)](https://travis-ci.org/mavidser/terraform-provider-linux)
========================

- [Documentation](https://github.com/mavidser/terraform-provider-linux/tree/master/docs)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.12.x
-	[Go](https://golang.org/doc/install) 1.12 (to build the provider plugin)

Usage
---------------------

```terraform
provider "linux" {
  host = "192.168.1.2"
  user = "user"
}
```

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/mavidser/terraform-provider-linux`

```sh
$ mkdir -p $GOPATH/src/github.com/terraform-providers; cd $GOPATH/src/github.com/terraform-providers
$ git clone git@github.com:mavidser/terraform-provider-linux
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/mavidser/terraform-provider-linux
$ make build
```

Using the provider
----------------------

Sample configuration for creating a few users/groups:

```terraform
resource "linux_group" "testgroup" {
  name = "testgroup"
  system = false
}

resource "linux_user" "testuser1" {
  name = "testuser1"
  gid = linux_group.testgroup.gid
  system = false
}

resource "linux_user" "testuser2" {
  name = "testuser2"
  gid = linux_group.testgroup.gid
  system = false
}
```

Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.11+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-linux
...
```

In order to test the provider, you can simply run `make test`.
Note: These tests will require docker installed to spin up a container with ssh access.

```sh
$ make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.

```sh
$ make testacc
```

In order to run only single Acceptance tests, execute the following steps:

```sh
# setup the testing environment
$ source ./scripts/tests_setup.sh

# run single tests
TF_LOG=INFO TF_ACC=1 go test -v ./linux -run TestAccUserCreation -timeout 360s

# cleanup the local testing resources
$ source ./scripts/tests_cleanup.sh
```