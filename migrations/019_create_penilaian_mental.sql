-- +goose Up

CREATE TABLE penilaian_mental (
    id BIGSERIAL PRIMARY KEY,

    serdik_id BIGINT NOT NULL,
    mental_component_id BIGINT NOT NULL,

    nilai DOUBLE PRECISION NOT NULL,
    catatan TEXT,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_penilaian_mental_serdik
        FOREIGN KEY (serdik_id)
        REFERENCES serdik(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_penilaian_mental_component
        FOREIGN KEY (mental_component_id)
        REFERENCES mental_components(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_penilaian_mental_serdik_id
    ON penilaian_mental(serdik_id);

CREATE INDEX idx_penilaian_mental_component_id
    ON penilaian_mental(mental_component_id);

CREATE UNIQUE INDEX uq_penilaian_mental_serdik_component
    ON penilaian_mental(serdik_id, mental_component_id);

CREATE TRIGGER trigger_penilaian_mental_updated_at
BEFORE UPDATE ON penilaian_mental
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- +goose Down

DROP TRIGGER IF EXISTS trigger_penilaian_mental_updated_at ON penilaian_mental;

DROP TABLE IF EXISTS penilaian_mental;