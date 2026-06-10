-- +goose NO TRANSACTION
-- +goose Up

-- nosis: display student ID used in GadikSubmissionModel.serdikNosis and profile screens
-- (may differ from no_serdik which is the internal system ID)
ALTER TABLE serdik ADD COLUMN IF NOT EXISTS nosis VARCHAR(100);

-- jabatan_senat: senate/student council position shown in Serdik profile screen
ALTER TABLE serdik ADD COLUMN IF NOT EXISTS jabatan_senat VARCHAR(255);

-- angkatan: cohort/generation, shown in profile and final recap report
ALTER TABLE serdik ADD COLUMN IF NOT EXISTS angkatan VARCHAR(100);

-- +goose Down

ALTER TABLE serdik DROP COLUMN IF EXISTS nosis;
ALTER TABLE serdik DROP COLUMN IF EXISTS jabatan_senat;
ALTER TABLE serdik DROP COLUMN IF EXISTS angkatan;
