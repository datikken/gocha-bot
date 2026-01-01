-- Create migrations table to track applied migrations
CREATE TABLE IF NOT EXISTS schema_migrations (
    version VARCHAR(255) PRIMARY KEY,
    applied_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    telegram_id BIGINT UNIQUE,
    username TEXT,
    role TEXT DEFAULT 'customer',
    invited_by INTEGER REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW()
);

-- Create visits table
CREATE TABLE IF NOT EXISTS visits (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    visit_date TIMESTAMP,
    status TEXT DEFAULT 'scheduled',
    created_at TIMESTAMP DEFAULT NOW()
);

-- Insert migration record
INSERT INTO schema_migrations (version) VALUES ('000001_initial_schema') ON CONFLICT (version) DO NOTHING;
