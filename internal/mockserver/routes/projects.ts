import { zValidator } from "@hono/zod-validator";
import { eq, or } from "drizzle-orm";
import { Hono } from "hono";
import z from "zod";
import { db, insertDefaultProjectRateLimits } from "../db";
import * as schema from "../db-schema";
import { now } from "../db-utils";
import { requireProject } from "../middleware/project";

async function externalKeyExists(id: string) {
  const key = await db.query.externalKeys.findFirst({
    where: eq(schema.externalKeys.id, id),
  });
  return !!key;
}

const route = new Hono();

route.get(
  "/",
  zValidator(
    "query",
    z.object({ include_archived: z.coerce.boolean().default(false) }),
  ),
  async (c) => {
    const include_archived = c.req.valid("query").include_archived;
    const projects = await db.query.projects.findMany({
      where: or(
        eq(schema.projects.status, "active"),
        include_archived ? eq(schema.projects.status, "archived") : undefined,
      ),
    });

    return c.json({
      object: "list",
      data: projects,
      has_more: false,
      first_id: projects.at(0)?.id,
      last_id: projects.at(-1)?.id,
    });
  },
);

route.post(
  "/",
  zValidator(
    "json",
    z.object({
      name: z.string(),
      geography: z
        .enum(["US", "EU", "JP", "IN", "KR", "CA", "AU", "SG"])
        .optional(),
      external_key_id: z.string().optional(),
    }),
  ),
  async (c) => {
    const { name, geography, external_key_id } = c.req.valid("json");

    if (external_key_id && !(await externalKeyExists(external_key_id))) {
      return c.json({ error: { message: "External key not found" } }, 400);
    }

    const [project] = await db
      .insert(schema.projects)
      .values({
        name,
        geography,
        external_key_id,
      })
      .returning();
    if (!project) {
      return c.json({ error: "Failed to create project" }, 500);
    }

    await insertDefaultProjectRateLimits({ projectId: project.id });

    return c.json(project);
  },
);

route.get("/:project_id", requireProject, async (c) => {
  const project = c.get("project");

  return c.json(project);
});

route.post(
  "/:project_id",
  requireProject,
  zValidator(
    "json",
    z.object({ name: z.string(), external_key_id: z.string().optional() }),
  ),
  async (c) => {
    const project = c.get("project");
    const { name, external_key_id } = c.req.valid("json");

    if (external_key_id && !(await externalKeyExists(external_key_id))) {
      return c.json({ error: { message: "External key not found" } }, 400);
    }

    const result = await db
      .update(schema.projects)
      .set({ name, external_key_id })
      .where(eq(schema.projects.id, project.id))
      .returning();

    return c.json(result[0]);
  },
);

route.post("/:project_id/archive", requireProject, async (c) => {
  const project = c.get("project");

  const result = await db
    .update(schema.projects)
    .set({ status: "archived", archived_at: now() })
    .where(eq(schema.projects.id, project.id))
    .returning();

  return c.json(result[0]);
});

export default route;
