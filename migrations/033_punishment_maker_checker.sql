-- +goose Up

-- Maker-Checker workflow for punishments (mirrors user_rewards.status).
-- Punishments created by Patun/Gadik must be approved by Korsis before they
-- count against a serdik's accumulated score.
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'punishment_log_status') THEN
        CREATE TYPE punishment_log_status AS ENUM (
            'pending',   -- created by maker (patun/gadik), awaiting korsis review
            'approved',  -- approved by korsis, now counts against the serdik
            'rejected'   -- rejected by korsis, no effect on score
        );
    END IF;
END$$;

ALTER TABLE punishment_logs
    ADD COLUMN IF NOT EXISTS status           punishment_log_status NOT NULL DEFAULT 'pending',
    ADD COLUMN IF NOT EXISTS approved_by       BIGINT,
    ADD COLUMN IF NOT EXISTS reviewed_at       TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS rejection_reason  TEXT;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_punishment_logs_approved_by'
    ) THEN
        ALTER TABLE punishment_logs
            ADD CONSTRAINT fk_punishment_logs_approved_by
            FOREIGN KEY (approved_by) REFERENCES users(id) ON DELETE SET NULL;
    END IF;
END$$;

-- Existing rows predate the workflow; treat them as already-approved so history
-- and accumulated scores remain intact after this migration.
UPDATE punishment_logs SET status = 'approved' WHERE status = 'pending' AND created_at < NOW();

CREATE INDEX IF NOT EXISTS idx_punishment_logs_status       ON punishment_logs(status);
CREATE INDEX IF NOT EXISTS idx_punishment_logs_user_status  ON punishment_logs(user_id, status);

-- +goose Down
ALTER TABLE punishment_logs
    DROP COLUMN IF EXISTS status,
    DROP COLUMN IF EXISTS approved_by,
    DROP COLUMN IF EXISTS reviewed_at,
    DROP COLUMN IF EXISTS rejection_reason;
DROP TYPE IF EXISTS punishment_log_status;
