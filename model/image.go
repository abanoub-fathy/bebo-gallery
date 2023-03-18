package model

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	uuid "github.com/satori/go.uuid"
)

type ImageService interface {
	CreateImage(reader io.ReadCloser, galleryID uuid.UUID, fileName string) error
	GetImagesByGalleryID(galleryID uuid.UUID) ([]string, error)
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
	imagePath, err := is.createImageDirPath(galleryID.String())
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

func (is *imageService) GetImagesByGalleryID(galleryID uuid.UUID) ([]string, error) {
	imagesDirPath := is.imagesPath(galleryID.String())
	fileNames, err := filepath.Glob(imagesDirPath + "*")
	if err != nil {
		return nil, err
	}
	for i := range fileNames {
		// replace the "\" char and add the "/" to each fileName
		fileNames[i] = "/" + strings.ReplaceAll(fileNames[i], "\\", "/")
	}
	return fileNames, nil
}

func (is *imageService) imagesPath(galleryID string) string {
	return fmt.Sprintf("images/galleries/%v/", galleryID)
}

func (is *imageService) createImageDirPath(galleryID string) (string, error) {
	imageDirPath := is.imagesPath(galleryID)
	err := os.MkdirAll(imageDirPath, 0755)
	if err != nil {
		return "", err
	}
	return imageDirPath, nil
}
