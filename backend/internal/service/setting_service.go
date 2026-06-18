package service

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"jewellery-billing/internal/domain"
)

// SettingService handles business logic for shop configuration.
type SettingService struct {
	repo domain.SettingRepository
}

func NewSettingService(repo domain.SettingRepository) *SettingService {
	return &SettingService{repo: repo}
}

// Get returns the current shop settings.
func (s *SettingService) Get(ctx context.Context) (*domain.ShopSettings, error) {
	return s.repo.Get(ctx)
}

// Update overrides the existing settings.
func (s *SettingService) Update(ctx context.Context, req domain.UpdateShopSettingsRequest) (*domain.ShopSettings, error) {
	if req.ShopName == "" {
		return nil, fmt.Errorf("shop name is required")
	}
	if req.InvoicePrefix == "" {
		req.InvoicePrefix = "INV"
	}

	settings := &domain.ShopSettings{
		ShopName:      req.ShopName,
		GSTIN:         req.GSTIN,
		Phone:         req.Phone,
		Address:       req.Address,
		InvoicePrefix: req.InvoicePrefix,
	}

	if err := s.repo.Update(ctx, settings); err != nil {
		return nil, fmt.Errorf("failed to update settings: %w", err)
	}

	return s.repo.Get(ctx)
}

// UploadLogo handles file saving locally and updates the path in the DB.
func (s *SettingService) UploadLogo(ctx context.Context, fileHeader *multipart.FileHeader) (string, error) {
	// Validate file type (basic extension check)
	ext := filepath.Ext(fileHeader.Filename)
	if ext != ".png" && ext != ".jpg" && ext != ".jpeg" {
		return "", fmt.Errorf("only .png, .jpg, and .jpeg files are allowed")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer file.Close()

	// Ensure uploads directory exists
	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %w", err)
	}

	// Generate safe filename
	filename := fmt.Sprintf("logo_%d%s", time.Now().Unix(), ext)
	dstPath := filepath.Join(uploadDir, filename)

	dst, err := os.Create(dstPath)
	if err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	// Store relative path in DB
	dbPath := fmt.Sprintf("/uploads/%s", filename)
	if err := s.repo.UpdateLogo(ctx, dbPath); err != nil {
		return "", fmt.Errorf("failed to update logo path in database: %w", err)
	}

	return dbPath, nil
}
