package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jianyuan/terraform-provider-openai/internal/acctest"
)

func TestAccOrganizationsDataSource(t *testing.T) {
	rn := "data.openai_organizations.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationsDataSourceConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("organizations"), knownvalue.SetPartial([]knownvalue.Check{
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"id":          knownvalue.StringExact(acctest.TestOrganizationId),
							"is_default":  knownvalue.Bool(true),
							"name":        knownvalue.NotNull(),
							"title":       knownvalue.NotNull(),
							"description": knownvalue.NotNull(),
						}),
					})),
				},
			},
		},
	})
}

const testAccOrganizationsDataSourceConfig = `
data "openai_organizations" "test" {
}
`
