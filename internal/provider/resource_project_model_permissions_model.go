package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/openai/openai-go/v3"
)

func (r *ProjectModelPermissionsResource) getNewParams(ctx context.Context, data ProjectModelPermissionsResourceModel) (*openai.AdminOrganizationProjectModelPermissionUpdateParams, diag.Diagnostics) {
	return r.getUpdateParams(ctx, data)
}

func (r *ProjectModelPermissionsResource) getUpdateParams(ctx context.Context, data ProjectModelPermissionsResourceModel) (*openai.AdminOrganizationProjectModelPermissionUpdateParams, diag.Diagnostics) {
	modelIds, diags := data.ModelIds.Get(ctx)
	if diags.HasError() {
		return nil, diags
	}

	return &openai.AdminOrganizationProjectModelPermissionUpdateParams{
		Mode:     openai.AdminOrganizationProjectModelPermissionUpdateParamsMode(data.Mode.ValueString()),
		ModelIDs: modelIds,
	}, nil
}
