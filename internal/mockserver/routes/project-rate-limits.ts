import { zValidator } from "@hono/zod-validator";
import { and, eq } from "drizzle-orm/sql";
import { Hono } from "hono";
import z from "zod";
import { db } from "../db";
import * as schema from "../db-schema";
import { requireProject, type ProjectEnv } from "../middleware/project";

const route = new Hono<ProjectEnv>();
route.use(requireProject);

route.get("/", async (c) => {
  const project = c.get("project");

  const rate_limits = await db.query.projectRateLimits.findMany({
    where: eq(schema.projectRateLimits.project_id, project.id),
  });

  return c.json({
    object: "list",
    data: rate_limits,
    has_more: false,
    first_id: rate_limits.at(0)?.id,
    last_id: rate_limits.at(-1)?.id,
  });
});

route.post(
  "/:rate_limit_id",
  zValidator(
    "json",
    z.object({
      batch_1_day_max_input_tokens: z.number().optional(),
      max_audio_megabytes_per_1_minute: z.number().optional(),
      max_images_per_1_minute: z.number().optional(),
      max_requests_per_1_day: z.number().optional(),
      max_requests_per_1_minute: z.number().optional(),
      max_tokens_per_1_minute: z.number().optional(),
    }),
  ),
  async (c) => {
    const project = c.get("project");
    const rate_limit_id = c.req.param("rate_limit_id");

    const [updatedRateLimit] = await db
      .update(schema.projectRateLimits)
      .set(c.req.valid("json"))
      .where(
        and(
          eq(schema.projectRateLimits.project_id, project.id),
          eq(schema.projectRateLimits.id, rate_limit_id),
        ),
      )
      .returning();
    if (!updatedRateLimit) {
      return c.json({ error: "Rate limit not found" }, 404);
    }

    return c.json(updatedRateLimit);
  },
);

export default route;
