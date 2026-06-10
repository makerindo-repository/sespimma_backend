-- +goose Up

-- Real-time GPS tracking table.
-- Mobile location_sync_service.dart syncs every 10 seconds with 2m distance filter.
-- Backend uses this for live geofence violation detection and zone presence verification.
CREATE TABLE IF NOT EXISTS user_location_logs (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT        NOT NULL,
    latitude    DOUBLE PRECISION NOT NULL,
    longitude   DOUBLE PRECISION NOT NULL,
    accuracy    DOUBLE PRECISION,           -- GPS accuracy in meters (optional)
    recorded_at TIMESTAMPTZ   NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_location_logs_user
        FOREIGN KEY (user_id) REFERENCES users(id)
        ON DELETE CASCADE ON UPDATE CASCADE
);

-- Compound index for querying latest position of a specific user
CREATE INDEX idx_location_logs_user_time ON user_location_logs(user_id, recorded_at DESC);

-- Index for time-window queries (zone violation sweeps)
CREATE INDEX idx_location_logs_recorded_at ON user_location_logs(recorded_at DESC);

-- +goose Down

DROP INDEX IF EXISTS idx_location_logs_recorded_at;
DROP INDEX IF EXISTS idx_location_logs_user_time;
DROP TABLE  IF EXISTS user_location_logs;
