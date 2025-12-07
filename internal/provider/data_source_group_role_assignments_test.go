package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jianyuan/terraform-provider-openai/internal/acctest"
)

func TestAccGroupRoleAssignmentsDataSource(t *testing.T) {
	rn := "data.openai_group_role_assignments.test"
	groupName := acctest.RandomWithPrefix("tf-group")
	roleName := acctest.RandomWithPrefix("tf-role")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGroupRoleAssignmentsDataSourceConfig(groupName, roleName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("roles"), knownvalue.SetExact([]knownvalue.Check{
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

func testAccGroupRoleAssignmentsDataSourceConfig(groupName, roleName string) string {
	return testAccGroupRoleAssignmentResourceConfig(groupName, roleName) + `
data "openai_group_role_assignments" "test" {
	group_id = openai_group_role_assignment.test.group_id
}
`
}
