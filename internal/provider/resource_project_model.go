package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/jianyuan/terraform-provider-openai/internal/openaiparam"
	"github.com/openai/openai-go/v3"
)

func (r *ProjectResource) getNewParams(ctx context.Context, data ProjectResourceModel) (*openai.AdminOrganizationProjectNewParams, diag.Diagnostics) {
	body := &openai.AdminOrganizationProjectNewParams{
		Name:          data.Name.ValueString(),
		ExternalKeyID: openaiparam.FromString(data.ExternalKeyId),
		Geography:     openaiparam.FromString(data.Geography),
	}

	return body, nil
}

func (r *ProjectResource) getUpdateParams(ctx context.Context, data ProjectResourceModel) (*openai.AdminOrganizationProjectUpdateParams, diag.Diagnostics) {
	body := &openai.AdminOrganizationProjectUpdateParams{
		Name:          openai.String(data.Name.ValueString()),
		ExternalKeyID: openaiparam.FromString(data.ExternalKeyId),
	}

	return body, nil
}
