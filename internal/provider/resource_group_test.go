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
	resource.AddTestSweepers("openai_group", &resource.Sweeper{
		Name: "openai_group",
		F: func(r string) error {
			ctx := context.Background()

			params := &apiclient.ListGroupsParams{
				Limit: ptr.Ptr(int64(100)),
			}

			for {
				httpResp, err := acctest.SharedClient.ListGroupsWithResponse(
					ctx,
					params,
				)

				if err != nil {
					return fmt.Errorf("[ERROR] Unable to read, got error: %s", err)
				} else if httpResp.StatusCode() != http.StatusOK || httpResp.JSON200 == nil {
					return fmt.Errorf("[ERROR] Unable to read, got status code %d: %s", httpResp.StatusCode(), string(httpResp.Body))
				}

				for _, group := range httpResp.JSON200.Data {
					if !strings.HasPrefix(group.Name, "tf-") {
						continue
					}

					log.Printf("[INFO] Found group %s (ID: %s)", group.Name, group.Id)

					httpResp, err := acctest.SharedClient.DeleteGroupWithResponse(
						ctx,
						group.Id,
					)

					if err != nil {
						log.Printf("[ERROR] Unable to delete, got error: %s", err)
						continue
					} else if httpResp.StatusCode() != http.StatusOK {
						log.Printf("[ERROR] Unable to delete, got status code %d: %s", httpResp.StatusCode(), string(httpResp.Body))
						continue
					}

					log.Printf("[INFO] Deleted group %s (ID: %s)", group.Name, group.Id)
				}

				if httpResp.JSON200.Next == nil {
					break
				}

				params.After = httpResp.JSON200.Next
			}

			return nil
		},
	})
}

func TestAccGroupResource(t *testing.T) {
	rn := "openai_group.test"
	groupName := acctest.RandomWithPrefix("tf-group")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGroupResourceConfig(groupName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(groupName)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("created_at"), knownvalue.NotNull()),
				},
			},
			{
				Config: testAccGroupResourceConfig(groupName + "-updated"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(groupName+"-updated")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("created_at"), knownvalue.NotNull()),
				},
			},
		},
	})
}

func testAccGroupResourceConfig(name string) string {
	return fmt.Sprintf(`
resource "openai_group" "test" {
	name = %[1]q
}
`, name)
}
