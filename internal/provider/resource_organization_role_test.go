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
	resource.AddTestSweepers("openai_organization_role", &resource.Sweeper{
		Name: "openai_organization_role",
		F: func(r string) error {
			ctx := context.Background()

			params := &apiclient.ListRolesParams{
				Limit: ptr.Ptr(int64(100)),
			}

			for {
				httpResp, err := acctest.SharedClient.ListRolesWithResponse(
					ctx,
					params,
				)

				if err != nil {
					return fmt.Errorf("[ERROR] Unable to read, got error: %s", err)
				} else if httpResp.StatusCode() != http.StatusOK || httpResp.JSON200 == nil {
					return fmt.Errorf("[ERROR] Unable to read, got status code %d: %s", httpResp.StatusCode(), string(httpResp.Body))
				}

				for _, role := range httpResp.JSON200.Data {
					if !strings.HasPrefix(role.Name, "tf-") {
						continue
					}

					log.Printf("[INFO] Found role %s (ID: %s)", role.Name, role.Id)

					httpResp, err := acctest.SharedClient.DeleteRoleWithResponse(
						ctx,
						role.Id,
					)

					if err != nil {
						log.Printf("[ERROR] Unable to delete, got error: %s", err)
						continue
					} else if httpResp.StatusCode() != http.StatusOK {
						log.Printf("[ERROR] Unable to delete, got status code %d: %s", httpResp.StatusCode(), string(httpResp.Body))
						continue
					}

					log.Printf("[INFO] Deleted role %s (ID: %s)", role.Name, role.Id)
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

func TestAccOrganizationRoleResource(t *testing.T) {
	rn := "openai_organization_role.test"
	roleName := acctest.RandomWithPrefix("tf-role")
	roleDescription := acctest.RandomWithPrefix("tf-role-description")

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
