package provider_test

import (
	"testing"

	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jianyuan/terraform-provider-openai/internal/acctest"
)

func TestAccProjectSpendLimitDataSource(t *testing.T) {
	rn := "data.openai_project_spend_limit.test"
	projectName := sdkacctest.RandomWithPrefix("tf-project")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectSpendLimitDataSourceConfig(projectName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("project_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("currency"), knownvalue.StringExact("USD")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("interval"), knownvalue.StringExact("month")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("threshold_amount"), knownvalue.Int64Exact(100)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("enforcement"), knownvalue.ObjectExact(map[string]knownvalue.Check{
						"status": knownvalue.StringExact("enforcing"),
					})),
				},
			},
		},
	})
}

func testAccProjectSpendLimitDataSourceConfig(projectName string) string {
	return testAccProjectResourceConfig(projectName) + `
resource "openai_project_spend_limit" "test" {
	project_id = openai_project.test.id
	currency = "USD"
	interval = "month"
	threshold_amount = 100
}

data "openai_project_spend_limit" "test" {
	project_id = openai_project_spend_limit.test.project_id
}
`
}
