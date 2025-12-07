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

func TestAccProjectRolesDataSource(t *testing.T) {
	rn := "data.openai_project_roles.test"
	projectName := acctest.RandomWithPrefix("tf-project")
	roleName := acctest.RandomWithPrefix("tf-role")
	roleDescription := acctest.RandomWithPrefix("tf-role-description")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectRolesDataSourceConfig(projectName, roleName, roleDescription),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("roles"), knownvalue.SetPartial([]knownvalue.Check{
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"id":          knownvalue.NotNull(),
							"name":        knownvalue.StringExact(roleName),
							"description": knownvalue.StringExact(roleDescription),
							"permissions": knownvalue.SetExact([]knownvalue.Check{
								knownvalue.StringExact("api.organization.projects.api_keys.read"),
								knownvalue.StringExact("api.organization.projects.api_keys.write"),
							}),
							"predefined_role": knownvalue.Bool(false),
							"resource_type":   knownvalue.StringExact("api.project"),
						}),
					})),
				},
			},
		},
	})
}

func testAccProjectRolesDataSourceConfig(projectName, roleName, roleDescription string) string {
	return fmt.Sprintf(`
resource "openai_project" "test" {
	name = %[1]q
}

resource "openai_project_role" "test" {
	project_id  = openai_project.test.id
	name        = %[2]q
	description = %[3]q
	permissions = [
		"api.organization.projects.api_keys.read",
		"api.organization.projects.api_keys.write",
	]
}

data "openai_project_roles" "test" {
	depends_on = [openai_project_role.test]
	project_id = openai_project.test.id
}
`, projectName, roleName, roleDescription)
}
