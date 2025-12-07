package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/jianyuan/terraform-provider-openai/internal/apiclient"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

func (m *ProjectRateLimitsDataSourceModel) Fill(ctx context.Context, rateLimits []apiclient.ProjectRateLimit) diag.Diagnostics {
	items := make([]ProjectRateLimitsDataSourceModelRateLimitsItem, len(rateLimits))
	for i, rl := range rateLimits {
		items[i] = ProjectRateLimitsDataSourceModelRateLimitsItem{
			Id:                          supertypes.NewStringValue(rl.Id),
			Model:                       supertypes.NewStringValue(rl.Model),
			MaxRequestsPer1Minute:       supertypes.NewInt64Value(rl.MaxRequestsPer1Minute),
			MaxTokensPer1Minute:         supertypes.NewInt64Value(rl.MaxTokensPer1Minute),
			MaxImagesPer1Minute:         supertypes.NewInt64PointerValue(rl.MaxImagesPer1Minute),
			MaxAudioMegabytesPer1Minute: supertypes.NewInt64PointerValue(rl.MaxAudioMegabytesPer1Minute),
			MaxRequestsPer1Day:          supertypes.NewInt64PointerValue(rl.MaxRequestsPer1Day),
			Batch1DayMaxInputTokens:     supertypes.NewInt64PointerValue(rl.Batch1DayMaxInputTokens),
		}
	}
	m.RateLimits = supertypes.NewSetNestedObjectValueOfValueSlice(ctx, items)
	return nil
}
