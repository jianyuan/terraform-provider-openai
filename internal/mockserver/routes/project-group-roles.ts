import { zValidator } from "@hono/zod-validator";
import { and, eq } from "drizzle-orm";
import { Hono } from "hono";
import z from "zod";
import { db } from "../db";
import * as schema from "../db-schema";
import { requireGroup, type GroupEnv } from "../middleware/group";
import { requireProject, type ProjectEnv } from "../middleware/project";

const route = new Hono<{
  Variables: ProjectEnv["Variables"] & GroupEnv["Variables"];
}>();
route.use(requireProject);
route.use(requireGroup);

route.get("/", async (c) => {
  const project = c.get("project");
  const group = c.get("group");

  const roles = await db.query.projectsToGroupsToRoles.findMany({
    where: and(
      eq(schema.projectsToGroupsToRoles.project_id, project.id),
      eq(schema.projectsToGroupsToRoles.group_id, group.id),
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
    const group = c.get("group");
    const { role_id } = c.req.valid("json");

    const role = await db.query.roles.findFirst({
      where: eq(schema.roles.id, role_id),
    });
    if (!role) {
      return c.json({ error: "Role not found" }, 404);
    }

    await db.insert(schema.projectsToGroupsToRoles).values({
      project_id: project.id,
      group_id: group.id,
      role_id: role.id,
    });

    return c.json({
      object: "group.role",
      group,
      role,
    });
  },
);

route.get("/:role_id", async (c) => {
  const project = c.get("project");
  const group = c.get("group");
  const role_id = c.req.param("role_id");

  const groupToRole = await db.query.projectsToGroupsToRoles.findFirst({
    where: and(
      eq(schema.projectsToGroupsToRoles.project_id, project.id),
      eq(schema.projectsToGroupsToRoles.group_id, group.id),
      eq(schema.projectsToGroupsToRoles.role_id, role_id),
    ),
    with: {
      role: true,
    },
  });
  if (!groupToRole) {
    return c.json({ error: "Group to role not found" }, 404);
  }

  return c.json(groupToRole.role);
});

route.delete("/:role_id", async (c) => {
  const project = c.get("project");
  const group = c.get("group");
  const role_id = c.req.param("role_id");

  const result = await db
    .delete(schema.projectsToGroupsToRoles)
    .where(
      and(
        eq(schema.projectsToGroupsToRoles.project_id, project.id),
        eq(schema.projectsToGroupsToRoles.group_id, group.id),
        eq(schema.projectsToGroupsToRoles.role_id, role_id),
      ),
    )
    .returning();
  if (!result[0]) {
    return c.json({ error: "Role not found" }, 404);
  }

  return c.json({
    object: "group.role.deleted",
    deleted: true,
  });
});

export default route;
