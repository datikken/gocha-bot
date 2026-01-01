-- Rollback initial schema migration

-- Drop tables in reverse order due to foreign key constraints
DROP TABLE IF EXISTS visits;
DROP TABLE IF EXISTS users;

-- Remove migration record
DELETE FROM schema_migrations WHERE version = '000001_initial_schema';
