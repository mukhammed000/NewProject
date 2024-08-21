CREATE TABLE IF NOT EXISTS users (
    user_id UUID PRIMARY KEY,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    is_email_verified Boolean DEFAULT false,
    gender VARCHAR(10) NOT NULL CHECK (gender IN ('male', 'female')),
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL,
    date_of_birth VARCHAR(20) NOT NULL,
    created_at VARCHAR(255) NOT NULL,
    updated_at VARCHAR(255),
    deleted_at BIGINT DEFAULT 0
);

CREATE TABLE IF NOT EXISTS tokens (
    token_id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    access_token TEXT NOT NULL,
    refresh_token TEXT NOT NULL,
    created_at VARCHAR(255) NOT NULL,
    updated_at VARCHAR(255),
    expires_at VARCHAR(255),
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);
