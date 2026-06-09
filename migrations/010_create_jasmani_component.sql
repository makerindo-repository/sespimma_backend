-- +goose Up

CREATE TYPE jasmani_age_group AS ENUM (
    'GOL_I',
    'GOL_II',
    'GOL_III',
    'GOL_IV'
);

CREATE TABLE jasmani_components (
    id BIGSERIAL PRIMARY KEY,
    parent_id BIGINT NULL,
    code VARCHAR(50) UNIQUE,
    name VARCHAR(255) NOT NULL,
    weight NUMERIC(5,2) NOT NULL,
    age_group jasmani_age_group NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_jasmani_components_parent
        FOREIGN KEY (parent_id)
        REFERENCES jasmani_components(id)
        ON DELETE SET NULL
);

CREATE INDEX idx_jasmani_components_parent_id
    ON jasmani_components(parent_id);

CREATE TRIGGER trigger_jasmani_components_updated_at
BEFORE UPDATE ON jasmani_components
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- +goose Down

DROP TRIGGER IF EXISTS trigger_jasmani_components_updated_at ON jasmani_components;

DROP TABLE IF EXISTS jasmani_components;

DROP TYPE IF EXISTS jasmani_age_group;