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
