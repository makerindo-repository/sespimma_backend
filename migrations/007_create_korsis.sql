-- +goose Up

CREATE TABLE korsis (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE,

    nama VARCHAR(150) NOT NULL,
    pangkat VARCHAR(100) NOT NULL,
    nrp_nip VARCHAR(50) UNIQUE,
    jabatan_struktural VARCHAR(150),
    peran_pengasuhan VARCHAR(150),

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_korsis_user
        FOREIGN KEY (user_id)
        REFERENCES users(id) 
        ON DELETE CASCADE
);

CREATE TRIGGER trg_korsis_updated_at
BEFORE UPDATE ON korsis
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();


CREATE INDEX idx_korsis_nrp_nip ON korsis(nrp_nip);

-- +goose Down

DROP TABLE IF EXISTS korsis;