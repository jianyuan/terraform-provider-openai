def fixes: {
  "Getorganizationspendlimit": "retrieve-organization-spend-limit",
  "Updateorganizationspendlimit": "update-organization-spend-limit",
  "Deleteorganizationspendlimit": "delete-organization-spend-limit",
  "Getprojectspendlimit": "retrieve-project-spend-limit",
  "Updateprojectspendlimit": "update-project-spend-limit",
  "Deleteprojectspendlimit": "delete-project-spend-limit"
};

# --- fix_any_of ---
def fix_any_of:
  if type == "object"
     and has("anyOf")
     and (.anyOf | type == "array")
     and (.anyOf | length > 0)
     and (.anyOf[-1] | type == "object")
     and (.anyOf[-1].type == "null")
  then
    (.anyOf[0:-1]) as $real_types |
    if ($real_types | length) == 1 then
      ($real_types[0] | fix_any_of) + {"nullable": true}
    else
      {"oneOf": [$real_types[] | fix_any_of], "nullable": true}
    end
  elif type == "object" then
    with_entries(.value |= fix_any_of)
  elif type == "array" then
    map(fix_any_of)
  else
    .
  end;

# --- fix_number_format ---
def fix_number_format:
  if type == "object" and (.type == "integer") then
    . + {"format": "int64"}
  elif type == "object" then
    with_entries(.value |= fix_number_format)
  elif type == "array" then
    map(fix_number_format)
  else
    .
  end;

# --- helper for remove_non_administrative_endpoints ---
def is_non_admin_op:
  (type == "object")
  and has("x-oaiMeta")
  and (."x-oaiMeta" | type == "object")
  and (."x-oaiMeta" | has("group"))
  and (."x-oaiMeta".group != "administration");

def should_remove_path:
  [.[] | select(is_non_admin_op)] | length > 0;

# --- pipeline ---
.openapi = "3.0.0"
| .paths |= with_entries(select(.value | should_remove_path | not))
| del(.webhooks)
| fix_any_of
| fix_number_format
| .paths |= with_entries(
    .value |= with_entries(
      if (.value.operationId != null) and (fixes[.value.operationId] != null) then
        .value.operationId = fixes[.value.operationId]
      else
        .
      end
    )
  )