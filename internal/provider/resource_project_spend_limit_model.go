package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/jianyuan/terraform-provider-openai/internal/apiclient"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

func (m *ProjectSpendLimitResourceModel) Fill(ctx context.Context, data apiclient.ProjectSpendLimitResource) diag.Diagnostics {
	var diags diag.Diagnostics

	if v, err := data.Currency.AsSpendLimitCurrency0(); err == nil {
		m.Currency = supertypes.NewStringValue(v)
	} else if v, err := data.Currency.AsSpendLimitCurrency1(); err == nil && v.Valid() {
		m.Currency = supertypes.NewStringValue(string(v))
	} else {
		diags.AddError("Unable to parse currency", "Unexpected value for currency")
	}

	if v, err := data.Interval.AsSpendLimitInterval0(); err == nil {
		m.Interval = supertypes.NewStringValue(v)
	} else if v, err := data.Interval.AsSpendLimitInterval1(); err == nil && v.Valid() {
		m.Interval = supertypes.NewStringValue(string(v))
	} else {
		diags.AddError("Unable to parse interval", "Unexpected value for interval")
	}

	m.ThresholdAmount = supertypes.NewInt64Value(data.ThresholdAmount)

	var enforcement ProjectSpendLimitResourceModelEnforcement
	if v, err := data.Enforcement.Status.AsSpendLimitEnforcementStatus0(); err == nil {
		enforcement.Status = supertypes.NewStringValue(v)
	} else if v, err := data.Enforcement.Status.AsSpendLimitEnforcementStatus1(); err == nil && v.Valid() {
		enforcement.Status = supertypes.NewStringValue(string(v))
	} else {
		diags.AddError("Unable to parse enforcement status", "Unexpected value for enforcement status")
	}
	m.Enforcement = supertypes.NewSingleNestedObjectValueOf(ctx, &enforcement)

	return diags
}

func (r *ProjectSpendLimitResource) getCreateJSONRequestBody(ctx context.Context, data ProjectSpendLimitResourceModel) (apiclient.UpdateProjectSpendLimitJSONRequestBody, diag.Diagnostics) {
	return apiclient.UpdateProjectSpendLimitJSONRequestBody{
		Currency:        apiclient.UpdateProjectSpendLimitBodyCurrency(data.Currency.ValueString()),
		Interval:        apiclient.UpdateProjectSpendLimitBodyInterval(data.Interval.ValueString()),
		ThresholdAmount: data.ThresholdAmount.ValueInt64(),
	}, nil
}

func (r *ProjectSpendLimitResource) getUpdateJSONRequestBody(ctx context.Context, data ProjectSpendLimitResourceModel) (apiclient.UpdateProjectSpendLimitJSONRequestBody, diag.Diagnostics) {
	return r.getCreateJSONRequestBody(ctx, data)
}
