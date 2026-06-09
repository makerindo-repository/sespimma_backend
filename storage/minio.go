package storage

import (
    "context"
    "github.com/minio/minio-go/v7"
    "github.com/minio/minio-go/v7/pkg/credentials"
)

func InitMinio() (*minio.Client, error) {
    client, err := minio.New("localhost:9000", &minio.Options{
        Creds:  credentials.NewStaticV4("minioadmin", "minioadmin", ""),
        Secure: false,
    })
    return client, err
}

func UploadFile(client *minio.Client, bucket, objectName, filePath string) error {
    _, err := client.FPutObject(
        context.Background(),
        bucket,
        objectName,
        filePath,
        minio.PutObjectOptions{ContentType: "image/jpeg"},
    )
    return err
}
