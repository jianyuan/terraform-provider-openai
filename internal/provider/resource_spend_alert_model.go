package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/shared/constant"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

func (m *SpendAlertResourceModel) Fill(ctx context.Context, data openai.OrganizationSpendAlert) diag.Diagnostics {
	m.Id = supertypes.NewStringValue(data.ID)
	m.Currency = supertypes.NewStringValue(string(data.Currency))
	m.Interval = supertypes.NewStringValue(string(data.Interval))
	m.ThresholdAmount = supertypes.NewInt64Value(data.ThresholdAmount)
	m.NotificationChannel = supertypes.NewSingleNestedObjectValueOf(ctx, &SpendAlertResourceModelNotificationChannel{
		Type:       supertypes.NewStringValue(string(data.NotificationChannel.Type)),
		Recipients: supertypes.NewSetValueOfSlice(ctx, data.NotificationChannel.Recipients),
		SubjectPrefix: (func() supertypes.StringValue {
			if data.NotificationChannel.JSON.SubjectPrefix.Valid() {
				return supertypes.NewStringValue(data.NotificationChannel.SubjectPrefix)
			}
			return supertypes.NewStringNull()
		})(),
	})
	return nil
}

func (r *SpendAlertResource) getCreateJSONRequestBody(ctx context.Context, data SpendAlertResourceModel) (*openai.AdminOrganizationSpendAlertNewParams, diag.Diagnostics) {
	notificationChannel, diags := data.NotificationChannel.Get(ctx)
	if diags.HasError() {
		return nil, diags
	}

	notificationChannelRecipients, diags := notificationChannel.Recipients.Get(ctx)
	if diags.HasError() {
		return nil, diags
	}

	body := &openai.AdminOrganizationSpendAlertNewParams{
		Currency:        openai.AdminOrganizationSpendAlertNewParamsCurrency(data.Currency.ValueString()),
		Interval:        openai.AdminOrganizationSpendAlertNewParamsInterval(data.Interval.ValueString()),
		ThresholdAmount: data.ThresholdAmount.ValueInt64(),
		NotificationChannel: openai.AdminOrganizationSpendAlertNewParamsNotificationChannel{
			Type:       constant.Email(notificationChannel.Type.ValueString()),
			Recipients: notificationChannelRecipients,
			SubjectPrefix: (func() param.Opt[string] {
				if notificationChannel.SubjectPrefix.IsKnown() {
					return openai.String(notificationChannel.SubjectPrefix.ValueString())
				}
				return param.Opt[string]{}
			})(),
		},
	}

	return body, nil
}

func (r *SpendAlertResource) getUpdateJSONRequestBody(ctx context.Context, data SpendAlertResourceModel) (*openai.AdminOrganizationSpendAlertUpdateParams, diag.Diagnostics) {
	notificationChannel, diags := data.NotificationChannel.Get(ctx)
	if diags.HasError() {
		return nil, diags
	}

	notificationChannelRecipients, diags := notificationChannel.Recipients.Get(ctx)
	if diags.HasError() {
		return nil, diags
	}

	return &openai.AdminOrganizationSpendAlertUpdateParams{
		Currency:        openai.AdminOrganizationSpendAlertUpdateParamsCurrency(data.Currency.ValueString()),
		Interval:        openai.AdminOrganizationSpendAlertUpdateParamsInterval(data.Interval.ValueString()),
		ThresholdAmount: data.ThresholdAmount.ValueInt64(),
		NotificationChannel: openai.AdminOrganizationSpendAlertUpdateParamsNotificationChannel{
			Type:       constant.Email(notificationChannel.Type.ValueString()),
			Recipients: notificationChannelRecipients,
			SubjectPrefix: (func() param.Opt[string] {
				if notificationChannel.SubjectPrefix.IsKnown() {
					return openai.String(notificationChannel.SubjectPrefix.ValueString())
				}
				return param.Opt[string]{}
			})(),
		},
	}, nil
}
