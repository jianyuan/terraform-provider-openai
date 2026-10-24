package provider

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/openai/openai-go/v3"
)

func TestMergeDiagnostics(t *testing.T) {
	var diags diag.Diagnostics

	do := func() (string, diag.Diagnostics) {
		var diags diag.Diagnostics
		diags.AddError("Error summary", "Error detail")
		diags.AddWarning("Warning summary", "Warning detail")
		return "Result", diags
	}

	v := mergeDiagnostics(do())(&diags)

	if v != "Result" {
		t.Errorf("Expected Result, got %s", v)
	}

	var expectedDiags diag.Diagnostics
	expectedDiags.AddError("Error summary", "Error detail")
	expectedDiags.AddWarning("Warning summary", "Warning detail")

	if !diags.Equal(expectedDiags) {
		t.Errorf("Expected %v, got %v", expectedDiags, diags)
	}
}

func TestClientErrorDetail(t *testing.T) {
	var withMessage openai.Error
	if err := withMessage.UnmarshalJSON([]byte(`{"message":"Bucket not found","type":"invalid_request_error","code":null,"param":null}`)); err != nil {
		t.Fatal(err)
	}
	withMessage.StatusCode = http.StatusBadRequest

	tests := []struct {
		name string
		err  error
		want string
	}{
		{"api error with message", fmt.Errorf("wrapped: %w", &withMessage), "OpenAI API error: 400 Bad Request: Bucket not found"},
		{"api error without message", &openai.Error{StatusCode: http.StatusNotFound}, "OpenAI API error: 404 Not Found"},
		{"other error", errors.New("connection refused"), "connection refused"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := clientErrorDetail(tt.err); got != tt.want {
				t.Errorf("clientErrorDetail() = %q, want %q", got, tt.want)
			}
		})
	}
}
