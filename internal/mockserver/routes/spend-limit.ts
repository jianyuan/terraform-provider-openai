import { zValidator } from "@hono/zod-validator";
import { Hono } from "hono";
import z from "zod";
import { db } from "../db";
import * as schema from "../db-schema";
import { buildConflictUpdateColumns } from "../db-utils";

const route = new Hono();

route.get("/", async (c) => {
  const spendLimit = await db.query.spendLimits.findFirst();
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
    const { currency, interval, threshold_amount } = c.req.valid("json");

    const [spendLimit] = await db
      .insert(schema.spendLimits)
      .values({
        currency,
        interval,
        threshold_amount,
        enforcement: {
          status: "enforcing",
        },
      })
      .onConflictDoUpdate({
        target: schema.spendLimits.object,
        set: buildConflictUpdateColumns(schema.spendLimits, [
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
  const [spendLimit] = await db.delete(schema.spendLimits).returning();
  if (!spendLimit) {
    return c.json({ error: "Spend limit not found" }, 404);
  }

  return c.json({
    object: "organization.spend_limit.deleted",
    deleted: true,
  });
});

export default route;
