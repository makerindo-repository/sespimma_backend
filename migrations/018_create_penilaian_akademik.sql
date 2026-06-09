-- +goose Up

CREATE TABLE penilaian_akademik (
    id BIGSERIAL PRIMARY KEY,

    serdik_id BIGINT NOT NULL,
    akademik_component_id BIGINT NOT NULL,

    nilai DOUBLE PRECISION NOT NULL,
    catatan TEXT,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_penilaian_akademik_serdik
        FOREIGN KEY (serdik_id)
        REFERENCES serdik(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_penilaian_akademik_component
        FOREIGN KEY (akademik_component_id)
        REFERENCES akademik_component(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_penilaian_akademik_serdik_id
    ON penilaian_akademik(serdik_id);

CREATE INDEX idx_penilaian_akademik_component_id
    ON penilaian_akademik(akademik_component_id);

CREATE UNIQUE INDEX uq_penilaian_akademik_serdik_component
    ON penilaian_akademik(serdik_id, akademik_component_id);

CREATE TRIGGER trg_penilaian_akademik_updated_at
BEFORE UPDATE ON penilaian_akademik
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- +goose Down

DROP TRIGGER IF EXISTS trg_penilaian_akademik_updated_at
ON penilaian_akademik;

DROP INDEX IF EXISTS uq_penilaian_akademik_serdik_component;
DROP INDEX IF EXISTS idx_penilaian_akademik_serdik_id;
DROP INDEX IF EXISTS idx_penilaian_akademik_component_id;

DROP TABLE IF EXISTS penilaian_akademik;