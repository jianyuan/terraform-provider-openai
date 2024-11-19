package acctest

import (
	"os"
	"testing"
)

var (
	TestAdminKey       = os.Getenv("OPENAI_ADMIN_KEY")
	TestOrganizationId = os.Getenv("OPENAI_TEST_ORGANIZATION_ID")
)

func PreCheck(t *testing.T) {
	if TestAdminKey == "" {
		t.Fatal("OPENAI_ADMIN_KEY must be set for acceptance tests")
	}

	if TestOrganizationId == "" {
		t.Fatal("OPENAI_TEST_ORGANIZATION_ID must be set for acceptance tests")
	}
}
