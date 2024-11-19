package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jianyuan/terraform-provider-openai/internal/acctest"
)

func TestAccInviteResource(t *testing.T) {
	rn := "openai_invite.test"
	email := fmt.Sprintf("tf-%d@example.com", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInviteResourceConfig(email, "reader"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("email"), knownvalue.StringExact(email)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("role"), knownvalue.StringExact("reader")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("invited_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("expires_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("accepted_at"), knownvalue.Null()),
				},
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccInviteResourceConfig(email, "owner"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("email"), knownvalue.StringExact(email)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("role"), knownvalue.StringExact("owner")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("invited_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("expires_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("accepted_at"), knownvalue.Null()),
				},
			},
			{
				Config: testAccInviteResourceConfig("changed-"+email, "owner"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("email"), knownvalue.StringExact("changed-"+email)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("role"), knownvalue.StringExact("owner")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("invited_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("expires_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("accepted_at"), knownvalue.Null()),
				},
			},
		},
	})
}

func testAccInviteResourceConfig(email, role string) string {
	return fmt.Sprintf(`
resource "openai_invite" "test" {
	email = %[1]q
	role  = %[2]q
}
`, email, role)
}
