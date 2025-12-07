package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/jianyuan/terraform-provider-openai/internal/apiclient"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

func (m *ProjectRateLimitResourceModel) Fill(ctx context.Context, data apiclient.ProjectRateLimit) diag.Diagnostics {
	m.Batch1DayMaxInputTokens = supertypes.NewInt64PointerValue(data.Batch1DayMaxInputTokens)
	m.MaxAudioMegabytesPer1Minute = supertypes.NewInt64PointerValue(data.MaxAudioMegabytesPer1Minute)
	m.MaxImagesPer1Minute = supertypes.NewInt64PointerValue(data.MaxImagesPer1Minute)
	m.MaxRequestsPer1Day = supertypes.NewInt64PointerValue(data.MaxRequestsPer1Day)
	m.MaxRequestsPer1Minute = supertypes.NewInt64Value(data.MaxRequestsPer1Minute)
	m.MaxTokensPer1Minute = supertypes.NewInt64Value(data.MaxTokensPer1Minute)
	return nil
}

func (r *ProjectRateLimitResource) resourceMatch(data ProjectRateLimitResourceModel, rateLimit apiclient.ProjectRateLimit) bool {
	return data.RateLimitId.ValueString() == rateLimit.Id
}

func (r *ProjectRateLimitResource) getCreateJSONRequestBody(ctx context.Context, data ProjectRateLimitResourceModel) (apiclient.UpdateProjectRateLimitsJSONRequestBody, diag.Diagnostics) {
	return apiclient.UpdateProjectRateLimitsJSONRequestBody{
		Batch1DayMaxInputTokens:     data.Batch1DayMaxInputTokens.GetInt64Ptr(),
		MaxAudioMegabytesPer1Minute: data.MaxAudioMegabytesPer1Minute.GetInt64Ptr(),
		MaxImagesPer1Minute:         data.MaxImagesPer1Minute.GetInt64Ptr(),
		MaxRequestsPer1Day:          data.MaxRequestsPer1Day.GetInt64Ptr(),
		MaxRequestsPer1Minute:       data.MaxRequestsPer1Minute.GetInt64Ptr(),
		MaxTokensPer1Minute:         data.MaxTokensPer1Minute.GetInt64Ptr(),
	}, nil
}

func (r *ProjectRateLimitResource) getUpdateJSONRequestBody(ctx context.Context, data ProjectRateLimitResourceModel) (apiclient.UpdateProjectRateLimitsJSONRequestBody, diag.Diagnostics) {
	return apiclient.UpdateProjectRateLimitsJSONRequestBody{
		Batch1DayMaxInputTokens:     data.Batch1DayMaxInputTokens.GetInt64Ptr(),
		MaxAudioMegabytesPer1Minute: data.MaxAudioMegabytesPer1Minute.GetInt64Ptr(),
		MaxImagesPer1Minute:         data.MaxImagesPer1Minute.GetInt64Ptr(),
		MaxRequestsPer1Day:          data.MaxRequestsPer1Day.GetInt64Ptr(),
		MaxRequestsPer1Minute:       data.MaxRequestsPer1Minute.GetInt64Ptr(),
		MaxTokensPer1Minute:         data.MaxTokensPer1Minute.GetInt64Ptr(),
	}, nil
}
