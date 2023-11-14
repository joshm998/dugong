import Database from "better-sqlite3";
import path from 'path';

export const db = Database(path.join(process.cwd(), "sqlite.db"), {});