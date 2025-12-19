import { zValidator } from "@hono/zod-validator";
import { and, eq, inArray, or } from "drizzle-orm";
import { Hono } from "hono";
import { bearerAuth } from "hono/bearer-auth";
import { logger } from "hono/logger";
import { prettyJSON } from "hono/pretty-json";
import z from "zod";
import { db, insertDefaultProjectRateLimits } from "./db";
import * as schema from "./db-schema";
import { now } from "./db-utils";

const app = new Hono();
app.use(logger());
app.use(prettyJSON());
app.use(
  "/*",
  bearerAuth({
    verifyToken: async (token) => {
      const apiKey = await db.query.adminApiKeys.findFirst({
        where: eq(schema.adminApiKeys.value, token),
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

app.get("/organization/admin_api_keys/:key_id", async (c) => {
  const key_id = c.req.param("key_id");

  const apiKey = await db.query.adminApiKeys.findFirst({
    where: eq(schema.adminApiKeys.id, key_id),
  });
  if (!apiKey) {
    return c.json({ error: "Admin API key not found" }, 404);
  }

  return c.json({
    ...apiKey,
    redacted_value: `sk-admin-***${apiKey.value.slice(-3)}`,
  });
});

app.delete("/organization/admin_api_keys/:key_id", async (c) => {
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

app.get("/organization/roles", async (c) => {
  const roles = await db.query.roles.findMany({
    where: eq(schema.roles.resource_type, "api.organization"),
  });

  return c.json({
    object: "list",
    data: roles,
    has_more: false,
    next: null,
  });
});

app.post(
  "/organization/roles",
  zValidator(
    "json",
    z.object({
      permissions: z.array(z.string()),
      role_name: z.string(),
      description: z.string().default(""),
    })
  ),
  async (c) => {
    const { permissions, role_name: name, description } = c.req.valid("json");

    const [role] = await db
      .insert(schema.roles)
      .values({
        name,
        description,
        permissions,
        resource_type: "api.organization",
      })
      .returning();

    return c.json(role);
  }
);

app.post(
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
    const { permissions, role_name: name, description } = c.req.valid("json");

    const [updatedRole] = await db
      .update(schema.roles)
      .set({
        name,
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
    if (!updatedRole) {
      return c.json({ error: "Role not found" }, 404);
    }

    return c.json(updatedRole);
  }
);

app.delete("/organization/roles/:role_id", async (c) => {
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
  if (!result[0]) {
    return c.json({ error: "Role not found" }, 404);
  }

  return c.json({
    object: "role.deleted",
    deleted: true,
  });
});

app.get(
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
);

app.post(
  "/organization/groups",
  zValidator("json", z.object({ name: z.string() })),
  async (c) => {
    const { name } = c.req.valid("json");

    const [group] = await db
      .insert(schema.groups)
      .values({
        name,
      })
      .returning();

    return c.json(group);
  }
);

app.post(
  "/organization/groups/:group_id",
  zValidator("json", z.object({ name: z.string() })),
  async (c) => {
    const group_id = c.req.param("group_id");
    const { name } = c.req.valid("json");

    const [group] = await db
      .update(schema.groups)
      .set({
        name,
      })
      .where(eq(schema.groups.id, group_id))
      .returning();
    if (!group) {
      return c.json({ error: "Group not found" }, 404);
    }

    return c.json(group);
  }
);

app.delete("/organization/groups/:group_id", async (c) => {
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
      where: eq(schema.groups.id, group_id),
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
      where: eq(schema.groups.id, group_id),
    });
    if (!group) {
      return c.json({ error: "Group not found" }, 404);
    }

    const role = await db.query.roles.findFirst({
      where: eq(schema.roles.id, role_id),
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
    where: and(
      eq(schema.groupsToRoles.group_id, group_id),
      eq(schema.groupsToRoles.role_id, role_id)
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

app.get("/organization/groups/:group_id/users", async (c) => {
  const group_id = c.req.param("group_id");

  const users = await db.query.groupsToUsers.findMany({
    where: eq(schema.groupsToUsers.group_id, group_id),
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
});

app.post(
  "/organization/groups/:group_id/users",
  zValidator("json", z.object({ user_id: z.string() })),
  async (c) => {
    const group_id = c.req.param("group_id")!;
    const { user_id } = c.req.valid("json");

    const group = await db.query.groups.findFirst({
      where: eq(schema.groups.id, group_id),
    });
    if (!group) {
      return c.json({ error: "Group not found" }, 404);
    }

    const user = await db.query.users.findFirst({
      where: eq(schema.users.id, user_id),
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
  }
);

app.delete("/organization/groups/:group_id/users/:user_id", async (c) => {
  const group_id = c.req.param("group_id");
  const user_id = c.req.param("user_id");

  const groupToUser = await db.query.groupsToUsers.findFirst({
    where: and(
      eq(schema.groupsToUsers.group_id, group_id),
      eq(schema.groupsToUsers.user_id, user_id)
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

app.get("/organization/users/:user_id", async (c) => {
  const user_id = c.req.param("user_id");

  const user = await db.query.users.findFirst({
    where: eq(schema.users.id, user_id),
  });
  if (!user) {
    return c.json({ error: "User not found" }, 404);
  }

  return c.json(user);
});

app.post(
  "/organization/users/:user_id",
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

app.get("/organization/users/:user_id/roles", async (c) => {
  const user_id = c.req.param("user_id");

  const roles = await db.query.usersToRoles.findMany({
    where: eq(schema.usersToRoles.user_id, user_id),
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
});

app.post(
  "/organization/users/:user_id/roles",
  zValidator("json", z.object({ role_id: z.string() })),
  async (c) => {
    const user_id = c.req.param("user_id")!;
    const { role_id } = c.req.valid("json");

    const user = await db.query.users.findFirst({
      where: eq(schema.users.id, user_id),
    });
    if (!user) {
      return c.json({ error: "User not found" }, 404);
    }

    const role = await db.query.roles.findFirst({
      where: and(
        eq(schema.roles.id, role_id),
        eq(schema.roles.resource_type, "api.organization")
      ),
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
  }
);

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

app.get("/organization/invites", async (c) => {
  const invites = await db.query.invites.findMany();

  return c.json({
    object: "list",
    data: invites,
    has_more: false,
    first_id: invites.at(0)?.id,
    last_id: invites.at(-1)?.id,
  });
});

app.post(
  "/organization/invites",
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

app.get("/organization/invites/:invite_id", async (c) => {
  const invite_id = c.req.param("invite_id");

  const invite = await db.query.invites.findFirst({
    where: eq(schema.invites.id, invite_id),
  });
  if (!invite) {
    return c.json({ error: "Invite not found" }, 404);
  }

  return c.json(invite);
});

app.delete("/organization/invites/:invite_id", async (c) => {
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

app.get(
  "/organization/projects",
  zValidator(
    "query",
    z.object({ include_archived: z.coerce.boolean().default(false) })
  ),
  async (c) => {
    const include_archived = c.req.valid("query").include_archived;
    const projects = await db.query.projects.findMany({
      where: or(
        eq(schema.projects.status, "active"),
        include_archived ? eq(schema.projects.status, "archived") : undefined
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
);

app.post(
  "/organization/projects",
  zValidator("json", z.object({ name: z.string() })),
  async (c) => {
    const { name } = c.req.valid("json");

    const [project] = await db
      .insert(schema.projects)
      .values({
        name,
      })
      .returning();
    if (!project) {
      return c.json({ error: "Failed to create project" }, 500);
    }

    await insertDefaultProjectRateLimits({ projectId: project.id });

    return c.json(project);
  }
);

app.get("/organization/projects/:project_id", async (c) => {
  const project_id = c.req.param("project_id");

  const project = await db.query.projects.findFirst({
    where: eq(schema.projects.id, project_id),
  });
  if (!project) {
    return c.json({ error: "Project not found" }, 404);
  }

  return c.json(project);
});

app.post(
  "/organization/projects/:project_id",
  zValidator("json", z.object({ name: z.string() })),
  async (c) => {
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
  }
);

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

app.get("/projects/:project_id/roles", async (c) => {
  const project_id = c.req.param("project_id");

  const roles = await db.query.projectsToRoles.findMany({
    where: eq(schema.projectsToRoles.project_id, project_id),
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
});

app.post(
  "/projects/:project_id/roles",
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
    const { permissions, role_name: name, description } = c.req.valid("json");

    const [role] = await db
      .insert(schema.roles)
      .values({
        name,
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

app.post(
  "/projects/:project_id/roles/:role_id",
  zValidator(
    "json",
    z.object({
      permissions: z.array(z.string()),
      role_name: z.string(),
      description: z.string().default(""),
    })
  ),
  async (c) => {
    const project_id = c.req.param("project_id");
    const role_id = c.req.param("role_id");
    const { permissions, role_name: name, description } = c.req.valid("json");

    const projectToRole = await db.query.projectsToRoles.findFirst({
      where: and(
        eq(schema.projectsToRoles.project_id, project_id),
        eq(schema.projectsToRoles.role_id, role_id)
      ),
    });
    if (!projectToRole) {
      return c.json({ error: "Role not found" }, 404);
    }

    const [role] = await db
      .update(schema.roles)
      .set({
        name,
        description,
        permissions,
      })
      .where(eq(schema.roles.id, projectToRole.role_id))
      .returning();

    return c.json(role);
  }
);

app.delete("/projects/:project_id/roles/:role_id", async (c) => {
  const project_id = c.req.param("project_id");
  const role_id = c.req.param("role_id");

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
    const project_id = c.req.param("project_id");
    const { user_id, role } = c.req.valid("json");

    const project = await db.query.projects.findFirst({
      where: eq(schema.projects.id, project_id),
    });
    if (!project) {
      return c.json({ error: "Project not found" }, 404);
    }

    const user = await db.query.users.findFirst({
      where: eq(schema.users.id, user_id),
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

app.get("/organization/projects/:project_id/users/:user_id", async (c) => {
  const project_id = c.req.param("project_id");
  const user_id = c.req.param("user_id");

  const projectToUser = await db.query.projectsToUsers.findFirst({
    where: and(
      eq(schema.projectsToUsers.project_id, project_id),
      eq(schema.projectsToUsers.user_id, user_id)
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
});

app.post(
  "/organization/projects/:project_id/users/:user_id",
  zValidator("json", z.object({ role: z.enum(["owner", "member"]) })),
  async (c) => {
    const project_id = c.req.param("project_id");
    const user_id = c.req.param("user_id");
    const { role } = c.req.valid("json");

    const projectToUser = await db.query.projectsToUsers.findFirst({
      where: and(
        eq(schema.projectsToUsers.project_id, project_id),
        eq(schema.projectsToUsers.user_id, user_id)
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
);

app.delete("/organization/projects/:project_id/users/:user_id", async (c) => {
  const project_id = c.req.param("project_id");
  const user_id = c.req.param("user_id");

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

app.post(
  "/organization/projects/:project_id/service_accounts",
  zValidator("json", z.object({ name: z.string() })),
  async (c) => {
    const project_id = c.req.param("project_id");
    const { name } = c.req.valid("json");

    const project = await db.query.projects.findFirst({
      where: eq(schema.projects.id, project_id),
    });
    if (!project) {
      return c.json({ error: "Project not found" }, 404);
    }

    const [service_account] = await db
      .insert(schema.projectServiceAccounts)
      .values({
        project_id,
        name,
      })
      .returning();
    if (!service_account) {
      return c.json({ error: "Failed to create service account" }, 500);
    }

    const [api_key] = await db
      .insert(schema.projectServiceAccountApiKeys)
      .values({
        project_service_account_id: service_account.id,
      })
      .returning();
    if (!api_key) {
      return c.json({ error: "Failed to create service account api key" }, 500);
    }

    return c.json({
      ...service_account,
      api_key,
    });
  }
);

app.get(
  "/organization/projects/:project_id/service_accounts/:service_account_id",
  async (c) => {
    const project_id = c.req.param("project_id");
    const service_account_id = c.req.param("service_account_id");

    const service_account = await db.query.projectServiceAccounts.findFirst({
      where: and(
        eq(schema.projectServiceAccounts.project_id, project_id),
        eq(schema.projectServiceAccounts.id, service_account_id)
      ),
    });
    if (!service_account) {
      return c.json({ error: "Service account not found" }, 404);
    }

    return c.json(service_account);
  }
);

app.delete(
  "/organization/projects/:project_id/service_accounts/:service_account_id",
  async (c) => {
    const project_id = c.req.param("project_id");
    const service_account_id = c.req.param("service_account_id");

    const result = await db
      .delete(schema.projectServiceAccounts)
      .where(
        and(
          eq(schema.projectServiceAccounts.project_id, project_id),
          eq(schema.projectServiceAccounts.id, service_account_id)
        )
      )
      .returning();
    if (!result[0]) {
      return c.json({ error: "Service account not found" }, 404);
    }

    return c.json({
      object: "organization.project.service_account.deleted",
      id: result[0].id,
      deleted: true,
    });
  }
);

app.get("/organization/projects/:project_id/rate_limits", async (c) => {
  const project_id = c.req.param("project_id");

  const rate_limits = await db.query.projectRateLimits.findMany({
    where: eq(schema.projectRateLimits.project_id, project_id),
  });

  return c.json({
    object: "list",
    data: rate_limits,
    has_more: false,
    first_id: rate_limits.at(0)?.id,
    last_id: rate_limits.at(-1)?.id,
  });
});

app.post(
  "/organization/projects/:project_id/rate_limits/:rate_limit_id",
  zValidator(
    "json",
    z.object({
      batch_1_day_max_input_tokens: z.number().optional(),
      max_audio_megabytes_per_1_minute: z.number().optional(),
      max_images_per_1_minute: z.number().optional(),
      max_requests_per_1_day: z.number().optional(),
      max_requests_per_1_minute: z.number().optional(),
      max_tokens_per_1_minute: z.number().optional(),
    })
  ),
  async (c) => {
    const project_id = c.req.param("project_id");
    const rate_limit_id = c.req.param("rate_limit_id");

    const [updatedRateLimit] = await db
      .update(schema.projectRateLimits)
      .set(c.req.valid("json"))
      .where(
        and(
          eq(schema.projectRateLimits.project_id, project_id),
          eq(schema.projectRateLimits.id, rate_limit_id)
        )
      )
      .returning();
    if (!updatedRateLimit) {
      return c.json({ error: "Rate limit not found" }, 404);
    }

    return c.json(updatedRateLimit);
  }
);

app.get("/projects/:project_id/groups/:group_id/roles", async (c) => {
  const project_id = c.req.param("project_id");
  const group_id = c.req.param("group_id");

  const roles = await db.query.projectsToGroupsToRoles.findMany({
    where: and(
      eq(schema.projectsToGroupsToRoles.project_id, project_id),
      eq(schema.projectsToGroupsToRoles.group_id, group_id)
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
});

app.post(
  "/projects/:project_id/groups/:group_id/roles",
  zValidator("json", z.object({ role_id: z.string() })),
  async (c) => {
    const project_id = c.req.param("project_id");
    const group_id = c.req.param("group_id");
    const { role_id } = c.req.valid("json");

    const project = await db.query.projects.findFirst({
      where: eq(schema.projects.id, project_id),
    });
    if (!project) {
      return c.json({ error: "Project not found" }, 404);
    }

    const group = await db.query.groups.findFirst({
      where: eq(schema.groups.id, group_id),
    });
    if (!group) {
      return c.json({ error: "Group not found" }, 404);
    }

    const role = await db.query.roles.findFirst({
      where: eq(schema.roles.id, role_id),
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
  }
);

app.delete(
  "/projects/:project_id/groups/:group_id/roles/:role_id",
  async (c) => {
    const project_id = c.req.param("project_id");
    const group_id = c.req.param("group_id");
    const role_id = c.req.param("role_id");

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

app.get("/projects/:project_id/users/:user_id/roles", async (c) => {
  const project_id = c.req.param("project_id");
  const user_id = c.req.param("user_id");

  const roles = await db.query.projectsToUsersToRoles.findMany({
    where: and(
      eq(schema.projectsToUsersToRoles.project_id, project_id),
      eq(schema.projectsToUsersToRoles.user_id, user_id)
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
});

app.post(
  "/projects/:project_id/users/:user_id/roles",
  zValidator("json", z.object({ role_id: z.string() })),
  async (c) => {
    const project_id = c.req.param("project_id");
    const user_id = c.req.param("user_id");
    const { role_id } = c.req.valid("json");

    const project = await db.query.projects.findFirst({
      where: eq(schema.projects.id, project_id),
    });
    if (!project) {
      return c.json({ error: "Project not found" }, 404);
    }

    const user = await db.query.users.findFirst({
      where: eq(schema.users.id, user_id),
    });
    if (!user) {
      return c.json({ error: "User not found" }, 404);
    }

    const role = await db.query.roles.findFirst({
      where: eq(schema.roles.id, role_id),
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
  }
);

app.delete("/projects/:project_id/users/:user_id/roles/:role_id", async (c) => {
  const project_id = c.req.param("project_id");
  const user_id = c.req.param("user_id");
  const role_id = c.req.param("role_id");

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
