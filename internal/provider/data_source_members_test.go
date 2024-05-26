package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jianyuan/terraform-provider-openai/internal/acctest"
)

func TestAccMembersDataSource(t *testing.T) {
	rn := "data.openai_members.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMembersDataSourceConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("organization_id"), knownvalue.StringExact(acctest.TestOrganizationId)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("invited_members"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("members"), knownvalue.SetPartial([]knownvalue.Check{
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"id":                 knownvalue.NotNull(),
							"email":              knownvalue.NotNull(),
							"name":               knownvalue.NotNull(),
							"picture":            knownvalue.NotNull(),
							"is_default":         knownvalue.Bool(true),
							"is_service_account": knownvalue.Bool(false),
							"role":               knownvalue.StringExact("owner"),
						}),
					})),
				},
			},
		},
	})
}

var testAccMembersDataSourceConfig = testAccOrganizationDataSourceConfig + `
data "openai_members" "test" {
  organization_id = data.openai_organization.test.id
}
`
