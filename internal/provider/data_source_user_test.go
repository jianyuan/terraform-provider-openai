package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jianyuan/terraform-provider-openai/internal/acctest"
)

func TestAccUserDataSource(t *testing.T) {
	rn := "data.openai_user.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUserDataSourceConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("email"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("role"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("added_at"), knownvalue.NotNull()),
				},
			},
		},
	})
}

var testAccUserDataSourceConfig = fmt.Sprintf(`
data "openai_user" "test" {
	id = %[1]q
}
`, acctest.TestUserId)
