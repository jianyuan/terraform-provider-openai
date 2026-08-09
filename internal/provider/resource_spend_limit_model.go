package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/openai/openai-go/v3"
)

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
