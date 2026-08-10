import { zValidator } from "@hono/zod-validator";
import { and, eq } from "drizzle-orm/sql";
import { Hono } from "hono";
import z from "zod";
import { db } from "../db";
import * as schema from "../db-schema";
import { type ProjectEnv, requireProject } from "../middleware/project";

const route = new Hono<ProjectEnv>();
route.use(requireProject);

route.post(
  "/",
  zValidator("json", z.object({ name: z.string() })),
  async (c) => {
    const project = c.get("project");
    const { name } = c.req.valid("json");

    const [service_account] = await db
      .insert(schema.projectServiceAccounts)
      .values({
        project_id: project.id,
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
  },
);

route.get("/:service_account_id", async (c) => {
  const project = c.get("project");
  const service_account_id = c.req.param("service_account_id");

  const service_account = await db.query.projectServiceAccounts.findFirst({
    where: and(
      eq(schema.projectServiceAccounts.project_id, project.id),
      eq(schema.projectServiceAccounts.id, service_account_id),
    ),
  });
  if (!service_account) {
    return c.json({ error: "Service account not found" }, 404);
  }

  return c.json(service_account);
});

route.delete("/:service_account_id", async (c) => {
  const project = c.get("project");
  const service_account_id = c.req.param("service_account_id");

  const result = await db
    .delete(schema.projectServiceAccounts)
    .where(
      and(
        eq(schema.projectServiceAccounts.project_id, project.id),
        eq(schema.projectServiceAccounts.id, service_account_id),
      ),
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
});

export default route;
