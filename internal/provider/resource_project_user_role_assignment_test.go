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

func TestAccProjectUserRoleAssignmentResource(t *testing.T) {
	rn := "openai_project_user_role_assignment.test"
	roleName := acctest.RandomWithPrefix("tf-role")
	projectName := acctest.RandomWithPrefix("tf-project")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectUserRoleAssignmentResourceConfig(projectName, roleName, acctest.TestUserId),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("project_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("user_id"), knownvalue.StringExact(acctest.TestUserId)),
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
					userId := rs.Primary.Attributes["user_id"]
					roleId := rs.Primary.Attributes["role_id"]
					return tfutils.BuildThreePartId(projectId, userId, roleId), nil
				},
			},
		},
	})
}

func testAccProjectUserRoleAssignmentResourceConfig(projectName, roleName, userId string) string {
	return testAccProjectRoleResourceConfig(projectName, roleName, "role description", `["api.organization.projects.api_keys.read"]`) + fmt.Sprintf(`
resource "openai_project_user" "test" {
	project_id = openai_project.test.id
	user_id    = %[1]q
	role       = "member"
}

resource "openai_project_user_role_assignment" "test" {
	project_id = openai_project_user.test.project_id
	user_id    = openai_project_user.test.user_id
	role_id    = openai_project_role.test.id
}
`, userId)
}
