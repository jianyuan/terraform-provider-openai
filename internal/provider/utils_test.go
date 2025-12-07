package provider

import (
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
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

func TestParseRateLimitHTTPResponse(t *testing.T) {
	resp := &http.Response{
		StatusCode: http.StatusTooManyRequests,
		Body: io.NopCloser(strings.NewReader(`{
			"error": {
				"message": "You've exceeded the rate limit, please slow down and try again after 6.937 seconds.",
				"type": "invalid_request_error",
				"param": null,
				"code": "rate_limit_exceeded"
			}
		}`)),
	}

	v, ok := parseRateLimitHTTPResponse(resp)

	if !ok {
		t.Errorf("Expected true, got %v", ok)
	}

	if v != time.Duration(6.937*float64(time.Second)) {
		t.Errorf("Expected 6s, got %v", v)
	}
}
