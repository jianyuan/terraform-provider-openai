package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/terraform-provider-openai/internal/apiclient"
)

var _ resource.Resource = &ProjectResource{}
var _ resource.ResourceWithImportState = &ProjectResource{}

func NewProjectResource() resource.Resource {
	return &ProjectResource{}
}

type ProjectResource struct {
	baseResource
}

// ProjectResourceModel describes the resource data model.
type ProjectResourceModel struct {
	Id             types.String `tfsdk:"id"`
	OrganizationId types.String `tfsdk:"organization_id"`
	Title          types.String `tfsdk:"title"`
}

func (m *ProjectResourceModel) Fill(p apiclient.Project) error {
	m.Id = types.StringValue(p.Id)
	m.OrganizationId = types.StringValue(p.OrganizationId)
	m.Title = types.StringValue(p.Title)
	return nil
}

func (r *ProjectResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

func (r *ProjectResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Project resource.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the project.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"organization_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the organization to which the project belongs.",
				Required:            true,
			},
			"title": schema.StringAttribute{
				MarkdownDescription: "Human-friendly label for the project, shown in user interfaces.",
				Required:            true,
			},
		},
	}
}

func (r *ProjectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ProjectResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.client.CreateOrganizationProjectWithResponse(
		ctx,
		data.OrganizationId.ValueString(),
		apiclient.CreateOrganizationProjectJSONRequestBody{
			Geography: "",
			Title:     data.Title.ValueString(),
		},
	)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read, got error: %s", err))
		return
	}

	if httpResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create, got status code %d: %s", httpResp.StatusCode(), string(httpResp.Body)))
		return
	}

	if err := data.Fill(*httpResp.JSON200); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to unmarshal response: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProjectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ProjectResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.client.GetOrganizationProjectWithResponse(
		ctx,
		data.OrganizationId.ValueString(),
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
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to unmarshal response: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProjectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ProjectResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.client.UpdateOrganizationProjectWithResponse(
		ctx,
		data.OrganizationId.ValueString(),
		data.Id.ValueString(),
		apiclient.UpdateOrganizationProjectJSONRequestBody{
			Title: data.Title.ValueStringPointer(),
		},
	)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update, got error: %s", err))
		return
	}

	if httpResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update, got status code %d: %s", httpResp.StatusCode(), string(httpResp.Body)))
		return
	}

	if err := data.Fill(*httpResp.JSON200); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to unmarshal response: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProjectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ProjectResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.client.UpdateOrganizationProjectWithResponse(
		ctx,
		data.OrganizationId.ValueString(),
		data.Id.ValueString(),
		apiclient.UpdateOrganizationProjectJSONRequestBody{
			Archive: Pointer(true),
		},
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

func (r *ProjectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	organizationId, id, err := SplitTwoPartId(req.ID, "organization-id", "id")
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Error parsing ID: %s", err.Error()))
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("organization_id"), organizationId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}
