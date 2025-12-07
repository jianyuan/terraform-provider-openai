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

func TestAccGroupRoleResource(t *testing.T) {
	rn := "openai_group_role.test"
	groupName := acctest.RandomWithPrefix("tf-group")
	roleName := acctest.RandomWithPrefix("tf-role")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGroupRoleResourceConfig(groupName, roleName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("group_id"), knownvalue.NotNull()),
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

func testAccGroupRoleResourceConfig(groupName, roleName string) string {
	return testAccGroupResourceConfig(groupName) + testAccOrganizationRoleResourceConfig(roleName, "role description", `["api.groups.read"]`) + `
resource "openai_group_role" "test" {
	group_id = openai_group.test.id
	role_id  = openai_organization_role.test.id
}
`
}
