package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/go-utils/ptr"
	"github.com/jianyuan/terraform-provider-openai/internal/apiclient"
)

type ProjectRateLimitModel struct {
	Id                          types.String `tfsdk:"id"`
	Model                       types.String `tfsdk:"model"`
	MaxRequestsPer1Minute       types.Int64  `tfsdk:"max_requests_per_1_minute"`
	MaxTokensPer1Minute         types.Int64  `tfsdk:"max_tokens_per_1_minute"`
	MaxImagesPer1Minute         types.Int64  `tfsdk:"max_images_per_1_minute"`
	MaxAudioMegabytesPer1Minute types.Int64  `tfsdk:"max_audio_megabytes_per_1_minute"`
	MaxRequestsPer1Day          types.Int64  `tfsdk:"max_requests_per_1_day"`
	Batch1DayMaxInputTokens     types.Int64  `tfsdk:"batch_1_day_max_input_tokens"`
}

func (m *ProjectRateLimitModel) Fill(ctx context.Context, rl apiclient.ProjectRateLimit) (diags diag.Diagnostics) {
	m.Id = types.StringValue(rl.Id)
	m.Model = types.StringValue(rl.Model)
	m.MaxRequestsPer1Minute = types.Int64Value(int64(rl.MaxRequestsPer1Minute))
	m.MaxTokensPer1Minute = types.Int64Value(int64(rl.MaxTokensPer1Minute))
	if rl.MaxImagesPer1Minute == nil {
		m.MaxImagesPer1Minute = types.Int64Null()
	} else {
		m.MaxImagesPer1Minute = types.Int64Value(int64(*rl.MaxImagesPer1Minute))
	}
	if rl.MaxAudioMegabytesPer1Minute == nil {
		m.MaxAudioMegabytesPer1Minute = types.Int64Null()
	} else {
		m.MaxAudioMegabytesPer1Minute = types.Int64Value(int64(*rl.MaxAudioMegabytesPer1Minute))
	}
	if rl.MaxRequestsPer1Day == nil {
		m.MaxRequestsPer1Day = types.Int64Null()
	} else {
		m.MaxRequestsPer1Day = types.Int64Value(int64(*rl.MaxRequestsPer1Day))
	}
	if rl.Batch1DayMaxInputTokens == nil {
		m.Batch1DayMaxInputTokens = types.Int64Null()
	} else {
		m.Batch1DayMaxInputTokens = types.Int64Value(int64(*rl.Batch1DayMaxInputTokens))
	}
	return
}

type ProjectRateLimitsDataSourceModel struct {
	ProjectId  types.String            `tfsdk:"project_id"`
	RateLimits []ProjectRateLimitModel `tfsdk:"rate_limits"`
}

func (m *ProjectRateLimitsDataSourceModel) Fill(ctx context.Context, rateLimits []apiclient.ProjectRateLimit) (diags diag.Diagnostics) {
	m.RateLimits = make([]ProjectRateLimitModel, len(rateLimits))
	for i, rl := range rateLimits {
		diags.Append(m.RateLimits[i].Fill(ctx, rl)...)
		if diags.HasError() {
			return
		}
	}
	return
}

var _ datasource.DataSource = &ProjectRateLimitsDataSource{}

func NewProjectRateLimitsDataSource() datasource.DataSource {
	return &ProjectRateLimitsDataSource{}
}

type ProjectRateLimitsDataSource struct {
	baseDataSource
}

func (d *ProjectRateLimitsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project_rate_limits"
}

func (d *ProjectRateLimitsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Returns the rate limits per model for a project.",

		Attributes: map[string]schema.Attribute{
			"project_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the project.",
				Required:            true,
			},
			"rate_limits": schema.SetNestedAttribute{
				MarkdownDescription: "List of users.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "The rate limit identifier.",
							Computed:            true,
						},
						"model": schema.StringAttribute{
							MarkdownDescription: "The model this rate limit applies to.",
							Computed:            true,
						},
						"max_requests_per_1_minute": schema.Int64Attribute{
							MarkdownDescription: "The maximum requests per minute.",
							Computed:            true,
						},
						"max_tokens_per_1_minute": schema.Int64Attribute{
							MarkdownDescription: "The maximum tokens per minute.",
							Computed:            true,
						},
						"max_images_per_1_minute": schema.Int64Attribute{
							MarkdownDescription: "The maximum images per minute. Only present for relevant models.",
							Computed:            true,
						},
						"max_audio_megabytes_per_1_minute": schema.Int64Attribute{
							MarkdownDescription: "The maximum audio megabytes per minute. Only present for relevant models.",
							Computed:            true,
						},
						"max_requests_per_1_day": schema.Int64Attribute{
							MarkdownDescription: "The maximum requests per day. Only present for relevant models.",
							Computed:            true,
						},
						"batch_1_day_max_input_tokens": schema.Int64Attribute{
							MarkdownDescription: "The maximum batch input tokens per day. Only present for relevant models.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *ProjectRateLimitsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ProjectRateLimitsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var rateLimits []apiclient.ProjectRateLimit
	params := &apiclient.ListProjectRateLimitsParams{
		Limit: ptr.Ptr(100),
	}

	for {
		httpResp, err := d.client.ListProjectRateLimitsWithResponse(
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

		rateLimits = append(rateLimits, httpResp.JSON200.Data...)

		if !httpResp.JSON200.HasMore {
			break
		}

		params.After = &httpResp.JSON200.LastId
	}

	resp.Diagnostics.Append(data.Fill(ctx, rateLimits)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
