package acctest

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/jianyuan/go-utils/must"
	"github.com/jianyuan/terraform-provider-openai/internal/apiclient"
	"github.com/jianyuan/terraform-provider-openai/internal/provider"
)

var (
	TestAdminKey = os.Getenv("OPENAI_ADMIN_KEY")
	TestUserId   = os.Getenv("OPENAI_TEST_USER_ID")

	SharedClient *apiclient.ClientWithResponses
)

func init() {
	SharedClient = must.Get(apiclient.NewClientWithResponses(
		"https://api.openai.com/v1",
		apiclient.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
			req.Header.Set("Authorization", "Bearer "+TestAdminKey)
			return nil
		}),
	))
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
