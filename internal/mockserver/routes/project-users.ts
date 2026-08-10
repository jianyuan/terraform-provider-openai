import { zValidator } from "@hono/zod-validator";
import { and, eq, inArray } from "drizzle-orm/sql";
import { Hono } from "hono";
import z from "zod";
import { db } from "../db";
import * as schema from "../db-schema";
import { type ProjectEnv, requireProject } from "../middleware/project";

const route = new Hono<ProjectEnv>();
route.use(requireProject);

route.post(
  "/",
  zValidator(
    "json",
    z.object({ user_id: z.string(), role: z.enum(["owner", "member"]) }),
  ),
  async (c) => {
    const project = c.get("project");
    const { user_id, role } = c.req.valid("json");

    const user = await db.query.users.findFirst({
      where: eq(schema.users.id, user_id),
    });
    if (!user) {
      return c.json({ error: "User not found" }, 404);
    }

    const [projectToUser] = await db
      .insert(schema.projectsToUsers)
      .values({
        project_id: project.id,
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
        project_id: project.id,
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
  },
);

route.get("/:user_id", async (c) => {
  const project = c.get("project");
  const user_id = c.req.param("user_id");

  const projectToUser = await db.query.projectsToUsers.findFirst({
    where: and(
      eq(schema.projectsToUsers.project_id, project.id),
      eq(schema.projectsToUsers.user_id, user_id),
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

route.post(
  "/:user_id",
  zValidator("json", z.object({ role: z.enum(["owner", "member"]) })),
  async (c) => {
    const project = c.get("project");
    const user_id = c.req.param("user_id");
    const { role } = c.req.valid("json");

    const projectToUser = await db.query.projectsToUsers.findFirst({
      where: and(
        eq(schema.projectsToUsers.project_id, project.id),
        eq(schema.projectsToUsers.user_id, user_id),
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
          eq(schema.projectsToUsers.project_id, project.id),
          eq(schema.projectsToUsers.user_id, user_id),
        ),
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
          eq(schema.projectsToUsersToRoles.project_id, project.id),
          eq(schema.projectsToUsersToRoles.user_id, user_id),
          inArray(schema.projectsToUsersToRoles.role_id, [
            "role_project_owner",
            "role_project_member",
          ]),
        ),
      );

    await db
      .insert(schema.projectsToUsersToRoles)
      .values({
        project_id: project.id,
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
  },
);

route.delete("/:user_id", async (c) => {
  const project = c.get("project");
  const user_id = c.req.param("user_id");

  const result = await db
    .delete(schema.projectsToUsers)
    .where(
      and(
        eq(schema.projectsToUsers.project_id, project.id),
        eq(schema.projectsToUsers.user_id, user_id),
      ),
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

export default route;
