import { eq } from "drizzle-orm";
import { Hono } from "hono";
import { bearerAuth } from "hono/bearer-auth";
import { logger } from "hono/logger";
import { prettyJSON } from "hono/pretty-json";
import { db } from "./db";
import * as schema from "./db-schema";
import adminApiKeys from "./routes/admin-api-keys";
import dataRetention from "./routes/data-retention";
import spendLimit from "./routes/spend-limit";
import projectSpendLimit from "./routes/project-spend-limit";
import projectSpendAlerts from "./routes/project-spend-alerts";
import spendAlerts from "./routes/spend-alerts";
import roles from "./routes/roles";
import groups from "./routes/groups";
import groupRoles from "./routes/group-roles";
import groupUsers from "./routes/group-users";
import users from "./routes/users";
import userRoles from "./routes/user-roles";
import invites from "./routes/invites";
import projects from "./routes/projects";
import projectRoles from "./routes/project-roles";
import projectUsers from "./routes/project-users";
import projectServiceAccounts from "./routes/project-service-accounts";
import projectRateLimits from "./routes/project-rate-limits";
import projectModelPermissions from "./routes/project-model-permissions";
import projectGroupRoles from "./routes/project-group-roles";
import projectUserRoles from "./routes/project-user-roles";

const app = new Hono();
app.use(logger());
app.use(prettyJSON());
app.use(
  "/*",
  bearerAuth({
    verifyToken: async (token) => {
      const apiKey = await db.query.adminApiKeys.findFirst({
        where: eq(schema.adminApiKeys.value, token),
      });
      return !!apiKey;
    },
  }),
);

app.get("/", (c) => c.text("Hello World"));

app.route("/organization/admin_api_keys", adminApiKeys);
app.route("/organization/data_retention", dataRetention);
app.route("/organization/spend_limit", spendLimit);
app.route("/organization/projects/:project_id/spend_limit", projectSpendLimit);
app.route(
  "/organization/projects/:project_id/spend_alerts",
  projectSpendAlerts,
);
app.route("/organization/spend_alerts", spendAlerts);
app.route("/organization/roles", roles);
app.route("/organization/groups", groups);
app.route("/organization/groups/:group_id/roles", groupRoles);
app.route("/organization/groups/:group_id/users", groupUsers);
app.route("/organization/users", users);
app.route("/organization/users/:user_id/roles", userRoles);
app.route("/organization/invites", invites);
app.route("/organization/projects", projects);
app.route("/projects/:project_id/roles", projectRoles);
app.route("/organization/projects/:project_id/users", projectUsers);
app.route(
  "/organization/projects/:project_id/service_accounts",
  projectServiceAccounts,
);
app.route("/organization/projects/:project_id/rate_limits", projectRateLimits);
app.route(
  "/organization/projects/:project_id/model_permissions",
  projectModelPermissions,
);
app.route("/projects/:project_id/groups/:group_id/roles", projectGroupRoles);
app.route("/projects/:project_id/users/:user_id/roles", projectUserRoles);

export default app;
