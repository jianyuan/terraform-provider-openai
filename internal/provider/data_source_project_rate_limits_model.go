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
		item := ProjectRateLimitsDataSourceModelRateLimitsItem{
			Id:                    supertypes.NewStringValue(rl.ID),
			Model:                 supertypes.NewStringValue(rl.Model),
			MaxRequestsPer1Minute: supertypes.NewInt64Value(rl.MaxRequestsPer1Minute),
			MaxTokensPer1Minute:   supertypes.NewInt64Value(rl.MaxTokensPer1Minute),
		}

		if rl.JSON.MaxImagesPer1Minute.Valid() {
			item.MaxImagesPer1Minute = supertypes.NewInt64Value(rl.MaxImagesPer1Minute)
		} else {
			item.MaxImagesPer1Minute = supertypes.NewInt64Null()
		}

		if rl.JSON.MaxAudioMegabytesPer1Minute.Valid() {
			item.MaxAudioMegabytesPer1Minute = supertypes.NewInt64Value(rl.MaxAudioMegabytesPer1Minute)
		} else {
			item.MaxAudioMegabytesPer1Minute = supertypes.NewInt64Null()
		}

		if rl.JSON.MaxRequestsPer1Day.Valid() {
			item.MaxRequestsPer1Day = supertypes.NewInt64Value(rl.MaxRequestsPer1Day)
		} else {
			item.MaxRequestsPer1Day = supertypes.NewInt64Null()
		}

		if rl.JSON.Batch1DayMaxInputTokens.Valid() {
			item.Batch1DayMaxInputTokens = supertypes.NewInt64Value(rl.Batch1DayMaxInputTokens)
		} else {
			item.Batch1DayMaxInputTokens = supertypes.NewInt64Null()
		}

		return item
	}))

	return nil
}
