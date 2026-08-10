import { zValidator } from "@hono/zod-validator";
import { eq } from "drizzle-orm";
import { Hono } from "hono";
import z from "zod";
import { db } from "../db";
import * as schema from "../db-schema";

const requestSchema = z.object({ name: z.string() });

const route = new Hono();

route.get(
  "/",
  zValidator("query", z.object({ limit: z.coerce.number().optional() })),
  async (c) => {
    const groups = await db.select().from(schema.groups);

    return c.json({
      object: "list",
      data: groups,
      has_more: false,
      next: null,
    });
  },
);

route.post("/", zValidator("json", requestSchema), async (c) => {
  const { name } = c.req.valid("json");

  const [group] = await db
    .insert(schema.groups)
    .values({
      name,
    })
    .returning();

  return c.json(group);
});

route.get("/:group_id", async (c) => {
  const group_id = c.req.param("group_id");
  const group = await db.query.groups.findFirst({
    where: eq(schema.groups.id, group_id),
  });
  if (!group) {
    return c.json({ error: "Group not found" }, 404);
  }

  return c.json(group);
});

route.post("/:group_id", zValidator("json", requestSchema), async (c) => {
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
});

route.delete("/:group_id", async (c) => {
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

export default route;
