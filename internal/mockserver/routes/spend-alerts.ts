import { zValidator } from "@hono/zod-validator";
import { eq } from "drizzle-orm";
import { Hono } from "hono";
import z from "zod";
import { db } from "../db";
import * as schema from "../db-schema";

const requestSchema = z.object({
  currency: z.enum(["USD"]),
  interval: z.enum(["month"]),
  threshold_amount: z.number(),
  notification_channel: z.object({
    type: z.literal("email"),
    recipients: z.array(z.string()),
    subject_prefix: z.string().optional(),
  }),
});

const route = new Hono();

route.get("/", async (c) => {
  const spendAlerts = await db.query.spendAlerts.findMany();
  return c.json({
    object: "list",
    data: spendAlerts,
    has_more: false,
    first_id: spendAlerts.at(0)?.id,
    last_id: spendAlerts.at(-1)?.id,
  });
});

route.post("/", zValidator("json", requestSchema), async (c) => {
  const { currency, interval, threshold_amount, notification_channel } =
    c.req.valid("json");

  const [spendAlert] = await db
    .insert(schema.spendAlerts)
    .values({
      currency,
      interval,
      threshold_amount,
      notification_channel,
    })
    .returning();

  return c.json(spendAlert);
});

route.get("/:alert_id", async (c) => {
  const alert_id = c.req.param("alert_id");
  const spendAlert = await db.query.spendAlerts.findFirst({
    where: eq(schema.spendAlerts.id, alert_id),
  });
  if (!spendAlert) {
    return c.json({ error: "Spend alert not found" }, 404);
  }

  return c.json(spendAlert);
});

route.post("/:alert_id", zValidator("json", requestSchema), async (c) => {
  const alert_id = c.req.param("alert_id");
  const { currency, interval, threshold_amount, notification_channel } =
    c.req.valid("json");

  const [updatedSpendAlert] = await db
    .update(schema.spendAlerts)
    .set({
      currency,
      interval,
      threshold_amount,
      notification_channel,
    })
    .where(eq(schema.spendAlerts.id, alert_id))
    .returning();
  if (!updatedSpendAlert) {
    return c.json({ error: "Spend alert not found" }, 404);
  }

  return c.json(updatedSpendAlert);
});

route.delete("/:alert_id", async (c) => {
  const alert_id = c.req.param("alert_id");

  const result = await db
    .delete(schema.spendAlerts)
    .where(eq(schema.spendAlerts.id, alert_id))
    .returning();
  if (!result[0]) {
    return c.json({ error: "Spend alert not found" }, 404);
  }

  return c.json({
    object: "spend_alert.deleted",
    id: result[0].id,
    deleted: true,
  });
});

export default route;
