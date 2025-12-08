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

func TestAccGroupRoleAssignmentResource(t *testing.T) {
	rn := "openai_group_role_assignment.test"
	roleName := acctest.RandomWithPrefix("tf-role")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGroupRoleAssignmentResourceConfig(acctest.TestGroupId, roleName),
				ConfigStateChecks: []statecheck.StateCheck{
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
					groupId := rs.Primary.Attributes["group_id"]
					roleId := rs.Primary.Attributes["role_id"]
					return tfutils.BuildTwoPartId(groupId, roleId), nil
				},
			},
		},
	})
}

func testAccGroupRoleAssignmentResourceConfig(groupId, roleName string) string {
	return testAccOrganizationRoleResourceConfig(roleName, "role description", `["api.groups.read"]`) + fmt.Sprintf(`
resource "openai_group_role_assignment" "test" {
	group_id = %[1]q
	role_id  = openai_organization_role.test.id
}	
`, groupId)
}
