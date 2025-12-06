export type ComputedOptionalRequired =
  | "computed"
  | "optional"
  | "computed_optional"
  | "required";

export type Attribute =
  | StringAttribute
  | IntAttribute
  | BoolAttribute
  | SetNestedAttribute
  | ObjectAttribute;

export interface BaseAttribute {
  name: string;
  description: string;
  computedOptionalRequired: ComputedOptionalRequired;
  sensitive?: boolean;
  planModifiers?: Array<string>;
  validators?: Array<string>;
  nullable?: boolean;
}

export interface StringAttribute extends BaseAttribute {
  type: "string";
}

export interface IntAttribute extends BaseAttribute {
  type: "int";
}

export interface BoolAttribute extends BaseAttribute {
  type: "bool";
}

export interface SetNestedAttribute extends BaseAttribute {
  type: "set_nested";
  attributes: Array<Attribute>;
}

export interface ObjectAttribute extends BaseAttribute {
  type: "object";
  attributes: Array<Attribute>;
}

export interface BaseDataSourceApiStrategy {
  model: string;
  readMethod: string;
  readRequestAttributes?: Array<string>;
}

export type DataSourceApiStrategy =
  | SimpleDataSourceApiStrategy
  | PaginateDataSourceApiStrategy;

export interface SimpleDataSourceApiStrategy extends BaseDataSourceApiStrategy {
  strategy: "simple";
}

export interface PaginateDataSourceApiStrategy
  extends BaseDataSourceApiStrategy {
  strategy: "paginate";
  hooks?: {
    readInitLoop?: string;
    readPreIterate?: string;
    readPostIterate?: string;
  };
}

export interface DataSource {
  name: string;
  description: string;
  api: DataSourceApiStrategy;
  attributes: Array<Attribute>;
}

export interface ResourceApiStrategy {
  createMethod: string;
  createRequestAttributes?: Array<string>;
  readMethod: string;
  readRequestAttributes?: Array<string>;
  updateMethod?: string;
  updateRequestAttributes?: Array<string>;
  deleteMethod?: string;
  deleteRequestAttributes?: Array<string>;
  hooks?: {};
}

export interface Resource {
  name: string;
  description: string;
  api: ResourceApiStrategy;
  importStateAttributes?: Array<string>;
  attributes: Array<Attribute>;
}
