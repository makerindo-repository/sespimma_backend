package models

import (
    "time"
    "github.com/google/uuid"
)

type File struct {
    ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    UserID    uuid.UUID
    FileType  string
    Path      string
    MimeType  string
    Size      int64
    CreatedAt time.Time `gorm:"autoCreateTime"`
}
