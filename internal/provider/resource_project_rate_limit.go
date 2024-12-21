package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/go-utils/ptr"
	"github.com/jianyuan/terraform-provider-openai/internal/apiclient"
)

type ProjectRateLimitResourceModel struct {
	ProjectId                   types.String `tfsdk:"project_id"`
	Model                       types.String `tfsdk:"model"`
	MaxRequestsPer1Minute       types.Int64  `tfsdk:"max_requests_per_1_minute"`
	MaxTokensPer1Minute         types.Int64  `tfsdk:"max_tokens_per_1_minute"`
	MaxImagesPer1Minute         types.Int64  `tfsdk:"max_images_per_1_minute"`
	MaxAudioMegabytesPer1Minute types.Int64  `tfsdk:"max_audio_megabytes_per_1_minute"`
	MaxRequestsPer1Day          types.Int64  `tfsdk:"max_requests_per_1_day"`
	Batch1DayMaxInputTokens     types.Int64  `tfsdk:"batch_1_day_max_input_tokens"`
}

func (m *ProjectRateLimitResourceModel) Fill(ctx context.Context, rl apiclient.ProjectRateLimit) (diags diag.Diagnostics) {
	m.Model = types.StringValue(rl.Model)
	m.MaxRequestsPer1Minute = types.Int64Value(int64(rl.MaxRequestsPer1Minute))
	m.MaxTokensPer1Minute = types.Int64Value(int64(rl.MaxTokensPer1Minute))
	if !m.MaxImagesPer1Minute.IsNull() {
		if rl.MaxImagesPer1Minute == nil {
			m.MaxImagesPer1Minute = types.Int64Null()
		} else {
			m.MaxImagesPer1Minute = types.Int64Value(int64(*rl.MaxImagesPer1Minute))
		}
	}
	if !m.MaxAudioMegabytesPer1Minute.IsNull() {
		if rl.MaxAudioMegabytesPer1Minute == nil {
			m.MaxAudioMegabytesPer1Minute = types.Int64Null()
		} else {
			m.MaxAudioMegabytesPer1Minute = types.Int64Value(int64(*rl.MaxAudioMegabytesPer1Minute))
		}
	}
	if !m.MaxRequestsPer1Day.IsNull() {
		if rl.MaxRequestsPer1Day == nil {
			m.MaxRequestsPer1Day = types.Int64Null()
		} else {
			m.MaxRequestsPer1Day = types.Int64Value(int64(*rl.MaxRequestsPer1Day))
		}
	}
	if !m.Batch1DayMaxInputTokens.IsNull() {
		if rl.Batch1DayMaxInputTokens == nil {
			m.Batch1DayMaxInputTokens = types.Int64Null()
		} else {
			m.Batch1DayMaxInputTokens = types.Int64Value(int64(*rl.Batch1DayMaxInputTokens))
		}
	}
	return
}

var _ resource.Resource = &ProjectRateLimitResource{}

func NewProjectRateLimitResource() resource.Resource {
	return &ProjectRateLimitResource{}
}

type ProjectRateLimitResource struct {
	baseResource
}

func (r *ProjectRateLimitResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project_rate_limit"
}

func (r *ProjectRateLimitResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manage rate limits per model for projects. Rate limits may be configured to be equal to or lower than the organization's rate limits.",

		Attributes: map[string]schema.Attribute{
			"project_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the project.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"model": schema.StringAttribute{
				MarkdownDescription: "The model to set rate limits for.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"max_requests_per_1_minute": schema.Int64Attribute{
				MarkdownDescription: "The maximum requests per minute.",
				Optional:            true,
			},
			"max_tokens_per_1_minute": schema.Int64Attribute{
				MarkdownDescription: "The maximum tokens per minute.",
				Optional:            true,
			},
			"max_images_per_1_minute": schema.Int64Attribute{
				MarkdownDescription: "The maximum images per minute. Only present for relevant models.",
				Optional:            true,
			},
			"max_audio_megabytes_per_1_minute": schema.Int64Attribute{
				MarkdownDescription: "The maximum audio megabytes per minute. Only present for relevant models.",
				Optional:            true,
			},
			"max_requests_per_1_day": schema.Int64Attribute{
				MarkdownDescription: "The maximum requests per day. Only present for relevant models.",
				Optional:            true,
			},
			"batch_1_day_max_input_tokens": schema.Int64Attribute{
				MarkdownDescription: "The maximum batch input tokens per day. Only present for relevant models.",
				Optional:            true,
			},
		},
	}
}

func (r *ProjectRateLimitResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ProjectRateLimitResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	body := apiclient.ProjectRateLimitUpdateRequest{}
	if !data.MaxRequestsPer1Minute.IsNull() {
		body.MaxRequestsPer1Minute = ptr.Ptr(int(data.MaxRequestsPer1Minute.ValueInt64()))
	}
	if !data.MaxTokensPer1Minute.IsNull() {
		body.MaxTokensPer1Minute = ptr.Ptr(int(data.MaxTokensPer1Minute.ValueInt64()))
	}
	if !data.MaxImagesPer1Minute.IsNull() {
		body.MaxImagesPer1Minute = ptr.Ptr(int(data.MaxImagesPer1Minute.ValueInt64()))
	}
	if !data.MaxAudioMegabytesPer1Minute.IsNull() {
		body.MaxAudioMegabytesPer1Minute = ptr.Ptr(int(data.MaxAudioMegabytesPer1Minute.ValueInt64()))
	}
	if !data.MaxRequestsPer1Day.IsNull() {
		body.MaxRequestsPer1Day = ptr.Ptr(int(data.MaxRequestsPer1Day.ValueInt64()))
	}
	if !data.Batch1DayMaxInputTokens.IsNull() {
		body.Batch1DayMaxInputTokens = ptr.Ptr(int(data.Batch1DayMaxInputTokens.ValueInt64()))
	}

	httpResp, err := r.client.UpdateProjectRateLimitsWithResponse(
		ctx,
		data.ProjectId.ValueString(),
		fmt.Sprintf("rl-%s", data.Model.ValueString()),
		body,
	)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create, got error: %s", err))
		return
	}

	if httpResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create, got status code %d: %s", httpResp.StatusCode(), string(httpResp.Body)))
		return
	}

	if httpResp.JSON200 == nil {
		resp.Diagnostics.AddError("Client Error", "Unable to create, got empty response")
		return
	}

	resp.Diagnostics.Append(data.Fill(ctx, *httpResp.JSON200)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProjectRateLimitResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ProjectRateLimitResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := &apiclient.ListProjectRateLimitsParams{
		Limit: ptr.Ptr(100),
	}

out:
	for {
		httpResp, err := r.client.ListProjectRateLimitsWithResponse(
			ctx,
			data.ProjectId.ValueString(),
			params,
		)

		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read, got error: %s", err))
			return
		}

		if httpResp.StatusCode() != http.StatusOK {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read, got status code %d: %s", httpResp.StatusCode(), string(httpResp.Body)))
			return
		}

		if httpResp.JSON200 == nil {
			resp.Diagnostics.AddError("Client Error", "Unable to read, got empty response")
			return
		}

		for _, rl := range httpResp.JSON200.Data {
			if rl.Model != data.Model.ValueString() {
				continue
			}

			resp.Diagnostics.Append(data.Fill(ctx, rl)...)
			if resp.Diagnostics.HasError() {
				return
			}

			break out
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProjectRateLimitResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ProjectRateLimitResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	body := apiclient.ProjectRateLimitUpdateRequest{}
	if !data.MaxRequestsPer1Minute.IsNull() {
		body.MaxRequestsPer1Minute = ptr.Ptr(int(data.MaxRequestsPer1Minute.ValueInt64()))
	}
	if !data.MaxTokensPer1Minute.IsNull() {
		body.MaxTokensPer1Minute = ptr.Ptr(int(data.MaxTokensPer1Minute.ValueInt64()))
	}
	if !data.MaxImagesPer1Minute.IsNull() {
		body.MaxImagesPer1Minute = ptr.Ptr(int(data.MaxImagesPer1Minute.ValueInt64()))
	}
	if !data.MaxAudioMegabytesPer1Minute.IsNull() {
		body.MaxAudioMegabytesPer1Minute = ptr.Ptr(int(data.MaxAudioMegabytesPer1Minute.ValueInt64()))
	}
	if !data.MaxRequestsPer1Day.IsNull() {
		body.MaxRequestsPer1Day = ptr.Ptr(int(data.MaxRequestsPer1Day.ValueInt64()))
	}
	if !data.Batch1DayMaxInputTokens.IsNull() {
		body.Batch1DayMaxInputTokens = ptr.Ptr(int(data.Batch1DayMaxInputTokens.ValueInt64()))
	}

	httpResp, err := r.client.UpdateProjectRateLimitsWithResponse(
		ctx,
		data.ProjectId.ValueString(),
		fmt.Sprintf("rl-%s", data.Model.ValueString()),
		body,
	)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update, got error: %s", err))
		return
	}

	if httpResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update, got status code %d: %s", httpResp.StatusCode(), string(httpResp.Body)))
		return
	}

	if httpResp.JSON200 == nil {
		resp.Diagnostics.AddError("Client Error", "Unable to update, got empty response")
		return
	}

	resp.Diagnostics.Append(data.Fill(ctx, *httpResp.JSON200)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProjectRateLimitResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddWarning("Delete not supported", "This resource does not support deletion.")
}
