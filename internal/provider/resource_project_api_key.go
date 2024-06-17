package provider

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/terraform-provider-openai/internal/apiclient"
)

var _ resource.Resource = &ProjectApiKeyResource{}
var _ resource.ResourceWithImportState = &ProjectApiKeyResource{}

func NewProjectApiKeyResource() resource.Resource {
	return &ProjectApiKeyResource{}
}

type ProjectApiKeyResource struct {
	baseResource
}

// ProjectApiKeyResourceModel describes the resource data model.
type ProjectApiKeyResourceModel struct {
	Id               types.String `tfsdk:"id"`
	OrganizationId   types.String `tfsdk:"organization_id"`
	ProjectId        types.String `tfsdk:"project_id"`
	ServiceAccountId types.String `tfsdk:"service_account_id"`
	Name             types.String `tfsdk:"name"`
	Scopes           types.Set    `tfsdk:"scopes"`
	Created          types.Int64  `tfsdk:"created"`
	RedactedKey      types.String `tfsdk:"redacted_key"`
}

func (m *ProjectApiKeyResourceModel) PartialFill(apiKey apiclient.ApiKey) {
	m.OrganizationId = types.StringValue(apiKey.Organization.Id)

	if apiKey.Project == nil {
		m.ProjectId = types.StringNull()
	} else {
		m.ProjectId = types.StringValue(apiKey.Project.Id)
	}

	if apiKey.Name == nil || *apiKey.Name == "" {
		m.Name = types.StringNull()
	} else {
		m.Name = types.StringPointerValue(apiKey.Name)

	}

	if len(apiKey.Scopes) == 0 {
		m.Scopes = types.SetNull(types.StringType)
	} else {
		scopeElements := make([]attr.Value, len(apiKey.Scopes))
		for i, scope := range apiKey.Scopes {
			scopeElements[i] = types.StringValue(scope)
		}
		m.Scopes = types.SetValueMust(types.StringType, scopeElements)
	}

	m.Created = types.Int64Value(apiKey.Created)
	m.RedactedKey = types.StringValue(apiKey.SensitiveId)
}

func (r *ProjectApiKeyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project_api_key"
}

func (r *ProjectApiKeyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Project API key resource.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The API key.",
				Computed:            true,
				Sensitive:           true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"organization_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the organization to which the project belongs.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"project_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the project. If not set, the default project will be used.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"service_account_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the service account to which the API key belongs. IDs can include letters, numbers, and hyphens.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the API key.",
				Optional:            true,
			},
			"scopes": schema.SetAttribute{
				MarkdownDescription: "The scopes of the API key. If not set, all scopes will be used.",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"created": schema.Int64Attribute{
				MarkdownDescription: "The timestamp when the API key was created.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"redacted_key": schema.StringAttribute{
				MarkdownDescription: "The redacted API key.",
				Computed:            true,
				Sensitive:           true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *ProjectApiKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ProjectApiKeyResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create the API key
	{
		httpResp, err := r.client.CreateServiceAccountKeyWithResponse(
			ctx,
			&apiclient.CreateServiceAccountKeyParams{
				OpenaiOrganization: data.OrganizationId.ValueStringPointer(),
				OpenaiProject:      data.ProjectId.ValueStringPointer(),
			},
			apiclient.CreateServiceAccountKeyJSONRequestBody{
				Id: data.ServiceAccountId.ValueString(),
			},
		)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read, got error: %s", err))
			return
		}
		if httpResp.StatusCode() != http.StatusOK {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create, got status code %d: %s", httpResp.StatusCode(), string(httpResp.Body)))
			return
		}
		if len(httpResp.JSON200.Secret) == 0 {
			resp.Diagnostics.AddError("Client Error", "Unable to create, got no secret")
			return
		}

		apiKey := httpResp.JSON200.Secret[0]

		// The rest of the data will be filled in the next step
		data.Id = types.StringValue(apiKey.SensitiveId)
		data.Created = types.Int64Value(apiKey.Created)
	}

	// Get the redacted key
	apiKey, err := r.readKey(
		ctx,
		data.OrganizationId.ValueString(),
		data.ProjectId.ValueStringPointer(),
		func(apiKey apiclient.ApiKey) bool {
			return data.Created.ValueInt64() == apiKey.Created && MatchStringWithMask(data.Id.ValueString(), apiKey.SensitiveId)
		},
	)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}

	data.RedactedKey = types.StringValue(apiKey.SensitiveId)

	// Update the name and scopes if they are set
	if !data.Name.IsNull() || !data.Scopes.IsNull() {
		var scopes []string

		if data.Scopes.IsNull() {
			scopes = []string{}
		} else {
			resp.Diagnostics.Append(data.Scopes.ElementsAs(ctx, &scopes, false)...)

			if resp.Diagnostics.HasError() {
				return
			}
		}

		name := data.Name.ValueString()

		if data.ProjectId.IsNull() {
			httpResp, err := r.client.UpdateOrganizationApiKeyWithResponse(
				ctx,
				data.OrganizationId.ValueString(),
				apiclient.UpdateOrganizationApiKeyJSONRequestBody{
					Action:      "update",
					CreatedAt:   data.Created.ValueInt64(),
					Name:        &name,
					RedactedKey: data.RedactedKey.ValueString(),
					Scopes:      &scopes,
				},
			)

			if err != nil {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update, got error: %s", err))
				return
			}

			apiKey = httpResp.JSON200.Key
		} else {
			httpResp, err := r.client.UpdateProjectApiKeyWithResponse(
				ctx,
				data.OrganizationId.ValueString(),
				data.ProjectId.ValueString(),
				apiclient.UpdateProjectApiKeyJSONRequestBody{
					Action:      "update",
					CreatedAt:   data.Created.ValueInt64(),
					Name:        &name,
					RedactedKey: data.RedactedKey.ValueString(),
					Scopes:      &scopes,
				},
			)

			if err != nil {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update, got error: %s", err))
				return
			}

			apiKey = httpResp.JSON200.Key
		}

	}

	data.PartialFill(*apiKey)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r ProjectApiKeyResource) readKey(ctx context.Context, organisationId string, projectId *string, matchFunc func(apiclient.ApiKey) bool) (*apiclient.ApiKey, error) {
	var apiKeys []apiclient.ApiKey

	if projectId == nil {
		httpResp, err := r.client.GetOrganizationApiKeysWithResponse(
			ctx,
			organisationId,
			&apiclient.GetOrganizationApiKeysParams{
				ExcludeProjectApiKeys: Pointer(true),
			},
		)
		if err != nil {
			return nil, fmt.Errorf("Unable to read, got error: %s", err)
		}
		if httpResp.StatusCode() != http.StatusOK {
			return nil, fmt.Errorf("Unable to read, got status code %d: %s", httpResp.StatusCode(), string(httpResp.Body))
		}

		apiKeys = httpResp.JSON200.Data
	} else {
		httpResp, err := r.client.GetProjectApiKeysWithResponse(
			ctx,
			organisationId,
			*projectId,
		)
		if err != nil {
			return nil, fmt.Errorf("Unable to read, got error: %s", err)
		}
		if httpResp.StatusCode() != http.StatusOK {
			return nil, fmt.Errorf("Unable to read, got status code %d: %s", httpResp.StatusCode(), string(httpResp.Body))
		}

		apiKeys = httpResp.JSON200.Data
	}

	for _, apiKey := range apiKeys {
		if matchFunc(apiKey) {
			return &apiKey, nil
		}
	}

	return nil, errors.New("Unable to read, got no matching key")
}

func (r *ProjectApiKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ProjectApiKeyResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	found, err := r.readKey(
		ctx,
		data.OrganizationId.ValueString(),
		data.ProjectId.ValueStringPointer(),
		func(apiKey apiclient.ApiKey) bool {
			return (data.Created.IsNull() || data.Created.ValueInt64() == apiKey.Created) && MatchStringWithMask(data.Id.ValueString(), apiKey.SensitiveId)
		},
	)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}

	data.ServiceAccountId = types.StringValue(found.User.Id)
	data.PartialFill(*found)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProjectApiKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state ProjectApiKeyResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.Name.Equal(state.Name) || !plan.Scopes.Equal(state.Scopes) {
		var scopes []string

		if plan.Scopes.IsNull() {
			scopes = []string{}
		} else {
			resp.Diagnostics.Append(plan.Scopes.ElementsAs(ctx, &scopes, false)...)
			if resp.Diagnostics.HasError() {
				return
			}
		}

		name := plan.Name.ValueString()

		var apiKey *apiclient.ApiKey
		if plan.ProjectId.IsNull() {
			httpResp, err := r.client.UpdateOrganizationApiKeyWithResponse(
				ctx,
				plan.OrganizationId.ValueString(),
				apiclient.UpdateOrganizationApiKeyJSONRequestBody{
					Action:      "update",
					CreatedAt:   plan.Created.ValueInt64(),
					Name:        &name,
					RedactedKey: plan.RedactedKey.ValueString(),
					Scopes:      &scopes,
				},
			)
			if err != nil {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update, got error: %s", err))
				return
			}
			if httpResp.StatusCode() != http.StatusOK {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update, got status code %d: %s", httpResp.StatusCode(), string(httpResp.Body)))
				return
			}

			apiKey = httpResp.JSON200.Key
		} else {
			httpResp, err := r.client.UpdateProjectApiKeyWithResponse(
				ctx,
				plan.OrganizationId.ValueString(),
				plan.ProjectId.ValueString(),
				apiclient.UpdateProjectApiKeyJSONRequestBody{
					Action:      "update",
					CreatedAt:   plan.Created.ValueInt64(),
					Name:        &name,
					RedactedKey: plan.RedactedKey.ValueString(),
					Scopes:      &scopes,
				},
			)
			if err != nil {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update, got error: %s", err))
				return
			}
			if httpResp.StatusCode() != http.StatusOK {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update, got status code %d: %s", httpResp.StatusCode(), string(httpResp.Body)))
				return
			}

			apiKey = httpResp.JSON200.Key
		}

		state.PartialFill(*apiKey)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ProjectApiKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ProjectApiKeyResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.ProjectId.IsNull() {
		httpResp, err := r.client.UpdateOrganizationApiKeyWithResponse(
			ctx,
			data.OrganizationId.ValueString(),
			apiclient.UpdateOrganizationApiKeyJSONRequestBody{
				Action:      "delete",
				CreatedAt:   data.Created.ValueInt64(),
				RedactedKey: data.RedactedKey.ValueString(),
			},
		)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update, got error: %s", err))
			return
		}
		if httpResp.StatusCode() != http.StatusOK {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update, got status code %d: %s", httpResp.StatusCode(), string(httpResp.Body)))
			return
		}
	} else {
		httpResp, err := r.client.UpdateProjectApiKeyWithResponse(
			ctx,
			data.OrganizationId.ValueString(),
			data.ProjectId.ValueString(),
			apiclient.UpdateProjectApiKeyJSONRequestBody{
				Action:      "delete",
				CreatedAt:   data.Created.ValueInt64(),
				RedactedKey: data.RedactedKey.ValueString(),
			},
		)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update, got error: %s", err))
			return
		}
		if httpResp.StatusCode() != http.StatusOK {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update, got status code %d: %s", httpResp.StatusCode(), string(httpResp.Body)))
			return
		}
	}
}

func (r *ProjectApiKeyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	organizationId, projectId, id, err := SplitThreePartId(req.ID, "organization-id", "project-id", "id")
	if err == nil {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("organization_id"), organizationId)...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), projectId)...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
		return

	}

	organizationId, id, err = SplitTwoPartId(req.ID, "organization-id", "id")
	if err == nil {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("organization_id"), organizationId)...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
		return
	}

	resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Error parsing ID: %s", err.Error()))
}
