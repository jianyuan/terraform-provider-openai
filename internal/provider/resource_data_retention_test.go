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

func TestAccDataRetentionResource(t *testing.T) {
	rn := "openai_data_retention.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataRetentionResourceConfig("zero_data_retention"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("type"), knownvalue.StringExact("zero_data_retention")),
				},
			},
			{
				Config: testAccDataRetentionResourceConfig("enhanced_zero_data_retention"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("type"), knownvalue.StringExact("enhanced_zero_data_retention")),
				},
			},
		},
	})
}

func testAccDataRetentionResourceConfig(retentionType string) string {
	return fmt.Sprintf(`
resource "openai_data_retention" "test" {
	type = %[1]q
}
`, retentionType)
}
