package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/openai/openai-go/v3"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

func (m *SpendLimitResourceModel) Fill(ctx context.Context, data openai.OrganizationSpendLimit) diag.Diagnostics {
	m.Currency = supertypes.NewStringValue(string(data.Currency))
	m.Interval = supertypes.NewStringValue(string(data.Interval))
	m.ThresholdAmount = supertypes.NewInt64Value(data.ThresholdAmount)
	m.Enforcement = supertypes.NewSingleNestedObjectValueOf(ctx, &SpendLimitResourceModelEnforcement{
		Status: supertypes.NewStringValue(data.Enforcement.Status),
	})
	return nil
}

func (r *SpendLimitResource) getNewParams(ctx context.Context, data SpendLimitResourceModel) (*openai.AdminOrganizationSpendLimitUpdateParams, diag.Diagnostics) {
	return r.getUpdateParams(ctx, data)
}

func (r *SpendLimitResource) getUpdateParams(ctx context.Context, data SpendLimitResourceModel) (*openai.AdminOrganizationSpendLimitUpdateParams, diag.Diagnostics) {
	return &openai.AdminOrganizationSpendLimitUpdateParams{
		Currency:        openai.AdminOrganizationSpendLimitUpdateParamsCurrency(data.Currency.ValueString()),
		Interval:        openai.AdminOrganizationSpendLimitUpdateParamsInterval(data.Interval.ValueString()),
		ThresholdAmount: data.ThresholdAmount.ValueInt64(),
	}, nil
}
