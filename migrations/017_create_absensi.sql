-- +goose Up

-- Create ENUM for status
CREATE TYPE attendance_status_enum AS ENUM (
    'hadir',
    'izin',
    'sakit',
    'tk'
);

CREATE TABLE absensi (
    id BIGSERIAL PRIMARY KEY,

    serdik_id BIGINT NOT NULL,
    kegiatan_id BIGINT NOT NULL,

    datetime TIMESTAMP NOT NULL,
    status attendance_status_enum NOT NULL,
    is_late BOOLEAN DEFAULT FALSE,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_absensi_serdik
        FOREIGN KEY (serdik_id)
        REFERENCES serdik(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT fk_absensi_kegiatan
        FOREIGN KEY (kegiatan_id)
        REFERENCES kegiatan(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

-- Indexes for faster lookups
CREATE INDEX idx_absensi_serdik_id ON absensi(serdik_id);
CREATE INDEX idx_absensi_kegiatan_id ON absensi(kegiatan_id);

-- Trigger for updated_at
CREATE TRIGGER trg_absensi_updated_at
BEFORE UPDATE ON absensi
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- +goose Down

DROP TRIGGER IF EXISTS trg_absensi_updated_at ON absensi;

DROP INDEX IF EXISTS idx_absensi_serdik_id;
DROP INDEX IF EXISTS idx_absensi_kegiatan_id;

DROP TABLE IF EXISTS absensi;

DROP TYPE IF EXISTS attendance_status_enum;
