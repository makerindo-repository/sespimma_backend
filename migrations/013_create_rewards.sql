-- +goose Up

-- ENUM types
CREATE TYPE reward_item_period_type AS ENUM (
    'daily',
    'weekly',
    'monthly',
    'semester',
    'unlimited'
);

CREATE TYPE user_reward_status AS ENUM (
    'pending',
    'approved',
    'rejected'
);

-- =========================
-- reward_categories
-- =========================
CREATE TABLE reward_categories (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(100) NOT NULL UNIQUE,
    name VARCHAR(150) NOT NULL,
    bobot NUMERIC(5,2) NOT NULL DEFAULT 0,
    sort_order INT NOT NULL DEFAULT 0,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_reward_categories_code ON reward_categories(code);
CREATE INDEX idx_reward_categories_sort_order ON reward_categories(sort_order);

CREATE TRIGGER trigger_reward_categories_updated_at
BEFORE UPDATE ON reward_categories
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- =========================
-- reward_items
-- =========================
CREATE TABLE reward_items (
    id BIGSERIAL PRIMARY KEY,
    category_id BIGINT NOT NULL,

    item_no INT NOT NULL,
    kegiatan TEXT NOT NULL,
    point NUMERIC(6,2) NOT NULL DEFAULT 0,

    max_count INT DEFAULT NULL,
    period_type reward_item_period_type DEFAULT 'unlimited',
    special_rule VARCHAR(255) DEFAULT NULL,

    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_reward_items_category
        FOREIGN KEY (category_id)
        REFERENCES reward_categories(id)
        ON DELETE CASCADE,

    CONSTRAINT uniq_reward_category_item_no UNIQUE (category_id, item_no)
);

CREATE INDEX idx_reward_items_category_id ON reward_items(category_id);
CREATE INDEX idx_reward_items_active ON reward_items(is_active);

CREATE TRIGGER trigger_reward_items_updated_at
BEFORE UPDATE ON reward_items
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- =========================
-- user_rewards
-- =========================
CREATE TABLE user_rewards (
    id BIGSERIAL PRIMARY KEY,

    user_id BIGINT NOT NULL,
    reward_item_id BIGINT NOT NULL,
    file_id UUID DEFAULT NULL,

    qty INT NOT NULL DEFAULT 1,
    point NUMERIC(8,2) NOT NULL,

    reward_date DATE NOT NULL,

    approved_by BIGINT DEFAULT NULL,
    created_by BIGINT DEFAULT NULL,

    notes TEXT DEFAULT NULL,

    status user_reward_status DEFAULT 'approved',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_user_rewards_item
        FOREIGN KEY (reward_item_id)
        REFERENCES reward_items(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_user_rewards_file
        FOREIGN KEY (file_id)
        REFERENCES files(id)
        ON DELETE SET NULL,

    CONSTRAINT fk_user_rewards_approved_by
        FOREIGN KEY (approved_by)
        REFERENCES users(id)
        ON DELETE SET NULL,

    CONSTRAINT fk_user_rewards_created_by
        FOREIGN KEY (created_by)
        REFERENCES users(id)
        ON DELETE SET NULL
);

CREATE INDEX idx_user_rewards_user_id ON user_rewards(user_id);
CREATE INDEX idx_user_rewards_reward_item_id ON user_rewards(reward_item_id);
CREATE INDEX idx_user_rewards_reward_date ON user_rewards(reward_date);
CREATE INDEX idx_user_rewards_status ON user_rewards(status);
CREATE INDEX idx_user_rewards_user_date ON user_rewards(user_id, reward_date);
CREATE INDEX idx_user_rewards_user_status ON user_rewards(user_id, status);

CREATE TRIGGER trigger_user_rewards_updated_at
BEFORE UPDATE ON user_rewards
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- +goose Down

DROP TRIGGER IF EXISTS trigger_user_rewards_updated_at ON user_rewards;
DROP TRIGGER IF EXISTS trigger_reward_items_updated_at ON reward_items;
DROP TRIGGER IF EXISTS trigger_reward_categories_updated_at ON reward_categories;

DROP TABLE IF EXISTS user_rewards;
DROP TABLE IF EXISTS reward_items;
DROP TABLE IF EXISTS reward_categories;

DROP TYPE IF EXISTS user_reward_status;
DROP TYPE IF EXISTS reward_item_period_type;
