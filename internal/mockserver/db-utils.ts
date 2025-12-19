import { createId } from "@paralleldrive/cuid2";

export function now() {
  return Math.floor(Date.now() / 1000);
}

export function idGenerator(prefix: string, separator: string = "_") {
  return () => `${prefix}${separator}${createId()}`;
}
