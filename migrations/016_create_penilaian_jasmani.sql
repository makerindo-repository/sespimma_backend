-- +goose Up

CREATE TABLE penilaian_jasmani (
    id BIGSERIAL PRIMARY KEY,

    serdik_id BIGINT NOT NULL,
    jasmani_component_id BIGINT NOT NULL,

    nilai DOUBLE PRECISION NOT NULL,
    catatan TEXT,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_penilaian_jasmani_serdik
        FOREIGN KEY (serdik_id)
        REFERENCES serdik(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_penilaian_jasmani_component
        FOREIGN KEY (jasmani_component_id)
        REFERENCES jasmani_components(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_penilaian_jasmani_serdik_id
    ON penilaian_jasmani(serdik_id);

CREATE INDEX idx_penilaian_jasmani_component_id
    ON penilaian_jasmani(jasmani_component_id);

CREATE TRIGGER trg_penilaian_jasmani_updated_at
BEFORE UPDATE ON penilaian_jasmani
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- +goose Down

DROP TRIGGER IF EXISTS trg_penilaian_jasmani_updated_at ON penilaian_jasmani;

DROP INDEX IF EXISTS idx_penilaian_jasmani_serdik_id;
DROP INDEX IF EXISTS idx_penilaian_jasmani_component_id;

DROP TABLE IF EXISTS penilaian_jasmani;
