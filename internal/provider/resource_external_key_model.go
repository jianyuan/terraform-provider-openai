package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/openai/openai-go/v3"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
	"github.com/samber/lo"
)

// The API returns the provider-specific fields flat, beside the common ones. Only the block
// for the key's type is filled; the others are null, as the ExactlyOneOf validators expect.
func (m *ExternalKeyResourceModel) Fill(ctx context.Context, data openai.ExternalKey) (diags diag.Diagnostics) {
	m.Id = supertypes.NewStringValue(data.ID)
	m.Name = supertypes.NewStringValue(data.Name)
	m.ApiProjectIds = supertypes.NewSetValueOfSlice(ctx, lo.Uniq(data.APIProjectIDs))
	m.CreatedAt = supertypes.NewInt64Value(data.CreatedAt)

	m.Aws = supertypes.NewSingleNestedObjectValueOfNull[ExternalKeyResourceModelAws](ctx)
	m.Azure = supertypes.NewSingleNestedObjectValueOfNull[ExternalKeyResourceModelAzure](ctx)
	m.Gcp = supertypes.NewSingleNestedObjectValueOfNull[ExternalKeyResourceModelGcp](ctx)

	switch data.Type {
	case "aws":
		m.Aws = supertypes.NewSingleNestedObjectValueOf(ctx, &ExternalKeyResourceModelAws{
			RoleArn:    supertypes.NewStringValue(data.RoleArn),
			KmsArn:     supertypes.NewStringValue(data.KmsArn),
			ExternalId: supertypes.NewStringValue(data.ExternalID),
		})
	case "azure":
		m.Azure = supertypes.NewSingleNestedObjectValueOf(ctx, &ExternalKeyResourceModelAzure{
			TenantId: supertypes.NewStringValue(data.TenantID),
			VaultUri: supertypes.NewStringValue(data.VaultUri),
			KeyName:  supertypes.NewStringValue(data.KeyName),
		})
	case "gcp":
		m.Gcp = supertypes.NewSingleNestedObjectValueOf(ctx, &ExternalKeyResourceModelGcp{
			WorkloadIdentityProjectNumber: supertypes.NewStringValue(data.WorkloadIdentityProjectNumber),
			WorkloadIdentityPoolId:        supertypes.NewStringValue(data.WorkloadIdentityPoolID),
			WorkloadIdentityProviderId:    supertypes.NewStringValue(data.WorkloadIdentityProviderID),
			Audience:                      supertypes.NewStringValue(data.Audience),
			KmsProjectId:                  supertypes.NewStringValue(data.KmsProjectID),
			KmsKeyRingName:                supertypes.NewStringValue(data.KmsKeyRingName),
			KmsKeyName:                    supertypes.NewStringValue(data.KmsKeyName),
			KmsKeyLocation:                supertypes.NewStringValue(data.KmsKeyLocation),
		})
	default:
		diags.AddError("Unsupported key type", fmt.Sprintf("External key %s has type %q, which this version of the provider does not support.", data.ID, data.Type))
	}

	return
}

func (r *ExternalKeyResource) getNewParams(ctx context.Context, data ExternalKeyResourceModel) (*openai.AdminOrganizationExternalKeyNewParams, diag.Diagnostics) {
	var diags diag.Diagnostics
	name := data.Name.ValueString()
	// The API only accepts keys registered for EKM.
	purpose := openai.String(string(openai.ExternalKeyPurposeEkm))

	switch {
	case !data.Aws.IsNull():
		aws, d := data.Aws.Get(ctx)
		if diags.Append(d...); diags.HasError() {
			return nil, diags
		}
		return &openai.AdminOrganizationExternalKeyNewParams{
			OfAws: &openai.AdminOrganizationExternalKeyNewParamsBodyAws{
				Name:       name,
				Purpose:    purpose,
				RoleArn:    aws.RoleArn.ValueString(),
				KmsArn:     aws.KmsArn.ValueString(),
				ExternalID: aws.ExternalId.ValueString(),
			},
		}, diags
	case !data.Azure.IsNull():
		azure, d := data.Azure.Get(ctx)
		if diags.Append(d...); diags.HasError() {
			return nil, diags
		}
		return &openai.AdminOrganizationExternalKeyNewParams{
			OfAzure: &openai.AdminOrganizationExternalKeyNewParamsBodyAzure{
				Name:     name,
				Purpose:  purpose,
				TenantID: azure.TenantId.ValueString(),
				VaultURI: azure.VaultUri.ValueString(),
				KeyName:  azure.KeyName.ValueString(),
			},
		}, diags
	case !data.Gcp.IsNull():
		gcp, d := data.Gcp.Get(ctx)
		if diags.Append(d...); diags.HasError() {
			return nil, diags
		}
		return &openai.AdminOrganizationExternalKeyNewParams{
			OfGcp: &openai.AdminOrganizationExternalKeyNewParamsBodyGcp{
				Name:                          name,
				Purpose:                       purpose,
				WorkloadIdentityProjectNumber: gcp.WorkloadIdentityProjectNumber.ValueString(),
				WorkloadIdentityPoolID:        gcp.WorkloadIdentityPoolId.ValueString(),
				WorkloadIdentityProviderID:    gcp.WorkloadIdentityProviderId.ValueString(),
				Audience:                      gcp.Audience.ValueString(),
				KmsProjectID:                  gcp.KmsProjectId.ValueString(),
				KmsKeyRingName:                gcp.KmsKeyRingName.ValueString(),
				KmsKeyName:                    gcp.KmsKeyName.ValueString(),
				KmsKeyLocation:                gcp.KmsKeyLocation.ValueString(),
			},
		}, diags
	}

	diags.AddError("Missing key", "One of `aws`, `azure` or `gcp` must be set.")
	return nil, diags
}
