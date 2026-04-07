CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS timescaledb;

-- Standard table for RBAC user management
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL CHECK (role IN ('admin', 'user')),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Metrics table optimized for time-series data
CREATE TABLE IF NOT EXISTS metrics (
    time TIMESTAMPTZ NOT NULL,
    node_name VARCHAR(255) NOT NULL,
    cpu_usage NUMERIC NOT NULL,
    memory_usage NUMERIC NOT NULL,
    disk_usage NUMERIC NOT NULL
);

-- Convert standard table into a TimescaleDB hypertable
SELECT create_hypertable('metrics', 'time', if_not_exists => TRUE);

-- Create default users with hashed passwords
-- pass for admin is 'admin123', pass for user is 'user123'
INSERT INTO users (username, password_hash, role) VALUES 
('admin', crypt('admin123', gen_salt('bf')), 'admin'),
('user1', crypt('user123', gen_salt('bf')), 'user')
ON CONFLICT (username) DO NOTHING;
