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

func TestAccSpendLimitResource(t *testing.T) {
	rn := "openai_spend_limit.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccSpendLimitResourceConfig(0),
				ExpectError: acctest.ExpectLiteralError(`Attribute threshold_amount value must be at least 1, got: 0`),
			},
			{
				Config: testAccSpendLimitResourceConfig(10),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("currency"), knownvalue.StringExact("USD")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("interval"), knownvalue.StringExact("month")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("threshold_amount"), knownvalue.Int32Exact(10)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("enforcement"), knownvalue.ObjectExact(map[string]knownvalue.Check{
						"status": knownvalue.StringExact("enforcing"),
					})),
				},
			},
			{
				Config: testAccSpendLimitResourceConfig(100),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("currency"), knownvalue.StringExact("USD")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("interval"), knownvalue.StringExact("month")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("threshold_amount"), knownvalue.Int32Exact(100)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("enforcement"), knownvalue.ObjectExact(map[string]knownvalue.Check{
						"status": knownvalue.StringExact("enforcing"),
					})),
				},
			},
		},
	})
}

func testAccSpendLimitResourceConfig(amount int) string {
	return fmt.Sprintf(`
resource "openai_spend_limit" "test" {
	currency         = "USD"
	interval         = "month"
	threshold_amount = %[1]d
}
`, amount)
}
