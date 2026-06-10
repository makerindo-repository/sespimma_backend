-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE EXTENSION IF NOT EXISTS postgis;

-- +goose Down
DROP FUNCTION IF EXISTS update_updated_at_column;
DROP EXTENSION IF EXISTS postgis;