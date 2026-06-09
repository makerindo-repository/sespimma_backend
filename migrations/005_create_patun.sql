-- +goose Up

CREATE TABLE patun (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE,

    nama VARCHAR(150) NOT NULL,
    pangkat VARCHAR(150),
    nrp_nip VARCHAR(50) UNIQUE NOT NULL,
    jabatan_struktural VARCHAR(150),
    peran_pengasuhan VARCHAR(150),
    pokjar VARCHAR(50),

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_patun_user
        FOREIGN KEY (user_id)
        REFERENCES users(id) 
        ON DELETE CASCADE
);

CREATE TRIGGER trg_patun_updated_at
BEFORE UPDATE ON patun
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE INDEX idx_patun_nrp_nip ON patun(nrp_nip);
CREATE INDEX idx_patun_pokjar ON patun(pokjar);

-- +goose Down

DROP TABLE IF EXISTS patun;