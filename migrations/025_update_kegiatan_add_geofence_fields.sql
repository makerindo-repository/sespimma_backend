-- +goose NO TRANSACTION
-- +goose Up

-- polygon_points: optional polygon boundary for complex geofence zones.
-- Mobile AttendanceZone has both radius (circle) and polygonPoints (List<LatLng>).
-- Stored as JSONB array: [{"lat": -6.967, "lng": 107.659}, ...]
ALTER TABLE kegiatan ADD COLUMN IF NOT EXISTS polygon_points JSONB;

-- cutoff_time: hard cutoff after which check-in is refused entirely.
-- Mobile AttendanceZone.cutoffTime is separate from batas_waktu_penugasan (soft deadline).
ALTER TABLE kegiatan ADD COLUMN IF NOT EXISTS cutoff_time TIME;

-- +goose Down

ALTER TABLE kegiatan DROP COLUMN IF EXISTS polygon_points;
ALTER TABLE kegiatan DROP COLUMN IF EXISTS cutoff_time;
