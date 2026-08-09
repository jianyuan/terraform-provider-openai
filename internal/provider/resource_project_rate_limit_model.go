package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/jianyuan/terraform-provider-openai/internal/openaiparam"
	"github.com/openai/openai-go/v3"
)

func (r *ProjectRateLimitResource) resourceMatch(data ProjectRateLimitResourceModel, rateLimit openai.ProjectRateLimit) bool {
	return data.RateLimitId.ValueString() == rateLimit.ID
}

func (r *ProjectRateLimitResource) getNewParams(ctx context.Context, data ProjectRateLimitResourceModel) (*openai.AdminOrganizationProjectRateLimitUpdateRateLimitParams, diag.Diagnostics) {
	return r.getUpdateParams(ctx, data)
}

func (r *ProjectRateLimitResource) getUpdateParams(ctx context.Context, data ProjectRateLimitResourceModel) (*openai.AdminOrganizationProjectRateLimitUpdateRateLimitParams, diag.Diagnostics) {
	return &openai.AdminOrganizationProjectRateLimitUpdateRateLimitParams{
		Batch1DayMaxInputTokens:     openaiparam.FromInt64(data.Batch1DayMaxInputTokens),
		MaxAudioMegabytesPer1Minute: openaiparam.FromInt64(data.MaxAudioMegabytesPer1Minute),
		MaxImagesPer1Minute:         openaiparam.FromInt64(data.MaxImagesPer1Minute),
		MaxRequestsPer1Day:          openaiparam.FromInt64(data.MaxRequestsPer1Day),
		MaxRequestsPer1Minute:       openaiparam.FromInt64(data.MaxRequestsPer1Minute),
		MaxTokensPer1Minute:         openaiparam.FromInt64(data.MaxTokensPer1Minute),
	}, nil
}
