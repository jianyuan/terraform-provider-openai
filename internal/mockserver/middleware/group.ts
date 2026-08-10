import { eq } from "drizzle-orm";
import { createMiddleware } from "hono/factory";
import { db } from "../db";
import * as schema from "../db-schema";

export type GroupEnv = {
  Variables: {
    group: typeof schema.groups.$inferSelect;
  };
};

export const requireGroup = createMiddleware<GroupEnv>(async (c, next) => {
  const group_id = c.req.param("group_id");
  if (!group_id) {
    return c.json({ error: "Group not found" }, 404);
  }

  const group = await db.query.groups.findFirst({
    where: eq(schema.groups.id, group_id),
  });
  if (!group) {
    return c.json({ error: "Group not found" }, 404);
  }

  c.set("group", group);
  await next();
});
