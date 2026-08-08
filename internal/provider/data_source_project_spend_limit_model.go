package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/openai/openai-go/v3"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

func (m *ProjectSpendLimitDataSourceModel) Fill(ctx context.Context, data openai.ProjectSpendLimit) diag.Diagnostics {
	m.Currency = supertypes.NewStringValue(string(data.Currency))
	m.Interval = supertypes.NewStringValue(string(data.Interval))
	m.ThresholdAmount = supertypes.NewInt64Value(data.ThresholdAmount)
	m.Enforcement = supertypes.NewSingleNestedObjectValueOf(ctx, &ProjectSpendLimitDataSourceModelEnforcement{
		Status: supertypes.NewStringValue(data.Enforcement.Status),
	})
	return nil
}
