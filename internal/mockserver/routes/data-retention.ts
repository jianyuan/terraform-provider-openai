import { zValidator } from "@hono/zod-validator";
import { Hono } from "hono";
import z from "zod";
import { db } from "../db";
import * as schema from "../db-schema";

const route = new Hono();

route.get("/", async (c) => {
  const dataRetention = await db.query.dataRetentions.findFirst();
  if (!dataRetention) {
    return c.json({ error: "Data retention not found" }, 404);
  }

  return c.json(dataRetention);
});

route.post(
  "/",
  zValidator(
    "json",
    z.object({
      retention_type: z.enum([
        "zero_data_retention",
        "modified_abuse_monitoring",
        "enhanced_zero_data_retention",
        "enhanced_modified_abuse_monitoring",
      ]),
    }),
  ),
  async (c) => {
    const { retention_type: type } = c.req.valid("json");

    const [dataRetention] = await db
      .update(schema.dataRetentions)
      .set({ type })
      .returning();

    return c.json(dataRetention);
  },
);

export default route;
