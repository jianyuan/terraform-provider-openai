import { zValidator } from "@hono/zod-validator";
import { and, eq } from "drizzle-orm";
import { Hono } from "hono";
import z from "zod";
import { db } from "../db";
import * as schema from "../db-schema";

const requestSchema = z.object({
  permissions: z.array(z.string()),
  role_name: z.string(),
  description: z.string().default(""),
});

const route = new Hono();

route.get("/", async (c) => {
  const roles = await db.query.roles.findMany({
    where: eq(schema.roles.resource_type, "api.organization"),
  });

  return c.json({
    object: "list",
    data: roles,
    has_more: false,
    next: null,
  });
});

route.post("/", zValidator("json", requestSchema), async (c) => {
  const { permissions, role_name: name, description } = c.req.valid("json");

  const [role] = await db
    .insert(schema.roles)
    .values({
      name,
      description,
      permissions,
      resource_type: "api.organization",
    })
    .returning();

  return c.json(role);
});

route.get("/:role_id", async (c) => {
  const role_id = c.req.param("role_id");
  const role = await db.query.roles.findFirst({
    where: eq(schema.roles.id, role_id),
  });
  if (!role) {
    return c.json({ error: "Role not found" }, 404);
  }

  return c.json(role);
});

route.post("/:role_id", zValidator("json", requestSchema), async (c) => {
  const role_id = c.req.param("role_id");
  const { permissions, role_name: name, description } = c.req.valid("json");

  const [updatedRole] = await db
    .update(schema.roles)
    .set({
      name,
      description,
      permissions,
    })
    .where(
      and(
        eq(schema.roles.id, role_id),
        eq(schema.roles.resource_type, "api.organization"),
      ),
    )
    .returning();
  if (!updatedRole) {
    return c.json({ error: "Role not found" }, 404);
  }

  return c.json(updatedRole);
});

route.delete("/:role_id", async (c) => {
  const role_id = c.req.param("role_id");

  const result = await db
    .delete(schema.roles)
    .where(
      and(
        eq(schema.roles.id, role_id),
        eq(schema.roles.resource_type, "api.organization"),
      ),
    )
    .returning();
  if (!result[0]) {
    return c.json({ error: "Role not found" }, 404);
  }

  return c.json({
    object: "role.deleted",
    deleted: true,
  });
});

export default route;
