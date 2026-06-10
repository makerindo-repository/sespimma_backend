package models

import (
    "time"
)

type Serdik struct {
    ID              uint      `gorm:"primaryKey;autoIncrement" json:"id"`

    UserID          uint      `gorm:"not null;index" json:"user_id"`
    User            *User     `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"user,omitempty"`

    NoSerdik        *string   `gorm:"type:varchar(30);unique" json:"no_serdik,omitempty"`
    NIK             *string   `gorm:"type:varchar(20)" json:"nik,omitempty"`
    NRP             *string   `gorm:"type:varchar(20)" json:"nrp,omitempty"`

    NamaLengkap     *string   `gorm:"type:varchar(150)" json:"nama_lengkap,omitempty"`
    TempatLahir     *string   `gorm:"type:varchar(100)" json:"tempat_lahir,omitempty"`
    TanggalLahir    *time.Time `json:"tanggal_lahir,omitempty"`

    JenisKelamin    string    `gorm:"type:gender_enum;default:'male'" json:"jenis_kelamin"`
    Agama           *string   `gorm:"type:varchar(20)" json:"agama,omitempty"`

    Email           *string   `gorm:"type:varchar(150)" json:"email,omitempty"`
    NoTelepon       *string   `gorm:"type:varchar(20)" json:"no_telepon,omitempty"`
    NoHandphone     *string   `gorm:"type:varchar(20)" json:"no_handphone,omitempty"`

    Alamat          *string   `gorm:"type:text" json:"alamat,omitempty"`

    PendidikanTerakhir *string `gorm:"type:varchar(50)" json:"pendidikan_terakhir,omitempty"`
    KelompokKelas   *string   `gorm:"type:varchar(50)" json:"kelompok_kelas,omitempty"`
    DiktukAwal      *string   `gorm:"type:varchar(50)" json:"diktuk_awal,omitempty"`
    TahunDiktuk     *int      `json:"tahun_diktuk,omitempty"`

    IsPersonel      bool      `gorm:"default:false" json:"is_personel"`

    Pangkat         *string   `gorm:"type:varchar(100)" json:"pangkat,omitempty"`
    Jabatan         *string   `gorm:"type:varchar(150)" json:"jabatan,omitempty"`
    Satker          *string   `gorm:"type:varchar(150)" json:"satker,omitempty"`

    Nosis           *string   `gorm:"type:varchar(100)" json:"nosis,omitempty"`
    JabatanSenat    *string   `gorm:"type:varchar(255)" json:"jabatan_senat,omitempty"`
    Angkatan        *string   `gorm:"type:varchar(100)" json:"angkatan,omitempty"`

    PokjarID        uint      `gorm:"not null;index" json:"pokjar_id"`
    Pokjar          *Pokjar   `gorm:"foreignKey:PokjarID;references:ID;constraint:OnDelete:RESTRICT" json:"pokjar,omitempty"`

    CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// Optional: custom table name
func (Serdik) TableName() string {
    return "serdik"
}
