package acctest

import (
	"os"
	"testing"
)

var (
	TestSessionKey     = os.Getenv("OPENAI_SESSION_KEY")
	TestOrganizationId = os.Getenv("OPENAI_TEST_ORGANIZATION_ID")
)

func PreCheck(t *testing.T) {
	if TestSessionKey == "" {
		t.Fatal("OPENAI_SESSION_KEY must be set for acceptance tests")
	}

	if TestOrganizationId == "" {
		t.Fatal("OPENAI_TEST_ORGANIZATION_ID must be set for acceptance tests")
	}
}
