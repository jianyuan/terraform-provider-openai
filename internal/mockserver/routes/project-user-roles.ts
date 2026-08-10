import { zValidator } from "@hono/zod-validator";
import { and, eq } from "drizzle-orm";
import { Hono } from "hono";
import z from "zod";
import { db } from "../db";
import * as schema from "../db-schema";
import { requireProject, type ProjectEnv } from "../middleware/project";

const route = new Hono<
  ProjectEnv,
  {},
  "/projects/:project_id/users/:user_id/roles"
>();
route.use(requireProject);

route.get("/", async (c) => {
  const project = c.get("project");
  const user_id = c.req.param("user_id");

  const roles = await db.query.projectsToUsersToRoles.findMany({
    where: and(
      eq(schema.projectsToUsersToRoles.project_id, project.id),
      eq(schema.projectsToUsersToRoles.user_id, user_id),
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

route.post(
  "/",
  zValidator("json", z.object({ role_id: z.string() })),
  async (c) => {
    const project = c.get("project");
    const user_id = c.req.param("user_id");
    const { role_id } = c.req.valid("json");

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
      project_id: project.id,
      user_id,
      role_id,
    });

    return c.json({
      object: "user.role",
      user,
      role,
    });
  },
);

route.get("/:role_id", async (c) => {
  const project = c.get("project");
  const user_id = c.req.param("user_id");
  const role_id = c.req.param("role_id");

  const userToRole = await db.query.projectsToUsersToRoles.findFirst({
    where: and(
      eq(schema.projectsToUsersToRoles.project_id, project.id),
      eq(schema.projectsToUsersToRoles.user_id, user_id),
      eq(schema.projectsToUsersToRoles.role_id, role_id),
    ),
    with: {
      role: true,
    },
  });
  if (!userToRole) {
    return c.json({ error: "User to role not found" }, 404);
  }

  return c.json(userToRole.role);
});

route.delete("/:role_id", async (c) => {
  const project = c.get("project");
  const user_id = c.req.param("user_id");
  const role_id = c.req.param("role_id");

  const result = await db
    .delete(schema.projectsToUsersToRoles)
    .where(
      and(
        eq(schema.projectsToUsersToRoles.project_id, project.id),
        eq(schema.projectsToUsersToRoles.user_id, user_id),
        eq(schema.projectsToUsersToRoles.role_id, role_id),
      ),
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

export default route;
