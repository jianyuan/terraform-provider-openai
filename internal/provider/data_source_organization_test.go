package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jianyuan/terraform-provider-openai/internal/acctest"
)

func TestAccOrganizationDataSource(t *testing.T) {
	rn := "data.openai_organization.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationDataSourceConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.StringExact(acctest.TestOrganizationId)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("is_default"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("title"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("description"), knownvalue.NotNull()),
				},
			},
		},
	})
}

func TestAccOrganizationDataSource_default(t *testing.T) {
	rn := "data.openai_organization.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationDataSourceConfig_default,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("is_default"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("title"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("description"), knownvalue.NotNull()),
				},
			},
		},
	})
}

var testAccOrganizationDataSourceConfig = fmt.Sprintf(`
data "openai_organization" "test" {
  id = "%s"
}
`, acctest.TestOrganizationId)

const testAccOrganizationDataSourceConfig_default = `
data "openai_organization" "test" {
}
`
