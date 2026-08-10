import { zValidator } from "@hono/zod-validator";
import { eq } from "drizzle-orm/sql";
import { Hono } from "hono";
import z from "zod";
import { db } from "../db";
import * as schema from "../db-schema";
import { buildConflictUpdateColumns } from "../db-utils";
import type { ProjectEnv } from "../middleware/project";
import { requireProject } from "../middleware/project";

const route = new Hono<ProjectEnv>();
route.use(requireProject);

route.get("/", async (c) => {
  const project = c.get("project");

  const modelPermissions = await db.query.projectModelPermissions.findFirst({
    where: eq(schema.projectModelPermissions.project_id, project.id),
  });
  if (!modelPermissions) {
    return c.json({ error: "Model permissions not found" }, 404);
  }

  return c.json(modelPermissions);
});

route.post(
  "/",
  zValidator(
    "json",
    z.object({
      mode: z.enum(["allow_list", "deny_list"]),
      model_ids: z.array(z.string()),
    }),
  ),
  async (c) => {
    const project = c.get("project");
    const { mode, model_ids } = c.req.valid("json");

    const [updatedModelPermissions] = await db
      .insert(schema.projectModelPermissions)
      .values({
        project_id: project.id,
        mode,
        model_ids,
      })
      .onConflictDoUpdate({
        target: schema.projectModelPermissions.project_id,
        set: buildConflictUpdateColumns(schema.projectModelPermissions, [
          "mode",
          "model_ids",
        ]),
      })
      .returning();
    if (!updatedModelPermissions) {
      return c.json({ error: "Model permissions not found" }, 404);
    }

    return c.json(updatedModelPermissions);
  },
);

route.delete("/", async (c) => {
  const project = c.get("project");

  const result = await db
    .delete(schema.projectModelPermissions)
    .where(eq(schema.projectModelPermissions.project_id, project.id))
    .returning();
  if (!result[0]) {
    return c.json({ error: "Model permissions not found" }, 404);
  }

  return c.json({
    object: "project.model_permissions.deleted",
    deleted: true,
  });
});

export default route;
