import type { DataSource, Resource } from "./schema";

export const DATASOURCES: Array<DataSource> = [
  {
    name: "groups",
    description: "Lists all groups in the organization.",
    api: {
      model: "GroupResponse",
      readStrategy: "paginate",
      readMethod: "ListGroups",
      readCursorParam: "Next",
    },
    attributes: [
      {
        name: "groups",
        type: "set_nested",
        description: "List of groups.",
        computedOptionalRequired: "computed",
        attributes: [
          {
            name: "id",
            type: "string",
            description: "Identifier for the group.",
            computedOptionalRequired: "computed",
          },
          {
            name: "name",
            type: "string",
            description: "Human readable name for the group.",
            computedOptionalRequired: "computed",
          },
          {
            name: "is_scim_managed",
            type: "bool",
            description: "Whether the group is managed through SCIM.",
            computedOptionalRequired: "computed",
          },
          {
            name: "created_at",
            type: "int",
            description:
              "Unix timestamp (in seconds) when the group was created.",
            computedOptionalRequired: "computed",
          },
        ],
      },
    ],
  },
  {
    name: "group_users",
    description: "Lists the users assigned to a group.",
    api: {
      model: "GroupUserAssignment",
      readMethod: "ListGroupUsers",
      readRequestAttributes: ["group_id"],
      readModel: "User",
      readStrategy: "paginate",
      readCursorParam: "Next",
    },
    attributes: [
      {
        name: "group_id",
        type: "string",
        description: "The ID of the group to update.",
        computedOptionalRequired: "required",
      },
      {
        name: "users",
        type: "set_nested",
        description: "List of users.",
        computedOptionalRequired: "computed",
        attributes: [
          {
            name: "id",
            type: "string",
            description: "User ID.",
            computedOptionalRequired: "computed",
          },
          {
            name: "email",
            type: "string",
            description: "The email address of the user.",
            computedOptionalRequired: "computed",
          },
          {
            name: "name",
            type: "string",
            description: "The name of the user.",
            computedOptionalRequired: "computed",
          },
          {
            name: "role",
            type: "string",
            description: "Role `owner` or `reader`.",
            computedOptionalRequired: "computed",
          },
          {
            name: "added_at",
            type: "int",
            description:
              "The Unix timestamp (in seconds) of when the user was added.",
            computedOptionalRequired: "computed",
          },
        ],
      },
    ],
  },
  {
    name: "group_role_assignments",
    description:
      "Lists the organization roles assigned to a group within the organization.",
    api: {
      model: "AssignedRoleDetails",
      readStrategy: "paginate",
      readMethod: "ListGroupRoleAssignments",
      readRequestAttributes: ["group_id"],
      readCursorParam: "Next",
    },
    attributes: [
      {
        name: "group_id",
        type: "string",
        description:
          "The ID of the group whose organization role assignments you want to list.",
        computedOptionalRequired: "required",
      },
      {
        name: "roles",
        type: "set_nested",
        description: "List of organization roles",
        computedOptionalRequired: "computed",
        attributes: [
          {
            name: "id",
            type: "string",
            description: "Identifier for the role.",
            computedOptionalRequired: "computed",
          },
          {
            name: "name",
            type: "string",
            description: "Unique name for the role.",
            computedOptionalRequired: "computed",
          },
          {
            name: "description",
            type: "string",
            description: "Description of the role.",
            computedOptionalRequired: "computed",
          },
          {
            name: "permissions",
            type: "set",
            description: "Permissions granted by the role.",
            computedOptionalRequired: "computed",
            elementType: "string",
          },
          {
            name: "predefined_role",
            type: "bool",
            description:
              "Whether the role is predefined and managed by OpenAI.",
            computedOptionalRequired: "computed",
          },
          {
            name: "resource_type",
            type: "string",
            description:
              "Resource type the role is bound to (for example `api.organization` or `api.project`).",
            computedOptionalRequired: "computed",
          },
        ],
      },
    ],
  },
  {
    name: "invite",
    description: "Retrieves an invite.",
    api: {
      model: "Invite",
      readStrategy: "simple",
      readMethod: "RetrieveInvite",
    },
    attributes: [
      {
        name: "id",
        type: "string",
        description: "Invite ID.",
        computedOptionalRequired: "required",
      },
      {
        name: "email",
        type: "string",
        description:
          "The email address of the individual to whom the invite was sent.",
        computedOptionalRequired: "computed",
      },
      {
        name: "role",
        type: "string",
        description: "`owner` or `reader`.",
        computedOptionalRequired: "computed",
      },
      {
        name: "status",
        type: "string",
        description: "`accepted`, `expired`, or `pending`.",
        computedOptionalRequired: "computed",
      },
      {
        name: "invited_at",
        type: "int",
        description:
          "The Unix timestamp (in seconds) of when the invite was sent.",
        computedOptionalRequired: "computed",
      },
      {
        name: "expires_at",
        type: "int",
        description:
          "The Unix timestamp (in seconds) of when the invite expires.",
        computedOptionalRequired: "computed",
      },
      {
        name: "accepted_at",
        type: "int",
        description:
          "The Unix timestamp (in seconds) of when the invite was accepted.",
        computedOptionalRequired: "computed",
      },
    ],
  },
  {
    name: "invites",
    description: "Lists all of the invites in the organization.",
    api: {
      model: "Invite",
      readStrategy: "paginate",
      readMethod: "ListInvites",
    },
    attributes: [
      {
        name: "invites",
        type: "set_nested",
        description: "List of invites.",
        computedOptionalRequired: "computed",
        attributes: [
          {
            name: "id",
            type: "string",
            description: "Invite ID.",
            computedOptionalRequired: "computed",
          },
          {
            name: "email",
            type: "string",
            description:
              "The email address of the individual to whom the invite was sent.",
            computedOptionalRequired: "computed",
          },
          {
            name: "role",
            type: "string",
            description: "`owner` or `reader`.",
            computedOptionalRequired: "computed",
          },
          {
            name: "status",
            type: "string",
            description: "`accepted`, `expired`, or `pending`.",
            computedOptionalRequired: "computed",
          },
          {
            name: "invited_at",
            type: "int",
            description:
              "The Unix timestamp (in seconds) of when the invite was sent.",
            computedOptionalRequired: "computed",
          },
          {
            name: "expires_at",
            type: "int",
            description:
              "The Unix timestamp (in seconds) of when the invite expires.",
            computedOptionalRequired: "computed",
          },
          {
            name: "accepted_at",
            type: "int",
            description:
              "The Unix timestamp (in seconds) of when the invite was accepted.",
            computedOptionalRequired: "computed",
          },
        ],
      },
    ],
  },
  {
    name: "organization_roles",
    description: "Lists the roles configured for the organization.",
    api: {
      model: "Role",
      readStrategy: "paginate",
      readMethod: "ListRoles",
      readCursorParam: "Next",
    },
    attributes: [
      {
        name: "roles",
        type: "set_nested",
        description: "List of roles.",
        computedOptionalRequired: "computed",
        attributes: [
          {
            name: "id",
            type: "string",
            description: "Identifier for the role.",
            computedOptionalRequired: "computed",
          },
          {
            name: "name",
            type: "string",
            description: "Unique name for the role.",
            computedOptionalRequired: "computed",
          },
          {
            name: "description",
            type: "string",
            description: "Description of the role.",
            computedOptionalRequired: "computed",
          },
          {
            name: "permissions",
            type: "set",
            description: "Permissions granted by the role.",
            computedOptionalRequired: "computed",
            elementType: "string",
          },
          {
            name: "predefined_role",
            type: "bool",
            description:
              "Whether the role is predefined and managed by OpenAI.",
            computedOptionalRequired: "computed",
          },
          {
            name: "resource_type",
            type: "string",
            description:
              "Resource type the role is bound to (for example `api.organization` or `api.project`).",
            computedOptionalRequired: "computed",
          },
        ],
      },
    ],
  },
  {
    name: "project",
    description: "Retrieve a project by ID.",
    api: {
      model: "Project",
      readStrategy: "simple",
      readMethod: "RetrieveProject",
    },
    attributes: [
      {
        name: "id",
        type: "string",
        description: "Project ID.",
        computedOptionalRequired: "required",
      },
      {
        name: "name",
        type: "string",
        description: "The name of the project. This appears in reporting.",
        computedOptionalRequired: "computed",
      },
      {
        name: "status",
        type: "string",
        description: "Status `active` or `archived`.",
        computedOptionalRequired: "computed",
      },
      {
        name: "external_key_id",
        type: "string",
        description:
          "The ID of the customer-managed encryption key used for Enterprise Key Management (EKM). EKM is only available on certain accounts. Refer to the [EKM (External Keys) in the Management API Article](https://help.openai.com/en/articles/20000953-ekm-external-keys-in-the-management-api).",
        computedOptionalRequired: "computed",
        nullable: true,
      },
      {
        name: "created_at",
        type: "int",
        description:
          "The Unix timestamp (in seconds) of when the project was created.",
        computedOptionalRequired: "computed",
      },
      {
        name: "archived_at",
        type: "int",
        description:
          "The Unix timestamp (in seconds) of when the project was archived or `null`.",
        computedOptionalRequired: "computed",
        nullable: true,
      },
    ],
  },
  {
    name: "projects",
    description: "List all projects in an organization.",
    api: {
      model: "Project",
      readStrategy: "paginate",
      readMethod: "ListProjects",
      readInitLoop: `
        params.IncludeArchived = data.IncludeArchived.ValueBoolPointer()

        // Set the limit for the API request
        if data.Limit.IsNull() {
          params.Limit = ptr.Ptr(int64(100))
        } else {
          requestLimit := data.Limit.ValueInt64()
          if requestLimit > 100 {
            params.Limit = ptr.Ptr(int64(100))
          } else {
            params.Limit = ptr.Ptr(requestLimit)
          }
        }
      `,
      readPreIterate: `
        // Recalculate the limit for each request to ensure we don't exceed the desired limit
        if !data.Limit.IsNull() {
          remainingLimit := data.Limit.ValueInt64() - int64(len(modelInstances))
          if remainingLimit <= 0 {
            break
          }
          if remainingLimit > 100 {
            params.Limit = ptr.Ptr(int64(100))
          } else {
            params.Limit = ptr.Ptr(remainingLimit)
          }
        }
      `,
      readPostIterate: `
        // If limit is set and we have enough projects, break.
        if !data.Limit.IsNull() && len(modelInstances) >= int(data.Limit.ValueInt64()) {
          modelInstances = modelInstances[:data.Limit.ValueInt64()]
          break
        }
      `,
    },
    attributes: [
      {
        name: "include_archived",
        type: "bool",
        description: "Include archived projects. Default is `false`.",
        computedOptionalRequired: "optional",
      },
      {
        name: "limit",
        type: "int",
        description:
          "Limit the number of projects to return. Default is to return all projects.",
        computedOptionalRequired: "optional",
        validators: ["int64validator.AtLeast(1)"],
      },
      {
        name: "projects",
        type: "set_nested",
        description: "List of projects.",
        computedOptionalRequired: "computed",
        attributes: [
          {
            name: "id",
            type: "string",
            description: "Project ID.",
            computedOptionalRequired: "computed",
          },
          {
            name: "name",
            type: "string",
            description: "The name of the project. This appears in reporting.",
            computedOptionalRequired: "computed",
          },
          {
            name: "status",
            type: "string",
            description: "Status `active` or `archived`.",
            computedOptionalRequired: "computed",
          },
          {
            name: "external_key_id",
            type: "string",
            description:
              "The ID of the customer-managed encryption key used for Enterprise Key Management (EKM). EKM is only available on certain accounts. Refer to the [EKM (External Keys) in the Management API Article](https://help.openai.com/en/articles/20000953-ekm-external-keys-in-the-management-api).",
            computedOptionalRequired: "computed",
            nullable: true,
          },
          {
            name: "created_at",
            type: "int",
            description:
              "The Unix timestamp (in seconds) of when the project was created.",
            computedOptionalRequired: "computed",
          },
          {
            name: "archived_at",
            type: "int",
            description:
              "The Unix timestamp (in seconds) of when the project was archived or `null`.",
            computedOptionalRequired: "computed",
            nullable: true,
          },
        ],
      },
    ],
  },
  {
    name: "project_rate_limits",
    description: "Returns the rate limits per model for a project.",
    api: {
      model: "ProjectRateLimit",
      readStrategy: "paginate",
      readMethod: "ListProjectRateLimits",
      readRequestAttributes: ["project_id"],
    },
    attributes: [
      {
        name: "project_id",
        type: "string",
        description: "The ID of the project.",
        computedOptionalRequired: "required",
      },
      {
        name: "rate_limits",
        type: "set_nested",
        description: "List of rate limits.",
        computedOptionalRequired: "computed",
        attributes: [
          {
            name: "id",
            type: "string",
            description: "The rate limit identifier.",
            computedOptionalRequired: "computed",
          },
          {
            name: "model",
            type: "string",
            description: "The model this rate limit applies to.",
            computedOptionalRequired: "computed",
          },
          {
            name: "max_requests_per_1_minute",
            type: "int",
            description: "The maximum requests per minute.",
            computedOptionalRequired: "computed",
          },
          {
            name: "max_tokens_per_1_minute",
            type: "int",
            description: "The maximum tokens per minute.",
            computedOptionalRequired: "computed",
          },
          {
            name: "max_images_per_1_minute",
            type: "int",
            description:
              "The maximum images per minute. Only present for relevant models.",
            computedOptionalRequired: "computed",
          },
          {
            name: "max_audio_megabytes_per_1_minute",
            type: "int",
            description:
              "The maximum audio megabytes per minute. Only present for relevant models.",
            computedOptionalRequired: "computed",
          },
          {
            name: "max_requests_per_1_day",
            type: "int",
            description:
              "The maximum requests per day. Only present for relevant models.",
            computedOptionalRequired: "computed",
          },
          {
            name: "batch_1_day_max_input_tokens",
            type: "int",
            description:
              "The maximum batch input tokens per day. Only present for relevant models.",
            computedOptionalRequired: "computed",
          },
        ],
      },
    ],
  },
  {
    name: "project_roles",
    description: "Lists the roles configured for a project.",
    api: {
      model: "Role",
      readStrategy: "paginate",
      readMethod: "ListProjectRoles",
      readRequestAttributes: ["project_id"],
      readCursorParam: "Next",
    },
    attributes: [
      {
        name: "project_id",
        type: "string",
        description: "The ID of the project to inspect.",
        computedOptionalRequired: "required",
      },
      {
        name: "roles",
        type: "set_nested",
        description: "List of roles configured for a project.",
        computedOptionalRequired: "computed",
        attributes: [
          {
            name: "id",
            type: "string",
            description: "Identifier for the role.",
            computedOptionalRequired: "computed",
          },
          {
            name: "name",
            type: "string",
            description: "Unique name for the role.",
            computedOptionalRequired: "computed",
          },
          {
            name: "description",
            type: "string",
            description: "Description of the role.",
            computedOptionalRequired: "computed",
          },
          {
            name: "permissions",
            type: "set",
            description: "Permissions granted by the role.",
            computedOptionalRequired: "computed",
            elementType: "string",
          },
          {
            name: "predefined_role",
            type: "bool",
            description:
              "Whether the role is predefined and managed by OpenAI.",
            computedOptionalRequired: "computed",
          },
          {
            name: "resource_type",
            type: "string",
            description:
              "Resource type the role is bound to (for example `api.organization` or `api.project`).",
            computedOptionalRequired: "computed",
          },
        ],
      },
    ],
  },
  {
    name: "user",
    description: "Retrieves a user by their identifier.",
    api: {
      model: "User",
      readStrategy: "simple",
      readMethod: "RetrieveUser",
    },
    attributes: [
      {
        name: "id",
        type: "string",
        description: "User ID.",
        computedOptionalRequired: "required",
      },
      {
        name: "email",
        type: "string",
        description: "The email address of the user.",
        computedOptionalRequired: "computed",
      },
      {
        name: "name",
        type: "string",
        description: "The name of the user.",
        computedOptionalRequired: "computed",
      },
      {
        name: "role",
        type: "string",
        description: "Role `owner` or `reader`.",
        computedOptionalRequired: "computed",
      },
      {
        name: "added_at",
        type: "int",
        description:
          "The Unix timestamp (in seconds) of when the user was added.",
        computedOptionalRequired: "computed",
      },
    ],
  },
  {
    name: "users",
    description: "Lists all of the users in the organization.",
    api: {
      model: "User",
      readStrategy: "paginate",
      readMethod: "ListUsers",
    },
    attributes: [
      {
        name: "users",
        type: "set_nested",
        description: "List of users.",
        computedOptionalRequired: "computed",
        attributes: [
          {
            name: "id",
            type: "string",
            description: "User ID.",
            computedOptionalRequired: "computed",
          },
          {
            name: "email",
            type: "string",
            description: "The email address of the user.",
            computedOptionalRequired: "computed",
          },
          {
            name: "name",
            type: "string",
            description: "The name of the user.",
            computedOptionalRequired: "computed",
          },
          {
            name: "role",
            type: "string",
            description: "Role `owner` or `reader`.",
            computedOptionalRequired: "computed",
          },
          {
            name: "added_at",
            type: "int",
            description:
              "The Unix timestamp (in seconds) of when the user was added.",
            computedOptionalRequired: "computed",
          },
        ],
      },
    ],
  },
  {
    name: "user_role_assignments",
    description:
      "Lists the organization roles assigned to a user within the organization.",
    api: {
      model: "AssignedRoleDetails",
      readStrategy: "paginate",
      readMethod: "ListUserRoleAssignments",
      readRequestAttributes: ["user_id"],
      readCursorParam: "Next",
    },
    attributes: [
      {
        name: "user_id",
        type: "string",
        description: "The ID of the user to inspect.",
        computedOptionalRequired: "required",
      },
      {
        name: "roles",
        type: "set_nested",
        description: "List of organization roles",
        computedOptionalRequired: "computed",
        attributes: [
          {
            name: "id",
            type: "string",
            description: "Identifier for the role.",
            computedOptionalRequired: "computed",
          },
          {
            name: "name",
            type: "string",
            description: "Unique name for the role.",
            computedOptionalRequired: "computed",
          },
          {
            name: "description",
            type: "string",
            description: "Description of the role.",
            computedOptionalRequired: "computed",
          },
          {
            name: "permissions",
            type: "set",
            description: "Permissions granted by the role.",
            computedOptionalRequired: "computed",
            elementType: "string",
          },
          {
            name: "predefined_role",
            type: "bool",
            description:
              "Whether the role is predefined and managed by OpenAI.",
            computedOptionalRequired: "computed",
          },
          {
            name: "resource_type",
            type: "string",
            description:
              "Resource type the role is bound to (for example `api.organization` or `api.project`).",
            computedOptionalRequired: "computed",
          },
        ],
      },
    ],
  },
];

export const RESOURCES: Array<Resource> = [
  {
    name: "admin_api_key",
    description: "Manages an organization admin API key.",
    api: {
      createMethod: "AdminApiKeysCreate",
      readMethod: "AdminApiKeysGet",
      readRequestAttributes: ["id"],
      deleteMethod: "AdminApiKeysDelete",
      deleteRequestAttributes: ["id"],
    },
    attributes: [
      {
        name: "name",
        type: "string",
        description: "The name of the organization admin API key.",
        computedOptionalRequired: "required",
        planModifiers: ["stringplanmodifier.RequiresReplace()"],
      },
      {
        name: "id",
        type: "string",
        description: "The ID of the organization admin API key.",
        computedOptionalRequired: "computed",
        planModifiers: ["stringplanmodifier.UseStateForUnknown()"],
      },
      {
        name: "created_at",
        type: "int",
        description:
          "The Unix timestamp (in seconds) of when the organization admin API key was created.",
        computedOptionalRequired: "computed",
        planModifiers: ["int64planmodifier.UseStateForUnknown()"],
      },
      {
        name: "api_key",
        type: "string",
        description:
          "The organization admin API key that can be used to authenticate with the API.",
        computedOptionalRequired: "computed",
        sensitive: true,
        planModifiers: ["stringplanmodifier.UseStateForUnknown()"],
      },
    ],
  },
  {
    name: "invite",
    description:
      "Invite and manage invitations for an organization. Invited users are automatically added to the Default project.",
    api: {
      createMethod: "InviteUser",
      readMethod: "RetrieveInvite",
      readRequestAttributes: ["id"],
      deleteMethod: "DeleteInvite",
      deleteRequestAttributes: ["id"],
    },
    importStateAttributes: ["id"],
    attributes: [
      {
        name: "id",
        type: "string",
        description: "Invite ID.",
        computedOptionalRequired: "computed",
        planModifiers: ["stringplanmodifier.UseStateForUnknown()"],
      },
      {
        name: "email",
        type: "string",
        description:
          "The email address of the individual to whom the invite was sent.",
        computedOptionalRequired: "required",
        planModifiers: ["stringplanmodifier.RequiresReplace()"],
      },
      {
        name: "role",
        type: "string",
        description: "`owner` or `reader`.",
        computedOptionalRequired: "required",
        planModifiers: ["stringplanmodifier.RequiresReplace()"],
        validators: ['stringvalidator.OneOf("owner", "reader")'],
      },
      {
        name: "status",
        type: "string",
        description: "`accepted`, `expired`, or `pending`.",
        computedOptionalRequired: "computed",
      },
      {
        name: "invited_at",
        type: "int",
        description:
          "The Unix timestamp (in seconds) of when the invite was sent.",
        computedOptionalRequired: "computed",
      },
      {
        name: "expires_at",
        type: "int",
        description:
          "The Unix timestamp (in seconds) of when the invite expires.",
        computedOptionalRequired: "computed",
      },
      {
        name: "accepted_at",
        type: "int",
        description:
          "The Unix timestamp (in seconds) of when the invite was accepted.",
        computedOptionalRequired: "computed",
      },
    ],
  },
  {
    name: "organization_role",
    description: "Creates a custom role for the organization.",
    api: {
      model: "Role",
      createMethod: "CreateRole",
      readMethod: "ListRoles",
      readStrategy: "paginate",
      readCursorParam: "Next",
      updateMethod: "UpdateRole",
      updateRequestAttributes: ["id"],
      deleteMethod: "DeleteRole",
      deleteRequestAttributes: ["id"],
    },
    importStateAttributes: ["id"],
    attributes: [
      {
        name: "id",
        type: "string",
        description: "Identifier for the role.",
        computedOptionalRequired: "computed",
        planModifiers: ["stringplanmodifier.UseStateForUnknown()"],
      },
      {
        name: "name",
        type: "string",
        description: "Unique name for the role.",
        computedOptionalRequired: "required",
      },
      {
        name: "description",
        type: "string",
        description: "Description of the role.",
        computedOptionalRequired: "optional",
      },
      {
        name: "permissions",
        type: "set",
        description: "Permissions to grant to the role.",
        computedOptionalRequired: "required",
        elementType: "string",
      },
    ],
  },
  {
    name: "project",
    description: "Project resource.",
    api: {
      createMethod: "CreateProject",
      readMethod: "RetrieveProject",
      readRequestAttributes: ["id"],
      updateMethod: "ModifyProject",
      updateRequestAttributes: ["id"],
      deleteMethod: "ArchiveProject",
      deleteRequestAttributes: ["id"],
    },
    importStateAttributes: ["id"],
    attributes: [
      {
        name: "id",
        type: "string",
        description: "The ID of the project.",
        computedOptionalRequired: "computed",
        planModifiers: ["stringplanmodifier.UseStateForUnknown()"],
      },
      {
        name: "name",
        type: "string",
        description:
          "The friendly name of the project, this name appears in reports.",
        computedOptionalRequired: "required",
      },
      {
        name: "status",
        type: "string",
        description: "Status `active` or `archived`.",
        computedOptionalRequired: "computed",
      },
      {
        name: "external_key_id",
        type: "string",
        description:
          "The ID of the customer-managed encryption key to use for Enterprise Key Management (EKM). EKM is only available on certain accounts. Refer to the [EKM (External Keys) in the Management API Article](https://help.openai.com/en/articles/20000953-ekm-external-keys-in-the-management-api).",
        computedOptionalRequired: "optional",
        planModifiers: ["stringplanmodifier.UseStateForUnknown()"],
      },
      {
        name: "created_at",
        type: "int",
        description:
          "The Unix timestamp (in seconds) of when the project was created.",
        computedOptionalRequired: "computed",
      },
      {
        name: "archived_at",
        type: "int",
        description:
          "The Unix timestamp (in seconds) of when the project was archived or `null`.",
        computedOptionalRequired: "computed",
      },
    ],
  },
  {
    name: "project_group_role_assignment",
    description: "Assigns a project role to a group within a project.",
    api: {
      model: "AssignedRoleDetails",
      createMethod: "AssignProjectGroupRole",
      createRequestAttributes: ["project_id", "group_id"],
      readMethod: "ListProjectGroupRoleAssignments",
      readRequestAttributes: ["project_id", "group_id"],
      readStrategy: "paginate",
      readCursorParam: "Next",
      deleteMethod: "UnassignProjectGroupRole",
      deleteRequestAttributes: ["project_id", "group_id", "role_id"],
    },
    attributes: [
      {
        name: "project_id",
        type: "string",
        description: "The ID of the project to update.",
        computedOptionalRequired: "required",
        planModifiers: ["stringplanmodifier.RequiresReplace()"],
      },
      {
        name: "group_id",
        type: "string",
        description: "Identifier of the group to add to the project.",
        computedOptionalRequired: "required",
        planModifiers: ["stringplanmodifier.RequiresReplace()"],
      },
      {
        name: "role_id",
        type: "string",
        description: "Identifier of the project role to grant to the group.",
        computedOptionalRequired: "required",
        planModifiers: ["stringplanmodifier.RequiresReplace()"],
      },
    ],
  },
  {
    name: "project_rate_limit",
    description:
      "Manage rate limits per model for projects. Rate limits may be configured to be equal to or lower than the organization's rate limits.",
    api: {
      model: "ProjectRateLimit",
      createMethod: "UpdateProjectRateLimits",
      createRequestAttributes: ["project_id", "rate_limit_id"],
      readMethod: "ListProjectRateLimits",
      readRequestAttributes: ["project_id"],
      readStrategy: "paginate",
      updateMethod: "UpdateProjectRateLimits",
      updateRequestAttributes: ["project_id", "rate_limit_id"],
    },
    attributes: [
      {
        name: "project_id",
        type: "string",
        description: "The ID of the project.",
        computedOptionalRequired: "required",
        planModifiers: ["stringplanmodifier.RequiresReplace()"],
      },
      {
        name: "rate_limit_id",
        type: "string",
        description:
          "The ID of the rate limit. This is typically in the format `rl-<model>`.",
        computedOptionalRequired: "required",
        planModifiers: ["stringplanmodifier.RequiresReplace()"],
      },
      {
        name: "max_requests_per_1_minute",
        type: "int",
        description: "The maximum requests per minute.",
        computedOptionalRequired: "computed_optional",
      },
      {
        name: "max_tokens_per_1_minute",
        type: "int",
        description: "The maximum tokens per minute.",
        computedOptionalRequired: "computed_optional",
      },
      {
        name: "max_images_per_1_minute",
        type: "int",
        description:
          "The maximum images per minute. Only relevant for certain models.",
        computedOptionalRequired: "computed_optional",
      },
      {
        name: "max_audio_megabytes_per_1_minute",
        type: "int",
        description:
          "The maximum audio megabytes per minute. Only relevant for certain models.",
        computedOptionalRequired: "computed_optional",
      },
      {
        name: "max_requests_per_1_day",
        type: "int",
        description:
          "The maximum requests per day. Only relevant for certain models.",
        computedOptionalRequired: "computed_optional",
      },
      {
        name: "batch_1_day_max_input_tokens",
        type: "int",
        description:
          "The maximum batch input tokens per day. Only relevant for certain models.",
        computedOptionalRequired: "computed_optional",
      },
    ],
  },
  {
    name: "project_role",
    description: "Creates a custom role for a project.",
    api: {
      model: "Role",
      createMethod: "CreateProjectRole",
      createRequestAttributes: ["project_id"],
      readMethod: "ListProjectRoles",
      readRequestAttributes: ["project_id"],
      readStrategy: "paginate",
      readCursorParam: "Next",
      updateMethod: "UpdateProjectRole",
      updateRequestAttributes: ["project_id", "id"],
      deleteMethod: "DeleteProjectRole",
      deleteRequestAttributes: ["project_id", "id"],
    },
    importStateAttributes: ["project_id", "id"],
    attributes: [
      {
        name: "id",
        type: "string",
        description: "Identifier for the role.",
        computedOptionalRequired: "computed",
        planModifiers: ["stringplanmodifier.UseStateForUnknown()"],
      },
      {
        name: "project_id",
        type: "string",
        description: "The ID of the project to create the role for.",
        computedOptionalRequired: "required",
      },
      {
        name: "name",
        type: "string",
        description: "Unique name for the role.",
        computedOptionalRequired: "required",
      },
      {
        name: "description",
        type: "string",
        description: "Description of the role.",
        computedOptionalRequired: "optional",
      },
      {
        name: "permissions",
        type: "set",
        description: "Permissions to grant to the role.",
        computedOptionalRequired: "required",
        elementType: "string",
      },
    ],
  },
  {
    name: "project_service_account",
    description:
      "Manage service accounts within a project. A service account is a bot user that is not associated with a user. If a user leaves an organization, their keys and membership in projects will no longer work. Service accounts do not have this limitation. However, service accounts can also be deleted from a project.",
    api: {
      createMethod: "CreateProjectServiceAccount",
      createRequestAttributes: ["project_id"],
      readMethod: "RetrieveProjectServiceAccount",
      readRequestAttributes: ["project_id", "id"],
      deleteMethod: "DeleteProjectServiceAccount",
      deleteRequestAttributes: ["project_id", "id"],
    },
    attributes: [
      {
        name: "project_id",
        type: "string",
        description: "The ID of the project.",
        computedOptionalRequired: "required",
        planModifiers: ["stringplanmodifier.RequiresReplace()"],
      },
      {
        name: "name",
        type: "string",
        description: "The name of the service account being created.",
        computedOptionalRequired: "required",
        planModifiers: ["stringplanmodifier.RequiresReplace()"],
      },
      {
        name: "id",
        type: "string",
        description: "The ID of the service account.",
        computedOptionalRequired: "computed",
        planModifiers: ["stringplanmodifier.UseStateForUnknown()"],
      },
      {
        name: "role",
        type: "string",
        description:
          "The role of the service account. Can be `owner` or `member`.",
        computedOptionalRequired: "computed",
      },
      {
        name: "created_at",
        type: "int",
        description:
          "The Unix timestamp (in seconds) of when the service account was created.",
        computedOptionalRequired: "computed",
        planModifiers: ["int64planmodifier.UseStateForUnknown()"],
      },
      {
        name: "api_key_id",
        type: "string",
        description:
          "Internal ID of the API key. This is a reference to the API key and not the actual key.",
        computedOptionalRequired: "computed",
        planModifiers: ["stringplanmodifier.UseStateForUnknown()"],
      },
      {
        name: "api_key",
        type: "string",
        description:
          "The API key that can be used to authenticate with the API.",
        computedOptionalRequired: "computed",
        sensitive: true,
        planModifiers: ["stringplanmodifier.UseStateForUnknown()"],
      },
    ],
  },
  {
    name: "project_user",
    description:
      "Adds a user to the project. Users must already be members of the organization to be added to a project.",
    api: {
      createMethod: "CreateProjectUser",
      createRequestAttributes: ["project_id"],
      readMethod: "RetrieveProjectUser",
      readRequestAttributes: ["project_id", "user_id"],
      updateMethod: "ModifyProjectUser",
      updateRequestAttributes: ["project_id", "user_id"],
      deleteMethod: "DeleteProjectUser",
      deleteRequestAttributes: ["project_id", "user_id"],
    },
    importStateAttributes: ["project_id", "user_id"],
    attributes: [
      {
        name: "project_id",
        type: "string",
        description: "The ID of the project.",
        computedOptionalRequired: "required",
        planModifiers: ["stringplanmodifier.RequiresReplace()"],
      },
      {
        name: "user_id",
        type: "string",
        description: "The ID of the user.",
        computedOptionalRequired: "required",
        planModifiers: ["stringplanmodifier.RequiresReplace()"],
      },
      {
        name: "role",
        type: "string",
        description: "`owner` or `member`.",
        computedOptionalRequired: "required",
        validators: ['stringvalidator.OneOf("owner", "member")'],
      },
    ],
  },
  {
    name: "user_role",
    description:
      "Modifies a user's role in the organization.\n\n**NOTE:** The new `openai_user_role_assignment` resource supports predefined roles like `owner` and `reader` as well as custom roles. This resource may be removed in a future release.",
    api: {
      createMethod: "ModifyUser",
      createRequestAttributes: ["user_id"],
      readMethod: "RetrieveUser",
      readRequestAttributes: ["user_id"],
      updateMethod: "ModifyUser",
      updateRequestAttributes: ["user_id"],
    },
    importStateAttributes: ["user_id"],
    attributes: [
      {
        name: "user_id",
        type: "string",
        description: "The ID of the user.",
        computedOptionalRequired: "required",
      },
      {
        name: "role",
        type: "string",
        description: "`owner` or `reader`.",
        computedOptionalRequired: "required",
        validators: ['stringvalidator.OneOf("owner", "reader")'],
      },
    ],
  },
  {
    name: "user_role_assignment",
    description:
      "Assigns an organization role to a user within the organization.\n\n**NOTE:** Predefined organization roles like `owner` and `reader` are in the format of `role-api-organization-<role_name>__api-organization__<org_id>`. You can use the `provider::openai::predefined_role_id(role, organization_id)` function to generate the role ID.",
    api: {
      model: "AssignedRoleDetails",
      createMethod: "AssignUserRole",
      createRequestAttributes: ["user_id"],
      readMethod: "ListUserRoleAssignments",
      readRequestAttributes: ["user_id"],
      readStrategy: "paginate",
      readCursorParam: "Next",
      updateMethod: "AssignUserRole",
      updateRequestAttributes: ["user_id"],
      deleteMethod: "UnassignUserRole",
      deleteRequestAttributes: ["user_id", "role_id"],
    },
    importStateAttributes: ["user_id", "role_id"],
    attributes: [
      {
        name: "user_id",
        type: "string",
        description:
          "The ID of the user that should receive the organization role.",
        computedOptionalRequired: "required",
      },
      {
        name: "role_id",
        type: "string",
        description: "Identifier of the role to assign.",
        computedOptionalRequired: "required",
      },
    ],
  },
  {
    name: "group",
    description: "Creates a new group in the organization.",
    api: {
      model: "GroupResponse",
      createMethod: "CreateGroup",
      readMethod: "ListGroups",
      readStrategy: "paginate",
      readCursorParam: "Next",
      updateMethod: "UpdateGroup",
      updateRequestAttributes: ["id"],
      deleteMethod: "DeleteGroup",
      deleteRequestAttributes: ["id"],
    },
    importStateAttributes: ["id"],
    attributes: [
      {
        name: "name",
        type: "string",
        description: "Human readable name for the group.",
        computedOptionalRequired: "required",
      },
      {
        name: "id",
        type: "string",
        description: "Identifier for the group.",
        computedOptionalRequired: "computed",
        planModifiers: ["stringplanmodifier.UseStateForUnknown()"],
      },
      {
        name: "created_at",
        type: "int",
        description: "Unix timestamp (in seconds) when the group was created.",
        computedOptionalRequired: "computed",
      },
    ],
  },
  {
    name: "group_user",
    description: "Adds a user to a group.",
    api: {
      model: "GroupUserAssignment",
      createMethod: "AddGroupUser",
      createRequestAttributes: ["group_id"],
      readMethod: "ListGroupUsers",
      readRequestAttributes: ["group_id"],
      readModel: "User",
      readStrategy: "paginate",
      readCursorParam: "Next",
      deleteMethod: "RemoveGroupUser",
      deleteRequestAttributes: ["group_id", "user_id"],
    },
    importStateAttributes: ["group_id", "user_id"],
    attributes: [
      {
        name: "group_id",
        type: "string",
        description: "The ID of the group to update.",
        computedOptionalRequired: "required",
        planModifiers: ["stringplanmodifier.RequiresReplace()"],
      },
      {
        name: "user_id",
        type: "string",
        description: "Identifier of the user to add to the group.",
        computedOptionalRequired: "required",
        planModifiers: ["stringplanmodifier.RequiresReplace()"],
      },
    ],
  },
  {
    name: "group_role_assignment",
    description:
      "Assigns an organization role to a group within the organization.",
    api: {
      model: "GroupRoleAssignment",
      createMethod: "AssignGroupRole",
      createRequestAttributes: ["group_id"],
      readMethod: "ListGroupRoleAssignments",
      readRequestAttributes: ["group_id"],
      readStrategy: "paginate",
      readModel: "AssignedRoleDetails",
      readCursorParam: "Next",
      deleteMethod: "UnassignGroupRole",
      deleteRequestAttributes: ["group_id", "role_id"],
    },
    importStateAttributes: ["group_id", "role_id"],
    attributes: [
      {
        name: "group_id",
        type: "string",
        description:
          "The ID of the group that should receive the organization role.",
        computedOptionalRequired: "required",
        planModifiers: ["stringplanmodifier.RequiresReplace()"],
      },
      {
        name: "role_id",
        type: "string",
        description: "Identifier of the role to assign.",
        computedOptionalRequired: "required",
        planModifiers: ["stringplanmodifier.RequiresReplace()"],
      },
    ],
  },
];
