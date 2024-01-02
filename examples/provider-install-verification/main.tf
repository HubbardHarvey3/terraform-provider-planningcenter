terraform {
  required_providers {
    planningcenter = {
      source = "github.com/HubbardHarvey3/planningcenter"
      version = "0.0.1"
    }
  }
}

provider "planningcenter" {
}

data "planningcenter_people" "test" {
  id = "138378248"
}

resource "planningcenter_people" "new_test" {
  first_name         = "Tester"
  last_name          = "McTesterson"
  site_administrator = false
  gender             = "Male"
  birthdate          = "1980-05-22"
}

resource "planningcenter_people" "testy" {
  first_name         = "Testy"
  last_name          = "McTesty"
  site_administrator = false
  gender             = "Female"
}

resource "planningcenter_people" "import-me" {
  first_name = "Importy"
  last_name  = "Importenson"
  gender     = "Male"
  birthdate  = "2012-02-10"
}

locals {
  addresses = ["testyUPDATE@hcubedcoder.com", "testerupdate@notgoogl.com"]
}

resource "planningcenter_email" "Tester" {
  count    = length(local.addresses)
  address  = local.addresses[count.index]
  primary  = count.index == 0 ? true : false
  location = count.index == 0 ? "Home" : "Work"
  relationships = {
    id = planningcenter_people.new_test.id
  }
}

resource "planningcenter_email" "import" {
  address  = "import@important.com"
  primary  = true
  location = "Other"
  relationships = {
    id = planningcenter_people.testy.id
  }

}

output "name_test" {
  value = planningcenter_people.testy.first_name
}
output "name" {
  value = planningcenter_people.new_test.first_name
}

output "gender" {
  value = planningcenter_people.new_test.gender
}

output "tester_id" {
  value = planningcenter_people.new_test.id
}

output "site_admin" {
  value = data.planningcenter_people.test.site_administrator
}
