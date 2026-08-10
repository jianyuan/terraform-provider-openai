import { zValidator } from "@hono/zod-validator";
import { and, eq } from "drizzle-orm";
import { Hono } from "hono";
import z from "zod";
import { db } from "../db";
import * as schema from "../db-schema";
import { type GroupEnv, requireGroup } from "../middleware/group";

const route = new Hono<GroupEnv>();
route.use(requireGroup);

route.get("/", async (c) => {
  const group = c.get("group");

  const users = await db.query.groupsToUsers.findMany({
    where: eq(schema.groupsToUsers.group_id, group.id),
    with: {
      user: {
        columns: {
          id: true,
          name: true,
          email: true,
        },
      },
    },
  });

  return c.json({
    object: "list",
    data: users.map((groupToUser) => groupToUser.user),
    has_more: false,
    next: null,
  });
});

route.post(
  "/",
  zValidator("json", z.object({ user_id: z.string() })),
  async (c) => {
    const group = c.get("group");
    const { user_id } = c.req.valid("json");

    const user = await db.query.users.findFirst({
      where: eq(schema.users.id, user_id),
    });
    if (!user) {
      return c.json({ error: "User not found" }, 404);
    }

    await db.insert(schema.groupsToUsers).values({
      group_id: group.id,
      user_id: user.id,
    });

    return c.json({
      object: "group.user",
      user_id: user.id,
      group_id: group.id,
    });
  },
);

route.get("/:user_id", async (c) => {
  const group = c.get("group");
  const user_id = c.req.param("user_id");

  const groupToUser = await db.query.groupsToUsers.findFirst({
    where: and(
      eq(schema.groupsToUsers.group_id, group.id),
      eq(schema.groupsToUsers.user_id, user_id),
    ),
    with: {
      user: true,
    },
  });
  if (!groupToUser) {
    return c.json({ error: "Group to user not found" }, 404);
  }

  return c.json(groupToUser.user);
});

route.delete("/:user_id", async (c) => {
  const group = c.get("group");
  const user_id = c.req.param("user_id");

  const groupToUser = await db.query.groupsToUsers.findFirst({
    where: and(
      eq(schema.groupsToUsers.group_id, group.id),
      eq(schema.groupsToUsers.user_id, user_id),
    ),
  });
  if (!groupToUser) {
    return c.json({ error: "Group to user not found" }, 404);
  }

  const result = await db
    .delete(schema.groupsToUsers)
    .where(
      and(
        eq(schema.groupsToUsers.group_id, group.id),
        eq(schema.groupsToUsers.user_id, user_id),
      ),
    )
    .returning();
  if (!result[0]) {
    return c.json({ error: "Group to user not found" }, 404);
  }

  return c.json({
    object: "group.user.deleted",
    deleted: true,
  });
});

export default route;
