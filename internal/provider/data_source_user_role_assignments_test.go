package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jianyuan/terraform-provider-openai/internal/acctest"
)

func TestAccUserRoleAssignmentsDataSource(t *testing.T) {
	rn := "data.openai_user_role_assignments.test"
	roleName := acctest.RandomWithPrefix("tf-role")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {
				Source: "hashicorp/time",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccUserRoleAssignmentsDataSourceConfig(acctest.TestUserId, roleName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("roles"), knownvalue.SetPartial([]knownvalue.Check{
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"id":              knownvalue.NotNull(),
							"name":            knownvalue.StringExact(roleName),
							"description":     knownvalue.NotNull(),
							"permissions":     knownvalue.NotNull(),
							"predefined_role": knownvalue.Bool(false),
							"resource_type":   knownvalue.StringExact("api.organization"),
						}),
					})),
				},
			},
		},
	})
}

func testAccUserRoleAssignmentsDataSourceConfig(userId, roleName string) string {
	return testAccUserRoleAssignmentResourceConfig(userId, roleName) + `
resource "time_sleep" "wait" {
	create_duration = "5s"

	triggers = {
		user_id = openai_user_role_assignment.test.user_id
	}
}

data "openai_user_role_assignments" "test" {
	user_id = time_sleep.wait.triggers.user_id
}
`
}
