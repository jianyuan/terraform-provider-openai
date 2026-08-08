package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/jianyuan/terraform-provider-openai/internal/openaiparam"
	"github.com/openai/openai-go/v3"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

func (m *ProjectRateLimitResourceModel) Fill(ctx context.Context, data openai.ProjectRateLimit) diag.Diagnostics {
	m.Batch1DayMaxInputTokens = (func() supertypes.Int64Value {
		if data.JSON.Batch1DayMaxInputTokens.Valid() {
			return supertypes.NewInt64Value(data.Batch1DayMaxInputTokens)
		}
		return supertypes.NewInt64Null()
	})()
	m.MaxAudioMegabytesPer1Minute = (func() supertypes.Int64Value {
		if data.JSON.MaxAudioMegabytesPer1Minute.Valid() {
			return supertypes.NewInt64Value(data.MaxAudioMegabytesPer1Minute)
		}
		return supertypes.NewInt64Null()
	})()
	m.MaxImagesPer1Minute = (func() supertypes.Int64Value {
		if data.JSON.MaxImagesPer1Minute.Valid() {
			return supertypes.NewInt64Value(data.MaxImagesPer1Minute)
		}
		return supertypes.NewInt64Null()
	})()
	m.MaxRequestsPer1Day = (func() supertypes.Int64Value {
		if data.JSON.MaxRequestsPer1Day.Valid() {
			return supertypes.NewInt64Value(data.MaxRequestsPer1Day)
		}
		return supertypes.NewInt64Null()
	})()
	m.MaxRequestsPer1Minute = supertypes.NewInt64Value(data.MaxRequestsPer1Minute)
	m.MaxTokensPer1Minute = supertypes.NewInt64Value(data.MaxTokensPer1Minute)
	return nil
}

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
