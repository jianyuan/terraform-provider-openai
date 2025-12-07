package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jianyuan/terraform-provider-openai/internal/acctest"
)

func TestAccGroupRoleAssignmentsDataSource(t *testing.T) {
	rn := "data.openai_group_role_assignments.test"
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
				Config: testAccGroupRoleAssignmentsDataSourceConfig(acctest.TestGroupId, roleName),
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

func testAccGroupRoleAssignmentsDataSourceConfig(groupId, roleName string) string {
	return testAccOrganizationRoleResourceConfig(roleName, "role description", `["api.groups.read"]`) + fmt.Sprintf(`
resource "openai_group_role_assignment" "test" {
	group_id = %[1]q
	role_id  = openai_organization_role.test.id
}

resource "time_sleep" "wait" {
	create_duration = "5s"

	triggers = {
		group_id = openai_group_role_assignment.test.group_id
	}
}

data "openai_group_role_assignments" "test" {
	group_id = time_sleep.wait.triggers.group_id
}
`, groupId)
}
