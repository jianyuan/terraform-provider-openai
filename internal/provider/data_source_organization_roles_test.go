package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jianyuan/terraform-provider-openai/internal/acctest"
)

func TestAccOrganizationRolesDataSource(t *testing.T) {
	rn := "data.openai_organization_roles.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationRolesDataSourceConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("roles"), knownvalue.SetPartial([]knownvalue.Check{
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"id":              knownvalue.NotNull(),
							"name":            knownvalue.StringExact("owner"),
							"description":     knownvalue.StringExact("Can modify billing information and manage organization members"),
							"permissions":     knownvalue.NotNull(),
							"predefined_role": knownvalue.Bool(true),
							"resource_type":   knownvalue.StringExact("api.organization"),
						}),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"id":              knownvalue.NotNull(),
							"name":            knownvalue.StringExact("reader"),
							"description":     knownvalue.StringExact("Can make standard API requests and read basic organizational data"),
							"permissions":     knownvalue.NotNull(),
							"predefined_role": knownvalue.Bool(true),
							"resource_type":   knownvalue.StringExact("api.organization"),
						}),
					})),
				},
			},
		},
	})
}

var testAccOrganizationRolesDataSourceConfig = `
data "openai_organization_roles" "test" {
}
`
