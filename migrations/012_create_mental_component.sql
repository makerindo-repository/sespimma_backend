-- +goose Up

CREATE TYPE mental_indicator_type AS ENUM (
    'component',
    'reward',
    'punishment'
);

CREATE TABLE mental_components (
    id BIGSERIAL PRIMARY KEY,
    parent_id BIGINT NULL,
    code VARCHAR(50) UNIQUE,
    name VARCHAR(255) NOT NULL,
    weight NUMERIC(5,2) NOT NULL,
    indicator_type mental_indicator_type DEFAULT 'component',
    point NUMERIC(6,2),
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_mental_components_parent
        FOREIGN KEY (parent_id)
        REFERENCES mental_components(id)
        ON DELETE SET NULL
);

CREATE INDEX idx_mental_components_parent_id
    ON mental_components(parent_id);

CREATE TRIGGER trigger_mental_components_updated_at
BEFORE UPDATE ON mental_components
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- +goose Down

DROP TRIGGER IF EXISTS trigger_mental_components_updated_at ON mental_components;

DROP TABLE IF EXISTS mental_components;

DROP TYPE IF EXISTS mental_indicator_type;