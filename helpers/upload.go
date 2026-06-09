package helpers

import (
    "github.com/disintegration/imaging"
)

func CreateThumbnail(srcPath, dstPath string) error {
    img, err := imaging.Open(srcPath)
    if err != nil {
        return err
    }
    thumb := imaging.Thumbnail(img, 200, 200, imaging.Lanczos)
    return imaging.Save(thumb, dstPath)
}
