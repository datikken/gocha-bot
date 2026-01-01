# Database Migrations

This project uses a custom migration system to manage database schema changes. Migrations are stored in the `migrations/` directory and are automatically applied when the application starts.

## Migration Files

Each migration consists of two files:
- `{version}_{name}.up.sql` - Contains the SQL to apply the migration
- `{version}_{name}.down.sql` - Contains the SQL to rollback the migration

## Migration Versioning

Versions should be in the format `000001`, `000002`, etc., with leading zeros to ensure proper sorting.

## Running Migrations

### Automatic (Default)
Migrations are automatically run when the application starts:
```bash
./gocha-bot
```

### Manual Commands

#### Run migrations only
```bash
./gocha-bot -migrate
```

#### Rollback to a specific version
```bash
./gocha-bot -rollback=000001
```

## Creating New Migrations

1. Create two files in the `migrations/` directory:
   - `000002_add_new_feature.up.sql`
   - `000002_add_new_feature.down.sql`

2. Write the appropriate SQL in each file

3. Test the migration by running:
   ```bash
   ./gocha-bot -migrate
   ```

4. Test rollback:
   ```bash
   ./gocha-bot -rollback=000002
   ```

## Migration Tracking

Applied migrations are tracked in the `schema_migrations` table with the following structure:
- `version` (VARCHAR(255) PRIMARY KEY)
- `applied_at` (TIMESTAMP WITH TIME ZONE)

## Example Migration Files

### 000002_add_user_email.up.sql
```sql
-- Add email column to users table
ALTER TABLE users ADD COLUMN email TEXT UNIQUE;

-- Insert migration record
INSERT INTO schema_migrations (version) VALUES ('000002_add_user_email') ON CONFLICT (version) DO NOTHING;
```

### 000002_add_user_email.down.sql
```sql
-- Remove email column from users table
ALTER TABLE users DROP COLUMN IF EXISTS email;

-- Remove migration record
DELETE FROM schema_migrations WHERE version = '000002_add_user_email';
```

## Best Practices

1. Always test migrations on a copy of production data first
2. Include both up and down migrations for every change
3. Use transactions in complex migrations when possible
4. Keep migration files small and focused on a single change
5. Never modify existing migration files after they've been applied in production
