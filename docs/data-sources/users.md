---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "openai_users Data Source - terraform-provider-openai"
subcategory: ""
description: |-
  Lists all of the users in the organization.
---

# openai_users (Data Source)

Lists all of the users in the organization.

## Example Usage

```terraform
data "openai_users" "example" {
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `users` (Attributes Set) List of users. (see [below for nested schema](#nestedatt--users))

<a id="nestedatt--users"></a>
### Nested Schema for `users`

Read-Only:

- `added_at` (Number) The Unix timestamp (in seconds) of when the user was added.
- `email` (String) The email address of the user.
- `id` (String) User identifier.
- `name` (String) The name of the user.
- `role` (String) Role `owner` or `reader`.