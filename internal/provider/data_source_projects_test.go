package provider_test

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
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectsDataSourceConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("projects"), knownvalue.SetPartial([]knownvalue.Check{
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"id":          knownvalue.NotNull(),
							"name":        knownvalue.StringExact("Default project"),
							"status":      knownvalue.StringExact("active"),
							"created_at":  knownvalue.NotNull(),
							"archived_at": knownvalue.Null(),
						}),
					})),
				},
			},
		},
	})
}

func TestAccProjectsDataSource_includeArchived(t *testing.T) {
	rn := "data.openai_projects.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectsDataSourceConfig_includeArchived,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("projects"), knownvalue.SetPartial([]knownvalue.Check{
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"id":          knownvalue.NotNull(),
							"name":        knownvalue.StringExact("Default project"),
							"status":      knownvalue.StringExact("archived"),
							"created_at":  knownvalue.NotNull(),
							"archived_at": knownvalue.NotNull(),
						}),
					})),
				},
			},
		},
	})
}

var testAccProjectsDataSourceConfig = `
data "openai_projects" "test" {
}
`

var testAccProjectsDataSourceConfig_includeArchived = `
data "openai_projects" "test" {
  include_archived = true

	limit = 10
}
`
