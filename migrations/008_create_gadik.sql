-- +goose Up

CREATE TABLE gadik (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE,

    nama VARCHAR(150) NOT NULL,
    pangkat VARCHAR(100) NOT NULL,
    nrp_nip VARCHAR(50) UNIQUE,
    jabatan_struktural VARCHAR(150),

    agama VARCHAR(50),
    eselon VARCHAR(20),
    golongan VARCHAR(20),

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_gadik_user
        FOREIGN KEY (user_id)
        REFERENCES users(id) 
        ON DELETE CASCADE
);

CREATE TRIGGER trg_gadik_updated_at
BEFORE UPDATE ON gadik
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE INDEX idx_gadik_nrp_nip ON gadik(nrp_nip);
CREATE INDEX idx_gadik_pangkat ON gadik(pangkat);

-- +goose Down

DROP TABLE IF EXISTS gadik;