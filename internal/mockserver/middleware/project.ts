import { eq } from "drizzle-orm";
import { createMiddleware } from "hono/factory";
import { db } from "../db";
import * as schema from "../db-schema";

export type ProjectEnv = {
  Variables: {
    project: typeof schema.projects.$inferSelect;
  };
};

export const requireProject = createMiddleware<ProjectEnv>(async (c, next) => {
  const project_id = c.req.param("project_id");
  if (!project_id) {
    return c.json({ error: "Project not found" }, 404);
  }

  const project = await db.query.projects.findFirst({
    where: eq(schema.projects.id, project_id),
  });
  if (!project) {
    return c.json({ error: "Project not found" }, 404);
  }

  c.set("project", project);
  await next();
});
