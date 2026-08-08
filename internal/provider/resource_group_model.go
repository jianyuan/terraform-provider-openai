package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/openai/openai-go/v3"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

func (m *GroupResourceModel) Fill(ctx context.Context, data any) diag.Diagnostics {
	switch v := data.(type) {
	case openai.Group:
		m.Id = supertypes.NewStringValue(v.ID)
		m.Name = supertypes.NewStringValue(v.Name)
		m.CreatedAt = supertypes.NewInt64Value(v.CreatedAt)
		return nil
	case openai.AdminOrganizationGroupUpdateResponse:
		m.Id = supertypes.NewStringValue(v.ID)
		m.Name = supertypes.NewStringValue(v.Name)
		m.CreatedAt = supertypes.NewInt64Value(v.CreatedAt)
		return nil
	default:
		var diags diag.Diagnostics
		diags.AddError("Unknown type", fmt.Sprintf("Unknown type: %T", data))
		return diags
	}
}

func (r *GroupResource) getNewParams(ctx context.Context, data GroupResourceModel) (*openai.AdminOrganizationGroupNewParams, diag.Diagnostics) {
	return &openai.AdminOrganizationGroupNewParams{
		Name: data.Name.ValueString(),
	}, nil
}

func (r *GroupResource) getUpdateParams(ctx context.Context, data GroupResourceModel) (*openai.AdminOrganizationGroupUpdateParams, diag.Diagnostics) {
	return &openai.AdminOrganizationGroupUpdateParams{
		Name: data.Name.ValueString(),
	}, nil
}
