-- Reset migration tracker so all migrations re-run cleanly.
-- This is a one-time fix after the cascading failure.
-- +goose Up

TRUNCATE _ci_migrations;

-- Drop ALL tables in reverse dependency order
DROP TABLE IF EXISTS health_records CASCADE;
DROP TABLE IF EXISTS serdik_health_data CASCADE;
DROP TABLE IF EXISTS sociometry_evaluations CASCADE;
DROP TABLE IF EXISTS sociometry_periods CASCADE;
DROP TABLE IF EXISTS izin_requests CASCADE;
DROP TABLE IF EXISTS submissions CASCADE;
DROP TABLE IF EXISTS assignments CASCADE;
DROP TABLE IF EXISTS user_location_logs CASCADE;
DROP TABLE IF EXISTS absensi CASCADE;
DROP TABLE IF EXISTS kegiatan CASCADE;
DROP TABLE IF EXISTS penilaian_mental CASCADE;
DROP TABLE IF EXISTS penilaian_akademik CASCADE;
DROP TABLE IF EXISTS penilaian_jasmani CASCADE;
DROP TABLE IF EXISTS mental_components CASCADE;
DROP TABLE IF EXISTS akademik_component CASCADE;
DROP TABLE IF EXISTS jasmani_components CASCADE;
DROP TABLE IF EXISTS punishment_logs CASCADE;
DROP TABLE IF EXISTS punishment_items CASCADE;
DROP TABLE IF EXISTS punishment_categories CASCADE;
DROP TABLE IF EXISTS user_rewards CASCADE;
DROP TABLE IF EXISTS reward_items CASCADE;
DROP TABLE IF EXISTS reward_categories CASCADE;
DROP TABLE IF EXISTS serdik CASCADE;
DROP TABLE IF EXISTS gadik CASCADE;
DROP TABLE IF EXISTS korsis CASCADE;
DROP TABLE IF EXISTS pokjar CASCADE;
DROP TABLE IF EXISTS patun CASCADE;
DROP TABLE IF EXISTS pimpinan CASCADE;
DROP TABLE IF EXISTS files CASCADE;
DROP TABLE IF EXISTS users CASCADE;

-- Drop ALL triggers
DROP TRIGGER IF EXISTS trg_users_updated_at ON users;
DROP TRIGGER IF EXISTS trg_files_updated_at ON files;
DROP TRIGGER IF EXISTS trg_serdik_updated_at ON serdik;

-- Drop ALL functions
DROP FUNCTION IF EXISTS update_updated_at_column CASCADE;

-- Drop ALL custom types
DROP TYPE IF EXISTS user_role CASCADE;
DROP TYPE IF EXISTS gender_enum CASCADE;
DROP TYPE IF EXISTS jasmani_age_group CASCADE;
DROP TYPE IF EXISTS mental_indicator_type CASCADE;
DROP TYPE IF EXISTS reward_item_period_type CASCADE;
DROP TYPE IF EXISTS user_reward_status CASCADE;
DROP TYPE IF EXISTS attendance_status_enum CASCADE;
DROP TYPE IF EXISTS assignment_status_enum CASCADE;
DROP TYPE IF EXISTS assignment_type CASCADE;
DROP TYPE IF EXISTS submission_status CASCADE;
DROP TYPE IF EXISTS submission_status_enum CASCADE;
DROP TYPE IF EXISTS izin_status CASCADE;
DROP TYPE IF EXISTS izin_status_enum CASCADE;
DROP TYPE IF EXISTS izin_type CASCADE;
DROP TYPE IF EXISTS sociometry_period_type CASCADE;

-- +goose Down
-- no-op: this is a one-time reset
