package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jianyuan/terraform-provider-openai/internal/acctest"
)

func TestAccProjectRateLimitResource(t *testing.T) {
	rn := "openai_project_rate_limit.test"
	projectName := acctest.RandomWithPrefix("tf-project")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectRateLimitResourceConfig(projectName, "text-embedding-3-small", 3, 150000),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs("openai_project.test", tfjsonpath.New("id"), rn, tfjsonpath.New("project_id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("model"), knownvalue.StringExact("text-embedding-3-small")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("max_requests_per_1_minute"), knownvalue.Int64Exact(3)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("max_tokens_per_1_minute"), knownvalue.Int64Exact(150000)),
				},
			},
			{
				Config: testAccProjectRateLimitResourceConfig(projectName, "text-embedding-3-small", 2, 75000),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs("openai_project.test", tfjsonpath.New("id"), rn, tfjsonpath.New("project_id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("model"), knownvalue.StringExact("text-embedding-3-small")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("max_requests_per_1_minute"), knownvalue.Int64Exact(2)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("max_tokens_per_1_minute"), knownvalue.Int64Exact(75000)),
				},
			},
		},
	})
}

func testAccProjectRateLimitResourceConfig(name, model string, maxRequestsPer1Minute, maxTokensPer1Minute int) string {
	return fmt.Sprintf(`
resource "openai_project" "test" {
	name = %[1]q
}

resource "openai_project_rate_limit" "test" {
	project_id = openai_project.test.id
	model      = %[2]q

	max_requests_per_1_minute = %[3]d
	max_tokens_per_1_minute   = %[4]d
}
`, name, model, maxRequestsPer1Minute, maxTokensPer1Minute)
}
