-- +goose Up

-- ============================================================================
-- Reward/Punishment catalog + Maker-Checker records, mirroring the mobile
-- lib/core/constants/reward_punishment_data.dart. Aspects map to the five
-- "Mental Kepribadian" components (MORAL, DISIPLIN, KEPEMIMPINAN,
-- PENGENDALIAN DIRI, PENAMPILAN). Points are signed (rewards +, punishments -).
-- ============================================================================

DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname='rp_type') THEN
        CREATE TYPE rp_type AS ENUM ('REWARD','PUNISHMENT');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname='rp_aspect') THEN
        CREATE TYPE rp_aspect AS ENUM ('MORAL','DISIPLIN','KEPEMIMPINAN','PENGENDALIAN DIRI','PENAMPILAN');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname='rp_record_status') THEN
        CREATE TYPE rp_record_status AS ENUM ('pending','approved','rejected');
    END IF;
END $$;

-- Catalog of rules (seeded below, mirrors the Flutter constant).
CREATE TABLE IF NOT EXISTS reward_punishment_rules (
    code        VARCHAR(20) PRIMARY KEY,
    type        rp_type      NOT NULL,
    aspect      rp_aspect    NOT NULL,
    description TEXT          NOT NULL,
    point       NUMERIC(5,2)  NOT NULL,
    note        VARCHAR(50),
    is_active   BOOLEAN       NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

-- Maker-Checker records: a penalty/reward assigned to a serdik.
CREATE TABLE IF NOT EXISTS reward_punishment_records (
    id               BIGSERIAL PRIMARY KEY,
    serdik_id        BIGINT          NOT NULL REFERENCES serdik(id) ON DELETE CASCADE,
    rule_code        VARCHAR(20)     NOT NULL REFERENCES reward_punishment_rules(code),
    type             rp_type         NOT NULL,
    aspect           rp_aspect       NOT NULL,
    point            NUMERIC(5,2)    NOT NULL,
    description      TEXT,
    status           rp_record_status NOT NULL DEFAULT 'pending',
    created_by       BIGINT          REFERENCES users(id) ON DELETE SET NULL,
    approved_by      BIGINT          REFERENCES users(id) ON DELETE SET NULL,
    reviewed_at      TIMESTAMPTZ,
    rejection_reason TEXT,
    attachment_path  VARCHAR(1000),
    created_at       TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_rpr_serdik        ON reward_punishment_records(serdik_id);
CREATE INDEX IF NOT EXISTS idx_rpr_status        ON reward_punishment_records(status);
CREATE INDEX IF NOT EXISTS idx_rpr_serdik_status ON reward_punishment_records(serdik_id, status);
CREATE INDEX IF NOT EXISTS idx_rpr_aspect        ON reward_punishment_records(aspect);

-- Seed / upsert the catalog (idempotent on re-run).
INSERT INTO reward_punishment_rules (code, type, aspect, description, point, note) VALUES
  ('R_M_01','REWARD','MORAL','Giat Keagamaan Mingguan',0.25,NULL),
  ('R_M_02','REWARD','MORAL','Memimpin Giat Keagamaan (Khutbah)',0.50,NULL),
  ('R_M_03','REWARD','MORAL','Bakti Lembaga',0.25,NULL),
  ('R_M_04','REWARD','MORAL','Memberikan Santunan (panti asuhan, warga dll)',0.40,'MAXS 3X'),
  ('R_M_05','REWARD','MORAL','Memberikan Sumbangan Ke Masjid, Gereja, Pure',0.38,'MAXS 3X'),
  ('R_M_06','REWARD','MORAL','Memberikan Sumbangan Buku Ke Perpustakaan',0.30,'MAXS 3X'),
  ('R_M_07','REWARD','MORAL','Reward Pejabat Upacara',0.30,'MAXS 3X'),
  ('R_M_08','REWARD','MORAL','Reward Pengucap Tribrata',0.30,NULL),
  ('R_M_09','REWARD','MORAL','Reward Pengucap Catur Prasetya',0.30,NULL),
  ('R_M_10','REWARD','MORAL','Reward Dirgen',0.30,NULL),
  ('R_M_11','REWARD','MORAL','Reward Pembaca Doa',0.25,NULL),
  ('R_D_01','REWARD','DISIPLIN','Selama Satu Bulan Tidak Ada Pelanggaran',0.25,NULL),
  ('R_D_02','REWARD','DISIPLIN','Danki Harian',0.35,NULL),
  ('R_D_03','REWARD','DISIPLIN','Pengibar Bendera',0.25,'NON SENAT'),
  ('R_K_01','REWARD','KEPEMIMPINAN','Peserta Didik Yang Menjabat Sebagai Perangkat Senat Mendapatkan Reward Sesuai Keaktifan Dan Kontribusinya Setiap Dua Minggu Sekali Ditandai Dengan Rekomendasi Dari Patun',0.25,NULL),
  ('R_K_02','REWARD','KEPEMIMPINAN','Mampu Menyelesaikan Tugas Dinas Sesuai Sprin Dengan Baik Dan Tuntas',0.25,NULL),
  ('R_K_03','REWARD','KEPEMIMPINAN','Menjadi Tim Perumus Seminar Sekolah',0.25,NULL),
  ('R_K_04','REWARD','KEPEMIMPINAN','Merumuskan Laporan Seminar Sekolah',0.25,NULL),
  ('R_K_05','REWARD','KEPEMIMPINAN','Kunjungan Perpustakaan',0.15,'2X SEMINGGU'),
  ('R_K_06','REWARD','KEPEMIMPINAN','Berani Tampil Menyelesaikan Permasalahan Yang Timbul Di Kalangan Teman-Temannya/Pokjarnya',0.30,NULL),
  ('R_PD_01','REWARD','PENGENDALIAN DIRI','Loyal Untuk Mendukung Segala Bentuk Kegiatan Di Lingkungan Teman-Temannya/Pokjarnya',0.30,NULL),
  ('R_P_01','REWARD','PENAMPILAN','Peduli lingkungan / merawat lingkungan lembaga',0.40,NULL),
  ('P_M_01','PUNISHMENT','MORAL','Tidak Mau Mengakui Kesalahan / Kekurangan',-0.70,NULL),
  ('P_M_02','PUNISHMENT','MORAL','Tidak Memberikan Keterangan Dengan Benar (Bohong)',-0.30,NULL),
  ('P_M_03','PUNISHMENT','MORAL','Tidak Melaksanakan Giat Keagamaan Mingguan',-0.20,NULL),
  ('P_M_04','PUNISHMENT','MORAL','Melanggar Norma Agama (Penistaan Agama)',-0.30,NULL),
  ('P_M_05','PUNISHMENT','MORAL','Tidak Tertib Dalam Melaksanakan / Mengikuti Upacara',-0.50,NULL),
  ('P_M_06','PUNISHMENT','MORAL','Tidak Melaksanakan Perintah Pimpinan',-0.50,NULL),
  ('P_M_07','PUNISHMENT','MORAL','Menjadi Profokator / Memperkeruh Suasana Dalam Kelompok',-0.30,NULL),
  ('P_M_08','PUNISHMENT','MORAL','Melakukan Pembiaran Terhadap Perselisihan Yang Terjadi Dalam Kelompok',-0.30,NULL),
  ('P_D_01','PUNISHMENT','DISIPLIN','Tidak mengisi Daftar Hadir (face id) Setiap Kegiatan Tepat Pada Waktunya',-0.50,NULL),
  ('P_D_02','PUNISHMENT','DISIPLIN','Terlambat Mengikuti Kegiatan / Kuliah Tanpa Alasan Yang Bisa Dipertanggungjawabkan',-0.53,NULL),
  ('P_D_03','PUNISHMENT','DISIPLIN','Terlambat Kembali Pada Waktu Ijin / IBL Tanpa Alasan Yang Dapat Dipertanggungjawabkan (Pengembalian Kartu IBL)',-0.50,NULL),
  ('P_D_04','PUNISHMENT','DISIPLIN','Terlambat Apel Pagi, Malam Dan Olga Pagi',-0.50,NULL),
  ('P_D_05','PUNISHMENT','DISIPLIN','Sengaja Tidak Mengikuti Kegiatan (Kuliah, Giat Pelatihan, Olah Raga, Giat Lain Yang Terjadwal) Tanpa Ijin / Bolos',-0.90,NULL),
  ('P_D_06','PUNISHMENT','DISIPLIN','Sengaja Meninggalkan Kegiatan Sewaktu Kuliah Berlangsung',-0.30,NULL),
  ('P_D_07','PUNISHMENT','DISIPLIN','Melanggar Ketentuan Menerima Tamu',-0.30,NULL),
  ('P_D_08','PUNISHMENT','DISIPLIN','Melanggar Ketentuan Parkir Kendaraan',-0.30,NULL),
  ('P_D_09','PUNISHMENT','DISIPLIN','Melanggar Ketentuan Ijin',-0.30,NULL),
  ('P_D_10','PUNISHMENT','DISIPLIN','Melanggar Ketentuan Di Dormitori, Kamar, Ruang Makan',-0.30,NULL),
  ('P_D_11','PUNISHMENT','DISIPLIN','Merokok Sambil Berjalan, Merokok Tidak Pada Tempatnya',-0.30,NULL),
  ('P_D_12','PUNISHMENT','DISIPLIN','Tidak Mengikuti Kegiatan Pelatihan',-0.80,NULL),
  ('P_D_13','PUNISHMENT','DISIPLIN','Berpakaian Tidak Sesuai Ketentuan',-0.50,NULL),
  ('P_D_14','PUNISHMENT','DISIPLIN','Tidak Membuat/Mengumpulkan Resume Mata Pelajaran',-0.50,NULL),
  ('P_D_15','PUNISHMENT','DISIPLIN','Apabila Setelah Dilaksanakan Pengulangan/Pembuatan Ulang Resume (Remidial) Isi Resume Masih Kurang Tepat/Tidak Sesuai Dengan Materi Mata Pelajaran',-0.50,NULL),
  ('P_K_01','PUNISHMENT','KEPEMIMPINAN','Membiarkan Kelompok Melakukan Pelanggaran',-0.50,NULL),
  ('P_K_02','PUNISHMENT','KEPEMIMPINAN','Mengajak Kelompok Untuk Melakukan Pelanggaran',-0.50,NULL),
  ('P_K_03','PUNISHMENT','KEPEMIMPINAN','Mengkritisi Tidak Proporsional (Asal Bunyi)',-1.00,NULL),
  ('P_K_04','PUNISHMENT','KEPEMIMPINAN','Menyuruh Orang Lain Untuk Mengerjakan Pekerjaan / Tugas Yang Dibebankan Kepadanya (Outsourching)',-0.50,NULL),
  ('P_K_05','PUNISHMENT','KEPEMIMPINAN','Waktu Jam Belajar Di Kelas Digunakan Untuk Kegiatan Lain Yang Tidak Bermanfaat',-0.40,NULL),
  ('P_K_06','PUNISHMENT','KEPEMIMPINAN','Mengganggu / Mencela Kegiatan Tentieur/Pencerahan Kepada Sesama Peserta Didik',-0.50,NULL),
  ('P_K_07','PUNISHMENT','KEPEMIMPINAN','Menggunakan Laptop Atau HP (Gadged) Yang Tidak Terkait Dengan Mata Pelajaran',-0.30,NULL),
  ('P_K_08','PUNISHMENT','KEPEMIMPINAN','Menulis Dan Atau Memviralkan Tulisan, Gambar Dan Video Yang Tidak Pantas Di Medsos',-0.50,NULL),
  ('P_K_09','PUNISHMENT','KEPEMIMPINAN','Acuh Tak Acuh / Masa Bodoh / Tidak Mau Bekerja Sama /Tidak Peduli Terhadap Kepentingan Organisasi',-0.31,NULL),
  ('P_K_10','PUNISHMENT','KEPEMIMPINAN','Melempar Tanggung Jawab',-0.50,NULL),
  ('P_K_11','PUNISHMENT','KEPEMIMPINAN','Menyampaikan Informasi Yang Belum Pasti Kebenarannya atau Tidak Berdasarkan',-0.50,NULL),
  ('P_K_12','PUNISHMENT','KEPEMIMPINAN','Tidak Mau Berpartisipasi Untuk Kepentingan Bersama',-0.50,NULL),
  ('P_K_13','PUNISHMENT','KEPEMIMPINAN','Tidak Mau Diingatkan / Dimotivasi Untuk Patuh',-0.50,NULL),
  ('P_PD_01','PUNISHMENT','PENGENDALIAN DIRI','Tidak Mampu Mengendalikan Amarah, Mudah Tersinggung Dan Sulit Memaafkan',-0.90,NULL),
  ('P_PD_02','PUNISHMENT','PENGENDALIAN DIRI','Berbicara Tidak Sesuai Tugasnya / Porsi Yang Ada Padanya/Celometan',-0.53,NULL),
  ('P_PD_03','PUNISHMENT','PENGENDALIAN DIRI','Mengeluarkan Kata-Kata Kotor / Menyinggung Perasaan Orang Lain Dan Bersikap Arogan',-0.50,NULL),
  ('P_PD_04','PUNISHMENT','PENGENDALIAN DIRI','Sering Mengantuk/Tidur Pada Saat Mengikuti Kegiatan Kuliah/Ceramah',-0.60,NULL),
  ('P_PD_05','PUNISHMENT','PENGENDALIAN DIRI','Selalu Menolak Kritik / Saran / Pendapat Yang Baik Dan Proporsional',-0.30,NULL),
  ('P_PD_06','PUNISHMENT','PENGENDALIAN DIRI','Pemborosan Terhadap Sarpras Dinas',-0.15,NULL),
  ('P_PD_07','PUNISHMENT','PENGENDALIAN DIRI','Tidak Membudayakan Diri Untuk Mandiri/Selalu Minta Dilayani (Membawa Ajudan / Staf Khusus)',-0.15,NULL),
  ('P_P_01','PUNISHMENT','PENAMPILAN','Rambut, Kumis, Jambang, Dan Jenggot Tidak Rapi',-0.25,NULL),
  ('P_P_02','PUNISHMENT','PENAMPILAN','Pakaian, Sepatu, Atribut Tidak Bersih, Tidak Rapi Dan Berwarna Mencolok Serta Tidak Sesuai Dengan Ketentuan',-0.20,NULL),
  ('P_P_03','PUNISHMENT','PENAMPILAN','Tidak Melengkapi Identitas Diri (KTA, KTP Dan SIM) Dan Tidak Menggunakan Kelengkapan Atribut Gampol',-0.25,NULL),
  ('P_P_04','PUNISHMENT','PENAMPILAN','Memakai Accesoris Badan Berlebihan (Akar Bahar, Kalung, Cincin, Anting Dll) Yang Tidak Sesuai Dengan Ketentuan Gampol',-0.30,NULL),
  ('P_P_05','PUNISHMENT','PENAMPILAN','Lantai, Dinding, Kamar Mandi Dan Kelengkapan Kamar Lainnya Kotor, Tidak Rapih Dan Tidak Terawat Dengan Baik',-0.53,NULL),
  ('P_P_06','PUNISHMENT','PENAMPILAN','Memasang / Menempel Foto, Pamplet, Gambar Porno, Dll Pada Dinding Kamar Yang Tidak Sesuai Dengan Ketentuan',-0.30,NULL),
  ('P_P_07','PUNISHMENT','PENAMPILAN','Menjemur Pakaian Tidak Pada Tempatnya',-0.25,NULL),
  ('P_P_08','PUNISHMENT','PENAMPILAN','Tidak Memelihara Dan Menjaga Kebersihan Ruang Kelas Dan Lembaga',-0.40,NULL)
ON CONFLICT (code) DO UPDATE SET
    type        = EXCLUDED.type,
    aspect      = EXCLUDED.aspect,
    description = EXCLUDED.description,
    point       = EXCLUDED.point,
    note        = EXCLUDED.note,
    updated_at  = NOW();

-- +goose Down
DROP TABLE IF EXISTS reward_punishment_records;
DROP TABLE IF EXISTS reward_punishment_rules;
DROP TYPE IF EXISTS rp_record_status;
DROP TYPE IF EXISTS rp_aspect;
DROP TYPE IF EXISTS rp_type;
