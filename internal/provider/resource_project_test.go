package provider_test

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jianyuan/go-utils/ptr"
	"github.com/jianyuan/terraform-provider-openai/internal/acctest"
	"github.com/jianyuan/terraform-provider-openai/internal/apiclient"
)

func init() {
	resource.AddTestSweepers("openai_project", &resource.Sweeper{
		Name: "openai_project",
		F: func(r string) error {
			ctx := context.Background()

			params := &apiclient.ListProjectsParams{
				Limit: ptr.Ptr(int64(100)),
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

				for _, project := range httpResp.JSON200.Data {
					if !strings.HasPrefix(project.Name, "tf-") {
						continue
					}

					log.Printf("[INFO] Destroying project %s", project.Id)

					_, err := acctest.SharedClient.ArchiveProjectWithResponse(
						ctx,
						project.Id,
					)

					if err != nil {
						log.Printf("[ERROR] Unable to archive project %s: %s", project.Id, err)
						continue
					}

					log.Printf("[INFO] Archived project %s", project.Id)
				}

				if !httpResp.JSON200.HasMore {
					break
				}

				params.After = &httpResp.JSON200.LastId
			}

			return nil
		},
	})
}

func TestAccProjectResource(t *testing.T) {
	rn := "openai_project.test"
	projectName := acctest.RandomWithPrefix("tf-project")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectResourceConfig(projectName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(projectName)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("status"), knownvalue.StringExact("active")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("archived_at"), knownvalue.Null()),
				},
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccProjectResourceConfig(projectName + "-updated"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(projectName+"-updated")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("status"), knownvalue.StringExact("active")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("archived_at"), knownvalue.Null()),
				},
			},
		},
	})
}

func testAccProjectResourceConfig(name string) string {
	return fmt.Sprintf(`
resource "openai_project" "test" {
	name = %[1]q
}
`, name)
}
