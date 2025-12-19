import { zValidator } from "@hono/zod-validator";
import { Hono } from "hono";
import { bearerAuth } from "hono/bearer-auth";
import { logger } from "hono/logger";
import { prettyJSON } from "hono/pretty-json";
import z from "zod";
import { db, insertDefaultProjectRateLimits } from "./db";
import * as schema from "./db-schema";
import { and, eq, inArray } from "drizzle-orm";
import { now } from "./db-utils";

const app = new Hono();
app.use(logger());
app.use(prettyJSON());
app.use(
  "/*",
  bearerAuth({
    verifyToken: async (token) => {
      const apiKey = await db.query.adminApiKeys.findFirst({
        where: (adminApiKeys, { eq }) => eq(adminApiKeys.value, token),
      });
      return !!apiKey;
    },
  })
);

app.get("/", (c) => c.text("Hello World"));

app.post(
  "/organization/admin_api_keys",
  zValidator("json", z.object({ name: z.string() })),
  async (c) => {
    const { name } = c.req.valid("json");

    const [apiKey] = await db
      .insert(schema.adminApiKeys)
      .values({
        name,
      })
      .returning();
    if (!apiKey) {
      return c.json({ error: "Failed to create admin API key" }, 500);
    }

    return c.json({
      ...apiKey,
      redacted_value: `sk-admin-***${apiKey.value.slice(-3)}`,
    });
  }
);

app
  .get("/organization/admin_api_keys/:key_id", async (c) => {
    const key_id = c.req.param("key_id");

    const apiKey = await db.query.adminApiKeys.findFirst({
      where: (adminApiKeys, { eq }) => eq(adminApiKeys.id, key_id),
    });
    if (!apiKey) {
      return c.json({ error: "Admin API key not found" }, 404);
    }

    return c.json({
      ...apiKey,
      redacted_value: `sk-admin-***${apiKey.value.slice(-3)}`,
    });
  })
  .delete(async (c) => {
    const key_id = c.req.param("key_id");

    const result = await db
      .delete(schema.adminApiKeys)
      .where(eq(schema.adminApiKeys.id, key_id))
      .returning();
    if (!result[0]) {
      return c.json({ error: "Admin API key not found" }, 404);
    }

    return c.json({
      object: "organization.admin_api_key.deleted",
      id: result[0].id,
      deleted: true,
    });
  });

app
  .get("/organization/roles", async (c) => {
    const roles = await db.query.roles.findMany({
      where: (roles, { eq }) => eq(roles.resource_type, "api.organization"),
    });

    return c.json({
      object: "list",
      data: roles,
      has_more: false,
      next: null,
    });
  })
  .post(
    zValidator(
      "json",
      z.object({
        permissions: z.array(z.string()),
        role_name: z.string(),
        description: z.string().default(""),
      })
    ),
    async (c) => {
      const { permissions, role_name, description } = c.req.valid("json");

      const [role] = await db
        .insert(schema.roles)
        .values({
          name: role_name,
          description,
          permissions,
          resource_type: "api.organization",
        })
        .returning();

      return c.json(role);
    }
  );

app
  .post(
    "/organization/roles/:role_id",
    zValidator(
      "json",
      z.object({
        permissions: z.array(z.string()),
        role_name: z.string(),
        description: z.string().default(""),
      })
    ),
    async (c) => {
      const role_id = c.req.param("role_id");
      const { permissions, role_name, description } = c.req.valid("json");

      const [role] = await db
        .update(schema.roles)
        .set({
          name: role_name,
          description,
          permissions,
        })
        .where(
          and(
            eq(schema.roles.id, role_id),
            eq(schema.roles.resource_type, "api.organization")
          )
        )
        .returning();

      return c.json(role);
    }
  )
  .delete(async (c) => {
    const role_id = c.req.param("role_id");
    const result = await db
      .delete(schema.roles)
      .where(
        and(
          eq(schema.roles.id, role_id),
          eq(schema.roles.resource_type, "api.organization")
        )
      )
      .returning();

    return c.json({
      object: "role.deleted",
      deleted: result.length > 0,
    });
  });

app
  .get(
    "/organization/groups",
    zValidator("query", z.object({ limit: z.coerce.number().optional() })),
    async (c) => {
      const groups = await db.select().from(schema.groups);

      return c.json({
        object: "list",
        data: groups,
        has_more: false,
        next: null,
      });
    }
  )
  .post(zValidator("json", z.object({ name: z.string() })), async (c) => {
    const data = c.req.valid("json");
    const [group] = await db
      .insert(schema.groups)
      .values({
        name: data.name,
      })
      .returning();

    return c.json(group);
  });

app
  .post(
    "/organization/groups/:group_id",
    zValidator("json", z.object({ name: z.string() })),
    async (c) => {
      const group_id = c.req.param("group_id");
      const data = c.req.valid("json");

      const [group] = await db
        .update(schema.groups)
        .set({
          name: data.name,
        })
        .where(eq(schema.groups.id, group_id))
        .returning();
      if (!group) {
        return c.json({ error: "Group not found" }, 404);
      }

      return c.json(group);
    }
  )
  .delete(async (c) => {
    const group_id = c.req.param("group_id");
    const result = await db
      .delete(schema.groups)
      .where(eq(schema.groups.id, group_id))
      .returning();
    if (!result[0]) {
      return c.json({ error: "Group not found" }, 404);
    }

    return c.json({
      object: "group.deleted",
      id: result[0].id,
      deleted: true,
    });
  });

app.get(
  "/organization/groups/:group_id/roles",
  zValidator("query", z.object({ limit: z.coerce.number().optional() })),
  async (c) => {
    const group_id = c.req.param("group_id");
    const group = await db.query.groups.findFirst({
      where: (groups, { eq }) => eq(groups.id, group_id),
      with: {
        groupsToRoles: {
          with: {
            role: true,
          },
        },
      },
    });
    if (!group) {
      return c.json({ error: "Group not found" }, 404);
    }

    return c.json({
      object: "list",
      data: group.groupsToRoles.map((groupToRole) => ({
        ...groupToRole.role,
      })),
      has_more: false,
      next: null,
    });
  }
);

app.post(
  "/organization/groups/:group_id/roles",
  zValidator("json", z.object({ role_id: z.string() })),
  async (c) => {
    const group_id = c.req.param("group_id");
    const { role_id } = c.req.valid("json");

    const group = await db.query.groups.findFirst({
      where: (groups, { eq }) => eq(groups.id, group_id),
    });
    if (!group) {
      return c.json({ error: "Group not found" }, 404);
    }

    const role = await db.query.roles.findFirst({
      where: (roles, { eq }) => eq(roles.id, role_id),
    });
    if (!role) {
      return c.json({ error: "Role not found" }, 404);
    }

    await db.insert(schema.groupsToRoles).values({
      group_id,
      role_id,
    });

    return c.json({
      object: "group.role",
      group,
      role,
    });
  }
);

app.delete("/organization/groups/:group_id/roles/:role_id", async (c) => {
  const group_id = c.req.param("group_id");
  const role_id = c.req.param("role_id");

  const groupToRole = await db.query.groupsToRoles.findFirst({
    where: (groupsToRoles, { eq }) =>
      and(
        eq(groupsToRoles.group_id, group_id),
        eq(groupsToRoles.role_id, role_id)
      ),
  });
  if (!groupToRole) {
    return c.json({ error: "Group to role not found" }, 404);
  }

  const result = await db
    .delete(schema.groupsToRoles)
    .where(
      and(
        eq(schema.groupsToRoles.group_id, group_id),
        eq(schema.groupsToRoles.role_id, role_id)
      )
    )
    .returning();
  if (!result[0]) {
    return c.json({ error: "Group to role not found" }, 404);
  }

  return c.json({
    object: "group.role.deleted",
    deleted: true,
  });
});

app
  .get("/organization/groups/:group_id/users", async (c) => {
    const group_id = c.req.param("group_id");

    const users = await db.query.groupsToUsers.findMany({
      where: (groupsToUsers, { eq }) => eq(groupsToUsers.group_id, group_id),
      with: {
        user: true,
      },
    });

    return c.json({
      object: "list",
      data: users.map((groupToUser) => groupToUser.user),
      has_more: false,
      next: null,
    });
  })
  .post(zValidator("json", z.object({ user_id: z.string() })), async (c) => {
    const group_id = c.req.param("group_id")!;
    const { user_id } = c.req.valid("json");

    const group = await db.query.groups.findFirst({
      where: (groups, { eq }) => eq(groups.id, group_id),
    });
    if (!group) {
      return c.json({ error: "Group not found" }, 404);
    }

    const user = await db.query.users.findFirst({
      where: (users, { eq }) => eq(users.id, user_id),
    });
    if (!user) {
      return c.json({ error: "User not found" }, 404);
    }

    await db.insert(schema.groupsToUsers).values({
      group_id,
      user_id,
    });

    return c.json({
      object: "group.user",
      user_id: user.id,
      group_id: group.id,
    });
  });

app.delete("/organization/groups/:group_id/users/:user_id", async (c) => {
  const group_id = c.req.param("group_id");
  const user_id = c.req.param("user_id");

  const groupToUser = await db.query.groupsToUsers.findFirst({
    where: (groupsToUsers, { eq }) =>
      and(
        eq(groupsToUsers.group_id, group_id),
        eq(groupsToUsers.user_id, user_id)
      ),
  });
  if (!groupToUser) {
    return c.json({ error: "Group to user not found" }, 404);
  }

  const result = await db
    .delete(schema.groupsToUsers)
    .where(
      and(
        eq(schema.groupsToUsers.group_id, group_id),
        eq(schema.groupsToUsers.user_id, user_id)
      )
    )
    .returning();
  if (!result[0]) {
    return c.json({ error: "Group to user not found" }, 404);
  }

  return c.json({
    object: "group.user.deleted",
    deleted: true,
  });
});

app.get("/organization/users", async (c) => {
  const users = await db.query.users.findMany();

  return c.json({
    object: "list",
    data: users,
    has_more: false,
    first_id: users.at(0)?.id,
    last_id: users.at(-1)?.id,
  });
});

app
  .get("/organization/users/:user_id", async (c) => {
    const user_id = c.req.param("user_id");

    const user = await db.query.users.findFirst({
      where: (users, { eq }) => eq(users.id, user_id),
    });
    if (!user) {
      return c.json({ error: "User not found" }, 404);
    }

    return c.json(user);
  })
  .post(
    zValidator("json", z.object({ role: z.enum(["owner", "reader"]) })),
    async (c) => {
      const user_id = c.req.param("user_id")!;
      const { role } = c.req.valid("json");

      const [user] = await db
        .update(schema.users)
        .set({
          role,
        })
        .where(eq(schema.users.id, user_id))
        .returning();

      // Built-in role
      await db
        .delete(schema.usersToRoles)
        .where(
          and(
            eq(schema.usersToRoles.user_id, user_id),
            inArray(schema.usersToRoles.role_id, [
              "role_organization_owner",
              "role_organization_reader",
            ])
          )
        );
      await db.insert(schema.usersToRoles).values({
        user_id,
        role_id:
          role === "owner"
            ? "role_organization_owner"
            : "role_organization_reader",
      });

      return c.json(user);
    }
  );

app
  .get("/organization/users/:user_id/roles", async (c) => {
    const user_id = c.req.param("user_id");

    const roles = await db.query.usersToRoles.findMany({
      where: (usersToRoles, { eq }) => eq(usersToRoles.user_id, user_id),
      with: {
        role: true,
      },
    });

    return c.json({
      object: "list",
      data: roles.map((userToRole) => userToRole.role),
      has_more: false,
      next: null,
    });
  })
  .post(zValidator("json", z.object({ role_id: z.string() })), async (c) => {
    const user_id = c.req.param("user_id")!;
    const { role_id } = c.req.valid("json");

    const user = await db.query.users.findFirst({
      where: (users, { eq }) => eq(users.id, user_id),
    });
    if (!user) {
      return c.json({ error: "User not found" }, 404);
    }

    const role = await db.query.roles.findFirst({
      where: (roles, { eq, and }) =>
        and(eq(roles.id, role_id), eq(roles.resource_type, "api.organization")),
    });
    if (!role) {
      return c.json({ error: "Role not found" }, 404);
    }

    await db.insert(schema.usersToRoles).values({
      user_id,
      role_id,
    });

    return c.json({
      object: "user.role",
      user,
      role,
    });
  });

app.delete("/organization/users/:user_id/roles/:role_id", async (c) => {
  const user_id = c.req.param("user_id");
  const role_id = c.req.param("role_id");

  const result = await db
    .delete(schema.usersToRoles)
    .where(
      and(
        eq(schema.usersToRoles.user_id, user_id),
        eq(schema.usersToRoles.role_id, role_id)
      )
    )
    .returning();
  if (!result[0]) {
    return c.json({ error: "User to role not found" }, 404);
  }

  return c.json({
    object: "user.role.deleted",
    deleted: true,
  });
});

app
  .get("/organization/invites", async (c) => {
    const invites = await db.query.invites.findMany();
    return c.json({
      object: "list",
      data: invites,
      has_more: false,
      first_id: invites.at(0)?.id,
      last_id: invites.at(-1)?.id,
    });
  })
  .post(
    zValidator(
      "json",
      z.object({ email: z.string(), role: z.enum(["owner", "reader"]) })
    ),
    async (c) => {
      const { email, role } = c.req.valid("json");

      const [invite] = await db
        .insert(schema.invites)
        .values({
          email,
          role,
        })
        .returning();

      return c.json(invite);
    }
  );

app
  .get("/organization/invites/:invite_id", async (c) => {
    const invite_id = c.req.param("invite_id");

    const invite = await db.query.invites.findFirst({
      where: (invites, { eq }) => eq(invites.id, invite_id),
    });
    if (!invite) {
      return c.json({ error: "Invite not found" }, 404);
    }

    return c.json(invite);
  })
  .delete(async (c) => {
    const invite_id = c.req.param("invite_id");
    const result = await db
      .delete(schema.invites)
      .where(eq(schema.invites.id, invite_id))
      .returning();
    if (!result[0]) {
      return c.json({ error: "Invite not found" }, 404);
    }
    return c.json({
      object: "organization.invite.deleted",
      id: result[0].id,
      deleted: true,
    });
  });

app
  .get(
    "/organization/projects",
    zValidator(
      "query",
      z.object({ include_archived: z.coerce.boolean().default(false) })
    ),
    async (c) => {
      const include_archived = c.req.valid("query").include_archived;
      const projects = await db.query.projects.findMany({
        where: (projects, { or, eq }) =>
          or(
            eq(projects.status, "active"),
            include_archived ? eq(projects.status, "archived") : undefined
          ),
      });

      return c.json({
        object: "list",
        data: projects,
        has_more: false,
        first_id: projects.at(0)?.id,
        last_id: projects.at(-1)?.id,
      });
    }
  )
  .post(zValidator("json", z.object({ name: z.string() })), async (c) => {
    const { name } = c.req.valid("json");

    const [project] = await db
      .insert(schema.projects)
      .values({
        name,
      })
      .returning();

    await insertDefaultProjectRateLimits({ projectId: project!.id });

    return c.json(project);
  });

app
  .get("/organization/projects/:project_id", async (c) => {
    const project_id = c.req.param("project_id");

    const project = await db.query.projects.findFirst({
      where: (projects, { eq }) => eq(projects.id, project_id),
    });
    if (!project) {
      return c.json({ error: "Project not found" }, 404);
    }

    return c.json(project);
  })
  .post(zValidator("json", z.object({ name: z.string() })), async (c) => {
    const project_id = c.req.param("project_id")!;
    const { name } = c.req.valid("json");

    const result = await db
      .update(schema.projects)
      .set({ name })
      .where(eq(schema.projects.id, project_id))
      .returning();
    if (!result[0]) {
      return c.json({ error: "Project not found" }, 404);
    }

    return c.json(result[0]);
  });

app.post("/organization/projects/:project_id/archive", async (c) => {
  const project_id = c.req.param("project_id");

  const result = await db
    .update(schema.projects)
    .set({ status: "archived", archived_at: now() })
    .where(eq(schema.projects.id, project_id))
    .returning();
  if (!result[0]) {
    return c.json({ error: "Project not found" }, 404);
  }

  return c.json(result[0]);
});

app
  .get("/projects/:project_id/roles", async (c) => {
    const project_id = c.req.param("project_id");

    const roles = await db.query.projectsToRoles.findMany({
      where: (projectsToRoles, { eq }) =>
        eq(projectsToRoles.project_id, project_id),
      with: {
        role: true,
      },
    });

    return c.json({
      object: "list",
      data: roles.map((role) => role.role),
      has_more: false,
      next: null,
    });
  })
  .post(
    zValidator(
      "json",
      z.object({
        permissions: z.array(z.string()),
        role_name: z.string(),
        description: z.string().default(""),
      })
    ),
    async (c) => {
      const project_id = c.req.param("project_id")!;
      const { permissions, role_name, description } = c.req.valid("json");

      const [role] = await db
        .insert(schema.roles)
        .values({
          name: role_name,
          description,
          permissions,
          resource_type: "api.project",
          predefined_role: false,
        })
        .returning();

      await db.insert(schema.projectsToRoles).values({
        project_id,
        role_id: role!.id,
      });

      return c.json(role);
    }
  );

app.delete("/projects/:project_id/roles/:role_id", async (c) => {
  const project_id = c.req.param("project_id")!;
  const role_id = c.req.param("role_id")!;

  const result = await db
    .delete(schema.projectsToRoles)
    .where(
      and(
        eq(schema.projectsToRoles.project_id, project_id),
        eq(schema.projectsToRoles.role_id, role_id)
      )
    )
    .returning();
  if (!result[0]) {
    return c.json({ error: "Role not found" }, 404);
  }

  return c.json({
    object: "role.deleted",
    id: result[0].role_id,
    deleted: true,
  });
});

app.post(
  "/organization/projects/:project_id/users",
  zValidator(
    "json",
    z.object({ user_id: z.string(), role: z.enum(["owner", "member"]) })
  ),
  async (c) => {
    const project_id = c.req.param("project_id")!;
    const { user_id, role } = c.req.valid("json");

    const project = await db.query.projects.findFirst({
      where: (projects, { eq }) => eq(projects.id, project_id),
    });
    if (!project) {
      return c.json({ error: "Project not found" }, 404);
    }

    const user = await db.query.users.findFirst({
      where: (users, { eq }) => eq(users.id, user_id),
    });
    if (!user) {
      return c.json({ error: "User not found" }, 404);
    }

    const [projectToUser] = await db
      .insert(schema.projectsToUsers)
      .values({
        project_id,
        user_id,
        role,
      })
      .returning();
    if (!projectToUser) {
      return c.json({ error: "Failed to add user to project" }, 500);
    }

    // Built-in role
    await db
      .insert(schema.projectsToUsersToRoles)
      .values({
        project_id,
        user_id,
        role_id:
          role === "owner" ? "role_project_owner" : "role_project_member",
      })
      .onConflictDoNothing();

    return c.json({
      object: "organization.project.user",
      id: user.id,
      email: user.email,
      role: projectToUser.role,
      added_at: projectToUser.added_at,
    });
  }
);

app
  .get("/organization/projects/:project_id/users/:user_id", async (c) => {
    const project_id = c.req.param("project_id");
    const user_id = c.req.param("user_id");

    const projectToUser = await db.query.projectsToUsers.findFirst({
      where: (projectsToUsers, { eq }) =>
        and(
          eq(projectsToUsers.project_id, project_id),
          eq(projectsToUsers.user_id, user_id)
        ),
      with: {
        user: true,
      },
    });
    if (!projectToUser) {
      return c.json({ error: "User not found" }, 404);
    }

    return c.json({
      object: "organization.project.user",
      id: projectToUser.user.id,
      email: projectToUser.user.email,
      role: projectToUser.role,
      added_at: projectToUser.added_at,
    });
  })
  .post(
    zValidator("json", z.object({ role: z.enum(["owner", "member"]) })),
    async (c) => {
      const project_id = c.req.param("project_id")!;
      const user_id = c.req.param("user_id")!;
      const { role } = c.req.valid("json");

      const projectToUser = await db.query.projectsToUsers.findFirst({
        where: (projectsToUsers, { eq }) =>
          and(
            eq(projectsToUsers.project_id, project_id),
            eq(projectsToUsers.user_id, user_id)
          ),
        with: {
          user: true,
        },
      });
      if (!projectToUser) {
        return c.json({ error: "User not found" }, 404);
      }

      const [updatedProjectToUser] = await db
        .update(schema.projectsToUsers)
        .set({ role })
        .where(
          and(
            eq(schema.projectsToUsers.project_id, project_id),
            eq(schema.projectsToUsers.user_id, user_id)
          )
        )
        .returning();
      if (!updatedProjectToUser) {
        return c.json({ error: "Failed to update user role" }, 500);
      }

      // Built-in role
      await db
        .delete(schema.projectsToUsersToRoles)
        .where(
          and(
            eq(schema.projectsToUsersToRoles.project_id, project_id),
            eq(schema.projectsToUsersToRoles.user_id, user_id),
            inArray(schema.projectsToUsersToRoles.role_id, [
              "role_project_owner",
              "role_project_member",
            ])
          )
        );

      await db
        .insert(schema.projectsToUsersToRoles)
        .values({
          project_id,
          user_id,
          role_id:
            role === "owner" ? "role_project_owner" : "role_project_member",
        })
        .onConflictDoNothing();

      return c.json({
        object: "organization.project.user",
        id: projectToUser.user.id,
        email: projectToUser.user.email,
        role: updatedProjectToUser.role,
        added_at: updatedProjectToUser.added_at,
      });
    }
  )
  .delete(async (c) => {
    const project_id = c.req.param("project_id")!;
    const user_id = c.req.param("user_id")!;

    const result = await db
      .delete(schema.projectsToUsers)
      .where(
        and(
          eq(schema.projectsToUsers.project_id, project_id),
          eq(schema.projectsToUsers.user_id, user_id)
        )
      )
      .returning();
    if (!result[0]) {
      return c.json({ error: "User not found" }, 404);
    }

    return c.json({
      object: "organization.project.user.deleted",
      id: result[0].user_id,
      deleted: true,
    });
  });

app.get("/organization/projects/:project_id/rate_limits", async (c) => {
  const project_id = c.req.param("project_id")!;

  const rate_limits = await db.query.projectRateLimits.findMany({
    where: (projectRateLimits, { eq }) =>
      eq(projectRateLimits.project_id, project_id),
  });

  return c.json({
    object: "list",
    data: rate_limits,
    has_more: false,
    first_id: rate_limits.at(0)?.id,
    last_id: rate_limits.at(-1)?.id,
  });
});

app
  .get("/projects/:project_id/groups/:group_id/roles", async (c) => {
    const project_id = c.req.param("project_id")!;
    const group_id = c.req.param("group_id")!;

    const roles = await db.query.projectsToGroupsToRoles.findMany({
      where: (projectsToGroupsToRoles, { eq }) =>
        and(
          eq(projectsToGroupsToRoles.project_id, project_id),
          eq(projectsToGroupsToRoles.group_id, group_id)
        ),
      with: {
        role: true,
      },
    });

    return c.json({
      object: "list",
      data: roles.map((role) => role.role),
      has_more: false,
      next: null,
    });
  })
  .post(zValidator("json", z.object({ role_id: z.string() })), async (c) => {
    const project_id = c.req.param("project_id")!;
    const group_id = c.req.param("group_id")!;
    const { role_id } = c.req.valid("json");

    const project = await db.query.projects.findFirst({
      where: (projects, { eq }) => eq(projects.id, project_id),
    });
    if (!project) {
      return c.json({ error: "Project not found" }, 404);
    }

    const group = await db.query.groups.findFirst({
      where: (groups, { eq }) => eq(groups.id, group_id),
    });
    if (!group) {
      return c.json({ error: "Group not found" }, 404);
    }

    const role = await db.query.roles.findFirst({
      where: (roles, { eq }) => eq(roles.id, role_id),
    });
    if (!role) {
      return c.json({ error: "Role not found" }, 404);
    }

    await db.insert(schema.projectsToGroupsToRoles).values({
      project_id,
      group_id,
      role_id,
    });

    return c.json({
      object: "group.role",
      group,
      role,
    });
  });

app.delete(
  "/projects/:project_id/groups/:group_id/roles/:role_id",
  async (c) => {
    const project_id = c.req.param("project_id")!;
    const group_id = c.req.param("group_id")!;
    const role_id = c.req.param("role_id")!;

    const result = await db
      .delete(schema.projectsToGroupsToRoles)
      .where(
        and(
          eq(schema.projectsToGroupsToRoles.project_id, project_id),
          eq(schema.projectsToGroupsToRoles.group_id, group_id),
          eq(schema.projectsToGroupsToRoles.role_id, role_id)
        )
      )
      .returning();
    if (!result[0]) {
      return c.json({ error: "Role not found" }, 404);
    }

    return c.json({
      object: "group.role.deleted",
      deleted: true,
    });
  }
);

app
  .get("/projects/:project_id/users/:user_id/roles", async (c) => {
    const project_id = c.req.param("project_id")!;
    const user_id = c.req.param("user_id")!;

    const roles = await db.query.projectsToUsersToRoles.findMany({
      where: (projectsToUsersToRoles, { eq }) =>
        and(
          eq(projectsToUsersToRoles.project_id, project_id),
          eq(projectsToUsersToRoles.user_id, user_id)
        ),
      with: {
        role: true,
      },
    });

    return c.json({
      object: "list",
      data: roles.map((role) => role.role),
      has_more: false,
      next: null,
    });
  })
  .post(zValidator("json", z.object({ role_id: z.string() })), async (c) => {
    const project_id = c.req.param("project_id")!;
    const user_id = c.req.param("user_id")!;
    const { role_id } = c.req.valid("json");

    const project = await db.query.projects.findFirst({
      where: (projects, { eq }) => eq(projects.id, project_id),
    });
    if (!project) {
      return c.json({ error: "Project not found" }, 404);
    }

    const user = await db.query.users.findFirst({
      where: (users, { eq }) => eq(users.id, user_id),
    });
    if (!user) {
      return c.json({ error: "User not found" }, 404);
    }

    const role = await db.query.roles.findFirst({
      where: (roles, { eq }) => eq(roles.id, role_id),
    });
    if (!role) {
      return c.json({ error: "Role not found" }, 404);
    }

    await db.insert(schema.projectsToUsersToRoles).values({
      project_id,
      user_id,
      role_id,
    });

    return c.json({
      object: "user.role",
      user,
      role,
    });
  });

app.delete("/projects/:project_id/users/:user_id/roles/:role_id", async (c) => {
  const project_id = c.req.param("project_id")!;
  const user_id = c.req.param("user_id")!;
  const role_id = c.req.param("role_id")!;

  const result = await db
    .delete(schema.projectsToUsersToRoles)
    .where(
      and(
        eq(schema.projectsToUsersToRoles.project_id, project_id),
        eq(schema.projectsToUsersToRoles.user_id, user_id),
        eq(schema.projectsToUsersToRoles.role_id, role_id)
      )
    )
    .returning();
  if (!result[0]) {
    return c.json({ error: "Role not found" }, 404);
  }

  return c.json({
    object: "user.role.deleted",
    deleted: true,
  });
});

export default app;
