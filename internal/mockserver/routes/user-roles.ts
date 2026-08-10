import { zValidator } from "@hono/zod-validator";
import { and, eq } from "drizzle-orm";
import { Hono } from "hono";
import z from "zod";
import { db } from "../db";
import * as schema from "../db-schema";

const route = new Hono<{}, {}, "/organization/users/:user_id/roles">();

route.get("/", async (c) => {
  const user_id = c.req.param("user_id");

  const roles = await db.query.usersToRoles.findMany({
    where: eq(schema.usersToRoles.user_id, user_id),
    with: {
      role: true,
    },
  });

  return c.json({
    object: "list",
    data: roles.map((userToRole) => userToRole.role),
    has_more: false,
    next: null,
  });
});

route.post(
  "/",
  zValidator("json", z.object({ role_id: z.string() })),
  async (c) => {
    const user_id = c.req.param("user_id")!;
    const { role_id } = c.req.valid("json");

    const user = await db.query.users.findFirst({
      where: eq(schema.users.id, user_id),
    });
    if (!user) {
      return c.json({ error: "User not found" }, 404);
    }

    const role = await db.query.roles.findFirst({
      where: and(
        eq(schema.roles.id, role_id),
        eq(schema.roles.resource_type, "api.organization"),
      ),
    });
    if (!role) {
      return c.json({ error: "Role not found" }, 404);
    }

    await db.insert(schema.usersToRoles).values({
      user_id,
      role_id,
    });

    return c.json({
      object: "user.role",
      user,
      role,
    });
  },
);

route.get("/:role_id", async (c) => {
  const user_id = c.req.param("user_id");
  const role_id = c.req.param("role_id");

  const userToRole = await db.query.usersToRoles.findFirst({
    where: and(
      eq(schema.usersToRoles.user_id, user_id),
      eq(schema.usersToRoles.role_id, role_id),
    ),
    with: {
      role: true,
    },
  });
  if (!userToRole) {
    return c.json({ error: "User to role not found" }, 404);
  }

  return c.json(userToRole.role);
});

route.delete("/:role_id", async (c) => {
  const user_id = c.req.param("user_id");
  const role_id = c.req.param("role_id");

  const result = await db
    .delete(schema.usersToRoles)
    .where(
      and(
        eq(schema.usersToRoles.user_id, user_id),
        eq(schema.usersToRoles.role_id, role_id),
      ),
    )
    .returning();
  if (!result[0]) {
    return c.json({ error: "User to role not found" }, 404);
  }

  return c.json({
    object: "user.role.deleted",
    deleted: true,
  });
});

export default route;
