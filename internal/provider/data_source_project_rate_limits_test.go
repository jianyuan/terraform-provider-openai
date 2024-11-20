package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jianyuan/terraform-provider-openai/internal/acctest"
)

func TestAccProjectRateLimitsDataSource(t *testing.T) {
	rn := "data.openai_project_rate_limits.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectRateLimitsDataSourceConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("project_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("rate_limits"), knownvalue.SetPartial([]knownvalue.Check{
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"id":                        knownvalue.NotNull(),
							"model":                     knownvalue.NotNull(),
							"max_requests_per_1_minute": knownvalue.NotNull(),
							"max_tokens_per_1_minute":   knownvalue.NotNull(),
						}),
					})),
				},
			},
		},
	})
}

var testAccProjectRateLimitsDataSourceConfig = `
data "openai_projects" "test" {
}

data "openai_project_rate_limits" "test" {
	project_id = tolist(data.openai_projects.test.projects)[0].id
}
`
