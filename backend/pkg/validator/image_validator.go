package validator

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
)

const (
	// MaxImageSize represents maximum allowed image size (2 MB)
	MaxImageSize = 2 * 1024 * 1024

	// MaxImageSizeMB is the human readable limit
	MaxImageSizeMB = 2
)

var (
	ErrImageTooLarge     = fmt.Errorf("ukuran gambar melebihi batas maksimum %dMB", MaxImageSizeMB)
	ErrInvalidImageType  = errors.New("format gambar tidak didukung: hanya file .webp, .png, dan .jpg/.jpeg yang diperbolehkan")
	ErrEmptyFile         = errors.New("file gambar tidak boleh kosong")
	ErrCannotReadFile    = errors.New("gagal membaca file gambar")

	// AllowedExtensions contains the supported file extensions
	AllowedExtensions = map[string]bool{
		".webp": true,
		".png":  true,
		".jpg":  true,
		".jpeg": true,
	}

	// AllowedMIMETypes contains the supported MIME types detected from content
	AllowedMIMETypes = map[string]bool{
		"image/webp": true,
		"image/png":  true,
		"image/jpeg": true,
	}
)

// ValidateImageFile validates a multipart.FileHeader for size (<= 2MB) and extension/MIME (webp, png, jpg).
func ValidateImageFile(fh *multipart.FileHeader) error {
	if fh == nil {
		return ErrEmptyFile
	}

	// 1. Validate file size (Max 2MB)
	if fh.Size <= 0 {
		return ErrEmptyFile
	}
	if fh.Size > MaxImageSize {
		return ErrImageTooLarge
	}

	// 2. Validate file extension
	ext := strings.ToLower(filepath.Ext(fh.Filename))
	if !AllowedExtensions[ext] {
		return ErrInvalidImageType
	}

	// 3. Validate actual content using magic bytes (MIME sniffing) to prevent spoofing
	file, err := fh.Open()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrCannotReadFile, err)
	}
	defer file.Close()

	headerBytes := make([]byte, 512)
	n, err := file.Read(headerBytes)
	if err != nil && err != io.EOF {
		return fmt.Errorf("%w: %v", ErrCannotReadFile, err)
	}

	detectedMIME := http.DetectContentType(headerBytes[:n])
	// Some systems may report image/x-png or variations; normalize
	detectedMIME = strings.ToLower(strings.Split(detectedMIME, ";")[0])

	if !AllowedMIMETypes[detectedMIME] {
		return fmt.Errorf("%w (terdeteksi: %s)", ErrInvalidImageType, detectedMIME)
	}

	return nil
}

// ValidateImageBytes validates raw image bytes and optional filename for size (<= 2MB) and MIME type.
func ValidateImageBytes(data []byte, filename string) error {
	if len(data) == 0 {
		return ErrEmptyFile
	}

	// 1. Validate size
	if len(data) > MaxImageSize {
		return ErrImageTooLarge
	}

	// 2. Validate filename extension if filename is provided
	if filename != "" {
		ext := strings.ToLower(filepath.Ext(filename))
		if !AllowedExtensions[ext] {
			return ErrInvalidImageType
		}
	}

	// 3. Validate MIME type using magic bytes
	headerLen := 512
	if len(data) < headerLen {
		headerLen = len(data)
	}
	detectedMIME := strings.ToLower(strings.Split(http.DetectContentType(data[:headerLen]), ";")[0])

	if !AllowedMIMETypes[detectedMIME] {
		return fmt.Errorf("%w (terdeteksi: %s)", ErrInvalidImageType, detectedMIME)
	}

	return nil
}
