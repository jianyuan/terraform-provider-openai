import { zValidator } from "@hono/zod-validator";
import { eq } from "drizzle-orm";
import { Hono } from "hono";
import z from "zod";
import { db } from "../db";
import * as schema from "../db-schema";

const route = new Hono();

route.get("/", async (c) => {
  const invites = await db.query.invites.findMany();

  return c.json({
    object: "list",
    data: invites,
    has_more: false,
    first_id: invites.at(0)?.id,
    last_id: invites.at(-1)?.id,
  });
});

route.post(
  "/",
  zValidator(
    "json",
    z.object({ email: z.string(), role: z.enum(["owner", "reader"]) }),
  ),
  async (c) => {
    const { email, role } = c.req.valid("json");

    const [invite] = await db
      .insert(schema.invites)
      .values({
        email,
        role,
      })
      .returning();

    return c.json(invite);
  },
);

route.get("/:invite_id", async (c) => {
  const invite_id = c.req.param("invite_id");

  const invite = await db.query.invites.findFirst({
    where: eq(schema.invites.id, invite_id),
  });
  if (!invite) {
    return c.json({ error: "Invite not found" }, 404);
  }

  return c.json(invite);
});

route.delete("/:invite_id", async (c) => {
  const invite_id = c.req.param("invite_id");

  const result = await db
    .delete(schema.invites)
    .where(eq(schema.invites.id, invite_id))
    .returning();
  if (!result[0]) {
    return c.json({ error: "Invite not found" }, 404);
  }

  return c.json({
    object: "organization.invite.deleted",
    id: result[0].id,
    deleted: true,
  });
});

export default route;
