package provider

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/openai/openai-go/v3"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

// Validation normally finishes within the validate call; these bound the wait when OpenAI
// still reports the configuration as pending.
const (
	externalStorageValidationPollInterval = 5 * time.Second
	externalStorageValidationTimeout      = 10 * time.Minute
)

// Only the block for the configuration's provider is filled; the others are null, as the
// ExactlyOneOf validators expect.
func (m *ExternalStorageResourceModel) Fill(ctx context.Context, data openai.ExternalStorageConfiguration) (diags diag.Diagnostics) {
	m.Id = supertypes.NewStringValue(data.ID)
	m.ProjectId = supertypes.NewStringValue(data.ProjectID)
	m.Status = supertypes.NewStringValue(string(data.Status))
	m.Geography = supertypes.NewStringValue(data.Geography)
	m.CreatedAt = supertypes.NewInt64Value(data.CreatedAt)

	m.Aws = supertypes.NewSingleNestedObjectValueOfNull[ExternalStorageResourceModelAws](ctx)
	m.Azure = supertypes.NewSingleNestedObjectValueOfNull[ExternalStorageResourceModelAzure](ctx)
	m.Gcp = supertypes.NewSingleNestedObjectValueOfNull[ExternalStorageResourceModelGcp](ctx)
	m.Oci = supertypes.NewSingleNestedObjectValueOfNull[ExternalStorageResourceModelOci](ctx)

	switch provider := data.Provider.AsAny().(type) {
	case openai.AwsExternalStorageProvider:
		m.Aws = supertypes.NewSingleNestedObjectValueOf(ctx, &ExternalStorageResourceModelAws{
			Bucket:     supertypes.NewStringValue(provider.Bucket),
			RoleArn:    supertypes.NewStringValue(provider.RoleArn),
			AccountId:  supertypes.NewStringValue(provider.AccountID),
			ExternalId: supertypes.NewStringValue(provider.ExternalID),
			Region:     supertypes.NewStringValue(provider.Region),
		})
	case openai.AzureExternalStorageProvider:
		m.Azure = supertypes.NewSingleNestedObjectValueOf(ctx, &ExternalStorageResourceModelAzure{
			TenantId:       supertypes.NewStringValue(provider.TenantID),
			SubscriptionId: supertypes.NewStringValue(provider.SubscriptionID),
			ResourceGroup:  supertypes.NewStringValue(provider.ResourceGroup),
			AccountName:    supertypes.NewStringValue(provider.AccountName),
			Container:      supertypes.NewStringValue(provider.Container),
			Region:         supertypes.NewStringValue(provider.Region),
		})
	case openai.GcpExternalStorageProvider:
		m.Gcp = supertypes.NewSingleNestedObjectValueOf(ctx, &ExternalStorageResourceModelGcp{
			Bucket:                        supertypes.NewStringValue(provider.Bucket),
			WorkloadIdentityProjectNumber: supertypes.NewStringValue(provider.WorkloadIdentityProjectNumber),
			WorkloadIdentityPoolId:        supertypes.NewStringValue(provider.WorkloadIdentityPoolID),
			WorkloadIdentityProviderId:    supertypes.NewStringValue(provider.WorkloadIdentityProviderID),
			Audience:                      supertypes.NewStringValue(provider.Audience),
			Region:                        supertypes.NewStringValue(provider.Region),
		})
	case openai.OciExternalStorageProvider:
		m.Oci = supertypes.NewSingleNestedObjectValueOf(ctx, &ExternalStorageResourceModelOci{
			TenancyOcid: supertypes.NewStringValue(provider.TenancyOcid),
			Region:      supertypes.NewStringValue(provider.Region),
			Bucket:      supertypes.NewStringValue(provider.Bucket),
		})
	default:
		diags.AddError("Unsupported provider", fmt.Sprintf("External storage %s uses provider type %q, which this version of the provider does not support.", data.ID, data.Provider.Type))
	}

	return
}

func (r *ExternalStorageResource) getNewParams(ctx context.Context, data ExternalStorageResourceModel) (*openai.AdminOrganizationExternalStorageNewParams, diag.Diagnostics) {
	var diags diag.Diagnostics
	params := &openai.AdminOrganizationExternalStorageNewParams{
		ProjectID: data.ProjectId.ValueString(),
	}

	switch {
	case !data.Aws.IsNull():
		aws, d := data.Aws.Get(ctx)
		if diags.Append(d...); diags.HasError() {
			return nil, diags
		}
		params.Provider.OfAws = &openai.AdminOrganizationExternalStorageNewParamsProviderAws{
			Bucket:  aws.Bucket.ValueString(),
			RoleArn: aws.RoleArn.ValueString(),
		}
	case !data.Azure.IsNull():
		azure, d := data.Azure.Get(ctx)
		if diags.Append(d...); diags.HasError() {
			return nil, diags
		}
		params.Provider.OfAzure = &openai.AdminOrganizationExternalStorageNewParamsProviderAzure{
			TenantID:       azure.TenantId.ValueString(),
			SubscriptionID: azure.SubscriptionId.ValueString(),
			ResourceGroup:  azure.ResourceGroup.ValueString(),
			AccountName:    azure.AccountName.ValueString(),
			Container:      azure.Container.ValueString(),
		}
	case !data.Gcp.IsNull():
		gcp, d := data.Gcp.Get(ctx)
		if diags.Append(d...); diags.HasError() {
			return nil, diags
		}
		params.Provider.OfGcp = &openai.AdminOrganizationExternalStorageNewParamsProviderGcp{
			Bucket:                        gcp.Bucket.ValueString(),
			WorkloadIdentityProjectNumber: gcp.WorkloadIdentityProjectNumber.ValueString(),
			WorkloadIdentityPoolID:        gcp.WorkloadIdentityPoolId.ValueString(),
			WorkloadIdentityProviderID:    gcp.WorkloadIdentityProviderId.ValueString(),
		}
	case !data.Oci.IsNull():
		oci, d := data.Oci.Get(ctx)
		if diags.Append(d...); diags.HasError() {
			return nil, diags
		}
		params.Provider.OfOci = &openai.AdminOrganizationExternalStorageNewParamsProviderOci{
			TenancyOcid: oci.TenancyOcid.ValueString(),
			Region:      oci.Region.ValueString(),
			Bucket:      oci.Bucket.ValueString(),
		}
	default:
		diags.AddError("Missing provider", "One of `aws`, `azure`, `gcp` or `oci` must be set.")
		return nil, diags
	}

	return params, diags
}

// validateAfterCreate validates a newly created configuration and waits for it to leave
// `pending`, so that a successful apply means OpenAI can use the storage. Any other outcome
// than `validated` deletes the configuration again and fails the create.
func (r *ExternalStorageResource) validateAfterCreate(ctx context.Context, created *openai.ExternalStorageConfiguration) (*openai.ExternalStorageConfiguration, diag.Diagnostics) {
	var diags diag.Diagnostics

	config, err := r.client.Admin.Organization.ExternalStorage.Validate(ctx, created.ID)
	if err == nil && config == nil {
		err = errors.New("got empty response body")
	}

	deadline := time.Now().Add(externalStorageValidationTimeout)
	for err == nil && config.Status == openai.ExternalStorageConfigurationStatusPending {
		if time.Now().After(deadline) {
			err = fmt.Errorf("still pending after %s", externalStorageValidationTimeout)
			break
		}

		select {
		case <-ctx.Done():
			err = ctx.Err()
		case <-time.After(externalStorageValidationPollInterval):
			config, err = r.client.Admin.Organization.ExternalStorage.Get(ctx, created.ID)
			if err == nil && config == nil {
				err = errors.New("got empty response body")
			}
		}
	}

	if err == nil && config.Status != openai.ExternalStorageConfigurationStatusValidated {
		err = fmt.Errorf("validation finished with status %q", config.Status)
	}

	if err == nil {
		return config, diags
	}

	diags.AddError("Validation Error", fmt.Sprintf("External storage %s failed validation, so it was deleted: %s", created.ID, clientErrorDetail(err)))

	// Clean up even if the create was cancelled, so a failed apply leaves nothing behind.
	cleanupCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), time.Minute)
	defer cancel()
	if _, err := r.client.Admin.Organization.ExternalStorage.Delete(cleanupCtx, created.ID); err != nil {
		if apiErr, ok := errors.AsType[*openai.Error](err); !ok || apiErr.StatusCode != http.StatusNotFound {
			diags.AddError("Client Error", fmt.Sprintf("Unable to delete external storage %s after it failed validation. Delete it with the API, or import it into Terraform and destroy it: %s", created.ID, clientErrorDetail(err)))
		}
	}

	return nil, diags
}
