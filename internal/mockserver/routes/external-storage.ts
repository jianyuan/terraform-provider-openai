import { zValidator } from "@hono/zod-validator";
import { eq } from "drizzle-orm";
import { Hono } from "hono";
import z from "zod";
import { db } from "../db";
import * as schema from "../db-schema";

const requestSchema = z.object({
  project_id: z.string(),
  provider: z.discriminatedUnion("type", [
    z.object({
      type: z.literal("aws"),
      bucket: z.string(),
      role_arn: z.string(),
    }),
    z.object({
      type: z.literal("azure"),
      tenant_id: z.string(),
      subscription_id: z.string(),
      resource_group: z.string(),
      account_name: z.string(),
      container: z.string(),
    }),
    z.object({
      type: z.literal("gcp"),
      bucket: z.string(),
      workload_identity_project_number: z.string(),
      workload_identity_pool_id: z.string(),
      workload_identity_provider_id: z.string(),
    }),
    z.object({
      type: z.literal("oci"),
      tenancy_ocid: z.string(),
      region: z.string(),
      bucket: z.string(),
    }),
  ]),
});

type Provider = z.infer<typeof requestSchema>["provider"];
type ExternalStorage = typeof schema.externalStorage.$inferSelect;

// Adds the fields OpenAI derives from the configuration.
function withDerivedFields(projectId: string, provider: Provider) {
  switch (provider.type) {
    case "aws":
      return {
        ...provider,
        account_id: provider.role_arn.split(":")[4] ?? "",
        external_id: "openai-external-id",
        region: "us-east-1",
      };
    case "azure":
      return { ...provider, region: "eastus" };
    case "gcp":
      return {
        ...provider,
        // OpenAI presents the project ID as the token audience. See "Allowed audiences" in
        // https://developers.openai.com/api/docs/guides/private-safety-processing
        audience: projectId,
        region: "us-central1",
      };
    case "oci":
      return provider;
  }
}

// The bucket or container name selects how validation behaves, so tests can exercise each
// outcome:
// - "invalid" fails validation with an error;
// - "unhealthy" and "degraded" finish validation with that status;
// - "pending" stays pending until the configuration is next read, then is validated;
// - "becomes-degraded" is validated, then reads as degraded afterwards.
// Any other name is validated.
function storageName(config: ExternalStorage) {
  return config.provider.bucket ?? config.provider.container ?? "";
}

async function setStatus(id: string, status: ExternalStorage["status"]) {
  const [config] = await db
    .update(schema.externalStorage)
    .set({ status })
    .where(eq(schema.externalStorage.id, id))
    .returning();
  return config;
}

const route = new Hono();

route.get("/", async (c) => {
  const project_id = c.req.query("project_id");

  const configs = await db.query.externalStorage.findMany({
    where: project_id
      ? eq(schema.externalStorage.project_id, project_id)
      : undefined,
  });

  return c.json({
    object: "list",
    data: configs,
    has_more: false,
    first_id: configs.at(0)?.id,
    last_id: configs.at(-1)?.id,
  });
});

route.post("/", zValidator("json", requestSchema), async (c) => {
  const { project_id, provider } = c.req.valid("json");

  const project = await db.query.projects.findFirst({
    where: eq(schema.projects.id, project_id),
  });
  if (!project) {
    return c.json({ error: { message: "Project not found" } }, 404);
  }

  const [config] = await db
    .insert(schema.externalStorage)
    .values({ project_id, provider: withDerivedFields(project_id, provider) })
    .returning();
  if (!config) {
    return c.json({ error: "Failed to create external storage" }, 500);
  }

  return c.json(config);
});

route.get("/:external_storage_id", async (c) => {
  const id = c.req.param("external_storage_id");

  let config = await db.query.externalStorage.findFirst({
    where: eq(schema.externalStorage.id, id),
  });
  if (!config) {
    return c.json({ error: { message: "External storage not found" } }, 404);
  }

  const name = storageName(config);
  if (config.status === "pending" && name.includes("pending")) {
    config = await setStatus(id, "validated");
  } else if (
    config.status === "validated" &&
    name.includes("becomes-degraded")
  ) {
    config = await setStatus(id, "degraded");
  }

  return c.json(config);
});

route.post("/:external_storage_id/validate", async (c) => {
  const id = c.req.param("external_storage_id");

  const config = await db.query.externalStorage.findFirst({
    where: eq(schema.externalStorage.id, id),
  });
  if (!config) {
    return c.json({ error: { message: "External storage not found" } }, 404);
  }

  const name = storageName(config);
  if (name.includes("invalid")) {
    return c.json(
      {
        error: {
          message: `OpenAI could not write to ${name}: access denied.`,
          type: "invalid_request_error",
          param: null,
          code: "external_storage_validation_failed",
        },
      },
      400,
    );
  }

  const status = name.includes("unhealthy")
    ? "unhealthy"
    : name.includes("becomes-degraded")
      ? "validated"
      : name.includes("degraded")
        ? "degraded"
        : name.includes("pending")
          ? "pending"
          : "validated";

  return c.json(await setStatus(id, status));
});

route.delete("/:external_storage_id", async (c) => {
  const id = c.req.param("external_storage_id");

  const [config] = await db
    .delete(schema.externalStorage)
    .where(eq(schema.externalStorage.id, id))
    .returning();
  if (!config) {
    return c.json({ error: { message: "External storage not found" } }, 404);
  }

  return c.json({
    object: "organization.external_storage.deleted",
    id: config.id,
    deleted: true,
  });
});

export default route;
