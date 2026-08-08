package sweep

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	intdiag "github.com/jianyuan/terraform-provider-openai/internal/diag"
	"github.com/openai/openai-go/v3"
)

var _ Sweepable = (*sweepResource)(nil)

type sweepResource struct {
	factory    func() resource.Resource
	client     *openai.Client
	attributes map[string]any
}

func NewSweepResource(factory func() resource.Resource, client *openai.Client, attributes map[string]any) *sweepResource {
	return &sweepResource{
		factory:    factory,
		client:     client,
		attributes: attributes,
	}
}

func (sr *sweepResource) Delete(ctx context.Context) error {
	res := sr.factory()

	if res, ok := res.(resource.ResourceWithConfigure); ok {
		var configureResp resource.ConfigureResponse
		res.Configure(ctx, resource.ConfigureRequest{ProviderData: sr.client}, &configureResp)
		if configureResp.Diagnostics.HasError() {
			return intdiag.DiagnosticsError(configureResp.Diagnostics)
		}
	}

	var schemaResp resource.SchemaResponse
	res.Schema(ctx, resource.SchemaRequest{}, &schemaResp)
	if schemaResp.Diagnostics.HasError() {
		return intdiag.DiagnosticsError(schemaResp.Diagnostics)
	}

	state := tfsdk.State{
		Raw:    tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), nil),
		Schema: schemaResp.Schema,
	}

	for k, v := range sr.attributes {
		if diags := state.SetAttribute(ctx, path.Root(k), v); diags.HasError() {
			return intdiag.DiagnosticsError(diags)
		}
	}

	log.Printf("[INFO] Deleting resource: %v", sr.attributes)
	var deleteResp resource.DeleteResponse
	res.Delete(ctx, resource.DeleteRequest{State: state}, &deleteResp)
	if deleteResp.Diagnostics.HasError() {
		return intdiag.DiagnosticsError(deleteResp.Diagnostics)
	}

	return nil
}
