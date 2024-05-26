package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jianyuan/terraform-provider-openai/internal/acctest"
)

func TestAccProjectApiKeyResource_defaultProject(t *testing.T) {
	rn := "openai_project_api_key.test"
	serviceAccountId := acctest.RandomWithPrefix("tf-service-account")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectApiKeyResourceConfig_defaultProject(serviceAccountId, `name = "tf-api-key"`),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("organization_id"), knownvalue.StringExact(acctest.TestOrganizationId)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("project_id"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("service_account_id"), knownvalue.StringExact(serviceAccountId)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact("tf-api-key")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("scopes"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("created"), knownvalue.NotNull()),
				},
			},
			{
				Config: testAccProjectApiKeyResourceConfig_defaultProject(serviceAccountId, ""),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("organization_id"), knownvalue.StringExact(acctest.TestOrganizationId)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("project_id"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("service_account_id"), knownvalue.StringExact(serviceAccountId)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("scopes"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("created"), knownvalue.NotNull()),
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
					organizationId := rs.Primary.Attributes["organization_id"]
					id := rs.Primary.ID
					return BuildTwoPartId(organizationId, id), nil
				},
				ImportStateVerify: true,
			},
			{
				Config: testAccProjectApiKeyResourceConfig_defaultProject(serviceAccountId, `name = "tf-api-key-updated"`),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("organization_id"), knownvalue.StringExact(acctest.TestOrganizationId)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("project_id"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("service_account_id"), knownvalue.StringExact(serviceAccountId)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact("tf-api-key-updated")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("scopes"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("created"), knownvalue.NotNull()),
				},
			},
		},
	})
}

func TestAccProjectApiKeyResource_namedProject(t *testing.T) {
	rn := "openai_project_api_key.test"
	projectTitle := acctest.RandomWithPrefix("tf-project")
	serviceAccountId := acctest.RandomWithPrefix("tf-service-account")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectApiKeyResourceConfig_namedProject(projectTitle, serviceAccountId, `name = "tf-api-key"`),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("organization_id"), knownvalue.StringExact(acctest.TestOrganizationId)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("project_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("service_account_id"), knownvalue.StringExact(serviceAccountId)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact("tf-api-key")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("scopes"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("created"), knownvalue.NotNull()),
				},
			},
			{
				Config: testAccProjectApiKeyResourceConfig_namedProject(projectTitle, serviceAccountId, ""),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("organization_id"), knownvalue.StringExact(acctest.TestOrganizationId)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("project_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("service_account_id"), knownvalue.StringExact(serviceAccountId)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("scopes"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("created"), knownvalue.NotNull()),
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
					organizationId := rs.Primary.Attributes["organization_id"]
					projectId := rs.Primary.Attributes["project_id"]
					id := rs.Primary.ID
					return BuildThreePartId(organizationId, projectId, id), nil
				},
				ImportStateVerify: true,
			},
			{
				Config: testAccProjectApiKeyResourceConfig_namedProject(projectTitle, serviceAccountId, `name = "tf-api-key-updated"`),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("organization_id"), knownvalue.StringExact(acctest.TestOrganizationId)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("project_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("service_account_id"), knownvalue.StringExact(serviceAccountId)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact("tf-api-key-updated")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("scopes"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("created"), knownvalue.NotNull()),
				},
			},
		},
	})
}

func TestAccProjectApiKeyResource_defaultProject_addNameAndScopes(t *testing.T) {
	rn := "openai_project_api_key.test"
	serviceAccountId := acctest.RandomWithPrefix("tf-service-account")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectApiKeyResourceConfig_defaultProject(serviceAccountId, ""),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("scopes"), knownvalue.Null()),
				},
			},
			{
				Config: testAccProjectApiKeyResourceConfig_defaultProject(serviceAccountId, `
					name   = "tf-api-key"
					scopes = ["model.read", "api.model.read"]
				`),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact("tf-api-key")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("scopes"), knownvalue.SetExact([]knownvalue.Check{
						knownvalue.StringExact("model.read"),
						knownvalue.StringExact("api.model.read"),
					})),
				},
			},
		},
	})
}

func TestAccProjectApiKeyResource_namedProject_addNameAndScopes(t *testing.T) {
	rn := "openai_project_api_key.test"
	projectTitle := acctest.RandomWithPrefix("tf-project")
	serviceAccountId := acctest.RandomWithPrefix("tf-service-account")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectApiKeyResourceConfig_namedProject(projectTitle, serviceAccountId, ""),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("scopes"), knownvalue.Null()),
				},
			},
			{
				Config: testAccProjectApiKeyResourceConfig_namedProject(projectTitle, serviceAccountId, `
					name   = "tf-api-key"
					scopes = ["model.read", "api.model.read"]
				`),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact("tf-api-key")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("scopes"), knownvalue.SetExact([]knownvalue.Check{
						knownvalue.StringExact("model.read"),
						knownvalue.StringExact("api.model.read"),
					})),
				},
			},
		},
	})
}

func TestAccProjectApiKeyResource_defaultProject_removeNameAndScopes(t *testing.T) {
	rn := "openai_project_api_key.test"
	serviceAccountId := acctest.RandomWithPrefix("tf-service-account")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectApiKeyResourceConfig_defaultProject(serviceAccountId, `
					name   = "tf-api-key"
					scopes = ["model.read", "api.model.read"]
				`),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact("tf-api-key")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("scopes"), knownvalue.SetExact([]knownvalue.Check{
						knownvalue.StringExact("model.read"),
						knownvalue.StringExact("api.model.read"),
					})),
				},
			},
			{
				Config: testAccProjectApiKeyResourceConfig_defaultProject(serviceAccountId, ""),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("scopes"), knownvalue.Null()),
				},
			},
		},
	})
}

func TestAccProjectApiKeyResource_namedProject_removeNameAndScopes(t *testing.T) {
	rn := "openai_project_api_key.test"
	projectTitle := acctest.RandomWithPrefix("tf-project")
	serviceAccountId := acctest.RandomWithPrefix("tf-service-account")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectApiKeyResourceConfig_namedProject(projectTitle, serviceAccountId, `
					name   = "tf-api-key"
					scopes = ["model.read", "api.model.read"]
				`),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact("tf-api-key")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("scopes"), knownvalue.SetExact([]knownvalue.Check{
						knownvalue.StringExact("model.read"),
						knownvalue.StringExact("api.model.read"),
					})),
				},
			},
			{
				Config: testAccProjectApiKeyResourceConfig_namedProject(projectTitle, serviceAccountId, ""),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("scopes"), knownvalue.Null()),
				},
			},
		},
	})
}

func testAccProjectApiKeyResourceConfig_defaultProject(serviceAccountId, extras string) string {
	return testAccOrganizationDataSourceConfig + fmt.Sprintf(`
resource "openai_project_api_key" "test" {
  organization_id    = data.openai_organization.test.id
  service_account_id = %[1]q
  %[2]s
}
`, serviceAccountId, extras)
}

func testAccProjectApiKeyResourceConfig_namedProject(projectTitle, serviceAccountId, extras string) string {
	return testAccOrganizationDataSourceConfig + fmt.Sprintf(`
resource "openai_project" "test" {
  organization_id = data.openai_organization.test.id
  title           = %[1]q
}

resource "openai_project_api_key" "test" {
  organization_id    = openai_project.test.organization_id
  project_id         = openai_project.test.id
  service_account_id = %[2]q
  %[3]s
}
`, projectTitle, serviceAccountId, extras)
}
