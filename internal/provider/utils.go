package provider

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/jianyuan/terraform-provider-openai/internal/apiclient"
)

// deduplicate removes duplicates from a slice of any comparable type T.
func deduplicate[T comparable](input []T) []T {
	seen := make(map[T]struct{})
	result := make([]T, 0, len(input))

	for _, v := range input {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}

func mergeDiagnostics[T any](v T, diagsOut diag.Diagnostics) func(diags *diag.Diagnostics) T {
	return func(diags *diag.Diagnostics) T {
		diags.Append(diagsOut...)
		return v
	}
}

func parseRateLimitHTTPResponse(resp *http.Response) (time.Duration, bool) {
	if resp == nil || resp.StatusCode != http.StatusTooManyRequests {
		return 0, false
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, false
	}
	resp.Body = io.NopCloser(bytes.NewBuffer(body))

	var errorResponse apiclient.ErrorResponse
	if err := json.Unmarshal(body, &errorResponse); err != nil {
		return 0, false
	}

	if errorResponse.Error.Code == nil || *errorResponse.Error.Code != "rate_limit_exceeded" {
		return 0, false
	}

	re := regexp.MustCompile(`after ([0-9.]+) seconds`)
	matches := re.FindStringSubmatch(errorResponse.Error.Message)
	if len(matches) != 2 {
		return 0, false
	}

	seconds, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return 0, false
	}

	return time.Duration(seconds * float64(time.Second)), true
}

func getBool(v any) bool {
	switch v := v.(type) {
	case bool:
		return v
	case *bool:
		if v == nil {
			return false
		}
		return *v
	default:
		panic("unknown type")
	}
}

func getString(v any) string {
	switch v := v.(type) {
	case string:
		return v
	case *string:
		if v == nil {
			return ""
		}
		return *v
	default:
		panic("unknown type")
	}
}
