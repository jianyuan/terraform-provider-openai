package provider

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jianyuan/go-utils/ptr"
	"github.com/jianyuan/terraform-provider-openai/internal/acctest"
	"github.com/jianyuan/terraform-provider-openai/internal/apiclient"
)

func init() {
	resource.AddTestSweepers("openai_project_service_account", &resource.Sweeper{
		Name: "openai_project_service_account",
		F: func(r string) error {
			ctx := context.Background()

			var projects []apiclient.Project

			// List all projects
			{
				params := &apiclient.ListProjectsParams{
					Limit: ptr.Ptr(100),
				}

				for {
					httpResp, err := acctest.SharedClient.ListProjectsWithResponse(
						ctx,
						params,
					)
					if err != nil {
						return fmt.Errorf("Unable to read, got error: %s", err)
					}

					if httpResp.StatusCode() != http.StatusOK {
						return fmt.Errorf("Unable to read, got status code %d: %s", httpResp.StatusCode(), string(httpResp.Body))
					}

					projects = append(projects, httpResp.JSON200.Data...)

					if !httpResp.JSON200.HasMore {
						break
					}

					params.After = &httpResp.JSON200.LastId
				}
			}

			for _, project := range projects {
				log.Printf("[INFO] Listing project service accounts for project %s", project.Id)

				var projectServiceAccounts []apiclient.ProjectServiceAccount
				params := &apiclient.ListProjectServiceAccountsParams{
					Limit: ptr.Ptr(100),
				}

				for {
					httpResp, err := acctest.SharedClient.ListProjectServiceAccountsWithResponse(
						ctx,
						project.Id,
						params,
					)

					if err != nil {
						return fmt.Errorf("Unable to read, got error: %s", err)
					}

					if httpResp.StatusCode() != http.StatusOK {
						return fmt.Errorf("Unable to read, got status code %d: %s", httpResp.StatusCode(), string(httpResp.Body))
					}

					for _, sa := range httpResp.JSON200.Data {
						if !strings.HasPrefix(sa.Id, "tf-") && !strings.HasPrefix(sa.Name, "test-") {
							continue
						}

						projectServiceAccounts = append(projectServiceAccounts, sa)
					}

					if !httpResp.JSON200.HasMore {
						break
					}

					params.After = &httpResp.JSON200.LastId
				}

				for _, sa := range projectServiceAccounts {
					log.Printf("[INFO] Destroying project service account %s", sa.Id)

					_, err := acctest.SharedClient.DeleteProjectServiceAccountWithResponse(
						ctx,
						project.Id,
						sa.Id,
					)

					if err != nil {
						log.Printf("[ERROR] Unable to delete project service account %s: %s", sa.Id, err)
						continue
					}

					log.Printf("[INFO] Deleted project service account %s", sa.Id)
				}
			}

			return nil
		},
	})
}

func TestAccProjectServiceAccountResource(t *testing.T) {
	rn := "openai_project_service_account.test"
	projectName := acctest.RandomWithPrefix("tf-project")
	projectServiceAccountName := acctest.RandomWithPrefix("tf-service-account")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectServiceAccountResourceConfig(projectName, projectServiceAccountName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs("openai_project.test", tfjsonpath.New("id"), rn, tfjsonpath.New("project_id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(projectServiceAccountName)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("role"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("api_key_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("api_key"), knownvalue.NotNull()),
					statecheck.ExpectKnownOutputValue("service_account_api_key", knownvalue.NotNull()),
				},
			},
			{
				Config: testAccProjectServiceAccountResourceConfig(projectName, projectServiceAccountName+"-changed"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs("openai_project.test", tfjsonpath.New("id"), rn, tfjsonpath.New("project_id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(projectServiceAccountName+"-changed")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("role"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("api_key_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("api_key"), knownvalue.NotNull()),
					statecheck.ExpectKnownOutputValue("service_account_api_key", knownvalue.NotNull()),
				},
			},
		},
	})
}

func testAccProjectServiceAccountResourceConfig(projectName, projectServiceAccountName string) string {
	return fmt.Sprintf(`
resource "openai_project" "test" {
	name = %[1]q
}

resource "openai_project_service_account" "test" {
	project_id = openai_project.test.id
	name       = %[2]q
}

output "service_account_api_key" {
	sensitive = true
	value     = openai_project_service_account.test.api_key
}

`, projectName, projectServiceAccountName)
}
