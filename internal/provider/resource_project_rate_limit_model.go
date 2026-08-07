package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/packages/param"
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

func (r *ProjectRateLimitResource) getCreateJSONRequestBody(ctx context.Context, data ProjectRateLimitResourceModel) (*openai.AdminOrganizationProjectRateLimitUpdateRateLimitParams, diag.Diagnostics) {
	return r.getUpdateJSONRequestBody(ctx, data)
}

func (r *ProjectRateLimitResource) getUpdateJSONRequestBody(ctx context.Context, data ProjectRateLimitResourceModel) (*openai.AdminOrganizationProjectRateLimitUpdateRateLimitParams, diag.Diagnostics) {
	return &openai.AdminOrganizationProjectRateLimitUpdateRateLimitParams{
		Batch1DayMaxInputTokens: (func() param.Opt[int64] {
			if data.Batch1DayMaxInputTokens.IsKnown() {
				return openai.Int(data.Batch1DayMaxInputTokens.ValueInt64())
			}
			return param.Opt[int64]{}
		})(),
		MaxAudioMegabytesPer1Minute: (func() param.Opt[int64] {
			if data.MaxAudioMegabytesPer1Minute.IsKnown() {
				return openai.Int(data.MaxAudioMegabytesPer1Minute.ValueInt64())
			}
			return param.Opt[int64]{}
		})(),
		MaxImagesPer1Minute: (func() param.Opt[int64] {
			if data.MaxImagesPer1Minute.IsKnown() {
				return openai.Int(data.MaxImagesPer1Minute.ValueInt64())
			}
			return param.Opt[int64]{}
		})(),
		MaxRequestsPer1Day: (func() param.Opt[int64] {
			if data.MaxRequestsPer1Day.IsKnown() {
				return openai.Int(data.MaxRequestsPer1Day.ValueInt64())
			}
			return param.Opt[int64]{}
		})(),
		MaxRequestsPer1Minute: (func() param.Opt[int64] {
			if data.MaxRequestsPer1Minute.IsKnown() {
				return openai.Int(data.MaxRequestsPer1Minute.ValueInt64())
			}
			return param.Opt[int64]{}
		})(),
		MaxTokensPer1Minute: (func() param.Opt[int64] {
			if data.MaxTokensPer1Minute.IsKnown() {
				return openai.Int(data.MaxTokensPer1Minute.ValueInt64())
			}
			return param.Opt[int64]{}
		})(),
	}, nil
}
