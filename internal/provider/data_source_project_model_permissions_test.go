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

func TestAccProjectModelPermissionsDataSource(t *testing.T) {
	rn := "data.openai_project_model_permissions.test"
	projectName := sdkacctest.RandomWithPrefix("tf-project")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectModelPermissionsDataSourceConfig(projectName, `
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
		},
	})
}

func testAccProjectModelPermissionsDataSourceConfig(projectName, extra string) string {
	return testAccProjectModelPermissionsResourceConfig(projectName, extra) + `
data "openai_project_model_permissions" "test" {
	project_id = openai_project_model_permissions.test.project_id
}
`
}
