// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	//"github.com/hashicorp/terraform-plugin-testing/config"
)

func TestAccEmailResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccEmailResourceConfig("one"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("planningcenter_people.testemail", "first_name", "acceptance"),
					resource.TestCheckResourceAttr("planningcenter_email.test", "address", "acceptance@testable.com"),
					resource.TestCheckResourceAttr("planningcenter_email.test", "primary", "true"),
					resource.TestCheckResourceAttr("planningcenter_email.test", "location", "Home"),
				),
			},
			// ImportState testing is skipped for now
			//			{
			//				ResourceName:      "planningcenter_email.test",
			//				ImportState:       true,
			//				ImportStateVerify: false,
			//			},
			// Update and Read testing
			{
				Config: testAccEmailResourceConfig("two"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("planningcenter_email.test", "address", "acceptance@testable.com"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccEmailResourceConfig(configurableAttribute string) string {
	return fmt.Sprintf(`
    resource "planningcenter_people" "testemail" {
      first_name="acceptance"
      last_name="tests"
      gender="Male"
      site_administrator=false
    }
    resource "planningcenter_email" "test" {
      address="acceptance@testable.com"
      primary=true
      location="Home"
      relationships = {
        id = planningcenter_people.testemail.id
      }
    }
`)
}
