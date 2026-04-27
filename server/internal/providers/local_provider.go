package providers

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
)

type LocalUploadProvider struct {
	basePath string
}

func NewLocalUploadProvider(basePath string) *LocalUploadProvider {
	return &LocalUploadProvider{
		basePath: basePath,
	}
}

func (p *LocalUploadProvider) UploadFile(file *multipart.FileHeader, path string) (string, error) {
	fullPath := filepath.Join(p.basePath, path)
	fullFilePath := filepath.Dir(fullPath)

	err := os.Mkdir(fullFilePath, 0750)
	if err != nil {
		return "", err
	}

	// Open the file source
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Create the destination
	dst, err := os.Create(fullFilePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// Read from destination
	_, err = dst.ReadFrom(src)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("/uploads/%s", path), nil
}

func (p *LocalUploadProvider) DeleteFile(path string) error {
	fullPath := filepath.Join(p.basePath, path)

	err := os.Remove(fullPath)
	if err != nil {
		return err
	}

	return nil
}
