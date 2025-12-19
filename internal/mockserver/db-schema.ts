import { relations } from "drizzle-orm";
import {
  integer,
  primaryKey,
  sqliteTable,
  text,
} from "drizzle-orm/sqlite-core";
import { idGenerator, now } from "./db-utils";

function objectColumn(key: string) {
  return text({ enum: [key] })
    .notNull()
    .default(key);
}

const created_at = integer().notNull().$defaultFn(now);
const project_id = text()
  .notNull()
  .references(() => projects.id, { onDelete: "cascade" });
const user_id = text()
  .notNull()
  .references(() => users.id, { onDelete: "cascade" });
const role_id = text()
  .notNull()
  .references(() => roles.id, { onDelete: "cascade" });

export const adminApiKeys = sqliteTable("admin_api_keys", {
  object: objectColumn("organization.admin_api_key"),
  id: text().primaryKey().$defaultFn(idGenerator("admin_api_key")),
  name: text().notNull(),
  value: text().notNull().$defaultFn(idGenerator("sk-admin", "-")),
  created_at,
});

export const users = sqliteTable("users", {
  object: objectColumn("organization.user"),
  id: text().primaryKey().$defaultFn(idGenerator("user")),
  name: text().notNull(),
  email: text().notNull(),
  role: text({ enum: ["owner", "reader"] }).notNull(),
  added_at: created_at,
});

export const usersRelation = relations(users, ({ many }) => ({
  groupsToUsers: many(groupsToUsers),
}));

export const usersToRoles = sqliteTable(
  "users_to_roles",
  {
    user_id,
    role_id,
  },
  (table) => [primaryKey({ columns: [table.user_id, table.role_id] })]
);

export const usersToRolesRelation = relations(usersToRoles, ({ one }) => ({
  user: one(users, {
    fields: [usersToRoles.user_id],
    references: [users.id],
  }),
  role: one(roles, {
    fields: [usersToRoles.role_id],
    references: [roles.id],
  }),
}));

export const groups = sqliteTable("groups", {
  object: objectColumn("group"),
  id: text().primaryKey().$defaultFn(idGenerator("group")),
  name: text().notNull(),
  is_scim_managed: integer({ mode: "boolean" }).notNull().default(false),
  created_at,
});

export const groupsRelation = relations(groups, ({ many }) => ({
  groupsToRoles: many(groupsToRoles),
  groupsToUsers: many(groupsToUsers),
}));

export const roles = sqliteTable("roles", {
  object: objectColumn("role"),
  id: text().primaryKey().$defaultFn(idGenerator("role")),
  name: text().notNull(),
  description: text().notNull(),
  permissions: text({ mode: "json" }).notNull().$type<string[]>(),
  resource_type: text({ enum: ["api.organization", "api.project"] }).notNull(),
  predefined_role: integer({ mode: "boolean" }).notNull().default(false),
  created_at,
});

export const rolesRelation = relations(roles, ({ many }) => ({
  groupsToRoles: many(groupsToRoles),
}));

export const groupsToRoles = sqliteTable(
  "groups_to_roles",
  {
    group_id: text().notNull(),
    role_id,
  },
  (table) => [primaryKey({ columns: [table.group_id, table.role_id] })]
);

export const groupsToRolesRelation = relations(groupsToRoles, ({ one }) => ({
  group: one(groups, {
    fields: [groupsToRoles.group_id],
    references: [groups.id],
  }),
  role: one(roles, {
    fields: [groupsToRoles.role_id],
    references: [roles.id],
  }),
}));

export const groupsToUsers = sqliteTable(
  "groups_to_users",
  {
    group_id: text().notNull(),
    user_id,
  },
  (table) => [primaryKey({ columns: [table.group_id, table.user_id] })]
);

export const groupsToUsersRelation = relations(groupsToUsers, ({ one }) => ({
  group: one(groups, {
    fields: [groupsToUsers.group_id],
    references: [groups.id],
  }),
  user: one(users, {
    fields: [groupsToUsers.user_id],
    references: [users.id],
  }),
}));

export const invites = sqliteTable("invites", {
  object: objectColumn("organization.invite"),
  id: text().primaryKey().$defaultFn(idGenerator("invite")),
  email: text().notNull(),
  role: text({ enum: ["owner", "reader"] }).notNull(),
  status: text({ enum: ["accepted", "expired", "pending"] })
    .notNull()
    .default("pending"),
  accepted_at: integer(),
  expires_at: integer().$defaultFn(() => now() + 60 * 60 * 24),
  invited_at: created_at,
});

export const projects = sqliteTable("projects", {
  object: objectColumn("organization.project"),
  id: text().primaryKey().$defaultFn(idGenerator("project")),
  name: text().notNull(),
  status: text({ enum: ["active", "archived"] })
    .notNull()
    .default("active"),
  archived_at: integer(),
  created_at,
});

export const projectsRelation = relations(projects, ({ many }) => ({
  projectsToRoles: many(projectsToRoles),
}));

export const projectsToRoles = sqliteTable(
  "projects_to_roles",
  {
    project_id,
    role_id,
  },
  (table) => [primaryKey({ columns: [table.project_id, table.role_id] })]
);

export const projectsToRolesRelation = relations(
  projectsToRoles,
  ({ one }) => ({
    project: one(projects, {
      fields: [projectsToRoles.project_id],
      references: [projects.id],
    }),
    role: one(roles, {
      fields: [projectsToRoles.role_id],
      references: [roles.id],
    }),
  })
);

export const projectsToGroupsToRoles = sqliteTable(
  "projects_to_groups_to_roles",
  {
    project_id,
    group_id: text().notNull(),
    role_id,
  },
  (table) => [
    primaryKey({ columns: [table.project_id, table.group_id, table.role_id] }),
  ]
);

export const projectsToGroupsToRolesRelation = relations(
  projectsToGroupsToRoles,
  ({ one }) => ({
    project: one(projects, {
      fields: [projectsToGroupsToRoles.project_id],
      references: [projects.id],
    }),
    group: one(groups, {
      fields: [projectsToGroupsToRoles.group_id],
      references: [groups.id],
    }),
    role: one(roles, {
      fields: [projectsToGroupsToRoles.role_id],
      references: [roles.id],
    }),
  })
);

export const projectsToUsers = sqliteTable(
  "projects_to_users",
  {
    project_id,
    user_id,
    role: text({ enum: ["owner", "member"] }).notNull(),
    added_at: created_at,
  },
  (table) => [primaryKey({ columns: [table.project_id, table.user_id] })]
);

export const projectsToUsersRelation = relations(
  projectsToUsers,
  ({ one }) => ({
    project: one(projects, {
      fields: [projectsToUsers.project_id],
      references: [projects.id],
    }),
    user: one(users, {
      fields: [projectsToUsers.user_id],
      references: [users.id],
    }),
  })
);

export const projectsToUsersToRoles = sqliteTable(
  "projects_to_users_to_roles",
  {
    project_id,
    user_id,
    role_id,
  },
  (table) => [
    primaryKey({ columns: [table.project_id, table.user_id, table.role_id] }),
  ]
);

export const projectsToUsersToRolesRelation = relations(
  projectsToUsersToRoles,
  ({ one }) => ({
    project: one(projects, {
      fields: [projectsToUsersToRoles.project_id],
      references: [projects.id],
    }),
    user: one(users, {
      fields: [projectsToUsersToRoles.user_id],
      references: [users.id],
    }),
    role: one(roles, {
      fields: [projectsToUsersToRoles.role_id],
      references: [roles.id],
    }),
  })
);

export const projectRateLimits = sqliteTable("project_rate_limits", {
  object: objectColumn("project.rate_limit"),
  id: text().primaryKey().$defaultFn(idGenerator("rl")),
  project_id,
  model: text().notNull(),
  batch_1_day_max_input_tokens: integer().notNull(),
  max_audio_megabytes_per_1_minute: integer().notNull(),
  max_images_per_1_minute: integer().notNull(),
  max_requests_per_1_day: integer().notNull(),
  max_requests_per_1_minute: integer().notNull(),
  max_tokens_per_1_minute: integer().notNull(),
});
