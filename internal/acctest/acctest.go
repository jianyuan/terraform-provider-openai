package acctest

import (
	"os"
	"testing"
)

var (
	TestOrganizationId = os.Getenv("OPENAI_TEST_ORGANIZATION_ID")
	TestApiKey         = os.Getenv("OPENAI_API_KEY")
)

func PreCheck(t *testing.T) {
	if TestApiKey == "" {
		t.Fatal("OPENAI_API_KEY must be set for acceptance tests")
	}

	if TestOrganizationId == "" {
		t.Fatal("OPENAI_TEST_ORGANIZATION_ID must be set for acceptance tests")
	}
}
