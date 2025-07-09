package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/terraform-provider-openai/internal/apiclient"
)

type AdminApiKeyResourceModel struct {
	Name      types.String `tfsdk:"name"`
	Id        types.String `tfsdk:"id"`
	CreatedAt types.Int64  `tfsdk:"created_at"`
	ApiKey    types.String `tfsdk:"api_key"`
}

func (m *AdminApiKeyResourceModel) Fill(ctx context.Context, apiKey apiclient.AdminApiKey) (diags diag.Diagnostics) {
	m.Id = types.StringValue(apiKey.Id)
	m.Name = types.StringValue(apiKey.Name)
	m.CreatedAt = types.Int64Value(apiKey.CreatedAt)
	return
}

func (m *AdminApiKeyResourceModel) FillFromCreate(ctx context.Context, apiKey apiclient.AdminApiKey) (diags diag.Diagnostics) {
	diags.Append(m.Fill(ctx, apiKey)...)
	if diags.HasError() {
		return
	}

	m.ApiKey = types.StringPointerValue(apiKey.Value)
	return
}

var _ resource.Resource = &AdminApiKeyResource{}

func NewAdminApiKeyResource() resource.Resource {
	return &AdminApiKeyResource{}
}

type AdminApiKeyResource struct {
	baseResource
}

func (r *AdminApiKeyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_admin_api_key"
}

func (r *AdminApiKeyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an organization admin API key.",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the organization admin API key.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the organization admin API key.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_at": schema.Int64Attribute{
				MarkdownDescription: "The Unix timestamp (in seconds) of when the organization admin API key was created.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "The organization admin API key that can be used to authenticate with the API.",
				Computed:            true,
				Sensitive:           true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *AdminApiKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data AdminApiKeyResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.client.AdminApiKeysCreateWithResponse(
		ctx,
		apiclient.AdminApiKeysCreateJSONRequestBody{
			Name: data.Name.ValueString(),
		},
	)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create, got error: %s", err))
		return
	} else if httpResp.StatusCode() != http.StatusOK || httpResp.JSON200 == nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create, got status code %d: %s", httpResp.StatusCode(), string(httpResp.Body)))
		return
	}

	resp.Diagnostics.Append(data.FillFromCreate(ctx, *httpResp.JSON200)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *AdminApiKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data AdminApiKeyResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.client.AdminApiKeysGetWithResponse(
		ctx,
		data.Id.ValueString(),
	)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read, got error: %s", err))
		return
	} else if httpResp.StatusCode() != http.StatusOK || httpResp.JSON200 == nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read, got status code %d: %s", httpResp.StatusCode(), string(httpResp.Body)))
		return
	}

	resp.Diagnostics.Append(data.Fill(ctx, *httpResp.JSON200)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *AdminApiKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Not Supported", "Update is not supported for this resource")
}

func (r *AdminApiKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data AdminApiKeyResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.client.AdminApiKeysDeleteWithResponse(
		ctx,
		data.Id.ValueString(),
	)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update, got error: %s", err))
		return
	} else if httpResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update, got status code %d: %s", httpResp.StatusCode(), string(httpResp.Body)))
		return
	}
}
