import { zValidator } from "@hono/zod-validator";
import { asc, eq, gt } from "drizzle-orm";
import { Hono } from "hono";
import z from "zod";
import { db } from "../db";
import * as schema from "../db-schema";

// The API only accepts keys registered for Enterprise Key Management.
const purpose = z.literal("ekm");

const requestSchema = z.discriminatedUnion("type", [
  z.object({
    type: z.literal("aws"),
    purpose,
    name: z.string(),
    role_arn: z.string(),
    kms_arn: z.string(),
    external_id: z.string(),
  }),
  z.object({
    type: z.literal("azure"),
    purpose,
    name: z.string(),
    tenant_id: z.string(),
    vault_uri: z.string(),
    key_name: z.string(),
  }),
  z.object({
    type: z.literal("gcp"),
    purpose,
    name: z.string(),
    workload_identity_project_number: z.string(),
    workload_identity_pool_id: z.string(),
    workload_identity_provider_id: z.string(),
    audience: z.string(),
    kms_project_id: z.string(),
    kms_key_ring_name: z.string(),
    kms_key_name: z.string(),
    kms_key_location: z.string(),
  }),
]);

async function serialize(key: typeof schema.externalKeys.$inferSelect) {
  const { config, ...rest } = key;
  const projects = await db.query.projects.findMany({
    where: eq(schema.projects.external_key_id, key.id),
  });
  return {
    ...rest,
    ...config,
    api_project_ids: projects.map((project) => project.id),
  };
}

const route = new Hono();

route.get("/", async (c) => {
  const limit = Number(c.req.query("limit") ?? 20);
  const after = c.req.query("after");

  const keys = await db.query.externalKeys.findMany({
    where: after ? gt(schema.externalKeys.id, after) : undefined,
    orderBy: asc(schema.externalKeys.id),
    limit: limit + 1,
  });
  const data = await Promise.all(keys.slice(0, limit).map(serialize));

  return c.json({
    object: "list",
    data,
    has_more: keys.length > limit,
    first_id: data.at(0)?.id,
    last_id: data.at(-1)?.id,
  });
});

route.post("/", zValidator("json", requestSchema), async (c) => {
  const { type, name, purpose, ...config } = c.req.valid("json");

  const [key] = await db
    .insert(schema.externalKeys)
    .values({ type, name, purpose, config })
    .returning();
  if (!key) {
    return c.json({ error: "Failed to create external key" }, 500);
  }

  return c.json(await serialize(key));
});

route.get("/:key_id", async (c) => {
  const key_id = c.req.param("key_id");

  const key = await db.query.externalKeys.findFirst({
    where: eq(schema.externalKeys.id, key_id),
  });
  if (!key) {
    return c.json({ error: "External key not found" }, 404);
  }

  return c.json(await serialize(key));
});

route.delete("/:key_id", async (c) => {
  const key_id = c.req.param("key_id");

  const [key] = await db
    .delete(schema.externalKeys)
    .where(eq(schema.externalKeys.id, key_id))
    .returning();
  if (!key) {
    return c.json({ error: "External key not found" }, 404);
  }

  return c.json(await serialize(key));
});

export default route;
