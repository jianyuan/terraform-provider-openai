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
	resource.AddTestSweepers("openai_organization_role", &resource.Sweeper{
		Name: "openai_organization_role",
		F: func(r string) error {
			ctx := context.Background()

			params := openai.AdminOrganizationRoleListParams{
				Limit: openai.Int(100),
			}

			var ids []string

			iter := acctest.SharedClient.Admin.Organization.Roles.ListAutoPaging(ctx, params)
			for iter.Next() {
				item := iter.Current()
				if strings.HasPrefix(item.Name, "tf-") {
					ids = append(ids, item.ID)
				}
			}

			for _, id := range ids {
				_, err := acctest.SharedClient.Admin.Organization.Roles.Delete(ctx, id)
				if err != nil {
					return fmt.Errorf("Unable to delete, got error: %w", err)
				}
			}

			return nil
		},
	})
}

func TestAccOrganizationRoleResource(t *testing.T) {
	rn := "openai_organization_role.test"
	roleName := sdkacctest.RandomWithPrefix("tf-role")
	roleDescription := sdkacctest.RandomWithPrefix("tf-role-description")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationRoleResourceConfig(roleName, roleDescription, `["api.groups.read"]`),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(roleName)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("description"), knownvalue.StringExact(roleDescription)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("permissions"), knownvalue.SetExact([]knownvalue.Check{
						knownvalue.StringExact("api.groups.read"),
					})),
				},
			},
			{
				Config: testAccOrganizationRoleResourceConfig(roleName+"-updated", roleDescription+"-updated", `["api.groups.read", "api.groups.write"]`),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(roleName+"-updated")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("description"), knownvalue.StringExact(roleDescription+"-updated")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("permissions"), knownvalue.SetExact([]knownvalue.Check{
						knownvalue.StringExact("api.groups.read"),
						knownvalue.StringExact("api.groups.write"),
					})),
				},
			},
		},
	})
}

func testAccOrganizationRoleResourceConfig(name, description, permissions string) string {
	return fmt.Sprintf(`
resource "openai_organization_role" "test" {
	name        = %[1]q
	description = %[2]q
	permissions = %[3]s
}
`, name, description, permissions)
}
