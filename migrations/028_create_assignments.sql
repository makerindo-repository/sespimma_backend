-- +goose Up

-- assignment_status_enum: lifecycle of a task from creation to archival
CREATE TYPE assignment_status_enum AS ENUM (
    'draft',
    'active',
    'closed',
    'archived'
);

-- assignments: tasks created by Gadik (instructors) and assigned to Pokjar groups.
-- Derived from mobile GadikAssignmentModel and TugasModel.
CREATE TABLE IF NOT EXISTS assignments (
    id               BIGSERIAL PRIMARY KEY,
    created_by       BIGINT        NOT NULL,  -- FK to users(id), must be gadik role
    judul            VARCHAR(500)  NOT NULL,
    deskripsi        TEXT,
    jenis_tugas      VARCHAR(100),            -- task type/category
    turunan_tugas    VARCHAR(100),            -- sub-type of task
    mapel            VARCHAR(255),            -- subject/course
    deadline         TIMESTAMPTZ   NOT NULL,
    target_pokjar_id INT,                     -- FK to pokjar(id); NULL = all pokjar
    instruksi        TEXT,
    status           assignment_status_enum NOT NULL DEFAULT 'active',
    file_name        VARCHAR(500),
    file_url         VARCHAR(1000),
    is_remedial      BOOLEAN       NOT NULL DEFAULT FALSE,
    created_at       TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ   NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_assignments_created_by
        FOREIGN KEY (created_by) REFERENCES users(id)
        ON DELETE RESTRICT ON UPDATE CASCADE,

    CONSTRAINT fk_assignments_pokjar
        FOREIGN KEY (target_pokjar_id) REFERENCES pokjar(id)
        ON DELETE SET NULL ON UPDATE CASCADE
);

CREATE TRIGGER trg_assignments_updated_at
    BEFORE UPDATE ON assignments
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE INDEX idx_assignments_created_by    ON assignments(created_by);
CREATE INDEX idx_assignments_pokjar        ON assignments(target_pokjar_id);
CREATE INDEX idx_assignments_status        ON assignments(status);
CREATE INDEX idx_assignments_deadline      ON assignments(deadline);

-- +goose Down

DROP TRIGGER  IF EXISTS trg_assignments_updated_at ON assignments;
DROP INDEX    IF EXISTS idx_assignments_deadline;
DROP INDEX    IF EXISTS idx_assignments_status;
DROP INDEX    IF EXISTS idx_assignments_pokjar;
DROP INDEX    IF EXISTS idx_assignments_created_by;
DROP TABLE    IF EXISTS assignments;
DROP TYPE     IF EXISTS assignment_status_enum;
