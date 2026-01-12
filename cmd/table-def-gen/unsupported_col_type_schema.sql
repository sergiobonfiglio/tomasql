-- Test schema with unsupported column types
CREATE TABLE test_table (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    -- Using an unsupported type (json) to test the ignore-unknown-types flag
    metadata JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);