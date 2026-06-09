-- +goose Up
CREATE TABLE akademik_component (
    id BIGSERIAL PRIMARY KEY,
    parent_id BIGINT NULL,
    code VARCHAR(50) UNIQUE,
    name VARCHAR(255) NOT NULL,
    weight NUMERIC(5,2) NOT NULL,
    level SMALLINT NOT NULL DEFAULT 1,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_akademik_component_parent
        FOREIGN KEY (parent_id)
        REFERENCES akademik_component(id)
        ON DELETE SET NULL
);
CREATE INDEX idx_akademik_component_parent_id
    ON akademik_component(parent_id);
CREATE TRIGGER trigger_akademik_component_updated_at
BEFORE UPDATE ON akademik_component
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- +goose Down

DROP TRIGGER IF EXISTS trigger_akademik_component_updated_at
ON akademik_component;
DROP INDEX IF EXISTS idx_akademik_component_parent_id;
DROP TABLE IF EXISTS akademik_component;