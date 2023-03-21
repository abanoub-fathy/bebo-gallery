package model

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	uuid "github.com/satori/go.uuid"
)

type Image struct {
	GalleryID string
	FileName  string
}

// Path method is used to return the full path to the image
func (i *Image) Path() string {
	return "/" + i.RelativePath()
}

func (i *Image) RelativePath() string {
	return fmt.Sprintf("images/galleries/%v/%v", i.GalleryID, i.FileName)
}

type ImageService interface {
	CreateImage(reader io.ReadCloser, galleryID uuid.UUID, fileName string) error
	GetImagesByGalleryID(galleryID uuid.UUID) ([]Image, error)
	DeleteImage(image *Image) error
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

func (is *imageService) GetImagesByGalleryID(galleryID uuid.UUID) ([]Image, error) {
	imagesDirPath := is.imagesPath(galleryID.String())
	fileNames, err := filepath.Glob(imagesDirPath + "*")
	if err != nil {
		return nil, err
	}
	images := make([]Image, len(fileNames))
	for i := range fileNames {
		fileNames[i] = strings.ReplaceAll(fileNames[i], "\\", "/")
		images[i] = Image{
			GalleryID: galleryID.String(),
			FileName:  strings.ReplaceAll(fileNames[i], is.imagesPath(galleryID.String()), ""),
		}
	}
	return images, nil
}

func (is *imageService) DeleteImage(image *Image) error {
	return os.Remove(image.RelativePath())
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
