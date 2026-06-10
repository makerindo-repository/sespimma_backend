-- +goose Up

-- submission_status_enum: tracks grading lifecycle
CREATE TYPE submission_status_enum AS ENUM (
    'pending',      -- submitted, not yet graded
    'graded',       -- graded by gadik
    'remedial',     -- failed, requires remedial resubmission
    'late'          -- submitted after deadline
);

-- submissions: serdik assignment submissions with full Gadik scoring breakdown.
-- All score fields derived from mobile GadikSubmissionModel.
CREATE TABLE IF NOT EXISTS submissions (
    id               BIGSERIAL PRIMARY KEY,
    assignment_id    BIGINT         NOT NULL,
    serdik_id        BIGINT         NOT NULL,

    submitted_at     TIMESTAMPTZ,
    file_url         VARCHAR(1000),
    file_name        VARCHAR(500),

    -- Grading state
    is_graded        BOOLEAN        NOT NULL DEFAULT FALSE,
    is_remedial      BOOLEAN        NOT NULL DEFAULT FALSE,
    status           submission_status_enum NOT NULL DEFAULT 'pending',
    catatan_pengajar TEXT,
    nilai_akhir      DOUBLE PRECISION,        -- final weighted score

    -- Detailed score breakdown (GadikSubmissionModel scoring fields)
    score_materi                    DOUBLE PRECISION,
    score_penulisan                 DOUBLE PRECISION,
    score_paparan                   DOUBLE PRECISION,
    score_keaktifan                 DOUBLE PRECISION,
    score_ujian                     DOUBLE PRECISION,
    score_keaktifan_perseorangan    DOUBLE PRECISION,
    score_produk_perseorangan       DOUBLE PRECISION,
    score_tata_ruang                DOUBLE PRECISION,

    created_at       TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ    NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_submissions_assignment
        FOREIGN KEY (assignment_id) REFERENCES assignments(id)
        ON DELETE CASCADE ON UPDATE CASCADE,

    CONSTRAINT fk_submissions_serdik
        FOREIGN KEY (serdik_id) REFERENCES serdik(id)
        ON DELETE CASCADE ON UPDATE CASCADE,

    -- One submission per serdik per assignment; remedials create a new record
    CONSTRAINT uq_submission_assignment_serdik
        UNIQUE (assignment_id, serdik_id)
);

CREATE TRIGGER trg_submissions_updated_at
    BEFORE UPDATE ON submissions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE INDEX idx_submissions_assignment_id ON submissions(assignment_id);
CREATE INDEX idx_submissions_serdik_id     ON submissions(serdik_id);
CREATE INDEX idx_submissions_status        ON submissions(status);

-- +goose Down

DROP TRIGGER  IF EXISTS trg_submissions_updated_at ON submissions;
DROP INDEX    IF EXISTS idx_submissions_status;
DROP INDEX    IF EXISTS idx_submissions_serdik_id;
DROP INDEX    IF EXISTS idx_submissions_assignment_id;
DROP TABLE    IF EXISTS submissions;
DROP TYPE     IF EXISTS submission_status_enum;
