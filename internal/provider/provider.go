package provider

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	inttflog "github.com/jianyuan/terraform-provider-openai/internal/tflog"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
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
	BaseUrl               types.String `tfsdk:"base_url"`
	AdminKey              types.String `tfsdk:"admin_key"`
	MaxRetries            types.Int64  `tfsdk:"max_retries"`
	RequestTimeoutSeconds types.Int64  `tfsdk:"request_timeout_seconds"`
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
				MarkdownDescription: "Base URL for the OpenAI API. It can also be set using the `OPENAI_BASE_URL` environment variable. Defaults to `https://api.openai.com/v1`.",
				Optional:            true,
			},
			"admin_key": schema.StringAttribute{
				MarkdownDescription: "The OpenAI admin key can be obtained through the [API Platform Organization](https://platform.openai.com/settings/organization/admin-keys) overview page. It can also be set using the `OPENAI_ADMIN_KEY` environment variable. Note that the admin key must begin with `sk-admin-`.",
				Optional:            true,
				Sensitive:           true,
			},
			"max_retries": schema.Int64Attribute{
				MarkdownDescription: "Maximum number of retries for failed requests. It can also be set using the `OPENAI_MAX_RETRIES` environment variable. Defaults to `3` retries.",
				Optional:            true,
			},
			"request_timeout_seconds": schema.Int64Attribute{
				MarkdownDescription: "Timeout for each request in seconds. It can also be set using the `OPENAI_REQUEST_TIMEOUT_SECONDS` environment variable. Defaults to `60` seconds.",
				Optional:            true,
			},
		},
	}
}

func (p *OpenAIProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data OpenAIProviderModel
	var err error

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var baseUrl string
	if !data.BaseUrl.IsNull() {
		baseUrl = data.BaseUrl.ValueString()
	} else if v := os.Getenv("OPENAI_BASE_URL"); v != "" {
		baseUrl = v
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

	maxRetries := 3
	if !data.MaxRetries.IsNull() {
		maxRetries = int(data.MaxRetries.ValueInt64())
	} else if v := os.Getenv("OPENAI_MAX_RETRIES"); v != "" {
		maxRetries, err = strconv.Atoi(v)
		if err != nil {
			resp.Diagnostics.AddError("invalid max_retries", "max_retries must be an integer")
			return
		}
	}

	requestTimeoutSeconds := 60
	if !data.RequestTimeoutSeconds.IsNull() {
		requestTimeoutSeconds = int(data.RequestTimeoutSeconds.ValueInt64())
	} else if v := os.Getenv("OPENAI_REQUEST_TIMEOUT_SECONDS"); v != "" {
		requestTimeoutSeconds, err = strconv.Atoi(v)
		if err != nil {
			resp.Diagnostics.AddError("invalid request_timeout_seconds", "request_timeout_seconds must be an integer")
			return
		}
	}

	client := new(openai.NewClient(
		option.WithBaseURL(baseUrl),
		option.WithAdminAPIKey(adminKey),
		option.WithHeader("User-Agent", fmt.Sprintf("Terraform/%s (+https://www.terraform.io) terraform-provider-openai/%s", req.TerraformVersion, p.version)),
		option.WithMaxRetries(maxRetries),
		option.WithRequestTimeout(time.Duration(requestTimeoutSeconds)*time.Second),
		option.WithDebugLog(inttflog.StandardLogger(ctx)),
	))

	pd := &providerData{
		client: client,
	}

	resp.DataSourceData = pd
	resp.ResourceData = pd
}

func (p *OpenAIProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{
		NewPredefinedProjectRoleIdFunction,
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
