package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/jianyuan/terraform-provider-openai/internal/apiclient"
)

var _ resource.Resource = &InviteResource{}
var _ resource.ResourceWithImportState = &InviteResource{}

func NewInviteResource() resource.Resource {
	return &InviteResource{}
}

type InviteResource struct {
	baseResource
}

func (r *InviteResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_invite"
}

func (r *InviteResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Invite and manage invitations for an organization. Invited users are automatically added to the Default project.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Invite ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"email": schema.StringAttribute{
				MarkdownDescription: "The email address of the individual to whom the invite was sent.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"role": schema.StringAttribute{
				MarkdownDescription: "`owner` or `reader`.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("owner", "reader"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "`accepted`, `expired`, or `pending`.",
				Computed:            true,
			},
			"invited_at": schema.Int64Attribute{
				MarkdownDescription: "The Unix timestamp (in seconds) of when the invite was sent.",
				Computed:            true,
			},
			"expires_at": schema.Int64Attribute{
				MarkdownDescription: "The Unix timestamp (in seconds) of when the invite expires.",
				Computed:            true,
			},
			"accepted_at": schema.Int64Attribute{
				MarkdownDescription: "The Unix timestamp (in seconds) of when the invite was accepted.",
				Computed:            true,
			},
		},
	}
}

func (r *InviteResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data InviteModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.client.InviteUserWithResponse(
		ctx,
		apiclient.InviteRequest{
			Email: data.Email.ValueString(),
			Role:  apiclient.InviteRequestRole(data.Role.ValueString()),
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

	if httpResp.JSON201 == nil {
		resp.Diagnostics.AddError("Client Error", "Unable to create, got empty response body")
		return
	}

	resp.Diagnostics.Append(data.Fill(ctx, *httpResp.JSON201)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *InviteResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data InviteModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.client.RetrieveInviteWithResponse(
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

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *InviteResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Not Supported", "Update is not supported for this resource")
}

func (r *InviteResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data InviteModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.client.DeleteInviteWithResponse(
		ctx,
		data.Id.ValueString(),
	)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete, got error: %s", err))
		return
	} else if httpResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete, got status code %d: %s", httpResp.StatusCode(), httpResp.Body))
		return
	}
}

func (r *InviteResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
