package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jianyuan/terraform-provider-openai/internal/acctest"
)

func TestAccProjectGroupRoleAssignmentsDataSource(t *testing.T) {
	rn := "data.openai_project_group_role_assignments.test"
	roleName := acctest.RandomWithPrefix("tf-role")
	projectName := acctest.RandomWithPrefix("tf-project")

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
				Config: testAccProjectGroupRoleAssignmentsDataSourceConfig(projectName, acctest.TestGroupId, roleName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("roles"), knownvalue.SetExact([]knownvalue.Check{
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"id":              knownvalue.NotNull(),
							"name":            knownvalue.StringExact(roleName),
							"description":     knownvalue.NotNull(),
							"permissions":     knownvalue.NotNull(),
							"predefined_role": knownvalue.Bool(false),
							"resource_type":   knownvalue.StringExact("api.project"),
						}),
					})),
				},
			},
		},
	})
}

func testAccProjectGroupRoleAssignmentsDataSourceConfig(projectName, groupId, roleName string) string {
	return testAccProjectGroupRoleAssignmentResourceConfig(projectName, groupId, roleName) + `
resource "time_sleep" "wait" {
	create_duration = "5s"

	triggers = {
		project_id = openai_project_group_role_assignment.test.project_id
		group_id   = openai_project_group_role_assignment.test.group_id
	}
}

data "openai_project_group_role_assignments" "test" {
	project_id = time_sleep.wait.triggers.project_id
	group_id   = time_sleep.wait.triggers.group_id
}
`
}
