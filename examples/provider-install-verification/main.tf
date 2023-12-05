terraform {
  required_providers {
    planningcenter = {
      source = "HubbardHarvey3/terraform-provider-planningcenter/planningcenter"
    }
  }
}

provider "planningcenter" {
}

//data "planningcenter_people" "test" {
//  id = "138378248"
//}

resource "planningcenter_people" "new_test" {
  first_name         = "Testy"
  last_name          = "McTesterson"
  site_administrator = false
  gender             = "Male"
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

//output "site_admin" {
//  value = data.planningcenter_people.test.site_administrator
//}
