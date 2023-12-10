// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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
					resource.TestCheckResourceAttr("planningcenter_email.test", "address", "acceptance@test.com"),
					resource.TestCheckResourceAttr("planningcenter_email.test", "primary", "true"),
					resource.TestCheckResourceAttr("planningcenter_email.test", "location", "Home"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "planningcenter_email.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccEmailResourceConfig("two"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("planningcenter_email.test", "location", "Work"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccEmailResourceConfig(configurableAttribute string) string {
	return fmt.Sprintf(`
resource "planningcenter_people" "test" {
  first_name="acceptance"
  last_name="ChangedTest"
  gender="Male"
}
resource "planningcenter_email" "test" {
  address="acceptance@test.com"
  primary=true
  location="Home"
  relationships = {
    id = planningcenter_people.test.id
  }
}
`)
}
