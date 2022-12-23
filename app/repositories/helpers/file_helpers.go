package helpers

import (
	"errors"
	"mime/multipart"
	"os"
	"outclass-api/app/models"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func UploadFile(c *fiber.Ctx, file *multipart.FileHeader, storedFilePath string) (*models.File, error) {
	fileExt := filepath.Ext(file.Filename)
	fileName := strings.ReplaceAll(uuid.NewString(), "-", "") + fileExt
	filePath := storedFilePath + fileName
	if _, err := os.Stat(storedFilePath); errors.Is(err, os.ErrNotExist) {
		if err := os.Mkdir(storedFilePath, os.ModePerm); err != nil {
			return nil, err
		}
	}

	if err := c.SaveFile(file, filePath); err != nil {
		return nil, err
	}
	serverScheme := os.Getenv("SERVER_SCHEME")
	serverHost := os.Getenv("SERVER_HOST")
	serverPort := os.Getenv("SERVER_PORT")
	server := serverScheme + "://" + serverHost + ":" + serverPort
	link := server + "/api/v1/files/" + fileName
	fileType := strings.Replace(fileExt, ".", "", 1)

	return &models.File{
		Id:   fileName,
		Name: file.Filename,
		Link: link,
		Type: &fileType,
		Size: &file.Size,
	}, nil
}

func DeleteFile(fileId string, storedFilePath string) error {
	filePath := storedFilePath + fileId

	if err := os.Remove(filePath); err != nil {
		return err
	}
	return nil
}
