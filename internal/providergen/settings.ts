import type { DataSource, Resource } from "./schema";

export const DATASOURCES: Array<DataSource> = [
  {
    name: "invite",
    description: "Retrieves an invite.",
    api: {
      strategy: "simple",
      readMethod: "RetrieveInvite",
      model: "Invite",
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
      strategy: "paginate",
      readMethod: "ListInvites",
      model: "Invite",
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
    name: "project",
    description: "Retrieve a project by ID.",
    api: {
      strategy: "simple",
      readMethod: "RetrieveProject",
      model: "Project",
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
      strategy: "paginate",
      readMethod: "ListProjects",
      model: "Project",
      hooks: {
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
      strategy: "paginate",
      readMethod: "ListProjectRateLimits",
      readRequestAttributes: ["project_id"],
      model: "ProjectRateLimit",
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
    name: "user",
    description: "Retrieves a user by their identifier.",
    api: {
      strategy: "simple",
      readMethod: "RetrieveUser",
      model: "User",
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
      strategy: "paginate",
      readMethod: "ListUsers",
      model: "User",
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
];
