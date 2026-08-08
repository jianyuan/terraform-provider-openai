package acctest

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/jianyuan/terraform-provider-openai/internal/provider"
	inttflog "github.com/jianyuan/terraform-provider-openai/internal/tflog"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

var (
	TestBaseUrl  = os.Getenv("OPENAI_BASE_URL")
	TestAdminKey = os.Getenv("OPENAI_ADMIN_KEY")
	TestUserId   = os.Getenv("OPENAI_TEST_USER_ID")
	TestGroupId  string

	SharedClient *openai.Client
)

func init() {
	if TestBaseUrl == "" {
		TestBaseUrl = "https://api.openai.com/v1"
	}

	SharedClient = new(openai.NewClient(
		option.WithBaseURL(TestBaseUrl),
		option.WithAdminAPIKey(TestAdminKey),
		option.WithHeader("User-Agent", fmt.Sprintf("Terraform/%s (+https://www.terraform.io) terraform-provider-openai/%s", "dev", "dev")),
		option.WithRequestTimeout(60*time.Second),
		option.WithMaxRetries(5),
		option.WithDebugLog(inttflog.StandardLogger(context.Background())),
	))

	TestGroupId = ensureTestGroupId(context.Background())
}

func PreCheck(t *testing.T) {
	if TestAdminKey == "" {
		t.Fatal("OPENAI_ADMIN_KEY must be set for acceptance tests")
	}

	if TestUserId == "" {
		t.Fatal("OPENAI_TEST_USER_ID must be set for acceptance tests")
	}
}

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"openai": providerserver.NewProtocol6WithError(provider.New("test")()),
}

func ensureTestGroupId(ctx context.Context) string {
	params := openai.AdminOrganizationGroupListParams{
		Limit: openai.Int(100),
	}

	iter := SharedClient.Admin.Organization.Groups.ListAutoPaging(ctx, params)
	for iter.Next() {
		group := iter.Current()
		if group.Name == "acc-tf-group" {
			return group.ID
		}
	}

	group, err := SharedClient.Admin.Organization.Groups.New(ctx, openai.AdminOrganizationGroupNewParams{
		Name: "acc-tf-group",
	})
	if err != nil {
		panic(err)
	}

	return group.ID
}
