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
	"github.com/jianyuan/terraform-provider-openai/internal/provider"
	"github.com/jianyuan/terraform-provider-openai/internal/sweep"
	"github.com/openai/openai-go/v3"
)

func init() {
	sweep.Register("openai_admin_api_key", func(ctx context.Context, client *openai.Client) ([]sweep.Sweepable, error) {
		params := openai.AdminOrganizationAdminAPIKeyListParams{
			Limit: openai.Int(100),
		}

		var sweepables []sweep.Sweepable

		iter := acctest.SharedClient.Admin.Organization.AdminAPIKeys.ListAutoPaging(ctx, params)
		for iter.Next() {
			item := iter.Current()
			if strings.HasPrefix(item.Name, "tf-") {
				sweepables = append(sweepables, sweep.NewSweepResource(provider.NewAdminApiKeyResource, acctest.SharedClient, map[string]any{
					"id": item.ID,
				}))
			}
		}

		return sweepables, nil
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
