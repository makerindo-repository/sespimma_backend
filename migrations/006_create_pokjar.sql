-- +goose Up

CREATE TABLE pokjar (
    id BIGSERIAL PRIMARY KEY,

    name VARCHAR(255) NOT NULL UNIQUE,
    grade VARCHAR(50) NOT NULL,

    patun_id BIGINT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_pokjar_patun
        FOREIGN KEY (patun_id)
        REFERENCES patun(id)
        ON DELETE SET NULL
);

-- index for foreign key-like field
CREATE INDEX idx_pokjar_patun_id ON pokjar(patun_id);

-- trigger for updated_at
CREATE TRIGGER trg_pokjar_updated_at
BEFORE UPDATE ON pokjar
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();


-- +goose Down

DROP TABLE IF EXISTS pokjar;
DROP TRIGGER IF EXISTS trg_pokjar_updated_at ON pokjar;
DROP INDEX IF EXISTS idx_pokjar_patun_id;
