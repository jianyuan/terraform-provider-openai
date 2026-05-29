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

func TestAccSpendAlertResource(t *testing.T) {
	rn := "openai_spend_alert.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSpendAlertResourceConfig(10, `
					recipients = ["a@example.com"]
				`),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("currency"), knownvalue.StringExact("USD")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("interval"), knownvalue.StringExact("month")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("threshold_amount"), knownvalue.Int32Exact(10)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("notification_channel").AtMapKey("type"), knownvalue.StringExact("email")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("notification_channel").AtMapKey("recipients"), knownvalue.SetExact([]knownvalue.Check{knownvalue.StringExact("a@example.com")})),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("notification_channel").AtMapKey("subject_prefix"), knownvalue.Null()),
				},
			},
			{
				Config: testAccSpendAlertResourceConfig(10, `
					recipients = ["a@example.com", "b@example.com"]
					subject_prefix = "OpenAI Terraform"
				`),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("currency"), knownvalue.StringExact("USD")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("interval"), knownvalue.StringExact("month")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("threshold_amount"), knownvalue.Int32Exact(10)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("notification_channel").AtMapKey("type"), knownvalue.StringExact("email")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("notification_channel").AtMapKey("recipients"), knownvalue.SetExact([]knownvalue.Check{knownvalue.StringExact("a@example.com"), knownvalue.StringExact("b@example.com")})),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("notification_channel").AtMapKey("subject_prefix"), knownvalue.StringExact("OpenAI Terraform")),
				},
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSpendAlertResourceConfig(amount int, notificationChannelExtra string) string {
	return fmt.Sprintf(`
resource "openai_spend_alert" "test" {
	currency         = "USD"
	interval         = "month"
	threshold_amount = %[1]d
	notification_channel = {
		type = "email"
		%[2]s
	}
}
`, amount, notificationChannelExtra)
}
