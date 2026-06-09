-- +goose Up
CREATE TABLE kegiatan (
    id BIGSERIAL PRIMARY KEY,
    created_id BIGINT NOT NULL,
    nama VARCHAR(255) NOT NULL,
    nama_lokasi VARCHAR(255) NOT NULL,
    tanggal_mulai DATE NOT NULL,
    waktu_mulai TIME NOT NULL,
    waktu_selesai TIME NOT NULL,
    batas_waktu_penugasan TIME NOT NULL,
    is_rutin BOOLEAN DEFAULT FALSE,
    is_pelatihan BOOLEAN DEFAULT FALSE,
    qr_code BOOLEAN DEFAULT FALSE,
    location GEOGRAPHY(Geometry, 4326) NOT NULL,
    radius DOUBLE PRECISION NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_kegiatan_created_id FOREIGN KEY (created_id)
        REFERENCES users(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE
);

CREATE TRIGGER trigger_kegiatan_updated_at
BEFORE UPDATE ON kegiatan
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- +goose Down
DROP TRIGGER IF EXISTS trigger_kegiatan_updated_at ON kegiatan;
DROP TABLE IF EXISTS kegiatan;
