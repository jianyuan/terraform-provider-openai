---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "openai_projects Data Source - terraform-provider-openai"
subcategory: ""
description: |-
  List all projects in an organization.
---

# openai_projects (Data Source)

List all projects in an organization.

## Example Usage

```terraform
data "openai_projects" "example" {
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `include_archived` (Boolean) Include archived projects. Default is `false`.
- `limit` (Number) Limit the number of projects to return. Default is to return all projects.

### Read-Only

- `projects` (Attributes Set) List of projects. (see [below for nested schema](#nestedatt--projects))

<a id="nestedatt--projects"></a>
### Nested Schema for `projects`

Read-Only:

- `archived_at` (Number) The Unix timestamp (in seconds) of when the project was archived or `null`.
- `created_at` (Number) The Unix timestamp (in seconds) of when the project was created.
- `id` (String) Project ID.
- `name` (String) The name of the project. This appears in reporting.
- `status` (String) Status `active` or `archived`.
