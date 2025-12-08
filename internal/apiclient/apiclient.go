package apiclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"github.com/jianyuan/go-utils/must"
)

func New(baseUrl string, terraformVersion string, providerVersion string, adminKey string) (*ClientWithResponses, error) {
	transport := cleanhttp.DefaultPooledClient().Transport
	transport = logging.NewLoggingHTTPTransport(transport)
	transport = NewRetryTransport(transport)

	httpClient := &http.Client{Transport: transport}

	client, err := NewClientWithResponses(
		baseUrl,
		WithHTTPClient(httpClient),
		WithRequestEditorFn(func(ctx context.Context, httpReq *http.Request) error {
			httpReq.Header.Set("Authorization", "Bearer "+adminKey)
			httpReq.Header.Set("User-Agent", fmt.Sprintf("Terraform/%s (+https://www.terraform.io) terraform-provider-openai/%s", terraformVersion, providerVersion))
			return nil
		}),
	)

	return client, err
}

func NewRetryTransport(transport http.RoundTripper) http.RoundTripper {
	retryClient := retryablehttp.NewClient()
	retryClient.HTTPClient.Transport = transport
	retryClient.ErrorHandler = retryablehttp.PassthroughErrorHandler
	retryClient.Logger = nil
	retryClient.RetryMax = 10

	projectServiceAccountPathPattern := regexp.MustCompile(`^/v1/organization/projects/[^/]+/service_accounts/[^/]+$`)
	retryClient.CheckRetry = func(ctx context.Context, resp *http.Response, err error) (bool, error) {
		// Retry on 404 for project service account. There's a small delay between creating a project service account and it being available.
		if resp.Request.Method == http.MethodGet && resp.StatusCode == http.StatusNotFound && projectServiceAccountPathPattern.MatchString(resp.Request.URL.Path) {
			return true, nil
		}

		if v, ok := parseRateLimitHTTPResponse(resp); ok {
			resp.Header.Set("X-Internal-Retry-After", v.String())
		} else {
			resp.Header.Set("X-Internal-Retry-After", "")
		}

		return retryablehttp.ErrorPropagatedRetryPolicy(ctx, resp, err)
	}

	retryClient.Backoff = func(durationMin, durationMax time.Duration, attemptNum int, resp *http.Response) time.Duration {
		var backoff time.Duration
		if resp != nil {
			retryAfter := resp.Header.Get("X-Internal-Retry-After")
			if retryAfter != "" {
				backoff = must.Get(time.ParseDuration(retryAfter))
			}
		}

		if backoff == 0 {
			backoff = retryablehttp.DefaultBackoff(durationMin, durationMax, attemptNum, resp)
		}

		tflog.Debug(
			resp.Request.Context(),
			fmt.Sprintf(
				"%s %s (status: %d): retrying in %s (%d left)",
				resp.Request.Method,
				resp.Request.URL.Redacted(),
				resp.StatusCode,
				backoff,
				retryClient.RetryMax-attemptNum,
			),
		)

		return backoff
	}

	return retryClient.StandardClient().Transport
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

	var errorResponse ErrorResponse
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
