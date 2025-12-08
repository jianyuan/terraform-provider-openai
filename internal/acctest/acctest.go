package acctest

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/jianyuan/go-utils/must"
	"github.com/jianyuan/go-utils/ptr"
	"github.com/jianyuan/terraform-provider-openai/internal/apiclient"
	"github.com/jianyuan/terraform-provider-openai/internal/provider"
)

var (
	TestAdminKey = os.Getenv("OPENAI_ADMIN_KEY")
	TestUserId   = os.Getenv("OPENAI_TEST_USER_ID")
	TestGroupId  string

	SharedClient *apiclient.ClientWithResponses
)

func init() {
	SharedClient = must.Get(apiclient.New("https://api.openai.com/v1", "", "", TestAdminKey))

	ctx := context.Background()
	TestGroupId = ensureTestGroupId(ctx)
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
	params := &apiclient.ListGroupsParams{
		Limit: ptr.Ptr(int64(100)),
	}

	for {
		httpResp := must.Get(SharedClient.ListGroupsWithResponse(ctx, params))

		if httpResp.StatusCode() != http.StatusOK || httpResp.JSON200 == nil {
			panic(fmt.Sprintf("failed to list groups: %v", string(httpResp.Body)))
		}

		for _, group := range httpResp.JSON200.Data {
			if group.Name == "acc-tf-group" {
				return group.Id
			}
		}

		if !httpResp.JSON200.HasMore {
			break
		}

		params.After = httpResp.JSON200.Next
	}

	httpResp := must.Get(SharedClient.CreateGroupWithResponse(ctx, apiclient.CreateGroupJSONRequestBody{
		Name: "acc-tf-group",
	}))

	if httpResp.StatusCode() != http.StatusOK || httpResp.JSON200 == nil {
		panic(fmt.Sprintf("failed to create group: %v", httpResp.JSON200))
	}

	return httpResp.JSON200.Id
}
