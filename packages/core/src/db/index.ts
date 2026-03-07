/**
 * @file index.ts
 * @module packages/core/src/db/index
 *
 * WHAT:
 * Database connection initialization for Drizzle ORM.
 *
 * WHY:
 * Provides a singleton `db` instance used by all repositories.
 * Supports switching between SQLite (dev) and PostgreSQL (prod) via env vars,
 * though predominantly targets SQLite for this local-first architecture.
 */

import { drizzle } from "drizzle-orm/better-sqlite3";
import Database from "better-sqlite3";
import * as schema from "./metamcp-schema.js";
import * as dotenv from "dotenv";
import path from "path";
import fs from "fs";

dotenv.config();

// Default to SQLite local file
const dbPath = process.env.DATABASE_URL || "metamcp.db";

// Ensure we are using absolute path if it is a local file
const resolvedDbPath = dbPath.startsWith("file:")
    ? dbPath.slice(5)
    : path.resolve(process.cwd(), dbPath);

// Ensure the parent directory exists before opening the database
const dbDir = path.dirname(resolvedDbPath);
if (!fs.existsSync(dbDir)) {
    fs.mkdirSync(dbDir, { recursive: true });
}

const sqlite = new Database(resolvedDbPath);
export const db = drizzle(sqlite, { schema });

// Export the schema for convenience
export { schema };
