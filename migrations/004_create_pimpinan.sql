-- +goose Up

CREATE TABLE pimpinan (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE,

    nama VARCHAR(150) NOT NULL,
    pangkat VARCHAR(100) NOT NULL,
    nrp_nip VARCHAR(50) NOT NULL UNIQUE,
    jabatan_struktural VARCHAR(255) NOT NULL,
    peran_pengasuhan VARCHAR(100) NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_pimpinan_user
        FOREIGN KEY (user_id)
        REFERENCES users(id) 
        ON DELETE CASCADE
);

CREATE TRIGGER trg_pimpinan_updated_at
BEFORE UPDATE ON pimpinan
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE INDEX idx_pimpinan_nrp_nip ON pimpinan(nrp_nip);

-- +goose Down

DROP TABLE IF EXISTS pimpinan;