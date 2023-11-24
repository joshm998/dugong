#! /usr/bin/env node 
import { drizzle } from "drizzle-orm/better-sqlite3";
import { migrate } from "drizzle-orm/better-sqlite3/migrator";
import Database from 'better-sqlite3';

const sqlite = new Database('../cms/sqlite.db');
const db = drizzle(sqlite);
 
const startMigration = async () => {
    await migrate(db, { migrationsFolder: "./migrations" });
}

startMigration();