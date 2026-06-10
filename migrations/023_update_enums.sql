-- +goose NO TRANSACTION
-- +goose Up

-- Mobile profile_screen.dart shows role 'kabag_bindik' (Head of Education division)
ALTER TYPE user_role ADD VALUE IF NOT EXISTS 'kabag_bindik';

-- Mobile AttendanceModel uses 'alpha' for absent; backend had 'tk' (tidak keterangan).
-- Adding 'alpha' as the canonical mobile-facing value. 'tk' is retained for legacy data.
ALTER TYPE attendance_status_enum ADD VALUE IF NOT EXISTS 'alpha';

-- +goose Down
-- Postgres does not support removing enum values; handled by full enum rebuild if needed.
