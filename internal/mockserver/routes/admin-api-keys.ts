import { zValidator } from "@hono/zod-validator";
import { eq } from "drizzle-orm";
import { Hono } from "hono";
import z from "zod";
import { db } from "../db";
import * as schema from "../db-schema";

const route = new Hono();

route.post(
  "/",
  zValidator("json", z.object({ name: z.string() })),
  async (c) => {
    const { name } = c.req.valid("json");

    const [apiKey] = await db
      .insert(schema.adminApiKeys)
      .values({
        name,
      })
      .returning();
    if (!apiKey) {
      return c.json({ error: "Failed to create admin API key" }, 500);
    }

    return c.json({
      ...apiKey,
      redacted_value: `sk-admin-***${apiKey.value.slice(-3)}`,
    });
  },
);

route.get("/:key_id", async (c) => {
  const key_id = c.req.param("key_id");

  const apiKey = await db.query.adminApiKeys.findFirst({
    where: eq(schema.adminApiKeys.id, key_id),
  });
  if (!apiKey) {
    return c.json({ error: "Admin API key not found" }, 404);
  }

  return c.json({
    ...apiKey,
    redacted_value: `sk-admin-***${apiKey.value.slice(-3)}`,
  });
});

route.delete("/:key_id", async (c) => {
  const key_id = c.req.param("key_id");

  const result = await db
    .delete(schema.adminApiKeys)
    .where(eq(schema.adminApiKeys.id, key_id))
    .returning();
  if (!result[0]) {
    return c.json({ error: "Admin API key not found" }, 404);
  }

  return c.json({
    object: "organization.admin_api_key.deleted",
    id: result[0].id,
    deleted: true,
  });
});

export default route;
