package provider_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jianyuan/terraform-provider-openai/internal/acctest"
	"github.com/openai/openai-go/v3"
)

func init() {
	resource.AddTestSweepers("openai_admin_api_key", &resource.Sweeper{
		Name: "openai_admin_api_key",
		F: func(r string) error {
			ctx := context.Background()

			params := openai.AdminOrganizationAdminAPIKeyListParams{
				Limit: openai.Int(100),
			}

			var ids []string

			iter := acctest.SharedClient.Admin.Organization.AdminAPIKeys.ListAutoPaging(ctx, params)
			for iter.Next() {
				item := iter.Current()
				if strings.HasPrefix(item.Name, "tf-") {
					ids = append(ids, item.ID)
				}
			}

			for _, apiKeyId := range ids {
				_, err := acctest.SharedClient.Admin.Organization.AdminAPIKeys.Delete(ctx, apiKeyId)
				if err != nil {
					return fmt.Errorf("Unable to delete, got error: %w", err)
				}
			}

			return nil
		},
	})
}

func TestAccAdminApiKeyResource(t *testing.T) {
	rn := "openai_admin_api_key.test"
	name := sdkacctest.RandomWithPrefix("tf")

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
