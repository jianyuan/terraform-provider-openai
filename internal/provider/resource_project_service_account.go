package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/terraform-provider-openai/internal/apiclient"
)

type ProjectServiceAccountResourceModel struct {
	ProjectId types.String `tfsdk:"project_id"`
	Name      types.String `tfsdk:"name"`
	Id        types.String `tfsdk:"id"`
	Role      types.String `tfsdk:"role"`
	CreatedAt types.Int64  `tfsdk:"created_at"`
	ApiKeyId  types.String `tfsdk:"api_key_id"`
	ApiKey    types.String `tfsdk:"api_key"`
}

func (m *ProjectServiceAccountResourceModel) Fill(sa apiclient.ProjectServiceAccount) error {
	m.Id = types.StringValue(sa.Id)
	m.Name = types.StringValue(sa.Name)
	m.Role = types.StringValue(string(sa.Role))
	m.CreatedAt = types.Int64Value(int64(sa.CreatedAt))
	return nil
}

func (m *ProjectServiceAccountResourceModel) FillFromCreate(sa apiclient.ProjectServiceAccountCreateResponse) error {
	m.Id = types.StringValue(sa.Id)
	m.Name = types.StringValue(sa.Name)
	m.Role = types.StringValue(string(sa.Role))
	m.CreatedAt = types.Int64Value(int64(sa.CreatedAt))
	m.ApiKeyId = types.StringValue(sa.ApiKey.Id)
	m.ApiKey = types.StringValue(sa.ApiKey.Value)
	return nil
}

var _ resource.Resource = &ProjectServiceAccountResource{}

func NewProjectServiceAccountResource() resource.Resource {
	return &ProjectServiceAccountResource{}
}

type ProjectServiceAccountResource struct {
	baseResource
}

func (r *ProjectServiceAccountResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project_service_account"
}

func (r *ProjectServiceAccountResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manage service accounts within a project. A service account is a bot user that is not associated with a user. If a user leaves an organization, their keys and membership in projects will no longer work. Service accounts do not have this limitation. However, service accounts can also be deleted from a project.",

		Attributes: map[string]schema.Attribute{
			"project_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the project.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the service account being created.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the service account.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"role": schema.StringAttribute{
				MarkdownDescription: "The role of the service account. Can be `owner` or `member`.",
				Computed:            true,
			},
			"created_at": schema.Int64Attribute{
				MarkdownDescription: "The Unix timestamp (in seconds) of when the service account was created.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"api_key_id": schema.StringAttribute{
				MarkdownDescription: "Internal ID of the API key. This is a reference to the API key and not the actual key.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "The API key that can be used to authenticate with the API.",
				Computed:            true,
				Sensitive:           true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *ProjectServiceAccountResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ProjectServiceAccountResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.client.CreateProjectServiceAccountWithResponse(
		ctx,
		data.ProjectId.ValueString(),
		apiclient.ProjectServiceAccountCreateRequest{
			Name: data.Name.ValueString(),
		},
	)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create, got error: %s", err))
		return
	}

	if httpResp.StatusCode() != http.StatusCreated {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create, got status code %d: %s", httpResp.StatusCode(), string(httpResp.Body)))
		return
	}

	if err := data.FillFromCreate(*httpResp.JSON201); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to fill data: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProjectServiceAccountResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ProjectServiceAccountResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.client.RetrieveProjectServiceAccountWithResponse(
		ctx,
		data.ProjectId.ValueString(),
		data.Id.ValueString(),
	)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read, got error: %s", err))
		return
	}

	if httpResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read, got status code %d: %s", httpResp.StatusCode(), string(httpResp.Body)))
		return
	}

	if err := data.Fill(*httpResp.JSON200); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to fill data: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProjectServiceAccountResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Not Supported", "Update is not supported for this resource")
}

func (r *ProjectServiceAccountResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ProjectServiceAccountResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.client.DeleteProjectServiceAccountWithResponse(
		ctx,
		data.ProjectId.ValueString(),
		data.Id.ValueString(),
	)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update, got error: %s", err))
		return
	}

	if httpResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update, got status code %d: %s", httpResp.StatusCode(), string(httpResp.Body)))
		return
	}
}
