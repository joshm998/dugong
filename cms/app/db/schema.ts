import { sql } from "drizzle-orm";
import { integer, sqliteTable, text } from "drizzle-orm/sqlite-core";

export const documents = sqliteTable('documents', {
    id: integer('id', { mode: 'number' }).primaryKey({ autoIncrement: true }),
    createdAt: text("created_at").default(sql`CURRENT_DATE`).notNull(),
    publishedAt: text("published_at"),
    name: text("name"),
    type: text("type")
});

export const settings = sqliteTable('settings', {
    id: integer('id', { mode: 'number' }).primaryKey({ autoIncrement: true }),
    createdAt: text("created_at").default(sql`CURRENT_DATE`).notNull(),
    publishedAt: text("published_at"),
    value: text("value", { mode: 'json' }),
    key: text("key"),
    name: text("name"),
    description: text("description")
});

export const fields = sqliteTable('fields', {
    id: integer('id', { mode: 'number' }).primaryKey({ autoIncrement: true }),
    createdAt: text("created_at").default(sql`CURRENT_DATE`).notNull(),
    value: text("value", { mode: 'json' }),
    key: text("key"),
    name: text("name"),
    fieldType: text("field_type"),
});

export const content = sqliteTable('content', {
    id: integer('id', { mode: 'number' }).primaryKey({ autoIncrement: true }),
    createdAt: text("created_at").default(sql`CURRENT_DATE`).notNull(),
    documentId: integer("document_id").references(() => documents.id),
    fieldId: integer("field_id").references(() => fields.id)
});


