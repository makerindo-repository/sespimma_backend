-- +goose Up

CREATE TABLE files (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    file_type TEXT NOT NULL,
    path TEXT NOT NULL,
    mime_type TEXT NOT NULL,
    size BIGINT NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP

);


-- trigger for updated_at
CREATE TRIGGER trg_files_updated_at
BEFORE UPDATE ON files
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();


-- +goose Down

DROP TABLE IF EXISTS files;
DROP TRIGGER IF EXISTS trg_files_updated_at ON files;