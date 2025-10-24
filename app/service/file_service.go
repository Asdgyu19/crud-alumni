package service

import (
	"crud-alumni/app/models"
	"crud-alumni/app/repository/mongo"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
)

type FileService struct {
	fileRepo *mongo.FileRepository
}

func NewFileService() *FileService {
	return &FileService{
		fileRepo: mongo.NewFileRepository(),
	}
}

// UploadFile handles file upload request
func (s *FileService) UploadFile(c *fiber.Ctx) error {
	// Get uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "No file uploaded",
		})
	}

	// Validasi ukuran file (maksimal 5MB)
	if file.Size > 5*1024*1024 {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "File size exceeds 5MB limit",
		})
	}

	// Validasi tipe file
	allowedTypes := map[string]bool{
		".pdf":  true,
		".doc":  true,
		".docx": true,
		".jpg":  true,
		".jpeg": true,
		".png":  true,
	}

	ext := filepath.Ext(file.Filename)
	if !allowedTypes[ext] {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "File type not allowed",
		})
	}

	// Buat direktori uploads jika belum ada
	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not create upload directory",
		})
	}

	// Generate nama file unik
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	filepath := fmt.Sprintf("%s/%s", uploadDir, filename)

	// Simpan file
	if err := c.SaveFile(file, filepath); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not save file",
		})
	}

	// Buat record file di database
	fileModel := &models.File{
		FileName:  file.Filename,
		FileType:  ext,
		FilePath:  filepath,
		FileSize:  file.Size,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.fileRepo.Create(fileModel); err != nil {
		// Hapus file jika gagal menyimpan ke database
		os.Remove(filepath)
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not save file information to database",
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "File uploaded successfully",
		"data":    fileModel,
	})
}

// GetAllFiles returns all files
func (s *FileService) GetAllFiles(c *fiber.Ctx) error {
	files, err := s.fileRepo.FindAll()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not retrieve files",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   files,
	})
}

// GetFileByID returns a specific file by ID
func (s *FileService) GetFileByID(c *fiber.Ctx) error {
	id := c.Params("id")

	file, err := s.fileRepo.FindByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "File not found",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   file,
	})
}

// DeleteFile removes a file by ID
func (s *FileService) DeleteFile(c *fiber.Ctx) error {
	id := c.Params("id")

	// Get file info first
	file, err := s.fileRepo.FindByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "File not found",
		})
	}

	// Delete physical file
	if err := os.Remove(file.FilePath); err != nil && !os.IsNotExist(err) {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not delete physical file",
		})
	}

	// Delete from database
	if err := s.fileRepo.Delete(id); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not delete file record from database",
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "File deleted successfully",
	})
}
