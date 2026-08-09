package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/openai/openai-go/v3"
)

func (r *DataRetentionResource) getNewParams(ctx context.Context, data DataRetentionResourceModel) (*openai.AdminOrganizationDataRetentionUpdateParams, diag.Diagnostics) {
	return r.getUpdateParams(ctx, data)
}

func (r *DataRetentionResource) getUpdateParams(ctx context.Context, data DataRetentionResourceModel) (*openai.AdminOrganizationDataRetentionUpdateParams, diag.Diagnostics) {
	return &openai.AdminOrganizationDataRetentionUpdateParams{
		RetentionType: openai.AdminOrganizationDataRetentionUpdateParamsRetentionType(data.Type.ValueString()),
	}, nil
}
