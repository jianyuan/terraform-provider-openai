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
	groupName := acctest.RandomWithPrefix("tf-group")
	projectName := acctest.RandomWithPrefix("tf-project")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectGroupRoleAssignmentsDataSourceConfig(roleName, groupName, projectName),
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

func testAccProjectGroupRoleAssignmentsDataSourceConfig(roleName, groupName, projectName string) string {
	return testAccProjectGroupRoleAssignmentResourceConfig(roleName, groupName, projectName) + `
data "openai_project_group_role_assignments" "test" {
	project_id = openai_project_group_role_assignment.test.project_id
	group_id = openai_project_group_role_assignment.test.group_id
}
`
}
