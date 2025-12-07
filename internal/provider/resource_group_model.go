package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/jianyuan/terraform-provider-openai/internal/apiclient"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

func (m *GroupResourceModel) Fill(ctx context.Context, data any) diag.Diagnostics {
	var diag diag.Diagnostics

	switch v := data.(type) {
	case apiclient.GroupResponse:
		m.Id = supertypes.NewStringValue(v.Id)
		m.Name = supertypes.NewStringValue(v.Name)
		m.CreatedAt = supertypes.NewInt64Value(v.CreatedAt)
		return nil
	case apiclient.GroupResourceWithSuccess:
		m.Id = supertypes.NewStringValue(v.Id)
		m.Name = supertypes.NewStringValue(v.Name)
		m.CreatedAt = supertypes.NewInt64Value(v.CreatedAt)
		return nil
	default:
		diag.AddError("Unknown type", fmt.Sprintf("Unknown type: %T", data))
		return diag
	}
}

func (r *GroupResource) resourceMatch(data GroupResourceModel, group apiclient.GroupResponse) bool {
	return data.Id.ValueString() == group.Id
}

func (r *GroupResource) getCreateJSONRequestBody(ctx context.Context, data GroupResourceModel) (apiclient.CreateGroupJSONRequestBody, diag.Diagnostics) {
	return apiclient.CreateGroupJSONRequestBody{
		Name: data.Name.ValueString(),
	}, nil
}

func (r *GroupResource) getUpdateJSONRequestBody(ctx context.Context, data GroupResourceModel) (apiclient.UpdateGroupJSONRequestBody, diag.Diagnostics) {
	return apiclient.UpdateGroupJSONRequestBody{
		Name: data.Name.ValueString(),
	}, nil
}
