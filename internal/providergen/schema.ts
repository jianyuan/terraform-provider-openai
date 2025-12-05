export type ComputedOptionalRequired =
  | "computed"
  | "optional"
  | "computed_optional"
  | "required";

export type Attribute = StringAttribute | IntAttribute | ObjectAttribute;

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

export interface ObjectAttribute extends BaseAttribute {
  type: "object";
  fields: Array<Attribute>;
}

export interface DataSource {
  name: string;
  description: string;
  api: {
    strategy: "simple";
    method: string;
    model: string;
  };
  attributes: Attribute[];
}
