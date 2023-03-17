package model

import (
	"fmt"
	"io"
	"os"

	uuid "github.com/satori/go.uuid"
)

type ImageService interface {
	CreateImage(reader io.ReadCloser, galleryID uuid.UUID, fileName string) error
}

type imageService struct{}

// make sure that imageService type implements ImageService interface
var _ ImageService = (*imageService)(nil)

func NewImageService() ImageService {
	return &imageService{}
}

func (is *imageService) CreateImage(reader io.ReadCloser, galleryID uuid.UUID, fileName string) error {
	defer reader.Close()

	// create image dir path
	imagePath, err := createImageDirPath(galleryID.String())
	if err != nil {
		return err
	}

	// create destination file
	destinationFile, err := os.Create(imagePath + fileName)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	// copy the uploaded file to destination file
	_, err = io.Copy(destinationFile, reader)
	if err != nil {
		return err
	}

	return nil
}

func createImageDirPath(galleryID string) (string, error) {
	imageDirPath := fmt.Sprintf("images/galleries/%v/", galleryID)
	err := os.MkdirAll(imageDirPath, 0755)
	if err != nil {
		return "", err
	}
	return imageDirPath, nil
}
