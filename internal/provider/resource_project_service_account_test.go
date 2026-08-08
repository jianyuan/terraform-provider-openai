package provider_test

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/compare"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jianyuan/terraform-provider-openai/internal/acctest"
	"github.com/jianyuan/terraform-provider-openai/internal/provider"
	"github.com/jianyuan/terraform-provider-openai/internal/sweep"
	"github.com/openai/openai-go/v3"
)

func init() {
	sweep.Register("openai_project_service_account", func(ctx context.Context, client *openai.Client) ([]sweep.Sweepable, error) {
		var sweepables []sweep.Sweepable

		params := openai.AdminOrganizationProjectListParams{
			Limit: openai.Int(100),
		}
		var projectIds []string
		iter := acctest.SharedClient.Admin.Organization.Projects.ListAutoPaging(ctx, params)
		for iter.Next() {
			item := iter.Current()
			projectIds = append(projectIds, item.ID)
		}

		for _, projectId := range projectIds {
			log.Printf("[INFO] Listing project service accounts for project %s", projectId)

			params := openai.AdminOrganizationProjectServiceAccountListParams{
				Limit: openai.Int(100),
			}

			var ids []string

			iter := acctest.SharedClient.Admin.Organization.Projects.ServiceAccounts.ListAutoPaging(ctx, projectId, params)
			for iter.Next() {
				item := iter.Current()
				if strings.HasPrefix(item.Name, "tf-") || strings.HasPrefix(item.Name, "test-") {
					ids = append(ids, item.ID)
				}
			}

			for _, id := range ids {
				sweepables = append(sweepables, sweep.NewSweepResource(provider.NewProjectServiceAccountResource, acctest.SharedClient, map[string]any{
					"project_id": projectId,
					"id":         id,
				}))
			}
		}

		return sweepables, nil
	})
}

func TestAccProjectServiceAccountResource(t *testing.T) {
	rn := "openai_project_service_account.test"
	projectName := sdkacctest.RandomWithPrefix("tf-project")
	projectServiceAccountName := sdkacctest.RandomWithPrefix("tf-service-account")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
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
