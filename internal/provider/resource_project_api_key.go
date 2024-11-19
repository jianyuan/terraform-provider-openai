package provider

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/iancoleman/orderedmap"
	"github.com/jianyuan/terraform-provider-openai/internal/apiclient"
	"github.com/jianyuan/terraform-provider-openai/internal/must"
)

//go:embed scopes.json
var rawApiKeyPermissions []byte

var apiKeyReadOnlyScope = "api.all.read"
var apiKeyPermissionAttributes map[string]schema.Attribute
var apiKeyPermissionScopes map[string]map[string][]string

func init() {
	type Permission struct {
		Name                string                `json:"name"`
		Description         string                `json:"description"`
		PermissionsToScopes orderedmap.OrderedMap `json:"permissions_to_scopes"`
		Endpoints           []string              `json:"endpoints"`
	}

	var permissions []Permission
	must.Do(json.Unmarshal(rawApiKeyPermissions, &permissions))

	apiKeyPermissionAttributes = make(map[string]schema.Attribute, len(permissions))
	apiKeyPermissionScopes = make(map[string]map[string][]string, len(permissions))

	for _, permission := range permissions {
		attribute := strings.ToLower(permission.Name)
		attribute = strings.ReplaceAll(attribute, " ", "_")
		attribute = strings.ReplaceAll(attribute, "-", "_")

		permissionKeys := make([]string, 0, len(permission.PermissionsToScopes.Keys()))
		permissionQuoted := make([]string, 0, len(permission.PermissionsToScopes.Keys()))
		apiKeyPermissionScopes[attribute] = make(map[string][]string, len(permission.PermissionsToScopes.Keys()))
		for _, key := range permission.PermissionsToScopes.Keys() {
			permissionKeys = append(permissionKeys, key)
			permissionQuoted = append(permissionQuoted, fmt.Sprintf("`%s`", key))

			scopesInterface, _ := permission.PermissionsToScopes.Get(key)
			scopes, _ := scopesInterface.([]interface{})
			for _, scopeInterface := range scopes {
				scope, _ := scopeInterface.(string)
				apiKeyPermissionScopes[attribute][key] = append(apiKeyPermissionScopes[attribute][key], scope)
			}
		}

		endpointQuoted := make([]string, 0, len(permission.Endpoints))
		for _, endpoint := range permission.Endpoints {
			endpointQuoted = append(endpointQuoted, fmt.Sprintf("`%s`", endpoint))
		}

		var valueNoun string
		if len(permissionKeys) == 1 {
			valueNoun = "value"
		} else {
			valueNoun = "values"
		}

		apiKeyPermissionAttributes[attribute] = schema.StringAttribute{
			MarkdownDescription: fmt.Sprintf(
				"%s. %s. Valid %s: %s. If omitted, the API key will not have access.",
				permission.Description,
				strings.Join(endpointQuoted, ", "),
				valueNoun,
				strings.Join(permissionQuoted, ", "),
			),
			Optional: true,
			Validators: []validator.String{
				stringvalidator.OneOf(permissionKeys...),
			},
		}
	}
}

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
	Id               types.String                  `tfsdk:"id"`
	ProjectId        types.String                  `tfsdk:"project_id"`
	ServiceAccountId types.String                  `tfsdk:"service_account_id"`
	Name             types.String                  `tfsdk:"name"`
	ReadOnly         types.Bool                    `tfsdk:"read_only"`
	Permissions      *ProjectApiKeyPermissionModel `tfsdk:"permissions"`
	Scopes           types.Set                     `tfsdk:"scopes"`
	Created          types.Int64                   `tfsdk:"created"`
	RedactedKey      types.String                  `tfsdk:"redacted_key"`
}

type ProjectApiKeyPermissionModel struct {
	Models            types.String `tfsdk:"models"`
	ModelCapabilities types.String `tfsdk:"model_capabilities"`
	Assistants        types.String `tfsdk:"assistants"`
	Threads           types.String `tfsdk:"threads"`
	FineTuning        types.String `tfsdk:"fine_tuning"`
	Files             types.String `tfsdk:"files"`
}

func (m *ProjectApiKeyPermissionModel) Fill(apiKey apiclient.ApiKey) {
	m.Models = types.StringNull()
	m.ModelCapabilities = types.StringNull()
	m.Assistants = types.StringNull()
	m.Threads = types.StringNull()
	m.FineTuning = types.StringNull()
	m.Files = types.StringNull()

	for key, permissions := range apiKeyPermissionScopes {
		for permission, scopes := range permissions {
			found := 0
			for _, scope := range scopes {
				for _, apiKeyScope := range apiKey.Scopes {
					if apiKeyScope == scope {
						found += 1
						break
					}
				}
			}

			if found == len(scopes) {
				switch key {
				case "models":
					m.Models = types.StringValue(permission)
				case "model_capabilities":
					m.ModelCapabilities = types.StringValue(permission)
				case "assistants":
					m.Assistants = types.StringValue(permission)
				case "threads":
					m.Threads = types.StringValue(permission)
				case "fine_tuning":
					m.FineTuning = types.StringValue(permission)
				case "files":
					m.Files = types.StringValue(permission)
				}
				break
			}
		}
	}
}

func (m *ProjectApiKeyPermissionModel) Equal(other *ProjectApiKeyPermissionModel) bool {
	return (m != nil && other != nil) &&
		m.Models.Equal(other.Models) &&
		m.ModelCapabilities.Equal(other.ModelCapabilities) &&
		m.Assistants.Equal(other.Assistants) &&
		m.Threads.Equal(other.Threads) &&
		m.FineTuning.Equal(other.FineTuning) &&
		m.Files.Equal(other.Files)
}

func (m *ProjectApiKeyPermissionModel) Scopes() []string {
	getScopes := func(key string, attr types.String) []string {
		if attr.IsNull() {
			return []string{}
		}
		return apiKeyPermissionScopes[key][attr.ValueString()]
	}

	var scopes []string
	scopes = append(scopes, getScopes("models", m.Models)...)
	scopes = append(scopes, getScopes("model_capabilities", m.ModelCapabilities)...)
	scopes = append(scopes, getScopes("assistants", m.Assistants)...)
	scopes = append(scopes, getScopes("threads", m.Threads)...)
	scopes = append(scopes, getScopes("fine_tuning", m.FineTuning)...)
	scopes = append(scopes, getScopes("files", m.Files)...)
	return scopes
}

func (m *ProjectApiKeyResourceModel) PartialFill(apiKey apiclient.ApiKey) {
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

	if len(apiKey.Scopes) == 0 {
		m.ReadOnly = types.BoolNull()
		m.Permissions = nil
	} else {
		if len(apiKey.Scopes) == 1 && apiKey.Scopes[0] == apiKeyReadOnlyScope {
			m.ReadOnly = types.BoolValue(true)
			m.Permissions = nil
		} else {
			m.ReadOnly = types.BoolNull()
			m.Permissions = &ProjectApiKeyPermissionModel{}
			m.Permissions.Fill(apiKey)
		}
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
			"read_only": schema.BoolAttribute{
				MarkdownDescription: "Whether the API key is read-only. If omitted, the API key will have full permissions.",
				Optional:            true,
				Validators: []validator.Bool{
					boolvalidator.ConflictsWith(path.MatchRoot("permissions")),
				},
			},
			"scopes": schema.SetAttribute{
				MarkdownDescription: "The scopes of the API key.",
				Computed:            true,
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

		Blocks: map[string]schema.Block{
			"permissions": schema.SingleNestedBlock{
				MarkdownDescription: "The permission of the API key. If omitted, the API key will have full permissions.",
				Attributes:          apiKeyPermissionAttributes,
				Validators: []validator.Object{
					objectvalidator.ConflictsWith(path.MatchRoot("read_only")),
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
				OpenaiProject: data.ProjectId.ValueStringPointer(),
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
		"TODO",
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
	if !data.Name.IsNull() || !data.ReadOnly.IsNull() || data.Permissions != nil {
		var scopes []string
		if !data.ReadOnly.IsNull() && data.ReadOnly.ValueBool() {
			scopes = []string{apiKeyReadOnlyScope}
		} else {
			if data.Permissions == nil {
				scopes = []string{}
			} else {
				scopes = data.Permissions.Scopes()
			}
		}

		name := data.Name.ValueString()

		if data.ProjectId.IsNull() {
			httpResp, err := r.client.UpdateOrganizationApiKeyWithResponse(
				ctx,
				"TODO",
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
				"TODO",
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
		"TODO",
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

	if !plan.Name.Equal(state.Name) || !plan.ReadOnly.Equal(state.ReadOnly) || !plan.Permissions.Equal(state.Permissions) {
		var scopes []string
		if !plan.ReadOnly.IsNull() && plan.ReadOnly.ValueBool() {
			scopes = []string{apiKeyReadOnlyScope}
		} else {
			if plan.Permissions == nil {
				scopes = []string{}
			} else {
				scopes = plan.Permissions.Scopes()
			}
		}

		name := plan.Name.ValueString()

		var apiKey *apiclient.ApiKey
		if plan.ProjectId.IsNull() {
			httpResp, err := r.client.UpdateOrganizationApiKeyWithResponse(
				ctx,
				"TODO",
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
				"TODO",
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
			"TODO",
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
			"TODO",
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
	projectId, id, err := SplitTwoPartId(req.ID, "project-id", "id")
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Error parsing ID: %s", err.Error()))
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), projectId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}
