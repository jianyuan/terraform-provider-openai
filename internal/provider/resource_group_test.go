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
	sweep.Register("openai_group", func(ctx context.Context, client *openai.Client) ([]sweep.Sweepable, error) {
		params := openai.AdminOrganizationGroupListParams{
			Limit: openai.Int(100),
		}

		var sweepables []sweep.Sweepable

		iter := acctest.SharedClient.Admin.Organization.Groups.ListAutoPaging(ctx, params)
		for iter.Next() {
			item := iter.Current()
			if strings.HasPrefix(item.Name, "tf-") {
				sweepables = append(sweepables, sweep.NewSweepResource(provider.NewGroupResource, acctest.SharedClient, map[string]any{
					"id": item.ID,
				}))
			}
		}

		return sweepables, nil
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
