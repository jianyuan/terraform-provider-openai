package provider

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/hashicorp/go-retryablehttp"
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
	BaseUrl  types.String `tfsdk:"base_url"`
	AdminKey types.String `tfsdk:"admin_key"`
}

func (p *OpenAIProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "openai"
	resp.Version = p.version
}

func (p *OpenAIProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The OpenAI provider enables you to configure resources and data sources for your OpenAI organization. It utilizes the official [Administration API](https://platform.openai.com/docs/api-reference/administration) to interact with the OpenAI platform.",
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
		resp.Diagnostics.AddError("admin_key is required", "admin_key is required")
		return
	} else if !strings.HasPrefix(adminKey, "sk-admin-") {
		resp.Diagnostics.AddError("admin_key must start with 'sk-admin-'", "admin_key must start with 'sk-admin-'")
		return
	}

	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 10

	client, err := apiclient.NewClientWithResponses(
		baseUrl,
		apiclient.WithHTTPClient(retryClient.StandardClient()),
		apiclient.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
			req.Header.Set("Authorization", "Bearer "+adminKey)
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
		NewInviteResource,
		NewProjectResource,
		NewProjectServiceAccountResource,
		NewProjectUserResource,
		NewUserRoleResource,
	}
}

func (p *OpenAIProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewInviteDataSource,
		NewInvitesDataSource,
		NewProjectDataSource,
		NewProjectsDataSource,
		NewUserDataSource,
		NewUsersDataSource,
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
