import { zValidator } from "@hono/zod-validator";
import { and, eq, inArray } from "drizzle-orm";
import { Hono } from "hono";
import z from "zod";
import { db } from "../db";
import * as schema from "../db-schema";

const route = new Hono();

route.get("/", async (c) => {
  const users = await db.query.users.findMany();

  return c.json({
    object: "list",
    data: users,
    has_more: false,
    first_id: users.at(0)?.id,
    last_id: users.at(-1)?.id,
  });
});

route.get("/:user_id", async (c) => {
  const user_id = c.req.param("user_id");

  const user = await db.query.users.findFirst({
    where: eq(schema.users.id, user_id),
  });
  if (!user) {
    return c.json({ error: "User not found" }, 404);
  }

  return c.json(user);
});

route.post(
  "/:user_id",
  zValidator("json", z.object({ role: z.enum(["owner", "reader"]) })),
  async (c) => {
    const user_id = c.req.param("user_id")!;
    const { role } = c.req.valid("json");

    const [user] = await db
      .update(schema.users)
      .set({
        role,
      })
      .where(eq(schema.users.id, user_id))
      .returning();

    // Built-in role
    await db
      .delete(schema.usersToRoles)
      .where(
        and(
          eq(schema.usersToRoles.user_id, user_id),
          inArray(schema.usersToRoles.role_id, [
            "role_organization_owner",
            "role_organization_reader",
          ]),
        ),
      );
    await db.insert(schema.usersToRoles).values({
      user_id,
      role_id:
        role === "owner"
          ? "role_organization_owner"
          : "role_organization_reader",
    });

    return c.json(user);
  },
);

export default route;
