package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/jianyuan/terraform-provider-openai/internal/openaiparam"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/shared/constant"
)

func (r *SpendAlertResource) getNewParams(ctx context.Context, data SpendAlertResourceModel) (*openai.AdminOrganizationSpendAlertNewParams, diag.Diagnostics) {
	notificationChannel, diags := data.NotificationChannel.Get(ctx)
	if diags.HasError() {
		return nil, diags
	}

	notificationChannelRecipients, diags := notificationChannel.Recipients.Get(ctx)
	if diags.HasError() {
		return nil, diags
	}

	return &openai.AdminOrganizationSpendAlertNewParams{
		Currency:        openai.AdminOrganizationSpendAlertNewParamsCurrency(data.Currency.ValueString()),
		Interval:        openai.AdminOrganizationSpendAlertNewParamsInterval(data.Interval.ValueString()),
		ThresholdAmount: data.ThresholdAmount.ValueInt64(),
		NotificationChannel: openai.AdminOrganizationSpendAlertNewParamsNotificationChannel{
			Type:          constant.Email(notificationChannel.Type.ValueString()),
			Recipients:    notificationChannelRecipients,
			SubjectPrefix: openaiparam.FromString(notificationChannel.SubjectPrefix),
		},
	}, nil
}

func (r *SpendAlertResource) getUpdateParams(ctx context.Context, data SpendAlertResourceModel) (*openai.AdminOrganizationSpendAlertUpdateParams, diag.Diagnostics) {
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
			Type:          constant.Email(notificationChannel.Type.ValueString()),
			Recipients:    notificationChannelRecipients,
			SubjectPrefix: openaiparam.FromString(notificationChannel.SubjectPrefix),
		},
	}, nil
}
