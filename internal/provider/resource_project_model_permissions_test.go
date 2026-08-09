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

func TestAccProjectModelPermissionsResource(t *testing.T) {
	rn := "openai_project_model_permissions.test"
	projectName := sdkacctest.RandomWithPrefix("tf-project")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectModelPermissionsResourceConfig(projectName, `
					mode      = "allow_list"
					model_ids = ["gpt-4", "gpt-4o"]
				`),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("project_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("mode"), knownvalue.StringExact("allow_list")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("model_ids"), knownvalue.SetExact([]knownvalue.Check{
						knownvalue.StringExact("gpt-4"),
						knownvalue.StringExact("gpt-4o"),
					})),
				},
			},
			{
				Config: testAccProjectModelPermissionsResourceConfig(projectName, `
					mode      = "deny_list"
					model_ids = ["o3"]
				`),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("project_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("mode"), knownvalue.StringExact("deny_list")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("model_ids"), knownvalue.SetExact([]knownvalue.Check{
						knownvalue.StringExact("o3"),
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
			},
		},
	})
}

func testAccProjectModelPermissionsResourceConfig(projectName, extra string) string {
	return fmt.Sprintf(`
resource "openai_project" "test" {
	name = %[1]q
}

resource "openai_project_model_permissions" "test" {
	project_id = openai_project.test.id
	%[2]s
}
`, projectName, extra)
}
