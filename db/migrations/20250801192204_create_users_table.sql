-- migrate:up
CREATE TABLE users (
                       id UUID PRIMARY KEY,
                       auth0_id TEXT UNIQUE NOT NULL,
                       email TEXT UNIQUE NOT NULL,
                       name TEXT,
                       created_at TIMESTAMP
);


-- migrate:down
DROP TABLE IF EXISTS users;