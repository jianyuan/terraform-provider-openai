package provider

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jianyuan/terraform-provider-openai/internal/acctest"
	"github.com/jianyuan/terraform-provider-openai/internal/apiclient"
	"github.com/jianyuan/terraform-provider-openai/internal/ptr"
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
