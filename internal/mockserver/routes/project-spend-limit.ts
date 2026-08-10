import { zValidator } from "@hono/zod-validator";
import { eq } from "drizzle-orm";
import { Hono } from "hono";
import z from "zod";
import { db } from "../db";
import * as schema from "../db-schema";
import { buildConflictUpdateColumns } from "../db-utils";
import { requireProject, type ProjectEnv } from "../middleware/project";

const route = new Hono<ProjectEnv>();
route.use(requireProject);

route.get("/", async (c) => {
  const project = c.get("project");

  const spendLimit = await db.query.projectSpendLimits.findFirst({
    where: eq(schema.projectSpendLimits.project_id, project.id),
  });
  if (!spendLimit) {
    return c.json({ error: "Spend limit not found" }, 404);
  }

  return c.json(spendLimit);
});

route.post(
  "/",
  zValidator(
    "json",
    z.object({
      currency: z.enum(["USD"]),
      interval: z.enum(["month"]),
      threshold_amount: z.number(),
    }),
  ),
  async (c) => {
    const project = c.get("project");
    const { currency, interval, threshold_amount } = c.req.valid("json");

    const [spendLimit] = await db
      .insert(schema.projectSpendLimits)
      .values({
        project_id: project.id,
        currency,
        interval,
        threshold_amount,
        enforcement: {
          status: "enforcing",
        },
      })
      .onConflictDoUpdate({
        target: schema.projectSpendLimits.project_id,
        set: buildConflictUpdateColumns(schema.projectSpendLimits, [
          "currency",
          "interval",
          "threshold_amount",
          "enforcement",
        ]),
      })
      .returning();

    return c.json(spendLimit);
  },
);

route.delete("/", async (c) => {
  const project = c.get("project");

  const [spendLimit] = await db
    .delete(schema.projectSpendLimits)
    .where(eq(schema.projectSpendLimits.project_id, project.id))
    .returning();
  if (!spendLimit) {
    return c.json({ error: "Spend limit not found" }, 404);
  }

  return c.json({
    object: "project.spend_limit.deleted",
    deleted: true,
  });
});

export default route;
