import { camelize } from "inflection";
import { DATASOURCES } from "./settings";
import type { DataSource, Attribute } from "./schema";
import { match } from "ts-pattern";

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

  return match(attribute)
    .with({ type: "string" }, () => {
      const parts: string[] = [];
      parts.push("schema.StringAttribute{");
      parts.push(...commonParts);
      parts.push("CustomType: supertypes.StringType{},");
      parts.push("}");
      return parts.join("\n");
    })
    .with({ type: "int" }, () => {
      const parts: string[] = [];
      parts.push("schema.Int64Attribute{");
      parts.push(...commonParts);
      parts.push("CustomType: supertypes.Int64Type{},");
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

  const params = ["ctx"];
  params.push(
    ...match(dataSource.api)
      .with({ strategy: "paginate" }, (api) => {
        const parts: string[] = [];
        if (api.params) {
          parts.push(
            ...api.params.map((param) => {
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

  const read = match(dataSource.api.strategy)
    .with(
      "paginate",
      () => `
    var modelInstances []apiclient.${dataSource.api.model}
    params := &apiclient.${dataSource.api.method}Params{
      Limit: ptr.Ptr(int64(100)),
    }

  done:
    for {
      httpResp, err := d.client.${
        dataSource.api.method
      }WithResponse(${params.join(",")})
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
    }

    resp.Diagnostics.Append(data.Fill(ctx, modelInstances)...)
    if resp.Diagnostics.HasError() {
      return
    }
    `
    )
    .with(
      "simple",
      () => `
    httpResp, err := d.client.${
      dataSource.api.method
    }WithResponse(${params.join(",")})
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

async function writeAndFormatGoFile(destination: URL, code: string) {
  await Bun.write(destination, code);
  await Bun.$`go fmt ${destination.pathname}`;
  await Bun.$`go tool goimports -w ${destination.pathname}`;
}

async function main() {
  console.log("Generating provider...");
  console.log("Generating resources...");
  console.log("Generating data sources...");
  for (const dataSource of DATASOURCES) {
    const code = generateDataSource({ dataSource });
    await writeAndFormatGoFile(
      new URL(`../provider/data_source_${dataSource.name}.go`, import.meta.url),
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
