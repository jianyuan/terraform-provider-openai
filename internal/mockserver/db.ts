import { Database } from "bun:sqlite";
import { pushSQLiteSchema } from "drizzle-kit/api";
import { drizzle } from "drizzle-orm/bun-sqlite";
import * as schema from "./db-schema";
import { now } from "./db-utils";

const sqlite = new Database(":memory:");
export const db = drizzle({ client: sqlite, schema });

// Apply migrations
async function applyMigrations() {
  const { apply } = await pushSQLiteSchema(schema, db as any);
  await apply();
}

// Seed
async function seedDatabase() {
  await db.insert(schema.adminApiKeys).values({
    name: "test",
    value: "sk-admin-test",
  });

  await db.insert(schema.users).values({
    id: "user_test",
    name: "John Doe",
    email: "john.doe@example.com",
    role: "owner",
  });

  await db.insert(schema.roles).values({
    id: "role_organization_owner",
    name: "owner",
    description:
      "Can modify billing information and manage organization members",
    permissions: ["*"],
    predefined_role: true,
    resource_type: "api.organization",
  });

  await db.insert(schema.roles).values({
    id: "role_organization_reader",
    name: "reader",
    description:
      "Can make standard API requests and read basic organizational data",
    permissions: ["*"],
    predefined_role: true,
    resource_type: "api.organization",
  });

  await db.insert(schema.roles).values({
    id: "role_project_owner",
    name: "owner",
    description: "",
    permissions: ["*"],
    predefined_role: true,
    resource_type: "api.project",
  });

  await db.insert(schema.roles).values({
    id: "role_project_member",
    name: "member",
    description: "",
    permissions: ["*"],
    predefined_role: true,
    resource_type: "api.project",
  });

  const [project] = await db
    .insert(schema.projects)
    .values({
      name: "Default project",
    })
    .returning();
  await insertDefaultProjectRateLimits({ projectId: project!.id });

  const [archivedProject] = await db
    .insert(schema.projects)
    .values({
      name: "Archived project",
      status: "archived",
      archived_at: now(),
    })
    .returning();
  await insertDefaultProjectRateLimits({ projectId: archivedProject!.id });
}

export async function insertDefaultProjectRateLimits({
  projectId,
}: {
  projectId: string;
}) {
  return await db
    .insert(schema.projectRateLimits)
    .values(
      schema.defaultModels.map((model) => ({
        id: `rl-${model.name}`,
        project_id: projectId,
        model: model.name,
        ...model.rateLimits,
      }))
    )
    .returning();
}

await applyMigrations();
await seedDatabase();
