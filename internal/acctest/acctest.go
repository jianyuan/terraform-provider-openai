package acctest

import (
	"os"
	"testing"
)

var (
	TestAdminKey = os.Getenv("OPENAI_ADMIN_KEY")
	TestUserId   = os.Getenv("OPENAI_TEST_USER_ID")
)

func PreCheck(t *testing.T) {
	if TestAdminKey == "" {
		t.Fatal("OPENAI_ADMIN_KEY must be set for acceptance tests")
	}

	if TestUserId == "" {
		t.Fatal("OPENAI_TEST_USER_ID must be set for acceptance tests")
	}
}
