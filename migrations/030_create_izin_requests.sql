-- +goose Up

-- izin_status_enum: approval workflow for leave requests
CREATE TYPE izin_status_enum AS ENUM (
    'pending',      -- submitted by serdik, awaiting review
    'disetujui',    -- approved by korsis/patun
    'ditolak'       -- rejected
);

-- izin_requests: leave/absence requests submitted by serdik.
-- Derived from mobile InboxItem (isIzin=true) and LeaveFormSheet widget.
-- Reviewed by korsis or patun role users.
CREATE TABLE IF NOT EXISTS izin_requests (
    id                   BIGSERIAL PRIMARY KEY,
    serdik_id            BIGINT         NOT NULL,
    kegiatan_id          BIGINT,                  -- linked activity; NULL for general leave

    start_time           TIMESTAMPTZ    NOT NULL,
    end_time             TIMESTAMPTZ    NOT NULL,
    description          TEXT           NOT NULL,
    attachment_path      VARCHAR(1000),
    attachment_file_name VARCHAR(500),

    -- Review fields
    reviewed_by          BIGINT,                  -- FK to users(id)
    reviewed_at          TIMESTAMPTZ,
    rejection_reason     TEXT,

    status               izin_status_enum NOT NULL DEFAULT 'pending',
    created_at           TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at           TIMESTAMPTZ    NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_izin_serdik
        FOREIGN KEY (serdik_id) REFERENCES serdik(id)
        ON DELETE CASCADE ON UPDATE CASCADE,

    CONSTRAINT fk_izin_kegiatan
        FOREIGN KEY (kegiatan_id) REFERENCES kegiatan(id)
        ON DELETE SET NULL ON UPDATE CASCADE,

    CONSTRAINT fk_izin_reviewer
        FOREIGN KEY (reviewed_by) REFERENCES users(id)
        ON DELETE SET NULL ON UPDATE CASCADE
);

CREATE TRIGGER trg_izin_requests_updated_at
    BEFORE UPDATE ON izin_requests
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE INDEX idx_izin_requests_serdik_id  ON izin_requests(serdik_id);
CREATE INDEX idx_izin_requests_status     ON izin_requests(status);
CREATE INDEX idx_izin_requests_start_time ON izin_requests(start_time);

-- +goose Down

DROP TRIGGER  IF EXISTS trg_izin_requests_updated_at ON izin_requests;
DROP INDEX    IF EXISTS idx_izin_requests_start_time;
DROP INDEX    IF EXISTS idx_izin_requests_status;
DROP INDEX    IF EXISTS idx_izin_requests_serdik_id;
DROP TABLE    IF EXISTS izin_requests;
DROP TYPE     IF EXISTS izin_status_enum;
