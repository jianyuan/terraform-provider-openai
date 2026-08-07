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
	resource.AddTestSweepers("openai_group", &resource.Sweeper{
		Name: "openai_group",
		F: func(r string) error {
			ctx := context.Background()

			params := openai.AdminOrganizationGroupListParams{
				Limit: openai.Int(100),
			}

			var ids []string

			iter := acctest.SharedClient.Admin.Organization.Groups.ListAutoPaging(ctx, params)
			for iter.Next() {
				item := iter.Current()
				if strings.HasPrefix(item.Name, "tf-") {
					ids = append(ids, item.ID)
				}
			}

			for _, apiKeyId := range ids {
				_, err := acctest.SharedClient.Admin.Organization.Groups.Delete(ctx, apiKeyId)
				if err != nil {
					return fmt.Errorf("Unable to delete, got error: %w", err)
				}
			}

			return nil
		},
	})
}

func TestAccGroupResource(t *testing.T) {
	rn := "openai_group.test"
	groupName := sdkacctest.RandomWithPrefix("tf-group")

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
