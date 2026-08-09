package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/openai/openai-go/v3"
)

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
