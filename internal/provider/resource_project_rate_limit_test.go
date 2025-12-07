package provider_test

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
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectRateLimitResourceConfig(projectName, "rl-text-embedding-3-small", 3, 3),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs("openai_project.test", tfjsonpath.New("id"), rn, tfjsonpath.New("project_id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("rate_limit_id"), knownvalue.StringExact("rl-text-embedding-3-small")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("max_requests_per_1_minute"), knownvalue.Int64Exact(3)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("max_tokens_per_1_minute"), knownvalue.Int64Exact(3)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("max_images_per_1_minute"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("max_audio_megabytes_per_1_minute"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("max_requests_per_1_day"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("batch_1_day_max_input_tokens"), knownvalue.Null()),
				},
			},
			{
				Config: testAccProjectRateLimitResourceConfig(projectName, "rl-text-embedding-3-small", 2, 2),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs("openai_project.test", tfjsonpath.New("id"), rn, tfjsonpath.New("project_id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("rate_limit_id"), knownvalue.StringExact("rl-text-embedding-3-small")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("max_requests_per_1_minute"), knownvalue.Int64Exact(2)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("max_tokens_per_1_minute"), knownvalue.Int64Exact(2)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("max_images_per_1_minute"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("max_audio_megabytes_per_1_minute"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("max_requests_per_1_day"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("batch_1_day_max_input_tokens"), knownvalue.Null()),
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
	project_id    = openai_project.test.id
	rate_limit_id = %[2]q

	max_requests_per_1_minute = %[3]d
	max_tokens_per_1_minute   = %[4]d
}
`, name, model, maxRequestsPer1Minute, maxTokensPer1Minute)
}
