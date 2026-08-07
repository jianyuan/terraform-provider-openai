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
	resource.AddTestSweepers("openai_project", &resource.Sweeper{
		Name: "openai_project",
		F: func(r string) error {
			ctx := context.Background()

			params := openai.AdminOrganizationProjectListParams{
				Limit: openai.Int(100),
			}

			var ids []string

			iter := acctest.SharedClient.Admin.Organization.Projects.ListAutoPaging(ctx, params)
			for iter.Next() {
				item := iter.Current()
				if strings.HasPrefix(item.Name, "tf-") {
					ids = append(ids, item.ID)
				}
			}

			for _, id := range ids {
				_, err := acctest.SharedClient.Admin.Organization.Projects.Archive(ctx, id)
				if err != nil {
					return fmt.Errorf("Unable to delete, got error: %w", err)
				}
			}

			return nil
		},
	})
}

func TestAccProjectResource(t *testing.T) {
	rn := "openai_project.test"
	projectName := sdkacctest.RandomWithPrefix("tf-project")

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
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("external_key_id"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("geography"), knownvalue.Null()),
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
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("external_key_id"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("geography"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("archived_at"), knownvalue.Null()),
				},
			},
		},
	})
}

func TestAccProjectResource_WithGeography(t *testing.T) {
	rn := "openai_project.test"
	projectName := sdkacctest.RandomWithPrefix("tf-project")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectResourceConfigWithGeography(projectName, "US"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(projectName)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("status"), knownvalue.StringExact("active")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("geography"), knownvalue.StringExact("US")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("archived_at"), knownvalue.Null()),
				},
			},
			{
				ResourceName:            rn,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"geography"},
			},
			{
				Config: testAccProjectResourceConfigWithGeography(projectName+"-updated", "US"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(projectName+"-updated")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("status"), knownvalue.StringExact("active")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("geography"), knownvalue.StringExact("US")),
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

func testAccProjectResourceConfigWithGeography(name, geography string) string {
	return fmt.Sprintf(`
resource "openai_project" "test" {
	name      = %[1]q
	geography = %[2]q
}
`, name, geography)
}
