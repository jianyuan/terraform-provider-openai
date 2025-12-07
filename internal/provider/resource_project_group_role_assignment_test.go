package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jianyuan/terraform-provider-openai/internal/acctest"
	"github.com/jianyuan/terraform-provider-openai/internal/tfutils"
)

func TestAccProjectGroupRoleAssignmentResource(t *testing.T) {
	rn := "openai_project_group_role_assignment.test"
	roleName := acctest.RandomWithPrefix("tf-role")
	projectName := acctest.RandomWithPrefix("tf-project")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectGroupRoleAssignmentResourceConfig(projectName, acctest.TestGroupId, roleName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("project_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("group_id"), knownvalue.StringExact(acctest.TestGroupId)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("role_id"), knownvalue.NotNull()),
				},
			},
			{
				ResourceName: rn,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[rn]
					if !ok {
						return "", fmt.Errorf("not found: %s", rn)
					}
					projectId := rs.Primary.Attributes["project_id"]
					groupId := rs.Primary.Attributes["group_id"]
					roleId := rs.Primary.Attributes["role_id"]
					return tfutils.BuildThreePartId(projectId, groupId, roleId), nil
				},
			},
		},
	})
}

func testAccProjectGroupRoleAssignmentResourceConfig(projectName, groupId, roleName string) string {
	return testAccProjectRoleResourceConfig(projectName, roleName, "role dscription", `["api.organization.projects.api_keys.read"]`) + fmt.Sprintf(`
resource "openai_project_group_role_assignment" "test" {
	project_id = openai_project.test.id
	group_id   = %[1]q
	role_id    = openai_project_role.test.id
}
`, groupId)
}
