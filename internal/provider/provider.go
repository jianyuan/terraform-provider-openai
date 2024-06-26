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
	SessionKey types.String `tfsdk:"session_key"`
}

func (p *OpenAIProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "openai"
	resp.Version = p.version
}

func (p *OpenAIProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The OpenAI provider allows you to configure resources and data sources for your OpenAI organization. It uses internal APIs, so breaking changes are expected.\n\n" +
			"Unfortunately, OpenAI's API keys do not allow some functionalities. Therefore, we need to obtain an OpenAI session key from the `Authorization` header of any requests to `https://api.openai.com/dashboard/*`. Log in to https://platform.openai.com, use Inspect Element to look for any requests to `https://api.openai.com/dashboard/*`, and grab the `Authorization` header value.",
		Attributes: map[string]schema.Attribute{
			"base_url": schema.StringAttribute{
				MarkdownDescription: "Base URL for the OpenAI API. Defaults to `https://api.openai.com`.",
				Optional:            true,
			},
			"session_key": schema.StringAttribute{
				MarkdownDescription: "The OpenAI session key can be obtained by accessing the dashboard in your browser. This can also be set via the `OPENAI_SESSION_KEY` environment variable. Note that the session key must start with `sess-`.",
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

	var sessionKey string
	if !data.SessionKey.IsNull() {
		sessionKey = data.SessionKey.ValueString()
	} else if v := os.Getenv("OPENAI_SESSION_KEY"); v != "" {
		sessionKey = v
	}

	if sessionKey == "" {
		resp.Diagnostics.AddError("session_key is required", "session_key is required")
		return
	} else if !strings.HasPrefix(sessionKey, "sess-") {
		resp.Diagnostics.AddError("session_key must start with 'sess-'", "session_key must start with 'sess-'")
		return
	}

	client, err := apiclient.NewClientWithResponses(
		baseUrl,
		apiclient.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
			req.Header.Set("Authorization", "Bearer "+sessionKey)
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
		NewProjectApiKeyResource,
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
