import { zValidator } from "@hono/zod-validator";
import { and, eq } from "drizzle-orm";
import { Hono } from "hono";
import z from "zod";
import { db } from "../db";
import * as schema from "../db-schema";
import { buildConflictUpdateColumns } from "../db-utils";
import { requireProject, type ProjectEnv } from "../middleware/project";

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

const route = new Hono<ProjectEnv>();
route.use(requireProject);

route.get("/", async (c) => {
  const project = c.get("project");

  const spendAlerts = await db.query.projectSpendAlerts.findMany({
    where: eq(schema.projectSpendAlerts.project_id, project.id),
  });

  return c.json({
    object: "list",
    data: spendAlerts,
    has_more: false,
    first_id: spendAlerts.at(0)?.id,
    last_id: spendAlerts.at(-1)?.id,
  });
});

route.post("/", zValidator("json", requestSchema), async (c) => {
  const project = c.get("project");
  const { currency, interval, threshold_amount, notification_channel } =
    c.req.valid("json");

  const [spendAlert] = await db
    .insert(schema.projectSpendAlerts)
    .values({
      project_id: project.id,
      currency,
      interval,
      threshold_amount,
      notification_channel,
    })
    .onConflictDoUpdate({
      target: schema.projectSpendAlerts.project_id,
      set: buildConflictUpdateColumns(schema.projectSpendAlerts, [
        "currency",
        "interval",
        "threshold_amount",
        "notification_channel",
      ]),
    })
    .returning();

  return c.json(spendAlert);
});

route.get("/:alert_id", async (c) => {
  const alert_id = c.req.param("alert_id");

  const spendAlert = await db.query.projectSpendAlerts.findFirst({
    where: eq(schema.projectSpendAlerts.id, alert_id),
  });
  if (!spendAlert) {
    return c.json({ error: "Spend alert not found" }, 404);
  }

  return c.json(spendAlert);
});

route.post("/:alert_id", zValidator("json", requestSchema), async (c) => {
  const project = c.get("project");
  const alert_id = c.req.param("alert_id");
  const { currency, interval, threshold_amount, notification_channel } =
    c.req.valid("json");

  const [spendAlert] = await db
    .update(schema.projectSpendAlerts)
    .set({
      currency,
      interval,
      threshold_amount,
      notification_channel,
    })
    .where(
      and(
        eq(schema.projectSpendAlerts.id, alert_id),
        eq(schema.projectSpendAlerts.project_id, project.id),
      ),
    )
    .returning();
  if (!spendAlert) {
    return c.json({ error: "Spend alert not found" }, 404);
  }

  return c.json(spendAlert);
});

route.delete("/:alert_id", async (c) => {
  const project = c.get("project");
  const alert_id = c.req.param("alert_id");

  const [spendAlert] = await db
    .delete(schema.projectSpendAlerts)
    .where(
      and(
        eq(schema.projectSpendAlerts.id, alert_id),
        eq(schema.projectSpendAlerts.project_id, project.id),
      ),
    )
    .returning();
  if (!spendAlert) {
    return c.json({ error: "Spend alert not found" }, 404);
  }

  return c.json({
    object: "project.spend_alert.deleted",
    deleted: true,
  });
});

export default route;
