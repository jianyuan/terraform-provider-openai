import { createId } from "@paralleldrive/cuid2";
import { SQL, getTableColumns, sql } from "drizzle-orm";
import type { SQLiteTable } from "drizzle-orm/sqlite-core";

export function now() {
  return Math.floor(Date.now() / 1000);
}

export function idGenerator(prefix: string) {
  return () => `${prefix}${createId()}`;
}

export function buildConflictUpdateColumns<
  T extends SQLiteTable,
  Q extends keyof T["_"]["columns"],
>(table: T, columns: Q[]): Record<Q, SQL> {
  const cls = getTableColumns(table);
  return columns.reduce(
    (acc, column) => {
      const colName = cls[column]!.name;
      acc[column] = sql.raw(`excluded.${colName}`);
      return acc;
    },
    {} as Record<Q, SQL>,
  );
}
