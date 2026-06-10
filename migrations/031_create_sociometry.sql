-- +goose Up

-- sociometry_period_type: awal = initial evaluation at start, akhir = final evaluation.
-- Used in NK (mental score) formula: NS = (sosiometri_awal + sosiometri_akhir) / 2
CREATE TYPE sociometry_period_type AS ENUM ('awal', 'akhir');

-- sociometry_periods: defines an active evaluation window for a pokjar group.
-- Managed by korsis role via mobile sosiometry period management screen.
CREATE TABLE IF NOT EXISTS sociometry_periods (
    id          BIGSERIAL PRIMARY KEY,
    pokjar_id   INT          NOT NULL,
    period_type sociometry_period_type NOT NULL,
    start_date  DATE         NOT NULL,
    end_date    DATE         NOT NULL,
    is_active   BOOLEAN      NOT NULL DEFAULT TRUE,
    created_by  BIGINT,                             -- korsis user who opened this period
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_sociometry_period_pokjar
        FOREIGN KEY (pokjar_id) REFERENCES pokjar(id)
        ON DELETE RESTRICT ON UPDATE CASCADE,

    CONSTRAINT fk_sociometry_period_created_by
        FOREIGN KEY (created_by) REFERENCES users(id)
        ON DELETE SET NULL ON UPDATE CASCADE,

    -- Only one active period of each type per pokjar at a time
    CONSTRAINT uq_sociometry_active_per_pokjar_type
        UNIQUE (pokjar_id, period_type, is_active)
);

CREATE TRIGGER trg_sociometry_periods_updated_at
    BEFORE UPDATE ON sociometry_periods
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE INDEX idx_sociometry_periods_pokjar    ON sociometry_periods(pokjar_id);
CREATE INDEX idx_sociometry_periods_is_active ON sociometry_periods(is_active);

-- sociometry_evaluations: peer-to-peer scores within a period.
-- Each serdik evaluates every other serdik in the same pokjar.
-- Mobile SociometryPeerModel.isEvaluated tracks completion.
CREATE TABLE IF NOT EXISTS sociometry_evaluations (
    id                   BIGSERIAL PRIMARY KEY,
    period_id            BIGINT       NOT NULL,
    evaluator_serdik_id  BIGINT       NOT NULL,
    evaluated_serdik_id  BIGINT       NOT NULL,
    score                DOUBLE PRECISION NOT NULL DEFAULT 0 CHECK (score >= 0 AND score <= 100),
    created_at           TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at           TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_socie_period
        FOREIGN KEY (period_id) REFERENCES sociometry_periods(id)
        ON DELETE CASCADE ON UPDATE CASCADE,

    CONSTRAINT fk_socie_evaluator
        FOREIGN KEY (evaluator_serdik_id) REFERENCES serdik(id)
        ON DELETE CASCADE ON UPDATE CASCADE,

    CONSTRAINT fk_socie_evaluated
        FOREIGN KEY (evaluated_serdik_id) REFERENCES serdik(id)
        ON DELETE CASCADE ON UPDATE CASCADE,

    -- One score per evaluator-evaluated pair per period
    CONSTRAINT uq_socie_evaluation
        UNIQUE (period_id, evaluator_serdik_id, evaluated_serdik_id),

    -- Cannot evaluate oneself
    CONSTRAINT chk_socie_no_self_eval
        CHECK (evaluator_serdik_id <> evaluated_serdik_id)
);

CREATE TRIGGER trg_sociometry_evaluations_updated_at
    BEFORE UPDATE ON sociometry_evaluations
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE INDEX idx_socie_evals_period    ON sociometry_evaluations(period_id);
CREATE INDEX idx_socie_evals_evaluated ON sociometry_evaluations(evaluated_serdik_id);
CREATE INDEX idx_socie_evals_evaluator ON sociometry_evaluations(evaluator_serdik_id);

-- +goose Down

DROP TRIGGER  IF EXISTS trg_sociometry_evaluations_updated_at ON sociometry_evaluations;
DROP TRIGGER  IF EXISTS trg_sociometry_periods_updated_at     ON sociometry_periods;

DROP INDEX    IF EXISTS idx_socie_evals_evaluator;
DROP INDEX    IF EXISTS idx_socie_evals_evaluated;
DROP INDEX    IF EXISTS idx_socie_evals_period;
DROP TABLE    IF EXISTS sociometry_evaluations;

DROP INDEX    IF EXISTS idx_sociometry_periods_is_active;
DROP INDEX    IF EXISTS idx_sociometry_periods_pokjar;
DROP TABLE    IF EXISTS sociometry_periods;

DROP TYPE     IF EXISTS sociometry_period_type;
