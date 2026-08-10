import { zValidator } from "@hono/zod-validator";
import { and, eq } from "drizzle-orm";
import { Hono } from "hono";
import z from "zod";
import { db } from "../db";
import * as schema from "../db-schema";
import { type GroupEnv, requireGroup } from "../middleware/group";

const route = new Hono<GroupEnv>();
route.use(requireGroup);

route.get(
  "/",
  zValidator("query", z.object({ limit: z.coerce.number().optional() })),
  async (c) => {
    const group = c.get("group");

    const groupsToRoles = await db.query.groupsToRoles.findMany({
      where: eq(schema.groupsToRoles.group_id, group.id),
      with: {
        role: true,
      },
    });

    return c.json({
      object: "list",
      data: groupsToRoles.map((groupToRole) => ({
        ...groupToRole.role,
      })),
      has_more: false,
      next: null,
    });
  },
);

route.post(
  "/",
  zValidator("json", z.object({ role_id: z.string() })),
  async (c) => {
    const group = c.get("group");
    const { role_id } = c.req.valid("json");

    const role = await db.query.roles.findFirst({
      where: eq(schema.roles.id, role_id),
    });
    if (!role) {
      return c.json({ error: "Role not found" }, 404);
    }

    await db.insert(schema.groupsToRoles).values({
      group_id: group.id,
      role_id: role.id,
    });

    return c.json({
      object: "group.role",
      group,
      role,
    });
  },
);

route.get("/:role_id", async (c) => {
  const group = c.get("group");
  const role_id = c.req.param("role_id");

  const groupToRole = await db.query.groupsToRoles.findFirst({
    where: and(
      eq(schema.groupsToRoles.group_id, group.id),
      eq(schema.groupsToRoles.role_id, role_id),
    ),
    with: {
      role: true,
    },
  });
  if (!groupToRole) {
    return c.json({ error: "Group to role not found" }, 404);
  }

  return c.json(groupToRole.role);
});

route.delete("/:role_id", async (c) => {
  const group = c.get("group");
  const role_id = c.req.param("role_id");

  const groupToRole = await db.query.groupsToRoles.findFirst({
    where: and(
      eq(schema.groupsToRoles.group_id, group.id),
      eq(schema.groupsToRoles.role_id, role_id),
    ),
  });
  if (!groupToRole) {
    return c.json({ error: "Group to role not found" }, 404);
  }

  const result = await db
    .delete(schema.groupsToRoles)
    .where(
      and(
        eq(schema.groupsToRoles.group_id, group.id),
        eq(schema.groupsToRoles.role_id, role_id),
      ),
    )
    .returning();
  if (!result[0]) {
    return c.json({ error: "Group to role not found" }, 404);
  }

  return c.json({
    object: "group.role.deleted",
    deleted: true,
  });
});

export default route;
