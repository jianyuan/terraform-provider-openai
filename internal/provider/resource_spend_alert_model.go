package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/jianyuan/terraform-provider-openai/internal/apiclient"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

func (m *SpendAlertResourceModel) Fill(ctx context.Context, data apiclient.OrganizationSpendAlert) diag.Diagnostics {
	m.Id = supertypes.NewStringValue(data.Id)
	m.Currency = supertypes.NewStringValue(string(data.Currency))
	m.Interval = supertypes.NewStringValue(string(data.Interval))
	m.ThresholdAmount = supertypes.NewInt64Value(data.ThresholdAmount)
	m.NotificationChannel = supertypes.NewSingleNestedObjectValueOf(ctx, &SpendAlertResourceModelNotificationChannel{
		Type:          supertypes.NewStringValue(string(data.NotificationChannel.Type)),
		Recipients:    supertypes.NewSetValueOfSlice(ctx, data.NotificationChannel.Recipients),
		SubjectPrefix: supertypes.NewStringPointerValue(data.NotificationChannel.SubjectPrefix),
	})
	return nil
}

func (r *SpendAlertResource) resourceMatch(data SpendAlertResourceModel, spendAlert apiclient.OrganizationSpendAlert) bool {
	return data.Id.ValueString() == spendAlert.Id
}

func (r *SpendAlertResource) getCreateJSONRequestBody(ctx context.Context, data SpendAlertResourceModel) (apiclient.CreateOrganizationSpendAlertJSONRequestBody, diag.Diagnostics) {
	notificationChannel, diags := data.NotificationChannel.Get(ctx)
	if diags.HasError() {
		return apiclient.CreateOrganizationSpendAlertJSONRequestBody{}, diags
	}

	notificationChannelRecipients, diags := notificationChannel.Recipients.Get(ctx)
	if diags.HasError() {
		return apiclient.CreateOrganizationSpendAlertJSONRequestBody{}, diags
	}

	return apiclient.CreateOrganizationSpendAlertJSONRequestBody{
		Currency:        apiclient.CreateSpendAlertBodyCurrency(data.Currency.ValueString()),
		Interval:        apiclient.CreateSpendAlertBodyInterval(data.Interval.ValueString()),
		ThresholdAmount: data.ThresholdAmount.ValueInt64(),
		NotificationChannel: apiclient.SpendAlertNotificationChannel{
			Type:          apiclient.SpendAlertNotificationChannelType(notificationChannel.Type.ValueString()),
			Recipients:    notificationChannelRecipients,
			SubjectPrefix: notificationChannel.SubjectPrefix.ValueStringPointer(),
		},
	}, nil
}

func (r *SpendAlertResource) getUpdateJSONRequestBody(ctx context.Context, data SpendAlertResourceModel) (apiclient.UpdateOrganizationSpendAlertJSONRequestBody, diag.Diagnostics) {
	notificationChannel, diags := data.NotificationChannel.Get(ctx)
	if diags.HasError() {
		return apiclient.UpdateOrganizationSpendAlertJSONRequestBody{}, diags
	}

	notificationChannelRecipients, diags := notificationChannel.Recipients.Get(ctx)
	if diags.HasError() {
		return apiclient.UpdateOrganizationSpendAlertJSONRequestBody{}, diags
	}

	return apiclient.UpdateOrganizationSpendAlertJSONRequestBody{
		Currency:        apiclient.CreateSpendAlertBodyCurrency(data.Currency.ValueString()),
		Interval:        apiclient.CreateSpendAlertBodyInterval(data.Interval.ValueString()),
		ThresholdAmount: data.ThresholdAmount.ValueInt64(),
		NotificationChannel: apiclient.SpendAlertNotificationChannel{
			Type:          apiclient.SpendAlertNotificationChannelType(notificationChannel.Type.ValueString()),
			Recipients:    notificationChannelRecipients,
			SubjectPrefix: notificationChannel.SubjectPrefix.ValueStringPointer(),
		},
	}, nil
}
