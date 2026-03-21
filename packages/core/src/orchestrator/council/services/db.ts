import { db, schema } from '../../../db/index.js';

class DatabaseService {
  constructor() {}

  getDb() {
    // In Borg, we use Drizzle 'db' for ORM and 'db.session.client' for raw sqlite if needed.
    // However, the migrated code expects a 'better-sqlite3' instance from 'getDb()'.
    return (db as any).session.client;
  }

  getDrizzle() {
    return db;
  }

  getSchema() {
    return schema;
  }

  close() {
    // Managed by central DB index
  }
}

export const dbService = new DatabaseService();
