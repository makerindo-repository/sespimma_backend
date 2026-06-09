-- +goose Up

CREATE TYPE user_role AS ENUM (
    'students',
    'admin',
    'teacher',
    'director',
    'patun'
);

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,

    email VARCHAR(150) NOT NULL UNIQUE,
    nrp VARCHAR(50) UNIQUE,
    nip VARCHAR(50) UNIQUE,
    password VARCHAR(255) NOT NULL,

    role user_role NOT NULL DEFAULT 'students',

    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    location GEOGRAPHY(Point, 4326),

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    current_token TEXT
);

CREATE TRIGGER trg_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();


CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_is_active ON users(is_active);

-- +goose Down

DROP TABLE IF EXISTS users;
DROP TYPE IF EXISTS user_role;