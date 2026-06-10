-- +goose NO TRANSACTION
-- +goose Up

-- FCM token for push notifications (mobile requires this on login)
ALTER TABLE users ADD COLUMN IF NOT EXISTS fcm_token VARCHAR(500);

-- Separate refresh token from current_token (mobile expects access_token + refresh_token pair)
ALTER TABLE users ADD COLUMN IF NOT EXISTS refresh_token TEXT;

-- NAK approval flag required by mobile login response (is_nak_approved)
ALTER TABLE users ADD COLUMN IF NOT EXISTS is_nak_approved BOOLEAN NOT NULL DEFAULT FALSE;

-- Avatar/profile photo URL (centralized on users, all roles can upload)
ALTER TABLE users ADD COLUMN IF NOT EXISTS profile_photo VARCHAR(500);

-- +goose Down

ALTER TABLE users DROP COLUMN IF EXISTS fcm_token;
ALTER TABLE users DROP COLUMN IF EXISTS refresh_token;
ALTER TABLE users DROP COLUMN IF EXISTS is_nak_approved;
ALTER TABLE users DROP COLUMN IF EXISTS profile_photo;
