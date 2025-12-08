package apiclient

import (
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

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
