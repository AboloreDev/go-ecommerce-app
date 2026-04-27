package services

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/aboloredev/armory/internal/interfaces"
	"github.com/google/uuid"
)

type UploadService struct {
	provider interfaces.UploadProvider
}

func NewUploadService(provider interfaces.UploadProvider) *UploadService {
	return &UploadService{
		provider: provider,
	}
}

func (s *UploadService) UploadProductImage(productID uint, file *multipart.FileHeader) (string, error) {
	// Get the file extension from filename and set to lowercase
	ext := strings.ToLower(filepath.Ext(file.Filename))

	if !isValidImageExt(ext) {
		return "", fmt.Errorf("This is not a valid image extension %s", ext)
	}

	newFileName := uuid.New().String() + ext

	path := fmt.Sprintf("product/%d/%s", productID, newFileName)

	return s.provider.UploadFile(file, path)
}

func (s *UploadService) DeleteProductImage(productID uint) error {
	path := fmt.Sprintf("product/%d/%s", productID)

	err := s.provider.DeleteFile(path)
	if err != nil {
		return err
	}

	return nil
}

func isValidImageExt(ext string) bool {
	var validExts = []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}

	for _, validExt := range validExts {
		if validExt == ext {
			return true
		}
	}
	return false
}
