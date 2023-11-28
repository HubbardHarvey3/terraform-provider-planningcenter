terraform {
  required_providers {
    planningcenter = {
      source = "HubbardHarvey3/terraform-provider-planningcenter/planningcenter"
    }
  }
}

provider "planningcenter" {
}

data "planningcenter_people" "test" {
  id = ""
}

output "name" {
   value =  data.planningcenter_people.test.name
}

output "gender" {
   value = data.planningcenter_people.test.gender
}
