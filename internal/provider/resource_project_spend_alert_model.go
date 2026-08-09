package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/jianyuan/terraform-provider-openai/internal/openaiparam"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/shared/constant"
)

func (r *ProjectSpendAlertResource) getNewParams(ctx context.Context, data ProjectSpendAlertResourceModel) (*openai.AdminOrganizationProjectSpendAlertNewParams, diag.Diagnostics) {
	notificationChannel, diags := data.NotificationChannel.Get(ctx)
	if diags.HasError() {
		return nil, diags
	}

	notificationChannelRecipients, diags := notificationChannel.Recipients.Get(ctx)
	if diags.HasError() {
		return nil, diags
	}

	return &openai.AdminOrganizationProjectSpendAlertNewParams{
		Currency:        openai.AdminOrganizationProjectSpendAlertNewParamsCurrency(data.Currency.ValueString()),
		Interval:        openai.AdminOrganizationProjectSpendAlertNewParamsInterval(data.Interval.ValueString()),
		ThresholdAmount: data.ThresholdAmount.ValueInt64(),
		NotificationChannel: openai.AdminOrganizationProjectSpendAlertNewParamsNotificationChannel{
			Type:          constant.Email(notificationChannel.Type.ValueString()),
			Recipients:    notificationChannelRecipients,
			SubjectPrefix: openaiparam.FromString(notificationChannel.SubjectPrefix),
		},
	}, nil
}

func (r *ProjectSpendAlertResource) getUpdateParams(ctx context.Context, data ProjectSpendAlertResourceModel) (*openai.AdminOrganizationProjectSpendAlertUpdateParams, diag.Diagnostics) {
	notificationChannel, diags := data.NotificationChannel.Get(ctx)
	if diags.HasError() {
		return nil, diags
	}

	notificationChannelRecipients, diags := notificationChannel.Recipients.Get(ctx)
	if diags.HasError() {
		return nil, diags
	}

	return &openai.AdminOrganizationProjectSpendAlertUpdateParams{
		Currency:        openai.AdminOrganizationProjectSpendAlertUpdateParamsCurrency(data.Currency.ValueString()),
		Interval:        openai.AdminOrganizationProjectSpendAlertUpdateParamsInterval(data.Interval.ValueString()),
		ThresholdAmount: data.ThresholdAmount.ValueInt64(),
		NotificationChannel: openai.AdminOrganizationProjectSpendAlertUpdateParamsNotificationChannel{
			Type:          constant.Email(notificationChannel.Type.ValueString()),
			Recipients:    notificationChannelRecipients,
			SubjectPrefix: openaiparam.FromString(notificationChannel.SubjectPrefix),
		},
	}, nil
}
