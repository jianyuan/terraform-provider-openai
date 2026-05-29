package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/jianyuan/terraform-provider-openai/internal/apiclient"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

func (m *OrganizationDataRetentionResourceModel) Fill(ctx context.Context, role apiclient.OrganizationDataRetention) diag.Diagnostics {
	m.Type = supertypes.NewStringValue(string(role.Type))
	return nil
}

func (r *OrganizationDataRetentionResource) getCreateJSONRequestBody(ctx context.Context, data OrganizationDataRetentionResourceModel) (apiclient.UpdateOrganizationDataRetentionJSONRequestBody, diag.Diagnostics) {
	var diags diag.Diagnostics
	return apiclient.UpdateOrganizationDataRetentionJSONRequestBody{
		RetentionType: apiclient.UpdateOrganizationDataRetentionBodyRetentionType(data.Type.ValueString()),
	}, diags
}

func (r *OrganizationDataRetentionResource) getUpdateJSONRequestBody(ctx context.Context, data OrganizationDataRetentionResourceModel) (apiclient.UpdateOrganizationDataRetentionJSONRequestBody, diag.Diagnostics) {
	var diags diag.Diagnostics
	return apiclient.UpdateOrganizationDataRetentionJSONRequestBody{
		RetentionType: apiclient.UpdateOrganizationDataRetentionBodyRetentionType(data.Type.ValueString()),
	}, diags
}
