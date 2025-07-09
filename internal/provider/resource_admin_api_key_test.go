package provider_test

import (
	"context"
	"fmt"
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
	resource.AddTestSweepers("openai_admin_api_key", &resource.Sweeper{
		Name: "openai_admin_api_key",
		F: func(r string) error {
			ctx := context.Background()

			params := &apiclient.AdminApiKeysListParams{
				Limit: ptr.Ptr(int64(100)),
			}

			for {
				httpResp, err := acctest.SharedClient.AdminApiKeysListWithResponse(
					ctx,
					params,
				)

				if err != nil {
					return fmt.Errorf("Unable to read, got error: %s", err)
				} else if httpResp.StatusCode() != http.StatusOK || httpResp.JSON200 == nil {
					return fmt.Errorf("Unable to read, got status code %d: %s", httpResp.StatusCode(), string(httpResp.Body))
				}

				for _, apiKey := range *httpResp.JSON200.Data {
					if strings.HasPrefix(apiKey.Name, "tf-") {
						httpResp, err := acctest.SharedClient.AdminApiKeysDeleteWithResponse(
							ctx,
							apiKey.Id,
						)

						if err != nil {
							return fmt.Errorf("Unable to delete, got error: %s", err)
						} else if httpResp.StatusCode() != http.StatusOK {
							return fmt.Errorf("Unable to delete, got status code %d: %s", httpResp.StatusCode(), string(httpResp.Body))
						}
					}
				}

				if !ptr.Value(httpResp.JSON200.HasMore) {
					break
				}

				params.After = httpResp.JSON200.LastId
			}

			return nil
		},
	})
}

func TestAccAdminApiKeyResource(t *testing.T) {
	rn := "openai_admin_api_key.test"
	name := acctest.RandomWithPrefix("tf")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAdminApiKeyResourceConfig(name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("api_key"), knownvalue.NotNull()),
				},
			},
			{
				Config: testAccAdminApiKeyResourceConfig(name + "-changed"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(name+"-changed")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("api_key"), knownvalue.NotNull()),
				},
			},
		},
	})
}

func testAccAdminApiKeyResourceConfig(name string) string {
	return fmt.Sprintf(`
resource "openai_admin_api_key" "test" {
	name = %[1]q
}
`, name)
}
