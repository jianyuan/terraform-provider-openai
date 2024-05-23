package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jianyuan/terraform-provider-openai/internal/acctest"
)

func TestAccProjectResource(t *testing.T) {
	rn := "openai_project.test"
	projectTitle := acctest.RandomWithPrefix("tf-project")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectResourceConfig(projectTitle),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("organization_id"), knownvalue.StringExact(acctest.TestOrganizationId)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("title"), knownvalue.StringExact(projectTitle)),
				},
			},
			{
				ResourceName: rn,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[rn]
					if !ok {
						return "", fmt.Errorf("not found: %s", rn)
					}
					organizationId := rs.Primary.Attributes["organization_id"]
					id := rs.Primary.ID
					return BuildTwoPartId(organizationId, id), nil
				},
				ImportStateVerify: true,
			},
			{
				Config: testAccProjectResourceConfig(projectTitle + "-updated"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("organization_id"), knownvalue.StringExact(acctest.TestOrganizationId)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("title"), knownvalue.StringExact(projectTitle+"-updated")),
				},
			},
		},
	})
}

func testAccProjectResourceConfig(title string) string {
	return testAccOrganizationDataSourceConfig + fmt.Sprintf(`
resource "openai_project" "test" {
  organization_id = data.openai_organization.test.id
  title           = %[1]q
}
`, title)
}
