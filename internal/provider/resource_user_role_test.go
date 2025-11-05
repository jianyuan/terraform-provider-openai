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

func TestAccUserRoleResource(t *testing.T) {
	rn := "openai_user_role.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Import existing user role
			{
				Config:        testAccUserRoleResourceConfig(acctest.TestUserId, "owner"),
				ResourceName:  rn,
				ImportState:   true,
				ImportStateId: acctest.TestUserId,
			},
			{
				Config: testAccUserRoleResourceConfig(acctest.TestUserId, "owner"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("user_id"), knownvalue.StringExact(acctest.TestUserId)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("role"), knownvalue.StringExact("owner")),
				},
			},
			{
				Config: testAccUserRoleResourceConfig(acctest.TestUserId, "reader"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("user_id"), knownvalue.StringExact(acctest.TestUserId)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("role"), knownvalue.StringExact("reader")),
				},
			},
			{
				// Detach state to prevent deletion of user
				Config: testAccUserRoleDetachState(),
			},
		},
	})
}

func testAccUserRoleResourceConfig(userId, role string) string {
	return fmt.Sprintf(`
resource "openai_user_role" "test" {
	user_id = %[1]q
	role    = %[2]q

	lifecycle {
		prevent_destroy = true
	}
}
`, userId, role)
}

func testAccUserRoleDetachState() string {
	return `
removed {
  from = openai_user_role.test

  lifecycle {
    destroy = false
  }
}
`
}
