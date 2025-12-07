package provider

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"github.com/jianyuan/go-utils/must"
	"github.com/jianyuan/terraform-provider-openai/internal/apiclient"
)

// Ensure OpenAIProvider satisfies various provider interfaces.
var _ provider.Provider = &OpenAIProvider{}
var _ provider.ProviderWithFunctions = &OpenAIProvider{}

// OpenAIProvider defines the provider implementation.
type OpenAIProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// OpenAIProviderModel describes the provider data model.
type OpenAIProviderModel struct {
	BaseUrl  types.String `tfsdk:"base_url"`
	AdminKey types.String `tfsdk:"admin_key"`
}

func (p *OpenAIProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "openai"
	resp.Version = p.version
}

func (p *OpenAIProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The OpenAI provider enables you to configure resources and data sources for your OpenAI organization. It utilizes the official [Administration API](https://platform.openai.com/docs/api-reference/administration) to interact with the OpenAI platform.\n\nIf you find this provider useful, please consider supporting me through GitHub Sponsorship or Ko-Fi to help with its development.\n\n[![Github-sponsors](https://img.shields.io/badge/sponsor-30363D?style=for-the-badge&logo=GitHub-Sponsors&logoColor=#EA4AAA)](https://github.com/sponsors/jianyuan)\n[![Ko-Fi](https://img.shields.io/badge/Ko--fi-F16061?style=for-the-badge&logo=ko-fi&logoColor=white)](https://ko-fi.com/L3L71DQEL)",
		Attributes: map[string]schema.Attribute{
			"base_url": schema.StringAttribute{
				MarkdownDescription: "Base URL for the OpenAI API. Defaults to `https://api.openai.com`.",
				Optional:            true,
			},
			"admin_key": schema.StringAttribute{
				MarkdownDescription: "The OpenAI admin key can be obtained through the [API Platform Organization](https://platform.openai.com/settings/organization/admin-keys) overview page. It can also be set using the `OPENAI_ADMIN_KEY` environment variable. Note that the admin key must begin with `sk-admin-`.",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *OpenAIProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data OpenAIProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var baseUrl string
	if !data.BaseUrl.IsNull() {
		baseUrl = data.BaseUrl.ValueString()
	} else {
		baseUrl = "https://api.openai.com/v1"
	}

	var adminKey string
	if !data.AdminKey.IsNull() {
		adminKey = data.AdminKey.ValueString()
	} else if v := os.Getenv("OPENAI_ADMIN_KEY"); v != "" {
		adminKey = v
	}

	if adminKey == "" {
		resp.Diagnostics.AddWarning("admin_key is required", "admin_key is required")
	} else if !strings.HasPrefix(adminKey, "sk-admin-") {
		resp.Diagnostics.AddError("admin_key must start with 'sk-admin-'", "admin_key must start with 'sk-admin-'")
		return
	}

	retryClient := retryablehttp.NewClient()
	retryClient.HTTPClient.Transport = logging.NewLoggingHTTPTransport(retryClient.HTTPClient.Transport)
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

	client, err := apiclient.NewClientWithResponses(
		baseUrl,
		apiclient.WithHTTPClient(retryClient.StandardClient()),
		apiclient.WithRequestEditorFn(func(ctx context.Context, httpReq *http.Request) error {
			httpReq.Header.Set("Authorization", "Bearer "+adminKey)
			httpReq.Header.Set("User-Agent", fmt.Sprintf("Terraform/%s (+https://www.terraform.io) terraform-provider-openai/%s", req.TerraformVersion, p.version))
			return nil
		}),
	)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create API client", err.Error())
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *OpenAIProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{
		NewPredefinedRoleIdFunction,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &OpenAIProvider{
			version: version,
		}
	}
}
