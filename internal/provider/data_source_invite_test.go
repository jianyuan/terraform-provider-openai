package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jianyuan/terraform-provider-openai/internal/acctest"
)

func TestAccInviteDataSource(t *testing.T) {
	email := fmt.Sprintf("tf-%d@example.com", acctest.RandInt())
	rn := "openai_invite.test"
	dn := "data.openai_invite.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInviteDataSourceConfig(email),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs(rn, tfjsonpath.New("id"), dn, tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.CompareValuePairs(rn, tfjsonpath.New("email"), dn, tfjsonpath.New("email"), compare.ValuesSame()),
					statecheck.CompareValuePairs(rn, tfjsonpath.New("role"), dn, tfjsonpath.New("role"), compare.ValuesSame()),
					statecheck.CompareValuePairs(rn, tfjsonpath.New("status"), dn, tfjsonpath.New("status"), compare.ValuesSame()),
					statecheck.CompareValuePairs(rn, tfjsonpath.New("invited_at"), dn, tfjsonpath.New("invited_at"), compare.ValuesSame()),
					statecheck.CompareValuePairs(rn, tfjsonpath.New("expires_at"), dn, tfjsonpath.New("expires_at"), compare.ValuesSame()),
					statecheck.CompareValuePairs(rn, tfjsonpath.New("accepted_at"), dn, tfjsonpath.New("accepted_at"), compare.ValuesSame()),
				},
			},
		},
	})
}

func testAccInviteDataSourceConfig(email string) string {
	return fmt.Sprintf(`
resource "openai_invite" "test" {
	email = %[1]q
	role  = "reader"
}

data "openai_invite" "test" {
	id = openai_invite.test.id
}
`, email)
}
