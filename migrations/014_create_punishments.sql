-- +goose Up

-- =========================
-- punishment_categories
-- =========================
CREATE TABLE IF NOT EXISTS punishment_categories (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    weight_percent NUMERIC(5,2) NOT NULL DEFAULT 0,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_punishment_categories_code
    ON punishment_categories(code);

CREATE TRIGGER trigger_punishment_categories_updated_at
BEFORE UPDATE ON punishment_categories
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();


-- =========================
-- punishment_items
-- =========================
CREATE TABLE IF NOT EXISTS punishment_items (
    id BIGSERIAL PRIMARY KEY,

    category_id BIGINT NOT NULL,
    item_no INT NOT NULL,

    activity TEXT NOT NULL,
    point NUMERIC(5,2) NOT NULL,

    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_punishment_category
        FOREIGN KEY (category_id)
        REFERENCES punishment_categories(id)
        ON DELETE CASCADE,

    CONSTRAINT uniq_category_item_no UNIQUE (category_id, item_no)
);

CREATE INDEX idx_punishment_items_category
    ON punishment_items(category_id);

CREATE INDEX idx_punishment_items_active
    ON punishment_items(is_active);

CREATE TRIGGER trigger_punishment_items_updated_at
BEFORE UPDATE ON punishment_items
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();


-- =========================
-- punishment_logs
-- =========================
CREATE TABLE IF NOT EXISTS punishment_logs (
    id BIGSERIAL PRIMARY KEY,

    user_id BIGINT NOT NULL,
    punishment_item_id BIGINT NOT NULL,
    file_id UUID NULL,

    qty INT NOT NULL DEFAULT 1,
    point NUMERIC(8,2) NOT NULL,

    violation_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    notes TEXT NULL,

    created_by BIGINT NULL,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_log_punishment_item
        FOREIGN KEY (punishment_item_id)
        REFERENCES punishment_items(id)
        ON DELETE RESTRICT,

    CONSTRAINT fk_log_file
        FOREIGN KEY (file_id)
        REFERENCES files(id)
        ON DELETE SET NULL,

    CONSTRAINT fk_log_created_by
        FOREIGN KEY (created_by)
        REFERENCES users(id)
        ON DELETE SET NULL
);

CREATE INDEX idx_punishment_logs_user
    ON punishment_logs(user_id);

CREATE INDEX idx_punishment_logs_item
    ON punishment_logs(punishment_item_id);

CREATE INDEX idx_punishment_logs_violation_date
    ON punishment_logs(violation_date);

CREATE INDEX idx_punishment_logs_user_date
    ON punishment_logs(user_id, violation_date);


-- +goose Down

DROP TABLE IF EXISTS punishment_logs;
DROP TABLE IF EXISTS punishment_items;
DROP TABLE IF EXISTS punishment_categories;
