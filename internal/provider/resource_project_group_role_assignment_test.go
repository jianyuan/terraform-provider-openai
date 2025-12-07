package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jianyuan/terraform-provider-openai/internal/acctest"
)

func TestAccProjectGroupRoleAssignmentResource(t *testing.T) {
	rn := "openai_project_group_role_assignment.test"
	roleName := acctest.RandomWithPrefix("tf-role")
	groupName := acctest.RandomWithPrefix("tf-group")
	projectName := acctest.RandomWithPrefix("tf-project")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectGroupRoleAssignmentResourceConfig(roleName, groupName, projectName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("project_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("group_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("role_id"), knownvalue.NotNull()),
				},
			},
		},
	})
}

func testAccProjectGroupRoleAssignmentResourceConfig(roleName, groupName, projectName string) string {
	return testAccProjectRoleResourceConfig(projectName, roleName, "role dscription", `["api.organization.projects.api_keys.read"]`) + testAccGroupResourceConfig(groupName) + `
resource "openai_project_group_role_assignment" "test" {
	project_id = openai_project.test.id
	group_id   = openai_group.test.id
	role_id    = openai_project_role.test.id
}
`
}
