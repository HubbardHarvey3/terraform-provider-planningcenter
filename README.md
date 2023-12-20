<a href="https://terraform.io">
    <img src=".github/tf.png" alt="Terraform logo" title="Terraform" align="left" height="50" />
</a>

# Terraform Provider for Planning Center

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.19

## Building The Provider

1. Pick your appropriate provider version from the release page
2. Extract the provider to the root of your Terraform directory
    <br> a. you should see a path similar to the one below in the root of your TF directory
```.terraform.d/plugins/github.com/HubbardHarvey3/planningcenter/<version>/<os>_<arch>```

If this project actually gets used, then I will probably setup future versions with the legit Terraform provider registry

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the provider

Before using the provider, ensure you have the following environment variables set by running the following commands:

```export PC_APP_ID=<your app_id>```<br>
```export PC_SECRET_TOKEN=<your secret>```

The Application ID and Secret can be found in your PlanningCenter developer account: https://api.planningcenteronline.com/oauth/applications
<br>
More info from the PlanningCenter Docs: https://developer.planning.center/docs/#/overview/authentication

```hcl

terraform {
  required_providers {
    planningcenter = {
      source = "github.com/HubbardHarvey3/planningcenter"
    }
  }
}

provider "planningcenter" {
}


resource "planningcenter_people" "test-people" {
  first_name         = "Tester"
  last_name          = "McTesterson"
  site_administrator = false
  gender             = "Male"
}

locals {
  addresses = ["testyUPDATE@hcubedcoder.com", "testerupdate@notgoogl.com"]
}

resource "planningcenter_email" "Tester" {
  count = length(local.addresses)
  address = local.addresses[count.index]
  primary = count.index == 0 ? true : false
  location = count.index == 0 ? "Home" : "Work"
  relationships = {
    id = planningcenter_people.new_test.id
  }
}

```
## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```

## TODO / Wishlist
- [x] Build People resource
  - [x] get acc tests for people working

- [x] Build Email resource
  - [ ] get acc tests for email working (Import not tested)

- [ ] getting proper versioning setup
- [ ] setup proper release process
- [x] setup docs

