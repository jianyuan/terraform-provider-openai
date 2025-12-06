export type ComputedOptionalRequired =
  | "computed"
  | "optional"
  | "computed_optional"
  | "required";

export type Attribute =
  | StringAttribute
  | IntAttribute
  | SetNestedAttribute
  | ObjectAttribute;

export interface BaseAttribute {
  name: string;
  description: string;
  computedOptionalRequired: ComputedOptionalRequired;
  nullable?: boolean;
}

export interface StringAttribute extends BaseAttribute {
  type: "string";
}

export interface IntAttribute extends BaseAttribute {
  type: "int";
}

export interface SetNestedAttribute extends BaseAttribute {
  type: "set_nested";
  attributes: Array<Attribute>;
}

export interface ObjectAttribute extends BaseAttribute {
  type: "object";
  attributes: Array<Attribute>;
}

export interface DataSource {
  name: string;
  description: string;
  api: {
    strategy: "simple" | "paginate";
    method: string;
    model: string;
    params?: Array<string>;
  };
  attributes: Array<Attribute>;
}
