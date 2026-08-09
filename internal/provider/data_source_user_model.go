package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/openai/openai-go/v3"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

func (m *UserDataSourceModel) Fill(ctx context.Context, user openai.OrganizationUser) diag.Diagnostics {
	m.Id = supertypes.NewStringValue(user.ID)
	m.Email = supertypes.NewStringValue(user.Email)
	m.Name = supertypes.NewStringValue(user.Name)
	m.Role = supertypes.NewStringValue(user.Role)
	m.AddedAt = supertypes.NewInt64Value(user.AddedAt)
	return nil
}
