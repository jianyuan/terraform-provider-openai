package provider_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/compare"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jianyuan/terraform-provider-openai/internal/acctest"
)

func TestAccExternalKeyResource_gcp(t *testing.T) {
	rn := "openai_external_key.test"
	name := sdkacctest.RandomWithPrefix("tf")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccExternalKeyResourceGcpConfig(name, "openai-key"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.StringRegexp(regexp.MustCompile(`^extkey_`))),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("api_project_ids"), knownvalue.SetSizeExact(0)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("gcp").AtMapKey("workload_identity_project_number"), knownvalue.StringExact("123456789012")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("gcp").AtMapKey("workload_identity_pool_id"), knownvalue.StringExact("openai-pool")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("gcp").AtMapKey("workload_identity_provider_id"), knownvalue.StringExact("openai-provider")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("gcp").AtMapKey("audience"), knownvalue.StringExact("org-test")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("gcp").AtMapKey("kms_project_id"), knownvalue.StringExact("kms-project")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("gcp").AtMapKey("kms_key_ring_name"), knownvalue.StringExact("openai-ring")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("gcp").AtMapKey("kms_key_name"), knownvalue.StringExact("openai-key")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("gcp").AtMapKey("kms_key_location"), knownvalue.StringExact("us")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("aws"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("azure"), knownvalue.Null()),
				},
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccExternalKeyResourceGcpConfig(name, "openai-key-rotated"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(rn, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("gcp").AtMapKey("kms_key_name"), knownvalue.StringExact("openai-key-rotated")),
				},
			},
		},
	})
}

func TestAccExternalKeyResource_aws(t *testing.T) {
	rn := "openai_external_key.test"
	name := sdkacctest.RandomWithPrefix("tf")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "openai_external_key" "test" {
	name = %[1]q
	aws = {
		role_arn    = "arn:aws:iam::123456789012:role/openai-ekm"
		kms_arn     = "arn:aws:kms:us-east-1:123456789012:key/00000000-0000-0000-0000-000000000000"
		external_id = "org-test"
	}
}
`, name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("aws").AtMapKey("role_arn"), knownvalue.StringExact("arn:aws:iam::123456789012:role/openai-ekm")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("aws").AtMapKey("kms_arn"), knownvalue.StringExact("arn:aws:kms:us-east-1:123456789012:key/00000000-0000-0000-0000-000000000000")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("aws").AtMapKey("external_id"), knownvalue.StringExact("org-test")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("azure"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("gcp"), knownvalue.Null()),
				},
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccExternalKeyResource_azure(t *testing.T) {
	rn := "openai_external_key.test"
	name := sdkacctest.RandomWithPrefix("tf")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "openai_external_key" "test" {
	name = %[1]q
	azure = {
		tenant_id = "00000000-0000-0000-0000-000000000000"
		vault_uri = "https://openai-ekm.vault.azure.net/"
		key_name  = "openai-key"
	}
}
`, name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("azure").AtMapKey("tenant_id"), knownvalue.StringExact("00000000-0000-0000-0000-000000000000")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("azure").AtMapKey("vault_uri"), knownvalue.StringExact("https://openai-ekm.vault.azure.net/")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("azure").AtMapKey("key_name"), knownvalue.StringExact("openai-key")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("aws"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("gcp"), knownvalue.Null()),
				},
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccExternalKeyResource_project(t *testing.T) {
	keyName := sdkacctest.RandomWithPrefix("tf")
	projectName := sdkacctest.RandomWithPrefix("tf-project")

	config := testAccExternalKeyResourceGcpConfig(keyName, "openai-key") + fmt.Sprintf(`
resource "openai_project" "test" {
	name            = %[1]q
	external_key_id = openai_external_key.test.id
}
`, projectName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs(
						"openai_project.test", tfjsonpath.New("external_key_id"),
						"openai_external_key.test", tfjsonpath.New("id"),
						compare.ValuesSame(),
					),
				},
			},
			{
				// The key lists the project once it is refreshed.
				RefreshState: true,
				Check: resource.TestCheckResourceAttrPair(
					"openai_external_key.test", "api_project_ids.0",
					"openai_project.test", "id",
				),
			},
		},
	})
}

func TestAccExternalKeyResource_requiresOneProvider(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "openai_external_key" "test" {
	name = "tf-invalid"
}
`,
				ExpectError: acctest.ExpectLiteralError(`No attribute specified when one (and only one) of`),
			},
		},
	})
}

func testAccExternalKeyResourceGcpConfig(name, keyName string) string {
	return fmt.Sprintf(`
resource "openai_external_key" "test" {
	name = %[1]q
	gcp = {
		workload_identity_project_number = "123456789012"
		workload_identity_pool_id        = "openai-pool"
		workload_identity_provider_id    = "openai-provider"
		audience                         = "org-test"
		kms_project_id                   = "kms-project"
		kms_key_ring_name                = "openai-ring"
		kms_key_name                     = %[2]q
		kms_key_location                 = "us"
	}
}
`, name, keyName)
}
