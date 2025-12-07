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

func TestAccGroupUserResource(t *testing.T) {
	rn := "openai_group_user.test"
	groupName := acctest.RandomWithPrefix("tf-group")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGroupUserResourceConfig(groupName, acctest.TestUserId),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("group_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("user_id"), knownvalue.StringExact(acctest.TestUserId)),
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
					userId := rs.Primary.Attributes["user_id"]
					return tfutils.BuildTwoPartId(groupId, userId), nil
				},
			},
		},
	})
}

func testAccGroupUserResourceConfig(groupName, userId string) string {
	return testAccGroupResourceConfig(groupName) + fmt.Sprintf(`
resource "openai_group_user" "test" {
	group_id = openai_group.test.id
	user_id  = %[1]q
}
`, userId)
}
