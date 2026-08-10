import { zValidator } from "@hono/zod-validator";
import { and, eq } from "drizzle-orm";
import { Hono } from "hono";
import z from "zod";
import { db } from "../db";
import * as schema from "../db-schema";
import { type ProjectEnv, requireProject } from "../middleware/project";

const requestSchema = z.object({
  permissions: z.array(z.string()),
  role_name: z.string(),
  description: z.string().default(""),
});

const route = new Hono<ProjectEnv>();
route.use(requireProject);

route.get("/", async (c) => {
  const project = c.get("project");

  const roles = await db.query.projectsToRoles.findMany({
    where: eq(schema.projectsToRoles.project_id, project.id),
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

route.post("/", zValidator("json", requestSchema), async (c) => {
  const project = c.get("project");
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
    project_id: project.id,
    role_id: role!.id,
  });

  return c.json(role);
});

route.get("/:role_id", async (c) => {
  const project = c.get("project");
  const role_id = c.req.param("role_id")!;

  const role = await db.query.projectsToRoles.findFirst({
    where: and(
      eq(schema.projectsToRoles.project_id, project.id),
      eq(schema.projectsToRoles.role_id, role_id),
    ),
    with: {
      role: true,
    },
  });
  if (!role) {
    return c.json({ error: "Role not found" }, 404);
  }

  return c.json(role.role);
});

route.post("/:role_id", zValidator("json", requestSchema), async (c) => {
  const project = c.get("project");
  const role_id = c.req.param("role_id");
  const { permissions, role_name: name, description } = c.req.valid("json");

  const projectToRole = await db.query.projectsToRoles.findFirst({
    where: and(
      eq(schema.projectsToRoles.project_id, project.id),
      eq(schema.projectsToRoles.role_id, role_id),
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
});

route.delete("/:role_id", async (c) => {
  const project = c.get("project");
  const role_id = c.req.param("role_id");

  const result = await db
    .delete(schema.projectsToRoles)
    .where(
      and(
        eq(schema.projectsToRoles.project_id, project.id),
        eq(schema.projectsToRoles.role_id, role_id),
      ),
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

export default route;
