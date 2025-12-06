import type { DataSource } from "./schema";

export const DATASOURCES: Array<DataSource> = [
  {
    name: "invite",
    description: "Retrieves an invite.",
    api: {
      strategy: "simple",
      method: "RetrieveInvite",
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
      method: "ListInvites",
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
      method: "RetrieveProject",
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
      method: "ListProjects",
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
      method: "ListProjectRateLimits",
      params: ["project_id"],
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
      method: "RetrieveUser",
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
      method: "ListUsers",
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
    ],
  },
];
