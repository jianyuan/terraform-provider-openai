package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/openai/openai-go/v3"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

func (m *ProjectSpendLimitResourceModel) Fill(ctx context.Context, data openai.ProjectSpendLimit) diag.Diagnostics {
	var diags diag.Diagnostics
	m.Currency = supertypes.NewStringValue(string(data.Currency))
	m.Interval = supertypes.NewStringValue(string(data.Interval))
	m.ThresholdAmount = supertypes.NewInt64Value(data.ThresholdAmount)
	m.Enforcement = supertypes.NewSingleNestedObjectValueOf(ctx, &ProjectSpendLimitResourceModelEnforcement{
		Status: supertypes.NewStringValue(data.Enforcement.Status),
	})
	return diags
}

func (r *ProjectSpendLimitResource) getNewParams(ctx context.Context, data ProjectSpendLimitResourceModel) (*openai.AdminOrganizationProjectSpendLimitUpdateParams, diag.Diagnostics) {
	return r.getUpdateParams(ctx, data)
}

func (r *ProjectSpendLimitResource) getUpdateParams(ctx context.Context, data ProjectSpendLimitResourceModel) (*openai.AdminOrganizationProjectSpendLimitUpdateParams, diag.Diagnostics) {
	return &openai.AdminOrganizationProjectSpendLimitUpdateParams{
		Currency:        openai.AdminOrganizationProjectSpendLimitUpdateParamsCurrency(data.Currency.ValueString()),
		Interval:        openai.AdminOrganizationProjectSpendLimitUpdateParamsInterval(data.Interval.ValueString()),
		ThresholdAmount: data.ThresholdAmount.ValueInt64(),
	}, nil
}
