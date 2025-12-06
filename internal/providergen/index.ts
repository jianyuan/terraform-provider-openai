import { camelize } from "inflection";
import { DATASOURCES, RESOURCES } from "./settings";
import type { DataSource, Attribute, Resource } from "./schema";
import { match, P } from "ts-pattern";
import { parseArgs } from "util";

function generateTerraformAttribute({
  parent,
  attribute,
}: {
  parent: string;
  attribute: Attribute;
}) {
  const commonParts: string[] = [];
  commonParts.push(
    `MarkdownDescription: ${JSON.stringify(attribute.description)},`
  );
  commonParts.push(
    match(attribute.computedOptionalRequired)
      .with("required", () => "Required: true,")
      .with("computed", () => "Computed: true,")
      .with("computed_optional", () => "Optional: true,\nComputed: true,")
      .with("optional", () => "Optional: true,")
      .exhaustive()
  );
  if (attribute.sensitive) {
    commonParts.push("Sensitive: true,");
  }

  return match(attribute)
    .with({ type: "string" }, () => {
      const parts: string[] = [];
      parts.push("schema.StringAttribute{");
      parts.push(...commonParts);
      parts.push("CustomType: supertypes.StringType{},");
      if (attribute.validators) {
        parts.push("Validators: []validator.String{");
        parts.push(...attribute.validators.map((validator) => `${validator},`));
        parts.push("},");
      }
      if (attribute.planModifiers) {
        parts.push("PlanModifiers: []planmodifier.String{");
        parts.push(
          ...attribute.planModifiers.map((modifier) => `${modifier},`)
        );
        parts.push("},");
      }
      parts.push("}");
      return parts.join("\n");
    })
    .with({ type: "int" }, (attribute) => {
      const parts: string[] = [];
      parts.push("schema.Int64Attribute{");
      parts.push(...commonParts);
      parts.push("CustomType: supertypes.Int64Type{},");
      if (attribute.validators) {
        parts.push("Validators: []validator.Int64{");
        parts.push(...attribute.validators.map((validator) => `${validator},`));
        parts.push("},");
      }
      if (attribute.planModifiers) {
        parts.push("PlanModifiers: []planmodifier.Int64{");
        parts.push(
          ...attribute.planModifiers.map((modifier) => `${modifier},`)
        );
        parts.push("},");
      }
      parts.push("}");
      return parts.join("\n");
    })
    .with({ type: "bool" }, () => {
      const parts: string[] = [];
      parts.push("schema.BoolAttribute{");
      parts.push(...commonParts);
      parts.push("CustomType: supertypes.BoolType{},");
      if (attribute.validators) {
        parts.push("Validators: []validator.Bool{");
        parts.push(...attribute.validators.map((validator) => `${validator},`));
        parts.push("},");
      }
      if (attribute.planModifiers) {
        parts.push("PlanModifiers: []planmodifier.Bool{");
        parts.push(
          ...attribute.planModifiers.map((modifier) => `${modifier},`)
        );
        parts.push("},");
      }
      parts.push("}");
      return parts.join("\n");
    })
    .with({ type: "set_nested" }, (attribute) => {
      const parts: string[] = [];
      parts.push("schema.SetNestedAttribute{");
      parts.push(...commonParts);
      parts.push(
        `CustomType: supertypes.NewSetNestedObjectTypeOf[${parent}${camelize(
          attribute.name
        )}Item](ctx),`
      );
      parts.push("NestedObject: schema.NestedAttributeObject{");
      parts.push("Attributes: map[string]schema.Attribute{");
      for (const nestedAttribute of attribute.attributes) {
        parts.push(
          `"${nestedAttribute.name}": ${generateTerraformAttribute({
            parent: attribute.name,
            attribute: nestedAttribute,
          })},`
        );
      }
      parts.push("},");
      parts.push("},");
      parts.push("}");
      return parts.join("\n");
    })
    .exhaustive();
}

function generateTerraformValueType({
  parent,
  attribute,
}: {
  parent: string;
  attribute: Attribute;
}) {
  return match(attribute)
    .with({ type: "string" }, () => "supertypes.StringValue")
    .with({ type: "int" }, () => "supertypes.Int64Value")
    .with({ type: "bool" }, () => "supertypes.BoolValue")
    .with(
      { type: "set_nested" },
      () =>
        `supertypes.SetNestedObjectValueOf[${parent}${camelize(
          attribute.name
        )}Item]`
    )
    .exhaustive();
}

function generateTerraformToPrimitive({
  attribute,
  srcVar,
}: {
  attribute: Attribute;
  srcVar: string;
}) {
  const srcVarName = `${srcVar}.${camelize(attribute.name)}`;
  return match(attribute)
    .with({ type: "string" }, () => `${srcVarName}.ValueString()`)
    .with({ type: "int" }, () => `${srcVarName}.ValueInt64()`)
    .with({ type: "bool" }, () => `${srcVarName}.ValueBool()`)
    .exhaustive();
}

function generatePrimitiveToTerraform({
  attribute,
  srcVar,
  destVar,
}: {
  attribute: Attribute;
  srcVar: string;
  destVar: string;
}) {
  const srcVarName = `${srcVar}.${camelize(attribute.name)}`;
  const destVarName = `${destVar}.${camelize(attribute.name)}`;
  return match(attribute)
    .with({ type: "string", nullable: true }, () =>
      `
      if ${srcVarName} != nil {
        ${destVarName} = supertypes.NewStringValue(string(*${srcVarName}))
      } else {
        ${destVarName} = supertypes.NewStringNull()
      }
      `.trim()
    )
    .with(
      { type: "string" },
      () => `${destVarName} = supertypes.NewStringValue(string(${srcVarName}))`
    )
    .with({ type: "int", nullable: true }, () =>
      `
      if ${srcVarName} != nil {
        ${destVarName} = supertypes.NewInt64Value(*${srcVarName})
      } else {
        ${destVarName} = supertypes.NewInt64Null()
      }
      `.trim()
    )
    .with(
      { type: "int" },
      () => `${destVarName} = supertypes.NewInt64Value(${srcVarName})`
    )
    .with(
      { type: "bool" },
      () => `${destVarName} = supertypes.NewBoolValue(${srcVarName})`
    )
    .with({ type: "set_nested" }, () => {
      const srcVarName = `${srcVar}.${camelize(attribute.name)}`;
      const destVarName = `${destVar}.${camelize(attribute.name)}`;
      return `
      // if ${srcVarName} != nil {
      //   ${destVarName} = supertypes.NewSetNestedObjectValueOf[${destVarName}](${srcVarName})
      // } else {
      //   ${destVarName} = supertypes.NewSetNestedObjectValueOf[${destVarName}](nil)
      // }
      // TODO
      `.trim();
    })
    .exhaustive();
}

function generateModel({
  name,
  attributes,
}: {
  name: string;
  attributes: Array<Attribute>;
}) {
  const structLines: string[] = [];
  const fillLines: string[] = [];
  const extras: string[] = [];

  for (const attribute of attributes) {
    structLines.push(
      `${camelize(attribute.name)} ${generateTerraformValueType({
        parent: name,
        attribute,
      })} \`tfsdk:"${attribute.name}"\``
    );
    fillLines.push(
      generatePrimitiveToTerraform({
        attribute,
        srcVar: "data",
        destVar: "m",
      })
    );

    extras.push(
      ...match(attribute)
        .with({ type: "set_nested" }, (attribute) => [
          generateModel({
            name: `${name}${camelize(attribute.name)}Item`,
            attributes: attribute.attributes,
          }),
        ])
        .otherwise(() => [])
    );
  }

  return `
type ${name} struct {
  ${structLines.join("\n")}
}

${extras.join("\n\n")}
`;
}

function generateDataSourceModel({ dataSource }: { dataSource: DataSource }) {
  const modelName = `${camelize(dataSource.name)}DataSourceModel`;
  return generateModel({
    name: modelName,
    attributes: dataSource.attributes,
  });
}

function generateDataSourceSchemaAttributes({
  dataSource,
}: {
  dataSource: DataSource;
}) {
  const lines: string[] = [];

  for (const attribute of dataSource.attributes) {
    lines.push(
      `"${attribute.name}": ${generateTerraformAttribute({
        parent: `${camelize(dataSource.name)}DataSourceModel`,
        attribute,
      })},`
    );
  }

  return lines.join("\n");
}

function generateDataSource({ dataSource }: { dataSource: DataSource }) {
  console.log(`Generating data source - ${dataSource.name}`);

  const dataSourceName = `${camelize(dataSource.name)}DataSource`;
  const modelName = `${camelize(dataSource.name)}DataSourceModel`;

  const readRequestParams = ["ctx"];
  readRequestParams.push(
    ...match(dataSource.api)
      .with({ strategy: "paginate" }, (api) => {
        const parts: string[] = [];
        if (api.readRequestAttributes) {
          parts.push(
            ...api.readRequestAttributes.map((param) => {
              const attribute = dataSource.attributes.find(
                (attribute) => attribute.name === param
              );
              if (!attribute) {
                throw new Error(
                  `Attribute ${param} not found in data source ${dataSource.name}`
                );
              }
              return generateTerraformToPrimitive({
                attribute,
                srcVar: "data",
              });
            })
          );
        }
        parts.push("params");
        return parts;
      })
      .with({ strategy: "simple" }, () => ["data.Id.ValueString()"])
      .exhaustive()
  );

  const read = match(dataSource.api)
    .with(
      { strategy: "paginate" },
      (api) => `
    var modelInstances []apiclient.${dataSource.api.model}
    params := &apiclient.${dataSource.api.readMethod}Params{
      Limit: ptr.Ptr(int64(100)),
    }

    ${api.hooks?.readInitLoop ?? ""}

  done:
    for {
      ${api.hooks?.readPreIterate ?? ""}

      httpResp, err := d.client.${
        dataSource.api.readMethod
      }WithResponse(${readRequestParams.join(",")})
      if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read, got error: %s", err))
        return
      } else if httpResp.StatusCode() != http.StatusOK {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read, got status code: %d", httpResp.StatusCode()))
        return
      } else if httpResp.JSON200 == nil {
        resp.Diagnostics.AddError("Client Error", "Unable to read, got empty response body")
        return
      }

      modelInstances = append(modelInstances, httpResp.JSON200.Data...)

      switch v := any(httpResp.JSON200.HasMore).(type) {
      case bool:
        if !v {
          break done
        }
      case *bool:
        if v == nil || !*v {
          break done
        }
      default:
        panic("unknown type")
      }

      switch v := any(httpResp.JSON200.LastId).(type) {
      case string:
        params.After = &v
      case *string:
        if v == nil {
          params.After = nil
        } else {
          params.After = v
        }
      default:
        panic("unknown type")
      }

      ${api.hooks?.readPostIterate ?? ""}
    }

    resp.Diagnostics.Append(data.Fill(ctx, modelInstances)...)
    if resp.Diagnostics.HasError() {
      return
    }
    `
    )
    .with(
      { strategy: "simple" },
      () => `
    httpResp, err := d.client.${
      dataSource.api.readMethod
    }WithResponse(${readRequestParams.join(",")})
    if err != nil {
      resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read, got error: %s", err))
      return
    } else if httpResp.StatusCode() != http.StatusOK {
      resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read, got status code: %d", httpResp.StatusCode()))
      return
    } else if httpResp.JSON200 == nil {
      resp.Diagnostics.AddError("Client Error", "Unable to read, got empty response body")
      return
    }

    resp.Diagnostics.Append(data.Fill(ctx, *httpResp.JSON200)...)
    if resp.Diagnostics.HasError() {
      return
    }
    `
    )
    .exhaustive();

  return `
// Code generated by providergen. DO NOT EDIT.
package provider

import (
  "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var _ datasource.DataSource = &${dataSourceName}{}

func New${dataSourceName}() datasource.DataSource {
  return &${dataSourceName}{}
}

type ${dataSourceName} struct {
  baseDataSource
}

func (d *${dataSourceName}) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
  resp.TypeName = req.ProviderTypeName + "_${dataSource.name}"
}

func (d *${dataSourceName}) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
  resp.Schema = schema.Schema{
    MarkdownDescription: ${JSON.stringify(dataSource.description)},
    Attributes: map[string]schema.Attribute{
      ${generateDataSourceSchemaAttributes({ dataSource })}
    },
  }
}

func (d *${dataSourceName}) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
  var data ${modelName}

  resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
  if resp.Diagnostics.HasError() {
    return
  }

  ${read}

  resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

${generateDataSourceModel({ dataSource })}
`;
}

function generateResourceModel({ resource }: { resource: Resource }) {
  const modelName = `${camelize(resource.name)}ResourceModel`;
  return generateModel({
    name: modelName,
    attributes: resource.attributes,
  });
}

function generateResourceSchemaAttributes({
  resource,
}: {
  resource: Resource;
}) {
  const lines: string[] = [];

  for (const attribute of resource.attributes) {
    lines.push(
      `"${attribute.name}": ${generateTerraformAttribute({
        parent: `${camelize(resource.name)}ResourceModel`,
        attribute,
      })},`
    );
  }

  return lines.join("\n");
}

function generateResource({ resource }: { resource: Resource }) {
  console.log(`Generating resource - ${resource.name}`);

  const resourceName = `${camelize(resource.name)}Resource`;
  const modelName = `${camelize(resource.name)}ResourceModel`;

  const createRequestParams = ["ctx"];
  if (resource.api.createRequestAttributes) {
    createRequestParams.push(
      ...resource.api.createRequestAttributes.map((param) => {
        const attribute = resource.attributes.find(
          (attribute) => attribute.name === param
        );
        if (!attribute) {
          throw new Error(
            `Attribute ${param} not found in resource ${resource.name}`
          );
        }
        return generateTerraformToPrimitive({
          attribute,
          srcVar: "data",
        });
      })
    );
  }
  createRequestParams.push("r.getCreateJSONRequestBody(data)");

  const readRequestParams = ["ctx"];
  if (resource.api.readRequestAttributes) {
    readRequestParams.push(
      ...resource.api.readRequestAttributes.map((param) => {
        const attribute = resource.attributes.find(
          (attribute) => attribute.name === param
        );
        if (!attribute) {
          throw new Error(
            `Attribute ${param} not found in resource ${resource.name}`
          );
        }
        return generateTerraformToPrimitive({
          attribute,
          srcVar: "data",
        });
      })
    );
  }

  const updateRequestParams = ["ctx"];
  if (resource.api.updateRequestAttributes) {
    updateRequestParams.push(
      ...resource.api.updateRequestAttributes.map((param) => {
        const attribute = resource.attributes.find(
          (attribute) => attribute.name === param
        );
        if (!attribute) {
          throw new Error(
            `Attribute ${param} not found in resource ${resource.name}`
          );
        }
        return generateTerraformToPrimitive({
          attribute,
          srcVar: "data",
        });
      })
    );
  }
  updateRequestParams.push("r.getUpdateJSONRequestBody(data)");

  const deleteRequestParams = ["ctx"];
  if (resource.api.deleteRequestAttributes) {
    deleteRequestParams.push(
      ...resource.api.deleteRequestAttributes.map((param) => {
        const attribute = resource.attributes.find(
          (attribute) => attribute.name === param
        );
        if (!attribute) {
          throw new Error(
            `Attribute ${param} not found in resource ${resource.name}`
          );
        }
        return generateTerraformToPrimitive({
          attribute,
          srcVar: "data",
        });
      })
    );
  }

  const importState = match(resource.importStateAttributes)
    .with([P.any], (attributes) => {
      return `
        func (r *${resourceName}) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
          resource.ImportStatePassthroughID(ctx, path.Root("${attributes[0]}"), req, resp)
        }
      `;
    })
    .with([P.any, P.any], (attributes) => {
      return `
        func (r *${resourceName}) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
          first, second, err := tfutils.SplitTwoPartId(req.ID, "${attributes[0]}", "${attributes[1]}")
          if err != nil {
            resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Error parsing ID: %s", err.Error()))
            return
          }

          resp.Diagnostics.Append(resp.State.SetAttribute(
            ctx, path.Root("${attributes[0]}"), first,
          )...)
          resp.Diagnostics.Append(resp.State.SetAttribute(
            ctx, path.Root("${attributes[1]}"), second,
          )...)
        }
      `;
    })
    .otherwise(() => "");

  return `
// Code generated by providergen. DO NOT EDIT.
package provider

import (
  "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var _ resource.Resource = &${resourceName}{}
${
  resource.importStateAttributes
    ? `var _ resource.ResourceWithImportState = &${resourceName}{}`
    : ""
}

func New${resourceName}() resource.Resource {
  return &${resourceName}{}
}

type ${resourceName} struct {
  baseResource
}

func (r *${resourceName}) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
  resp.TypeName = req.ProviderTypeName + "_${resource.name}"
}

func (r *${resourceName}) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
  resp.Schema = schema.Schema{
    MarkdownDescription: ${JSON.stringify(resource.description)},
    Attributes: map[string]schema.Attribute{
      ${generateResourceSchemaAttributes({ resource })}
    },
  }
}

func (r *${resourceName}) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
  var data ${modelName}

  resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
  if resp.Diagnostics.HasError() {
    return
  }

  httpResp, err := r.client.${
    resource.api.createMethod
  }WithResponse(${createRequestParams.join(",")})
  if err != nil {
    resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create, got error: %s", err))
    return
  } else if httpResp.StatusCode() != http.StatusOK {
    resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create, got status code: %d", httpResp.StatusCode()))
    return
  } else if httpResp.JSON200 == nil {
    resp.Diagnostics.AddError("Client Error", "Unable to create, got empty response body")
    return
  }

  resp.Diagnostics.Append(data.Fill(ctx, *httpResp.JSON200)...)
  if resp.Diagnostics.HasError() {
    return
  }

  resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *${resourceName}) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
  var data ${modelName}

  resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
  if resp.Diagnostics.HasError() {
    return
  }

  httpResp, err := r.client.${
    resource.api.readMethod
  }WithResponse(${readRequestParams.join(",")})
  if err != nil {
    resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read, got error: %s", err))
    return
  } else if httpResp.StatusCode() != http.StatusOK {
    resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read, got status code: %d", httpResp.StatusCode()))
    return
  } else if httpResp.JSON200 == nil {
    resp.Diagnostics.AddError("Client Error", "Unable to read, got empty response body")
    return
  }

  resp.Diagnostics.Append(data.Fill(ctx, *httpResp.JSON200)...)
  if resp.Diagnostics.HasError() {
    return
  }

  resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *${resourceName}) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
  ${
    resource.api.updateMethod
      ? `
      var data ${modelName}

      resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
      if resp.Diagnostics.HasError() {
        return
      }

      httpResp, err := r.client.${
        resource.api.updateMethod
      }WithResponse(${updateRequestParams.join(",")})
      if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update, got error: %s", err))
        return
      } else if httpResp.StatusCode() != http.StatusOK {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update, got status code: %d", httpResp.StatusCode()))
        return
      } else if httpResp.JSON200 == nil {
        resp.Diagnostics.AddError("Client Error", "Unable to update, got empty response body")
        return
      }

      resp.Diagnostics.Append(data.Fill(ctx, *httpResp.JSON200)...)
      if resp.Diagnostics.HasError() {
        return
      }

      resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
      `
      : `
      resp.Diagnostics.AddError("Not Supported", "Update is not supported for this resource")
      `
  }
}

func (r *${resourceName}) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
  ${
    resource.api.deleteMethod
      ? `
      var data ${modelName}

      resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
      if resp.Diagnostics.HasError() {
        return
      }

      httpResp, err := r.client.${
        resource.api.deleteMethod
      }WithResponse(${deleteRequestParams.join(",")})
      if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete, got error: %s", err))
        return
      } else if httpResp.StatusCode() == http.StatusNotFound {
        return
      } else if httpResp.StatusCode() != http.StatusOK {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete, got status code: %d", httpResp.StatusCode()))
        return
      }
      `
      : `
      resp.Diagnostics.AddWarning("Not Supported", "Delete is not supported for this resource. Please manually delete the resource.")
      `
  }
}

${importState}


${generateResourceModel({ resource })}
`;
}

async function writeAndFormatGoFile(destination: URL, code: string) {
  await Bun.write(destination, code);
  await Bun.$`go fmt ${destination.pathname}`;
  await Bun.$`go tool goimports -w ${destination.pathname}`;
}

async function main() {
  const { values } = parseArgs({
    args: Bun.argv,
    options: {
      filter: {
        type: "string",
      },
    },
    strict: true,
    allowPositionals: true,
  });

  console.log("Generating provider...");

  console.log("Generating data sources...");
  for (const dataSource of DATASOURCES) {
    if (values.filter && values.filter !== dataSource.name) {
      continue;
    }

    const code = generateDataSource({ dataSource });
    await writeAndFormatGoFile(
      new URL(`../provider/data_source_${dataSource.name}.go`, import.meta.url),
      code
    );
  }

  console.log("Generating resources...");
  for (const resource of RESOURCES) {
    if (values.filter && values.filter !== resource.name) {
      continue;
    }

    const code = generateResource({ resource });
    await writeAndFormatGoFile(
      new URL(`../provider/resource_${resource.name}.go`, import.meta.url),
      code
    );
  }
}

await main()
  .then(() => {
    console.log("✨ Done");
    process.exit(0);
  })
  .catch((err) => {
    console.error(err);
    process.exit(1);
  });
