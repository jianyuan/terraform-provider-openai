package acctest

import (
	"os"
	"testing"
)

var (
	TestAdminKey = os.Getenv("OPENAI_ADMIN_KEY")
)

func PreCheck(t *testing.T) {
	if TestAdminKey == "" {
		t.Fatal("OPENAI_ADMIN_KEY must be set for acceptance tests")
	}
}
