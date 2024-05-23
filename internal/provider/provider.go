package provider

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
	BaseUrl    types.String `tfsdk:"base_url"`
	ApiKey     types.String `tfsdk:"api_key"`
	SessionKey types.String `tfsdk:"session_key"`
}

func (p *OpenAIProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "openai"
	resp.Version = p.version
}

func (p *OpenAIProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"base_url": schema.StringAttribute{
				MarkdownDescription: "Base URL for the OpenAI API. Defaults to `https://api.openai.com`.",
				Optional:            true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "API key for the OpenAI API.",
				Optional:            true,
				Sensitive:           true,
			},
			"session_key": schema.StringAttribute{
				MarkdownDescription: "Session key for the OpenAI API.",
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
		baseUrl = "https://api.openai.com"
	}

	var apiKey string
	if !data.ApiKey.IsNull() {
		apiKey = data.ApiKey.ValueString()
	} else if v := os.Getenv("OPENAI_API_KEY"); v != "" {
		apiKey = v
	}

	if apiKey == "" {
		resp.Diagnostics.AddError("api_key is required", "api_key is required")
		return
	}

	var sessionKey string
	if !data.SessionKey.IsNull() {
		sessionKey = data.SessionKey.ValueString()
	} else if v := os.Getenv("OPENAI_SESSION_KEY"); v != "" {
		sessionKey = v
	}

	if sessionKey == "" {
		resp.Diagnostics.AddError("session_key is required", "session_key is required")
		return
	}

	client, err := apiclient.NewClientWithResponses(
		baseUrl,
		apiclient.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
			if strings.HasPrefix(req.URL.Path, "/v1") {
				req.Header.Set("Authorization", "Bearer "+apiKey)
			} else if strings.HasPrefix(req.URL.Path, "/dashboard") {
				req.Header.Set("Authorization", "Bearer "+sessionKey)
			}

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

func (p *OpenAIProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewProjectResource,
	}
}

func (p *OpenAIProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewMembersDataSource,
		NewOrganizationDataSource,
		NewOrganizationsDataSource,
		NewProjectDataSource,
		NewProjectsDataSource,
	}
}

func (p *OpenAIProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &OpenAIProvider{
			version: version,
		}
	}
}
