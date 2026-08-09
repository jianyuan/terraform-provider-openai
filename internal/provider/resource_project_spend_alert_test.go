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
	"github.com/jianyuan/terraform-provider-openai/internal/tfutils"
)

func TestAccProjectSpendAlertResource(t *testing.T) {
	rn := "openai_project_spend_alert.test"
	projectName := sdkacctest.RandomWithPrefix("tf-project")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectSpendAlertResourceConfig(projectName, 10, `
					recipients = ["a@example.com"]
				`),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("project_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("currency"), knownvalue.StringExact("USD")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("interval"), knownvalue.StringExact("month")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("threshold_amount"), knownvalue.Int32Exact(10)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("notification_channel").AtMapKey("type"), knownvalue.StringExact("email")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("notification_channel").AtMapKey("recipients"), knownvalue.SetExact([]knownvalue.Check{knownvalue.StringExact("a@example.com")})),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("notification_channel").AtMapKey("subject_prefix"), knownvalue.Null()),
				},
			},
			{
				Config: testAccProjectSpendAlertResourceConfig(projectName, 10, `
					recipients = ["a@example.com", "b@example.com"]
					subject_prefix = "OpenAI Terraform"
				`),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("project_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("currency"), knownvalue.StringExact("USD")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("interval"), knownvalue.StringExact("month")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("threshold_amount"), knownvalue.Int32Exact(10)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("notification_channel").AtMapKey("type"), knownvalue.StringExact("email")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("notification_channel").AtMapKey("recipients"), knownvalue.SetExact([]knownvalue.Check{knownvalue.StringExact("a@example.com"), knownvalue.StringExact("b@example.com")})),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("notification_channel").AtMapKey("subject_prefix"), knownvalue.StringExact("OpenAI Terraform")),
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
					alertId := rs.Primary.Attributes["id"]
					return tfutils.BuildTwoPartId(projectId, alertId), nil
				},
				ImportStateVerify: true,
			},
		},
	})
}

func testAccProjectSpendAlertResourceConfig(projectName string, amount int, notificationChannelExtra string) string {
	return fmt.Sprintf(`
resource "openai_project" "test" {
	name = %[1]q
}

resource "openai_project_spend_alert" "test" {
	project_id       = openai_project.test.id
	currency         = "USD"
	interval         = "month"
	threshold_amount = %[2]d
	notification_channel = {
		type = "email"
		%[3]s
	}
}
`, projectName, amount, notificationChannelExtra)
}
