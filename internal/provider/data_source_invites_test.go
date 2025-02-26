package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jianyuan/terraform-provider-openai/internal/acctest"
)

func TestAccInvitesDataSource(t *testing.T) {
	rn := "data.openai_invites.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInvitesDataSourceConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("invites"), knownvalue.SetPartial([]knownvalue.Check{
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"id":          knownvalue.NotNull(),
							"email":       knownvalue.NotNull(),
							"role":        knownvalue.StringExact("reader"),
							"status":      knownvalue.StringExact("expired"),
							"invited_at":  knownvalue.NotNull(),
							"expires_at":  knownvalue.NotNull(),
							"accepted_at": knownvalue.Null(),
						}),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"id":          knownvalue.NotNull(),
							"email":       knownvalue.NotNull(),
							"role":        knownvalue.StringExact("reader"),
							"status":      knownvalue.StringExact("accepted"),
							"invited_at":  knownvalue.NotNull(),
							"expires_at":  knownvalue.NotNull(),
							"accepted_at": knownvalue.NotNull(),
						}),
					})),
				},
			},
		},
	})
}

var testAccInvitesDataSourceConfig = `
data "openai_invites" "test" {
}
`
