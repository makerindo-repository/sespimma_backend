-- +goose NO TRANSACTION
-- +goose Up

-- Mobile AttendanceModel requires latitude/longitude for GPS-based check-in validation
ALTER TABLE absensi ADD COLUMN IF NOT EXISTS latitude  DOUBLE PRECISION;
ALTER TABLE absensi ADD COLUMN IF NOT EXISTS longitude DOUBLE PRECISION;

-- is_location_valid: set by backend after verifying coords against kegiatan geofence radius
ALTER TABLE absensi ADD COLUMN IF NOT EXISTS is_location_valid BOOLEAN NOT NULL DEFAULT FALSE;

-- type: distinguishes check-in vs check-out records within same kegiatan
ALTER TABLE absensi ADD COLUMN IF NOT EXISTS type VARCHAR(20) NOT NULL DEFAULT 'checkin'
    CHECK (type IN ('checkin', 'checkout', 'absent', 'permission'));

-- method: how the attendance was recorded (gps, qr_code, manual)
ALTER TABLE absensi ADD COLUMN IF NOT EXISTS method VARCHAR(20) NOT NULL DEFAULT 'gps';

-- note: optional free-text note attached to the attendance record
ALTER TABLE absensi ADD COLUMN IF NOT EXISTS note TEXT;

-- user_id: non-serdik roles (gadik, korsis, etc.) also need attendance records.
-- serdik_id remains for serdik-specific queries; user_id covers all roles.
ALTER TABLE absensi ADD COLUMN IF NOT EXISTS user_id BIGINT;
ALTER TABLE absensi
    ADD CONSTRAINT fk_absensi_user
        FOREIGN KEY (user_id) REFERENCES users(id)
        ON DELETE CASCADE ON UPDATE CASCADE;

-- Make serdik_id and kegiatan_id nullable so generic user attendance can be recorded
-- without being tied to a specific serdik profile or kegiatan event.
ALTER TABLE absensi ALTER COLUMN serdik_id  DROP NOT NULL;
ALTER TABLE absensi ALTER COLUMN kegiatan_id DROP NOT NULL;

CREATE INDEX IF NOT EXISTS idx_absensi_user_id   ON absensi(user_id);
CREATE INDEX IF NOT EXISTS idx_absensi_type       ON absensi(type);
CREATE INDEX IF NOT EXISTS idx_absensi_datetime   ON absensi(datetime);

-- +goose Down

DROP INDEX IF EXISTS idx_absensi_datetime;
DROP INDEX IF EXISTS idx_absensi_type;
DROP INDEX IF EXISTS idx_absensi_user_id;

ALTER TABLE absensi DROP CONSTRAINT IF EXISTS fk_absensi_user;
ALTER TABLE absensi DROP COLUMN IF EXISTS user_id;
ALTER TABLE absensi DROP COLUMN IF EXISTS note;
ALTER TABLE absensi DROP COLUMN IF EXISTS method;
ALTER TABLE absensi DROP COLUMN IF EXISTS type;
ALTER TABLE absensi DROP COLUMN IF EXISTS is_location_valid;
ALTER TABLE absensi DROP COLUMN IF EXISTS longitude;
ALTER TABLE absensi DROP COLUMN IF EXISTS latitude;

ALTER TABLE absensi ALTER COLUMN serdik_id  SET NOT NULL;
ALTER TABLE absensi ALTER COLUMN kegiatan_id SET NOT NULL;
