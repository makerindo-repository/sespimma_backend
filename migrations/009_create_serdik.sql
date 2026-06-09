-- +goose Up

CREATE TYPE gender_enum AS ENUM (
    'male',
    'female'
);

CREATE TABLE serdik (
    id BIGSERIAL PRIMARY KEY,

    user_id BIGINT NOT NULL,

    no_serdik VARCHAR(30) UNIQUE,
    nik VARCHAR(20),
    nrp VARCHAR(20),

    nama_lengkap VARCHAR(150),
    tempat_lahir VARCHAR(100),
    tanggal_lahir DATE,

    jenis_kelamin gender_enum DEFAULT 'male',
    agama VARCHAR(20),

    email VARCHAR(150),
    no_telepon VARCHAR(20),
    no_handphone VARCHAR(20),

    alamat TEXT,

    pendidikan_terakhir VARCHAR(50),
    kelompok_kelas VARCHAR(50),
    diktuk_awal VARCHAR(50),
    tahun_diktuk INTEGER,

    is_personel BOOLEAN DEFAULT FALSE,

    pangkat VARCHAR(100),
    jabatan VARCHAR(150),
    satker VARCHAR(150),

    pokjar_id BIGINT NOT NULL,

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),

    CONSTRAINT fk_serdik_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_serdik_pokajar
        FOREIGN KEY (pokjar_id)
        REFERENCES pokjar(id)
        ON DELETE RESTRICT
);

CREATE INDEX idx_serdik_user_id ON serdik(user_id);
CREATE INDEX idx_serdik_pokjar_id ON serdik(pokjar_id);

CREATE TRIGGER trg_serdik_updated_at
BEFORE UPDATE ON serdik
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- +goose Down

DROP TABLE IF EXISTS serdik;
DROP TYPE IF EXISTS gender_enum;

DROP TRIGGER IF EXISTS trg_serdik_updated_at ON serdik;

DROP INDEX IF EXISTS idx_serdik_user_id;
DROP INDEX IF EXISTS idx_serdik_pokjar_id;