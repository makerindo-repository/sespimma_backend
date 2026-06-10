-- +goose Up

-- serdik_health_data: aggregated health assessment data per serdik.
-- Derived from mobile HealthMonitoringData / SerdikHealthData model.
-- Nilai A/B are doctor-assigned scores; Nilai C starts at 80 and decreases per health_records.
-- Final NKes = (nilaiA + nilaiB + nilaiC) / 3, used in physical score (NKJ) calculation.
CREATE TABLE IF NOT EXISTS serdik_health_data (
    id                BIGSERIAL PRIMARY KEY,
    serdik_id         BIGINT         NOT NULL UNIQUE,

    -- Nilai A: initial health examination score
    nilai_a           DOUBLE PRECISION,
    catatan_dokter_a  TEXT,

    -- Nilai B: final/follow-up health examination score
    nilai_b           DOUBLE PRECISION,
    catatan_dokter_b  TEXT,

    -- Nilai C base score; decremented by minus_points from individual health_records
    base_nilai_c      INT            NOT NULL DEFAULT 80,

    created_at        TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ    NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_health_data_serdik
        FOREIGN KEY (serdik_id) REFERENCES serdik(id)
        ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TRIGGER trg_serdik_health_data_updated_at
    BEFORE UPDATE ON serdik_health_data
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- health_records: individual medical incidents recorded by medis role users.
-- Each record carries a minus_points deduction applied to nilai_c (base_nilai_c - sum(minus_points)).
CREATE TABLE IF NOT EXISTS health_records (
    id                      BIGSERIAL PRIMARY KEY,
    serdik_health_data_id   BIGINT         NOT NULL,
    medis_user_id           BIGINT         NOT NULL,   -- medis role user who recorded this
    type                    VARCHAR(100)   NOT NULL,   -- e.g. 'rawat_jalan', 'rawat_inap'
    description             TEXT           NOT NULL,
    photo_path              VARCHAR(1000),
    minus_points            INT            NOT NULL DEFAULT 0 CHECK (minus_points >= 0),
    recorded_at             TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    created_at              TIMESTAMPTZ    NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_health_record_data
        FOREIGN KEY (serdik_health_data_id) REFERENCES serdik_health_data(id)
        ON DELETE CASCADE ON UPDATE CASCADE,

    CONSTRAINT fk_health_record_medis
        FOREIGN KEY (medis_user_id) REFERENCES users(id)
        ON DELETE RESTRICT ON UPDATE CASCADE
);

CREATE INDEX idx_health_records_health_data_id ON health_records(serdik_health_data_id);
CREATE INDEX idx_health_records_medis_user_id  ON health_records(medis_user_id);
CREATE INDEX idx_health_records_recorded_at    ON health_records(recorded_at DESC);

-- +goose Down

DROP INDEX    IF EXISTS idx_health_records_recorded_at;
DROP INDEX    IF EXISTS idx_health_records_medis_user_id;
DROP INDEX    IF EXISTS idx_health_records_health_data_id;
DROP TABLE    IF EXISTS health_records;

DROP TRIGGER  IF EXISTS trg_serdik_health_data_updated_at ON serdik_health_data;
DROP TABLE    IF EXISTS serdik_health_data;
