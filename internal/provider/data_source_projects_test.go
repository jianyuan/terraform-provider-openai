package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jianyuan/terraform-provider-openai/internal/acctest"
)

func TestAccProjectsDataSource(t *testing.T) {
	rn := "data.openai_projects.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectsDataSourceConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("organization_id"), knownvalue.StringExact(acctest.TestOrganizationId)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("projects"), knownvalue.SetPartial([]knownvalue.Check{
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"id":    knownvalue.NotNull(),
							"title": knownvalue.NotNull(),
						}),
					})),
				},
			},
		},
	})
}

var testAccProjectsDataSourceConfig = testAccOrganizationDataSourceConfig + `
data "openai_projects" "test" {
  organization_id = data.openai_organization.test.id
}
`
