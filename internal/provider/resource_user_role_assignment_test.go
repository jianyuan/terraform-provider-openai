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

func TestAccUserRoleAssignmentResource(t *testing.T) {
	rn := "openai_user_role_assignment.test"
	roleName := acctest.RandomWithPrefix("tf-role")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUserRoleAssignmentResourceConfig(acctest.TestUserId, roleName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("user_id"), knownvalue.NotNull()),
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
					userId := rs.Primary.Attributes["user_id"]
					roleId := rs.Primary.Attributes["role_id"]
					return tfutils.BuildTwoPartId(userId, roleId), nil
				},
			},
		},
	})
}

func testAccUserRoleAssignmentResourceConfig(userId, roleName string) string {
	return testAccOrganizationRoleResourceConfig(roleName, "role description", `["api.groups.read"]`) + fmt.Sprintf(`
resource "openai_user_role_assignment" "test" {
	user_id = %q
	role_id  = openai_organization_role.test.id
}
`, userId)
}
