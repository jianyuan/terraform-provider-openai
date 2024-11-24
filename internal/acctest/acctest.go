package acctest

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/jianyuan/go-utils/must"
	"github.com/jianyuan/terraform-provider-openai/internal/apiclient"
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
