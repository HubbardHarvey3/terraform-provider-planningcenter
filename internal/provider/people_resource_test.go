// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccExampleResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccPeopleResourceConfig("one"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("planningcenter_people.test", "first_name", "acceptance"),
					resource.TestCheckResourceAttr("planningcenter_people.test", "site_administrator", "false"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "planningcenter_people.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccPeopleResourceConfig("two"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("planningcenter_people.test", "last_name", "ChangedTest"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccPeopleResourceConfig(configurableAttribute string) string {
	return fmt.Sprintf(`
resource "planningcenter_people" "test" {
  first_name="acceptance"
  last_name="ChangedTest"
  gender="Male"
}
`)
}
