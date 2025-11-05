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

func TestAccInvitesDataSource(t *testing.T) {
	email := fmt.Sprintf("tf-%d@example.com", acctest.RandInt())
	rn := "data.openai_invites.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInvitesDataSourceConfig(email),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("invites"), knownvalue.SetPartial([]knownvalue.Check{
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"id":          knownvalue.NotNull(),
							"email":       knownvalue.StringExact(email),
							"role":        knownvalue.StringExact("reader"),
							"status":      knownvalue.NotNull(),
							"invited_at":  knownvalue.NotNull(),
							"expires_at":  knownvalue.NotNull(),
							"accepted_at": knownvalue.Null(),
						}),
					})),
				},
			},
		},
	})
}

func testAccInvitesDataSourceConfig(email string) string {
	return fmt.Sprintf(`
resource "openai_invite" "test" {
	email = %[1]q
	role  = "reader"
}

data "openai_invites" "test" {
	depends_on = [openai_invite.test]
}
`, email)

}
