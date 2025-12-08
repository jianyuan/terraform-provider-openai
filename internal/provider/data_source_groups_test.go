package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jianyuan/terraform-provider-openai/internal/acctest"
)

func TestAccGroupsDataSource(t *testing.T) {
	rn := "data.openai_groups.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGroupsDataSourceConfig(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("groups"), knownvalue.SetPartial([]knownvalue.Check{
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"id":              knownvalue.StringExact(acctest.TestGroupId),
							"name":            knownvalue.StringExact("acc-tf-group"),
							"is_scim_managed": knownvalue.Bool(false),
							"created_at":      knownvalue.NotNull(),
						}),
					})),
				},
			},
		},
	})
}

func testAccGroupsDataSourceConfig() string {
	return `
data "openai_groups" "test" {
}
`
}
