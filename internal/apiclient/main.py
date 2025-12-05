from pathlib import Path

import httpx
import yaml

ROOT = Path(__file__).parent
SPEC_ORIGINAL_PATH = ROOT / "api-original.yaml"
SPEC_OUTPUT_PATH = ROOT / "api.yaml"

OPENAPI_SPEC_URL = (
    "https://app.stainless.com/api/spec/documented/openai/openapi.documented.yml"
)


def fix_openapi_version(spec):
    spec["openapi"] = "3.0.0"

    return spec


def fix_remove_non_administrative_endpoints(spec):
    paths_to_remove = set()

    for path, operations in spec["paths"].items():
        for operation in operations.values():
            if (
                "x-oaiMeta" in operation
                and operation["x-oaiMeta"]["group"] != "administration"
            ):
                paths_to_remove.add(path)

    for path in paths_to_remove:
        del spec["paths"][path]

    return spec


def fix_certificate_id_parameter(spec):
    certificate_id_parameter = {
        "name": "certificate_id",
        "in": "path",
        "description": "Unique ID of the certificate to retrieve.",
        "required": True,
        "schema": {"type": "string"},
    }

    path = spec["paths"]["/organization/certificates/{certificate_id}"]
    path["post"]["parameters"] = [certificate_id_parameter]
    path["delete"]["parameters"] = [certificate_id_parameter]

    return spec


def fix_any_of(spec):
    match spec:
        case {"anyOf": [real_type, {"type": "null"}]}:
            return {
                **real_type,
                "nullable": True,
            }
        case dict():
            return {k: fix_any_of(v) for k, v in spec.items()}
        case list():
            return [fix_any_of(v) for v in spec]
        case _:
            return spec


def fix_number_format(spec):
    match spec:
        case {"type": "integer", **rest}:
            return {
                **rest,
                "type": "integer",
                "format": "int64",
            }
        case dict():
            return {k: fix_number_format(v) for k, v in spec.items()}
        case list():
            return [fix_number_format(v) for v in spec]
        case _:
            return spec


def add_external_key_id_to_project(spec):
    property_spec = {
        "type": "string",
        "description": "The ID of the customer-managed encryption key for Enterprise Key Management (EKM).",
    }
    spec["components"]["schemas"]["Project"]["properties"]["external_key_id"] = (
        property_spec
    )
    spec["components"]["schemas"]["ProjectCreateRequest"]["properties"][
        "external_key_id"
    ] = property_spec

    return spec


def rename_user_role_assignment_object(spec):
    spec["components"]["schemas"]["UserRoleAssignment"]["properties"]["object"][
        "x-enum-varnames"
    ] = ["UserRoleAssignmentObjectUserRole"]

    return spec


fix_funcs = [
    fix_openapi_version,
    fix_remove_non_administrative_endpoints,
    fix_certificate_id_parameter,
    fix_any_of,
    fix_number_format,
    add_external_key_id_to_project,
    rename_user_role_assignment_object,
]


def main() -> None:
    response = httpx.get(OPENAPI_SPEC_URL)
    response.raise_for_status()

    SPEC_ORIGINAL_PATH.write_text(response.text)

    with SPEC_ORIGINAL_PATH.open() as f:
        spec = yaml.safe_load(stream=f)

    for fix_func in fix_funcs:
        spec = fix_func(spec)

    SPEC_OUTPUT_PATH.write_text(yaml.dump(spec))


if __name__ == "__main__":
    main()
