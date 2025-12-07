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

func TestAccProjectRoleResource(t *testing.T) {
	rn := "openai_project_role.test"
	projectName := acctest.RandomWithPrefix("tf-project")
	roleName := acctest.RandomWithPrefix("tf-role")
	roleDescription := acctest.RandomWithPrefix("tf-role-description")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectRoleResourceConfig(projectName, roleName, roleDescription, `["api.organization.projects.api_keys.read"]`),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("project_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(roleName)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("description"), knownvalue.StringExact(roleDescription)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("permissions"), knownvalue.SetExact([]knownvalue.Check{
						knownvalue.StringExact("api.organization.projects.api_keys.read"),
					})),
				},
			},
			{
				Config: testAccProjectRoleResourceConfig(projectName, roleName+"-updated", roleDescription+"-updated", `["api.organization.projects.api_keys.read", "api.organization.projects.api_keys.write"]`),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("project_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(roleName+"-updated")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("description"), knownvalue.StringExact(roleDescription+"-updated")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("permissions"), knownvalue.SetExact([]knownvalue.Check{
						knownvalue.StringExact("api.organization.projects.api_keys.read"),
						knownvalue.StringExact("api.organization.projects.api_keys.write"),
					})),
				},
			},
		},
	})
}

func testAccProjectRoleResourceConfig(projectName, name, description, permissions string) string {
	return fmt.Sprintf(`
resource "openai_project" "test" {
	name = %[1]q
}

resource "openai_project_role" "test" {
	project_id = openai_project.test.id
	name        = %[2]q
	description = %[3]q
	permissions = %[4]s
}
`, projectName, name, description, permissions)
}
