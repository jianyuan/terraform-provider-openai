import type { DataSource, Resource } from "./schema";

export const DATASOURCES: Array<DataSource> = [
  {
    name: "groups",
    description: "Lists all groups in the organization.",
    api: {
      readStrategy: "paginate",
      readModel: "Group",
      readMethod: "Admin.Organization.Groups.ListAutoPaging",
      readRequestParamsStruct: "AdminOrganizationGroupListParams",
    },
    filler: {
      model: "[]openai.Group",
    },
    attributes: [
      {
        name: "groups",
        type: "set_nested",
        description: "List of groups.",
        computedOptionalRequired: "computed",
        filler: {
          model: "openai.Group",
          sourceAttribute: [],
        },
        attributes: [
          {
            name: "id",
            type: "string",
            description: "Identifier for the group.",
            computedOptionalRequired: "computed",
            filler: {
              sourceAttribute: ["ID"],
            },
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
            type: "int64",
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
      readStrategy: "paginate",
      readModel: "OrganizationGroupUser",
      readMethod: "Admin.Organization.Groups.Users.ListAutoPaging",
      readRequestParamsStruct: "AdminOrganizationGroupUserListParams",
      readRequestAttributes: ["group_id"],
    },
    filler: {
      model: "[]openai.OrganizationGroupUser",
    },
    attributes: [
      {
        name: "group_id",
        type: "string",
        description: "The ID of the group to update.",
        computedOptionalRequired: "required",
        filler: { skip: true },
      },
      {
        name: "users",
        type: "set_nested",
        description: "List of users.",
        computedOptionalRequired: "computed",
        filler: {
          model: "openai.OrganizationGroupUser",
          sourceAttribute: [],
        },
        attributes: [
          {
            name: "id",
            type: "string",
            description: "User ID.",
            computedOptionalRequired: "computed",
            filler: {
              sourceAttribute: ["ID"],
            },
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
        ],
      },
    ],
  },
  {
    name: "group_role_assignments",
    description:
      "Lists the organization roles assigned to a group within the organization.",
    api: {
      readStrategy: "paginate",
      readModel: "AdminOrganizationGroupRoleListResponse",
      readMethod: "Admin.Organization.Groups.Roles.ListAutoPaging",
      readRequestParamsStruct: "AdminOrganizationGroupRoleListParams",
      readRequestAttributes: ["group_id"],
    },
    filler: {
      model: "[]openai.AdminOrganizationGroupRoleListResponse",
    },
    attributes: [
      {
        name: "group_id",
        type: "string",
        description:
          "The ID of the group whose organization role assignments you want to list.",
        computedOptionalRequired: "required",
        filler: { skip: true },
      },
      {
        name: "roles",
        type: "set_nested",
        description: "List of organization roles",
        computedOptionalRequired: "computed",
        filler: {
          model: "openai.AdminOrganizationGroupRoleListResponse",
          sourceAttribute: [],
        },
        attributes: [
          {
            name: "id",
            type: "string",
            description: "Identifier for the role.",
            computedOptionalRequired: "computed",
            filler: {
              sourceAttribute: ["ID"],
            },
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
      readStrategy: "simple",
      readMethod: "Admin.Organization.Invites.Get",
    },
    filler: {
      model: "openai.Invite",
    },
    attributes: [
      {
        name: "id",
        type: "string",
        description: "Invite ID.",
        computedOptionalRequired: "required",
        filler: {
          sourceAttribute: ["ID"],
        },
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
        name: "created_at",
        type: "int64",
        description:
          "The Unix timestamp (in seconds) of when the invite was sent.",
        computedOptionalRequired: "computed",
      },
      {
        name: "expires_at",
        type: "int64",
        description:
          "The Unix timestamp (in seconds) of when the invite expires.",
        computedOptionalRequired: "computed",
        nullable: true,
      },
      {
        name: "accepted_at",
        type: "int64",
        description:
          "The Unix timestamp (in seconds) of when the invite was accepted.",
        computedOptionalRequired: "computed",
        nullable: true,
      },
    ],
  },
  {
    name: "invites",
    description: "Lists all of the invites in the organization.",
    api: {
      readStrategy: "paginate",
      readModel: "Invite",
      readMethod: "Admin.Organization.Invites.ListAutoPaging",
      readRequestParamsStruct: "AdminOrganizationInviteListParams",
    },
    filler: {
      model: "[]openai.Invite",
    },
    attributes: [
      {
        name: "invites",
        type: "set_nested",
        description: "List of invites.",
        computedOptionalRequired: "computed",
        filler: {
          model: "openai.Invite",
          sourceAttribute: [],
        },
        attributes: [
          {
            name: "id",
            type: "string",
            description: "Invite ID.",
            computedOptionalRequired: "computed",
            filler: {
              sourceAttribute: ["ID"],
            },
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
            name: "created_at",
            type: "int64",
            description:
              "The Unix timestamp (in seconds) of when the invite was sent.",
            computedOptionalRequired: "computed",
          },
          {
            name: "expires_at",
            type: "int64",
            description:
              "The Unix timestamp (in seconds) of when the invite expires.",
            computedOptionalRequired: "computed",
            nullable: true,
          },
          {
            name: "accepted_at",
            type: "int64",
            description:
              "The Unix timestamp (in seconds) of when the invite was accepted.",
            computedOptionalRequired: "computed",
            nullable: true,
          },
        ],
      },
    ],
  },
  {
    name: "organization_roles",
    description: "Lists the roles configured for the organization.",
    api: {
      readStrategy: "paginate",
      readModel: "Role",
      readMethod: "Admin.Organization.Roles.ListAutoPaging",
      readRequestParamsStruct: "AdminOrganizationRoleListParams",
    },
    filler: {
      model: "[]openai.Role",
    },
    attributes: [
      {
        name: "roles",
        type: "set_nested",
        description: "List of roles.",
        computedOptionalRequired: "computed",
        filler: {
          model: "openai.Role",
          sourceAttribute: [],
        },
        attributes: [
          {
            name: "id",
            type: "string",
            description: "Identifier for the role.",
            computedOptionalRequired: "computed",
            filler: {
              sourceAttribute: ["ID"],
            },
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
      readStrategy: "simple",
      readMethod: "Admin.Organization.Projects.Get",
    },
    filler: {
      model: "openai.Project",
    },
    attributes: [
      {
        name: "id",
        type: "string",
        description: "Project ID.",
        computedOptionalRequired: "required",
        filler: {
          sourceAttribute: ["ID"],
        },
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
        filler: {
          sourceAttribute: ["ExternalKeyID"],
        },
      },
      {
        name: "created_at",
        type: "int64",
        description:
          "The Unix timestamp (in seconds) of when the project was created.",
        computedOptionalRequired: "computed",
      },
      {
        name: "archived_at",
        type: "int64",
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
      readStrategy: "paginate",
      readModel: "Project",
      readMethod: "Admin.Organization.Projects.ListAutoPaging",
      readRequestParamsStruct: "AdminOrganizationProjectListParams",
      readInitLoop: `
        if data.IncludeArchived.IsKnown() {
          params.IncludeArchived = openai.Bool(data.IncludeArchived.ValueBool())
        }

        // Set the limit for the API request
        if data.Limit.IsKnown() {
          requestLimit := data.Limit.ValueInt64()
          if requestLimit > 100 {
            params.Limit = openai.Int(100)
          } else {
            params.Limit = openai.Int(requestLimit)
          }
        } else {
          params.Limit = openai.Int(100)
        }
      `,
      readPostIterate: `
        // If limit is set and we have enough projects, break.
        if data.Limit.IsKnown() && len(modelInstances) >= int(data.Limit.ValueInt64()) {
          modelInstances = modelInstances[:data.Limit.ValueInt64()]
          break
        }
      `,
    },
    filler: {
      model: "[]openai.Project",
    },
    attributes: [
      {
        name: "include_archived",
        type: "bool",
        description: "Include archived projects. Default is `false`.",
        computedOptionalRequired: "optional",
        filler: { skip: true },
      },
      {
        name: "limit",
        type: "int64",
        description:
          "Limit the number of projects to return. Default is to return all projects.",
        computedOptionalRequired: "optional",
        validators: ["int64validator.AtLeast(1)"],
        filler: { skip: true },
      },
      {
        name: "projects",
        type: "set_nested",
        description: "List of projects.",
        computedOptionalRequired: "computed",
        filler: {
          model: "openai.Project",
          sourceAttribute: [],
        },
        attributes: [
          {
            name: "id",
            type: "string",
            description: "Project ID.",
            computedOptionalRequired: "computed",
            filler: {
              sourceAttribute: ["ID"],
            },
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
            filler: {
              sourceAttribute: ["ExternalKeyID"],
            },
          },
          {
            name: "created_at",
            type: "int64",
            description:
              "The Unix timestamp (in seconds) of when the project was created.",
            computedOptionalRequired: "computed",
          },
          {
            name: "archived_at",
            type: "int64",
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
      readStrategy: "paginate",
      readModel: "ProjectRateLimit",
      readMethod:
        "Admin.Organization.Projects.RateLimits.ListRateLimitsAutoPaging",
      readRequestAttributes: ["project_id"],
      readRequestParamsStruct:
        "AdminOrganizationProjectRateLimitListRateLimitsParams",
    },
    filler: {
      model: "[]openai.ProjectRateLimit",
    },
    attributes: [
      {
        name: "project_id",
        type: "string",
        description: "The ID of the project.",
        computedOptionalRequired: "required",
        filler: { skip: true },
      },
      {
        name: "rate_limits",
        type: "set_nested",
        description: "List of rate limits.",
        computedOptionalRequired: "computed",
        filler: {
          model: "openai.ProjectRateLimit",
          sourceAttribute: [],
        },
        attributes: [
          {
            name: "id",
            type: "string",
            description: "The rate limit identifier.",
            computedOptionalRequired: "computed",
            filler: {
              sourceAttribute: ["ID"],
            },
          },
          {
            name: "model",
            type: "string",
            description: "The model this rate limit applies to.",
            computedOptionalRequired: "computed",
          },
          {
            name: "max_requests_per_1_minute",
            type: "int64",
            description: "The maximum requests per minute.",
            computedOptionalRequired: "computed",
          },
          {
            name: "max_tokens_per_1_minute",
            type: "int64",
            description: "The maximum tokens per minute.",
            computedOptionalRequired: "computed",
          },
          {
            name: "max_images_per_1_minute",
            type: "int64",
            description:
              "The maximum images per minute. Only present for relevant models.",
            computedOptionalRequired: "computed",
            nullable: true,
          },
          {
            name: "max_audio_megabytes_per_1_minute",
            type: "int64",
            description:
              "The maximum audio megabytes per minute. Only present for relevant models.",
            computedOptionalRequired: "computed",
            nullable: true,
          },
          {
            name: "max_requests_per_1_day",
            type: "int64",
            description:
              "The maximum requests per day. Only present for relevant models.",
            computedOptionalRequired: "computed",
            nullable: true,
          },
          {
            name: "batch_1_day_max_input_tokens",
            type: "int64",
            description:
              "The maximum batch input tokens per day. Only present for relevant models.",
            computedOptionalRequired: "computed",
            nullable: true,
          },
        ],
      },
    ],
  },
  {
    name: "project_roles",
    description: "Lists the roles configured for a project.",
    api: {
      readStrategy: "paginate",
      readModel: "Role",
      readMethod: "Admin.Organization.Projects.Roles.ListAutoPaging",
      readRequestAttributes: ["project_id"],
      readRequestParamsStruct: "AdminOrganizationProjectRoleListParams",
    },
    filler: {
      model: "[]openai.Role",
    },
    attributes: [
      {
        name: "project_id",
        type: "string",
        description: "The ID of the project to inspect.",
        computedOptionalRequired: "required",
        filler: { skip: true },
      },
      {
        name: "roles",
        type: "set_nested",
        description: "List of roles configured for a project.",
        computedOptionalRequired: "computed",
        filler: {
          model: "openai.Role",
          sourceAttribute: [],
        },
        attributes: [
          {
            name: "id",
            type: "string",
            description: "Identifier for the role.",
            computedOptionalRequired: "computed",
            filler: {
              sourceAttribute: ["ID"],
            },
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
    name: "project_group_role_assignments",
    description:
      "Lists the project roles assigned to a group within a project.",
    api: {
      readStrategy: "paginate",
      readModel: "AdminOrganizationProjectGroupRoleListResponse",
      readMethod: "Admin.Organization.Projects.Groups.Roles.ListAutoPaging",
      readRequestParamsStruct: "AdminOrganizationProjectGroupRoleListParams",
      readRequestAttributes: ["project_id", "group_id"],
    },
    filler: {
      model: "[]openai.AdminOrganizationProjectGroupRoleListResponse",
    },
    attributes: [
      {
        name: "project_id",
        type: "string",
        description: "The ID of the project to inspect.",
        computedOptionalRequired: "required",
        filler: { skip: true },
      },
      {
        name: "group_id",
        type: "string",
        description: "The ID of the group to inspect.",
        computedOptionalRequired: "required",
        filler: { skip: true },
      },
      {
        name: "roles",
        type: "set_nested",
        description:
          "List of project roles assigned to the group within the project.",
        computedOptionalRequired: "computed",
        filler: {
          model: "openai.AdminOrganizationProjectGroupRoleListResponse",
          sourceAttribute: [],
        },
        attributes: [
          {
            name: "id",
            type: "string",
            description: "Identifier for the role.",
            computedOptionalRequired: "computed",
            filler: {
              sourceAttribute: ["ID"],
            },
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
    name: "project_user_role_assignments",
    description: "Lists the project roles assigned to a user within a project.",
    api: {
      readStrategy: "paginate",
      readModel: "AdminOrganizationProjectUserRoleListResponse",
      readMethod: "Admin.Organization.Projects.Users.Roles.ListAutoPaging",
      readRequestParamsStruct: "AdminOrganizationProjectUserRoleListParams",
      readRequestAttributes: ["project_id", "user_id"],
    },
    filler: {
      model: "[]openai.AdminOrganizationProjectUserRoleListResponse",
    },
    attributes: [
      {
        name: "project_id",
        type: "string",
        description: "The ID of the project to inspect.",
        computedOptionalRequired: "required",
        filler: { skip: true },
      },
      {
        name: "user_id",
        type: "string",
        description: "The ID of the user to inspect.",
        computedOptionalRequired: "required",
        filler: { skip: true },
      },
      {
        name: "roles",
        type: "set_nested",
        description:
          "List of project roles assigned to the group within the project.",
        computedOptionalRequired: "computed",
        filler: {
          model: "openai.AdminOrganizationProjectUserRoleListResponse",
          sourceAttribute: [],
        },
        attributes: [
          {
            name: "id",
            type: "string",
            description: "Identifier for the role.",
            computedOptionalRequired: "computed",
            filler: {
              sourceAttribute: ["ID"],
            },
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
      readStrategy: "simple",
      readMethod: "Admin.Organization.Users.Get",
    },
    filler: {
      model: "openai.OrganizationUser",
    },
    attributes: [
      {
        name: "id",
        type: "string",
        description: "User ID.",
        computedOptionalRequired: "required",
        filler: {
          sourceAttribute: ["ID"],
        },
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
        type: "int64",
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
      readStrategy: "paginate",
      readModel: "OrganizationUser",
      readMethod: "Admin.Organization.Users.ListAutoPaging",
      readRequestParamsStruct: "AdminOrganizationUserListParams",
    },
    filler: {
      model: "[]openai.OrganizationUser",
    },
    attributes: [
      {
        name: "users",
        type: "set_nested",
        description: "List of users.",
        computedOptionalRequired: "computed",
        filler: {
          model: "openai.OrganizationUser",
          sourceAttribute: [],
        },
        attributes: [
          {
            name: "id",
            type: "string",
            description: "User ID.",
            computedOptionalRequired: "computed",
            filler: {
              sourceAttribute: ["ID"],
            },
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
            type: "int64",
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
      readStrategy: "paginate",
      readModel: "AdminOrganizationUserRoleListResponse",
      readMethod: "Admin.Organization.Users.Roles.ListAutoPaging",
      readRequestParamsStruct: "AdminOrganizationUserRoleListParams",
      readRequestAttributes: ["user_id"],
    },
    filler: {
      model: "[]openai.AdminOrganizationUserRoleListResponse",
    },
    attributes: [
      {
        name: "user_id",
        type: "string",
        description: "The ID of the user to inspect.",
        computedOptionalRequired: "required",
        filler: { skip: true },
      },
      {
        name: "roles",
        type: "set_nested",
        description: "List of organization roles",
        computedOptionalRequired: "computed",
        filler: {
          model: "openai.AdminOrganizationUserRoleListResponse",
          sourceAttribute: [],
        },
        attributes: [
          {
            name: "id",
            type: "string",
            description: "Identifier for the role.",
            computedOptionalRequired: "computed",
            filler: {
              sourceAttribute: ["ID"],
            },
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
    name: "spend_limit",
    description: "Retrieves organization spend limit.",
    api: {
      readStrategy: "simple",
      readMethod: "Admin.Organization.SpendLimit.Get",
      readRequestAttributes: [],
    },
    filler: {
      model: "openai.OrganizationSpendLimit",
    },
    attributes: [
      {
        name: "currency",
        type: "string",
        description:
          "The currency for the threshold amount. Currently, only `USD` is supported.",
        computedOptionalRequired: "computed",
      },
      {
        name: "interval",
        type: "string",
        description:
          "The time interval for evaluating spend against the threshold. Currently, only `month` is supported.",
        computedOptionalRequired: "computed",
      },
      {
        name: "threshold_amount",
        type: "int64",
        description: "The hard spend limit amount, in cents.",
        computedOptionalRequired: "computed",
      },
      {
        name: "enforcement",
        type: "single_nested",
        description: "The current enforcement state of the hard spend limit.",
        computedOptionalRequired: "computed",
        filler: {
          model: "openai.OrganizationSpendLimitEnforcement",
        },
        attributes: [
          {
            name: "status",
            type: "string",
            description: "Whether the hard spend limit is currently enforcing.",
            computedOptionalRequired: "computed",
          },
        ],
      },
    ],
  },
  {
    name: "project_spend_limit",
    description: "Retrieves project spend limit.",
    api: {
      readStrategy: "simple",
      readMethod: "Admin.Organization.Projects.SpendLimit.Get",
      readRequestAttributes: ["project_id"],
    },
    filler: {
      model: "openai.ProjectSpendLimit",
    },
    attributes: [
      {
        name: "project_id",
        type: "string",
        description: "The ID of the project.",
        computedOptionalRequired: "required",
        filler: { skip: true },
      },
      {
        name: "currency",
        type: "string",
        description:
          "The currency for the threshold amount. Currently, only `USD` is supported.",
        computedOptionalRequired: "computed",
      },
      {
        name: "interval",
        type: "string",
        description:
          "The time interval for evaluating spend against the threshold. Currently, only `month` is supported.",
        computedOptionalRequired: "computed",
      },
      {
        name: "threshold_amount",
        type: "int64",
        description: "The hard spend limit amount, in cents.",
        computedOptionalRequired: "computed",
      },
      {
        name: "enforcement",
        type: "single_nested",
        description: "The current enforcement state of the hard spend limit.",
        computedOptionalRequired: "computed",
        filler: {
          model: "openai.ProjectSpendLimitEnforcement",
        },
        attributes: [
          {
            name: "status",
            type: "string",
            description: "Whether the hard spend limit is currently enforcing.",
            computedOptionalRequired: "computed",
          },
        ],
      },
    ],
  },
  {
    name: "project_model_permissions",
    description: "Retrieves model access permissions for a project.",
    api: {
      readStrategy: "simple",
      readMethod: "Admin.Organization.Projects.ModelPermissions.Get",
      readRequestAttributes: ["project_id"],
    },
    filler: {
      model: "openai.ProjectModelPermissions",
    },
    attributes: [
      {
        name: "project_id",
        type: "string",
        description:
          "The ID of the project for which model permissions are being set.",
        computedOptionalRequired: "required",
        filler: { skip: true },
      },
      {
        name: "mode",
        type: "string",
        description:
          "The model permissions mode to apply. One of `allow_list` or `deny_list`.",
        computedOptionalRequired: "computed",
      },
      {
        name: "model_ids",
        type: "set",
        elementType: "string",
        description: "The model IDs included in this permissions policy.",
        computedOptionalRequired: "computed",
        filler: {
          sourceAttribute: ["ModelIDs"],
        },
      },
    ],
  },
];

export const RESOURCES: Array<Resource> = [
  {
    name: "admin_api_key",
    description: "Manages an organization admin API key.",
    api: {
      method: "Admin.Organization.AdminAPIKeys",
      createMethod: "New",
      readMethod: "Get",
      readRequestAttributes: ["id"],
      deleteMethod: "Delete",
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
        type: "int64",
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
      method: "Admin.Organization.Invites",
      createMethod: "New",
      readMethod: "Get",
      readRequestAttributes: ["id"],
      deleteMethod: "Delete",
      deleteRequestAttributes: ["id"],
    },
    importStateAttributes: ["id"],
    filler: {
      model: "openai.Invite",
    },
    attributes: [
      {
        name: "id",
        type: "string",
        description: "Invite ID.",
        computedOptionalRequired: "computed",
        planModifiers: ["stringplanmodifier.UseStateForUnknown()"],
        filler: {
          sourceAttribute: ["ID"],
        },
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
        name: "created_at",
        type: "int64",
        description:
          "The Unix timestamp (in seconds) of when the invite was sent.",
        computedOptionalRequired: "computed",
      },
      {
        name: "expires_at",
        type: "int64",
        description:
          "The Unix timestamp (in seconds) of when the invite expires.",
        computedOptionalRequired: "computed",
        nullable: true,
      },
      {
        name: "accepted_at",
        type: "int64",
        description:
          "The Unix timestamp (in seconds) of when the invite was accepted.",
        computedOptionalRequired: "computed",
        nullable: true,
      },
    ],
  },
  {
    name: "organization_role",
    description: "Creates a custom role for the organization.",
    api: {
      method: "Admin.Organization.Roles",
      createMethod: "New",
      readMethod: "Get",
      readRequestAttributes: ["id"],
      updateMethod: "Update",
      updateRequestAttributes: ["id"],
      deleteMethod: "Delete",
      deleteRequestAttributes: ["id"],
    },
    importStateAttributes: ["id"],
    filler: {
      model: "openai.Role",
    },
    attributes: [
      {
        name: "id",
        type: "string",
        description: "Identifier for the role.",
        computedOptionalRequired: "computed",
        planModifiers: ["stringplanmodifier.UseStateForUnknown()"],
        filler: {
          sourceAttribute: ["ID"],
        },
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
      method: "Admin.Organization.Projects",
      createMethod: "New",
      readMethod: "Get",
      readRequestAttributes: ["id"],
      updateMethod: "Update",
      updateRequestAttributes: ["id"],
      deleteMethod: "Archive",
      deleteRequestAttributes: ["id"],
    },
    importStateAttributes: ["id"],
    filler: {
      model: "openai.Project",
    },
    attributes: [
      {
        name: "id",
        type: "string",
        description: "The ID of the project.",
        computedOptionalRequired: "computed",
        planModifiers: ["stringplanmodifier.UseStateForUnknown()"],
        filler: {
          sourceAttribute: ["ID"],
        },
      },
      {
        name: "name",
        type: "string",
        description:
          "The friendly name of the project, this name appears in reports.",
        computedOptionalRequired: "required",
      },
      {
        name: "geography",
        type: "string",
        description:
          "Create the project with the specified data residency region. Your organization must have access to Data residency functionality in order to use.",
        computedOptionalRequired: "optional",
        planModifiers: ["stringplanmodifier.RequiresReplace()"],
        validators: [
          'stringvalidator.OneOf("US", "EU", "JP", "IN", "KR", "CA", "AU", "SG")',
        ],
        filler: { skip: true },
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
        nullable: true,
        filler: {
          sourceAttribute: ["ExternalKeyID"],
        },
      },
      {
        name: "created_at",
        type: "int64",
        description:
          "The Unix timestamp (in seconds) of when the project was created.",
        computedOptionalRequired: "computed",
      },
      {
        name: "archived_at",
        type: "int64",
        description:
          "The Unix timestamp (in seconds) of when the project was archived or `null`.",
        computedOptionalRequired: "computed",
        nullable: true,
      },
    ],
  },
  {
    name: "project_group_role_assignment",
    description: "Assigns a project role to a group within a project.",
    api: {
      method: "Admin.Organization.Projects.Groups.Roles",
      createMethod: "New",
      createRequestAttributes: ["project_id", "group_id"],
      readMethod: "Get",
      readRequestAttributes: ["project_id", "group_id", "role_id"],
      deleteMethod: "Delete",
      deleteRequestAttributes: ["project_id", "group_id", "role_id"],
    },
    importStateAttributes: ["project_id", "group_id", "role_id"],
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
    name: "project_user_role_assignment",
    description: "Assigns a project role to a user within a project.",
    api: {
      method: "Admin.Organization.Projects.Users.Roles",
      createMethod: "New",
      createRequestAttributes: ["project_id", "user_id"],
      readMethod: "Get",
      readRequestAttributes: ["project_id", "user_id", "role_id"],
      deleteMethod: "Delete",
      deleteRequestAttributes: ["project_id", "user_id", "role_id"],
    },
    importStateAttributes: ["project_id", "user_id", "role_id"],
    attributes: [
      {
        name: "project_id",
        type: "string",
        description: "The ID of the project to update.",
        computedOptionalRequired: "required",
        planModifiers: ["stringplanmodifier.RequiresReplace()"],
      },
      {
        name: "user_id",
        type: "string",
        description: "The ID of the user that should receive the project role.",
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
  {
    name: "project_rate_limit",
    description:
      "Manage rate limits per model for projects. Rate limits may be configured to be equal to or lower than the organization's rate limits.",
    api: {
      method: "Admin.Organization.Projects.RateLimits",
      createMethod: "UpdateRateLimit",
      createRequestAttributes: ["project_id", "rate_limit_id"],
      readStrategy: "paginate",
      readModel: "ProjectRateLimit",
      readMethod: "ListRateLimitsAutoPaging",
      readRequestAttributes: ["project_id"],
      readRequestParamsStruct:
        "AdminOrganizationProjectRateLimitListRateLimitsParams",
      updateMethod: "UpdateRateLimit",
      updateRequestAttributes: ["project_id", "rate_limit_id"],
    },
    filler: {
      model: "openai.ProjectRateLimit",
    },
    attributes: [
      {
        name: "project_id",
        type: "string",
        description: "The ID of the project.",
        computedOptionalRequired: "required",
        planModifiers: ["stringplanmodifier.RequiresReplace()"],
        filler: { skip: true },
      },
      {
        name: "rate_limit_id",
        type: "string",
        description:
          "The ID of the rate limit. This is typically in the format `rl-<model>`.",
        computedOptionalRequired: "required",
        planModifiers: ["stringplanmodifier.RequiresReplace()"],
        filler: { skip: true },
      },
      {
        name: "max_requests_per_1_minute",
        type: "int64",
        description: "The maximum requests per minute.",
        computedOptionalRequired: "computed_optional",
      },
      {
        name: "max_tokens_per_1_minute",
        type: "int64",
        description: "The maximum tokens per minute.",
        computedOptionalRequired: "computed_optional",
      },
      {
        name: "max_images_per_1_minute",
        type: "int64",
        description:
          "The maximum images per minute. Only relevant for certain models.",
        computedOptionalRequired: "computed_optional",
        nullable: true,
      },
      {
        name: "max_audio_megabytes_per_1_minute",
        type: "int64",
        description:
          "The maximum audio megabytes per minute. Only relevant for certain models.",
        computedOptionalRequired: "computed_optional",
        nullable: true,
      },
      {
        name: "max_requests_per_1_day",
        type: "int64",
        description:
          "The maximum requests per day. Only relevant for certain models.",
        computedOptionalRequired: "computed_optional",
        nullable: true,
      },
      {
        name: "batch_1_day_max_input_tokens",
        type: "int64",
        description:
          "The maximum batch input tokens per day. Only relevant for certain models.",
        computedOptionalRequired: "computed_optional",
        nullable: true,
      },
    ],
  },
  {
    name: "project_role",
    description: "Creates a custom role for a project.",
    api: {
      method: "Admin.Organization.Projects.Roles",
      createMethod: "New",
      createRequestAttributes: ["project_id"],
      readMethod: "Get",
      readRequestAttributes: ["project_id", "id"],
      updateMethod: "Update",
      updateRequestAttributes: ["project_id", "id"],
      deleteMethod: "Delete",
      deleteRequestAttributes: ["project_id", "id"],
    },
    importStateAttributes: ["project_id", "id"],
    filler: {
      model: "openai.Role",
    },
    attributes: [
      {
        name: "id",
        type: "string",
        description: "Identifier for the role.",
        computedOptionalRequired: "computed",
        planModifiers: ["stringplanmodifier.UseStateForUnknown()"],
        filler: {
          sourceAttribute: ["ID"],
        },
      },
      {
        name: "project_id",
        type: "string",
        description: "The ID of the project to create the role for.",
        computedOptionalRequired: "required",
        filler: { skip: true },
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
      method: "Admin.Organization.Projects.ServiceAccounts",
      createMethod: "New",
      createRequestAttributes: ["project_id"],
      readMethod: "Get",
      readRequestAttributes: ["project_id", "id"],
      updateMethod: "Update",
      updateRequestAttributes: ["project_id", "id"],
      deleteMethod: "Delete",
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
        type: "int64",
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
      method: "Admin.Organization.Projects.Users",
      createMethod: "New",
      createRequestAttributes: ["project_id"],
      readMethod: "Get",
      readRequestAttributes: ["project_id", "user_id"],
      updateMethod: "Update",
      updateRequestAttributes: ["project_id", "user_id"],
      deleteMethod: "Delete",
      deleteRequestAttributes: ["project_id", "user_id"],
    },
    importStateAttributes: ["project_id", "user_id"],
    filler: {
      model: "openai.ProjectUser",
    },
    attributes: [
      {
        name: "project_id",
        type: "string",
        description: "The ID of the project.",
        computedOptionalRequired: "required",
        planModifiers: ["stringplanmodifier.RequiresReplace()"],
        filler: { skip: true },
      },
      {
        name: "user_id",
        type: "string",
        description: "The ID of the user.",
        computedOptionalRequired: "required",
        planModifiers: ["stringplanmodifier.RequiresReplace()"],
        filler: {
          sourceAttribute: ["ID"],
        },
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
      method: "Admin.Organization.Users",
      createMethod: "Update",
      createRequestAttributes: ["user_id"],
      readMethod: "Get",
      readRequestAttributes: ["user_id"],
      updateMethod: "Update",
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
      method: "Admin.Organization.Users.Roles",
      createMethod: "New",
      createRequestAttributes: ["user_id"],
      readMethod: "Get",
      readRequestAttributes: ["user_id", "role_id"],
      updateMethod: "New",
      updateRequestAttributes: ["user_id"],
      deleteMethod: "Delete",
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
      method: "Admin.Organization.Groups",
      createMethod: "New",
      readMethod: "Get",
      readRequestAttributes: ["id"],
      updateMethod: "Update",
      updateRequestAttributes: ["id"],
      deleteMethod: "Delete",
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
        type: "int64",
        description: "Unix timestamp (in seconds) when the group was created.",
        computedOptionalRequired: "computed",
      },
    ],
  },
  {
    name: "group_user",
    description: "Adds a user to a group.",
    api: {
      method: "Admin.Organization.Groups.Users",
      createMethod: "New",
      createRequestAttributes: ["group_id"],
      readMethod: "Get",
      readRequestAttributes: ["group_id", "user_id"],
      deleteMethod: "Delete",
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
      method: "Admin.Organization.Groups.Roles",
      createMethod: "New",
      createRequestAttributes: ["group_id"],
      readMethod: "Get",
      readRequestAttributes: ["group_id", "role_id"],
      deleteMethod: "Delete",
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
  {
    name: "data_retention",
    description: "Updates organization data retention controls.",
    api: {
      method: "Admin.Organization.DataRetention",
      createMethod: "Update",
      readMethod: "Get",
      updateMethod: "Update",
    },
    filler: {
      model: "openai.OrganizationDataRetention",
    },
    attributes: [
      {
        name: "type",
        type: "string",
        description:
          "The desired organization data retention type. Must be one of `zero_data_retention`, `enhanced_zero_data_retention`, `modified_abuse_monitoring`, or `enhanced_modified_abuse_monitoring`.",
        computedOptionalRequired: "required",
        validators: [
          'stringvalidator.OneOf("zero_data_retention", "modified_abuse_monitoring", "enhanced_zero_data_retention", "enhanced_modified_abuse_monitoring")',
        ],
      },
    ],
  },
  {
    name: "spend_limit",
    description: "Updates organization spend limit.",
    api: {
      method: "Admin.Organization.SpendLimit",
      createMethod: "Update",
      readMethod: "Get",
      updateMethod: "Update",
      deleteMethod: "Delete",
    },
    filler: {
      model: "openai.OrganizationSpendLimit",
    },
    attributes: [
      {
        name: "currency",
        type: "string",
        description:
          "The currency for the threshold amount. Currently, only `USD` is supported.",
        computedOptionalRequired: "required",
        validators: ['stringvalidator.OneOf("USD")'],
      },
      {
        name: "interval",
        type: "string",
        description:
          "The time interval for evaluating spend against the threshold. Currently, only `month` is supported.",
        computedOptionalRequired: "required",
        validators: ['stringvalidator.OneOf("month")'],
      },
      {
        name: "threshold_amount",
        type: "int64",
        description: "The hard spend limit amount, in cents.",
        computedOptionalRequired: "required",
        validators: ["int64validator.AtLeast(1)"],
      },
      {
        name: "enforcement",
        type: "single_nested",
        description: "The current enforcement state of the hard spend limit.",
        computedOptionalRequired: "computed",
        filler: {
          model: "openai.OrganizationSpendLimitEnforcement",
        },
        attributes: [
          {
            name: "status",
            type: "string",
            description: "Whether the hard spend limit is currently enforcing.",
            computedOptionalRequired: "computed",
          },
        ],
      },
    ],
  },
  {
    name: "project_spend_limit",
    description: "Updates project spend limit.",
    api: {
      method: "Admin.Organization.Projects.SpendLimit",
      createMethod: "Update",
      createRequestAttributes: ["project_id"],
      readMethod: "Get",
      readRequestAttributes: ["project_id"],
      updateMethod: "Update",
      updateRequestAttributes: ["project_id"],
      deleteMethod: "Delete",
      deleteRequestAttributes: ["project_id"],
    },
    importStateAttributes: ["project_id"],
    filler: {
      model: "openai.ProjectSpendLimit",
    },
    attributes: [
      {
        name: "project_id",
        type: "string",
        description:
          "The ID of the project for which the spend limit is being set.",
        computedOptionalRequired: "required",
        planModifiers: ["stringplanmodifier.RequiresReplace()"],
        filler: { skip: true },
      },
      {
        name: "currency",
        type: "string",
        description:
          "The currency for the threshold amount. Currently, only `USD` is supported.",
        computedOptionalRequired: "required",
        validators: ['stringvalidator.OneOf("USD")'],
      },
      {
        name: "interval",
        type: "string",
        description:
          "The time interval for evaluating spend against the threshold. Currently, only `month` is supported.",
        computedOptionalRequired: "required",
        validators: ['stringvalidator.OneOf("month")'],
      },
      {
        name: "threshold_amount",
        type: "int64",
        description: "The hard spend limit amount, in cents.",
        computedOptionalRequired: "required",
        validators: ["int64validator.AtLeast(1)"],
      },
      {
        name: "enforcement",
        type: "single_nested",
        description: "The current enforcement state of the hard spend limit.",
        computedOptionalRequired: "computed",
        filler: {
          model: "openai.ProjectSpendLimitEnforcement",
        },
        attributes: [
          {
            name: "status",
            type: "string",
            description: "Whether the hard spend limit is currently enforcing.",
            computedOptionalRequired: "computed",
          },
        ],
      },
    ],
  },
  {
    name: "project_model_permissions",
    description: "Updates model access permissions for a project.",
    api: {
      method: "Admin.Organization.Projects.ModelPermissions",
      createMethod: "Update",
      createRequestAttributes: ["project_id"],
      updateMethod: "Update",
      updateRequestAttributes: ["project_id"],
      readMethod: "Get",
      readRequestAttributes: ["project_id"],
      deleteMethod: "Delete",
      deleteRequestAttributes: ["project_id"],
    },
    importStateAttributes: ["project_id"],
    filler: {
      model: "openai.ProjectModelPermissions",
    },
    attributes: [
      {
        name: "project_id",
        type: "string",
        description:
          "The ID of the project for which model permissions are being set.",
        computedOptionalRequired: "required",
        planModifiers: ["stringplanmodifier.RequiresReplace()"],
        filler: { skip: true },
      },
      {
        name: "mode",
        type: "string",
        description:
          "The model permissions mode to apply. One of `allow_list` or `deny_list`.",
        computedOptionalRequired: "required",
        validators: ['stringvalidator.OneOf("allow_list","deny_list")'],
      },
      {
        name: "model_ids",
        type: "set",
        elementType: "string",
        description: "The model IDs included in this permissions policy.",
        computedOptionalRequired: "required",
        filler: {
          sourceAttribute: ["ModelIDs"],
        },
      },
    ],
  },
  {
    name: "project_spend_alert",
    description: "Creates a project spend alert.",
    api: {
      method: "Admin.Organization.Projects.SpendAlerts",
      createMethod: "New",
      createRequestAttributes: ["project_id"],
      readMethod: "Get",
      readRequestAttributes: ["project_id", "id"],
      updateMethod: "Update",
      updateRequestAttributes: ["project_id", "id"],
      deleteMethod: "Delete",
      deleteRequestAttributes: ["project_id", "id"],
    },
    importStateAttributes: ["project_id", "id"],
    filler: {
      model: "openai.ProjectSpendAlert",
    },
    attributes: [
      {
        name: "project_id",
        type: "string",
        description:
          "The ID of the project for which spend alert is being set.",
        computedOptionalRequired: "required",
        planModifiers: ["stringplanmodifier.RequiresReplace()"],
        filler: { skip: true },
      },
      {
        name: "id",
        type: "string",
        description: "Spend alert ID.",
        computedOptionalRequired: "computed",
        planModifiers: ["stringplanmodifier.UseStateForUnknown()"],
        filler: {
          sourceAttribute: ["ID"],
        },
      },
      {
        name: "currency",
        type: "string",
        description: "The currency for the threshold amount (e.g. `USD`).",
        computedOptionalRequired: "required",
      },
      {
        name: "interval",
        type: "string",
        description: "The interval for the spend alert (e.g. `month`).",
        computedOptionalRequired: "required",
      },
      {
        name: "notification_channel",
        type: "single_nested",
        description: "Email notification settings for a spend alert.",
        computedOptionalRequired: "required",
        filler: {
          model: "openai.ProjectSpendAlertNotificationChannel",
        },
        attributes: [
          {
            name: "type",
            type: "string",
            description:
              "The notification channel type. Currently only `email` is supported.",
            computedOptionalRequired: "required",
            validators: ['stringvalidator.OneOf("email")'],
          },
          {
            name: "recipients",
            type: "set",
            description:
              "Email addresses that receive the spend alert notification.",
            computedOptionalRequired: "required",
            elementType: "string",
          },
          {
            name: "subject_prefix",
            type: "string",
            description: "Optional subject prefix for alert emails.",
            computedOptionalRequired: "optional",
            nullable: true,
          },
        ],
      },
      {
        name: "threshold_amount",
        type: "int64",
        description: "The alert threshold amount, in cents.",
        computedOptionalRequired: "required",
      },
    ],
  },
  {
    name: "spend_alert",
    description: "Creates an organization spend alert.",
    api: {
      method: "Admin.Organization.SpendAlerts",
      createMethod: "New",
      readMethod: "Get",
      readRequestAttributes: ["id"],
      updateMethod: "Update",
      updateRequestAttributes: ["id"],
      deleteMethod: "Delete",
      deleteRequestAttributes: ["id"],
    },
    importStateAttributes: ["id"],
    filler: {
      model: "openai.OrganizationSpendAlert",
    },
    attributes: [
      {
        name: "id",
        type: "string",
        description: "Spend alert ID.",
        computedOptionalRequired: "computed",
        planModifiers: ["stringplanmodifier.UseStateForUnknown()"],
        filler: {
          sourceAttribute: ["ID"],
        },
      },
      {
        name: "currency",
        type: "string",
        description: "The currency for the threshold amount (e.g. `USD`).",
        computedOptionalRequired: "required",
      },
      {
        name: "interval",
        type: "string",
        description: "The interval for the spend alert (e.g. `month`).",
        computedOptionalRequired: "required",
      },
      {
        name: "notification_channel",
        type: "single_nested",
        description: "Email notification settings for a spend alert.",
        computedOptionalRequired: "required",
        filler: {
          model: "openai.OrganizationSpendAlertNotificationChannel",
        },
        attributes: [
          {
            name: "type",
            type: "string",
            description:
              "The notification channel type. Currently only `email` is supported.",
            computedOptionalRequired: "required",
            validators: ['stringvalidator.OneOf("email")'],
          },
          {
            name: "recipients",
            type: "set",
            description:
              "Email addresses that receive the spend alert notification.",
            computedOptionalRequired: "required",
            elementType: "string",
          },
          {
            name: "subject_prefix",
            type: "string",
            description: "Optional subject prefix for alert emails.",
            computedOptionalRequired: "optional",
            nullable: true,
          },
        ],
      },
      {
        name: "threshold_amount",
        type: "int64",
        description: "The alert threshold amount, in cents.",
        computedOptionalRequired: "required",
      },
    ],
  },
];
