package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jianyuan/terraform-provider-openai/internal/acctest"
)

func TestAccGroupUsersDataSource(t *testing.T) {
	rn := "data.openai_group_users.test"
	groupName := acctest.RandomWithPrefix("tf-group")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGroupUsersDataSourceConfig(groupName, acctest.TestUserId),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("group_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("users"), knownvalue.SetPartial([]knownvalue.Check{
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"id":       knownvalue.StringExact(acctest.TestUserId),
							"email":    knownvalue.NotNull(),
							"name":     knownvalue.NotNull(),
							"role":     knownvalue.NotNull(),
							"added_at": knownvalue.NotNull(),
						}),
					})),
				},
			},
		},
	})
}

func testAccGroupUsersDataSourceConfig(groupName, userId string) string {
	return testAccGroupUserResourceConfig(groupName, userId) + `
data "openai_group_users" "test" {
	depends_on = [openai_group_user.test]
	group_id = openai_group.test.id
}
`
}
