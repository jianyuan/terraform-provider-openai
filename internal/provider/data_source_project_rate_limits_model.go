package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/openai/openai-go/v3"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
	"github.com/samber/lo"
)

func (m *ProjectRateLimitsDataSourceModel) Fill(ctx context.Context, rateLimits []openai.ProjectRateLimit) diag.Diagnostics {
	m.RateLimits = supertypes.NewSetNestedObjectValueOfValueSlice(ctx, lo.Map(rateLimits, func(rl openai.ProjectRateLimit, _ int) ProjectRateLimitsDataSourceModelRateLimitsItem {
		return ProjectRateLimitsDataSourceModelRateLimitsItem{
			Id:                    supertypes.NewStringValue(rl.ID),
			Model:                 supertypes.NewStringValue(rl.Model),
			MaxRequestsPer1Minute: supertypes.NewInt64Value(rl.MaxRequestsPer1Minute),
			MaxTokensPer1Minute:   supertypes.NewInt64Value(rl.MaxTokensPer1Minute),
			MaxImagesPer1Minute: (func() supertypes.Int64Value {
				if rl.JSON.MaxImagesPer1Minute.Valid() {
					return supertypes.NewInt64Value(rl.MaxImagesPer1Minute)
				}
				return supertypes.NewInt64Null()
			})(),
			MaxAudioMegabytesPer1Minute: (func() supertypes.Int64Value {
				if rl.JSON.MaxAudioMegabytesPer1Minute.Valid() {
					return supertypes.NewInt64Value(rl.MaxAudioMegabytesPer1Minute)
				}
				return supertypes.NewInt64Null()
			})(),
			MaxRequestsPer1Day: (func() supertypes.Int64Value {
				if rl.JSON.MaxRequestsPer1Day.Valid() {
					return supertypes.NewInt64Value(rl.MaxRequestsPer1Day)
				}
				return supertypes.NewInt64Null()
			})(),
			Batch1DayMaxInputTokens: (func() supertypes.Int64Value {
				if rl.JSON.Batch1DayMaxInputTokens.Valid() {
					return supertypes.NewInt64Value(rl.Batch1DayMaxInputTokens)
				}
				return supertypes.NewInt64Null()
			})(),
		}
	}))

	return nil
}
