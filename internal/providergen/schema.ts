export type ComputedOptionalRequired =
  | "computed"
  | "optional"
  | "computed_optional"
  | "required";

export type Attribute =
  | StringAttribute
  | Int64Attribute
  | Float64Attribute
  | BoolAttribute
  | ListAttribute
  | ListNestedAttribute
  | SetAttribute
  | SetNestedAttribute
  | MapAttribute
  | ObjectAttribute
  | SingleNestedAttribute;

export interface BaseAttribute {
  name: string;
  customType?: {
    type: string;
    value: string;
  };
  description: string;
  computedOptionalRequired: ComputedOptionalRequired;
  sensitive?: boolean;
  planModifiers?: Array<string>;
  validators?: Array<string>;
  nullable?: boolean;
  filler?: {
    skip?: boolean;
    sourceAttribute?: Array<string>;
    destinationAttribute?: Array<string>;
  };
}

export interface StringAttribute extends BaseAttribute {
  type: "string";
}

export interface Int64Attribute extends BaseAttribute {
  type: "int64";
}

export interface Float64Attribute extends BaseAttribute {
  type: "float64";
}

export interface BoolAttribute extends BaseAttribute {
  type: "bool";
}

export interface ListAttribute extends BaseAttribute {
  type: "list";
  elementType: "string";
}

export interface ListNestedAttribute extends BaseAttribute {
  type: "list_nested";
  attributes: Array<Attribute>;
  filler?: BaseAttribute["filler"] & {
    model: string;
  };
}

export interface SetAttribute extends BaseAttribute {
  type: "set";
  elementType: "string";
}

export interface SetNestedAttribute extends BaseAttribute {
  type: "set_nested";
  attributes: Array<Attribute>;
  filler?: BaseAttribute["filler"] & {
    model: string;
  };
}

export interface MapAttribute extends BaseAttribute {
  type: "map";
  elementType: "string";
}

export interface SingleNestedAttribute extends BaseAttribute {
  type: "single_nested";
  attributes: Array<Attribute>;
  filler?: BaseAttribute["filler"] & {
    model: string;
  };
}

export interface ObjectAttribute extends BaseAttribute {
  type: "object";
  attributes: Array<Attribute>;
}

export interface SingleNestedAttribute extends BaseAttribute {
  type: "single_nested";
  attributes: Array<Attribute>;
  filler?: BaseAttribute["filler"] & {
    model: string;
  };
}

export interface BaseDataSourceApiStrategy {
  readMethod: string;
  readRequestAttributes?: Array<string>;
}

export type DataSourceApiStrategy =
  | SimpleDataSourceApiStrategy
  | PaginateDataSourceApiStrategy;

export interface SimpleDataSourceApiStrategy extends BaseDataSourceApiStrategy {
  readStrategy: "simple";
}

export interface PaginateDataSourceApiStrategy extends BaseDataSourceApiStrategy {
  readStrategy: "paginate";
  readRequestParamsStruct: string;
  readModel: string;
  readInitLoop?: string;
  readPreIterate?: string;
  readPostIterate?: string;
}

export interface DataSource {
  name: string;
  description: string;
  api: DataSourceApiStrategy;
  filler?: {
    model: string;
  };
  attributes: Array<Attribute>;
}

export interface BaseResourceApiStrategy {
  method?: string;
  createMethod: string;
  createRequestAttributes?: Array<string>;
  readMethod: string;
  readRequestAttributes?: Array<string>;
  readRequestParamsStruct?: string;
  updateMethod?: string;
  updateRequestAttributes?: Array<string>;
  deleteMethod?: string;
  deleteRequestAttributes?: Array<string>;
}

export interface SimpleResourceApiStrategy extends BaseResourceApiStrategy {
  readStrategy?: never;
}

export interface PaginateResourceApiStrategy extends BaseResourceApiStrategy {
  readStrategy: "paginate";
  readModel: string;
}

export type ResourceApiStrategy =
  | SimpleResourceApiStrategy
  | PaginateResourceApiStrategy;

export interface Resource {
  name: string;
  description: string;
  api: ResourceApiStrategy;
  importStateAttributes?: Array<string>;
  filler?: {
    model: string;
  };
  attributes: Array<Attribute>;
}
