package provider_test

import (
	"fmt"
	"testing"

	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jianyuan/terraform-provider-openai/internal/acctest"
)

func TestAccProjectSpendLimitResource(t *testing.T) {
	rn := "openai_project_spend_limit.test"
	projectName := sdkacctest.RandomWithPrefix("tf-project")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccProjectSpendLimitResourceConfig(projectName, 0),
				ExpectError: acctest.ExpectLiteralError(`Attribute threshold_amount value must be at least 1, got: 0`),
			},
			{
				Config: testAccProjectSpendLimitResourceConfig(projectName, 10),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("project_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("currency"), knownvalue.StringExact("USD")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("interval"), knownvalue.StringExact("month")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("threshold_amount"), knownvalue.Int32Exact(10)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("enforcement"), knownvalue.ObjectExact(map[string]knownvalue.Check{
						"status": knownvalue.StringExact("enforcing"),
					})),
				},
			},
			{
				Config: testAccProjectSpendLimitResourceConfig(projectName, 100),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("project_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("currency"), knownvalue.StringExact("USD")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("interval"), knownvalue.StringExact("month")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("threshold_amount"), knownvalue.Int32Exact(100)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("enforcement"), knownvalue.ObjectExact(map[string]knownvalue.Check{
						"status": knownvalue.StringExact("enforcing"),
					})),
				},
			},
			{
				ResourceName: rn,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[rn]
					if !ok {
						return "", fmt.Errorf("not found: %s", rn)
					}
					projectId := rs.Primary.Attributes["project_id"]
					return projectId, nil
				},
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "project_id",
			},
		},
	})
}

func testAccProjectSpendLimitResourceConfig(projectName string, amount int) string {
	return testAccProjectResourceConfig(projectName) + fmt.Sprintf(`
resource "openai_project_spend_limit" "test" {
	project_id       = openai_project.test.id
	currency         = "USD"
	interval         = "month"
	threshold_amount = %[1]d
}
`, amount)
}
