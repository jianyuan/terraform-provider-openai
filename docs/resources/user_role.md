---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "openai_user_role Resource - terraform-provider-openai"
subcategory: ""
description: |-
  Modifies a user's role in the organization.
---

# openai_user_role (Resource)

Modifies a user's role in the organization.

## Example Usage

```terraform
resource "openai_user_role" "example" {
  user_id = "user-000000000000000000000000"
  role    = "owner"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `role` (String) `owner` or `reader`.
- `user_id` (String) The ID of the user.