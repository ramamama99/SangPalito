package validator_test

import (
	"bytes"
	"mime/multipart"
	"net/textproto"
	"testing"

	"sangpalito/backend/pkg/validator"
)

// Minimal magic bytes for image formats
var (
	// PNG signature: 89 50 4E 47 0D 0A 1A 0A
	pngMagic = []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52}
	// JPEG signature: FF D8 FF
	jpegMagic = []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46, 0x49, 0x46, 0x00, 0x01}
	// WebP signature: RIFF....WEBP
	webpMagic = []byte{
		0x52, 0x49, 0x46, 0x46, // "RIFF"
		0x24, 0x00, 0x00, 0x00, // file length
		0x57, 0x45, 0x42, 0x50, // "WEBP"
		0x56, 0x50, 0x38, 0x20, // "VP8 "
	}
)

func createMultipartHeader(filename string, content []byte) *multipart.FileHeader {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", `form-data; name="image"; filename="`+filename+`"`)
	part, _ := writer.CreatePart(h)
	_, _ = part.Write(content)
	_ = writer.Close()

	reader := multipart.NewReader(&body, writer.Boundary())
	form, _ := reader.ReadForm(int64(len(content) + 1024))
	return form.File["image"][0]
}

func TestValidateImageBytes(t *testing.T) {
	t.Run("Valid PNG", func(t *testing.T) {
		err := validator.ValidateImageBytes(pngMagic, "avatar.png")
		if err != nil {
			t.Fatalf("expected nil error for valid png, got %v", err)
		}
	})

	t.Run("Valid JPG", func(t *testing.T) {
		err := validator.ValidateImageBytes(jpegMagic, "photo.jpg")
		if err != nil {
			t.Fatalf("expected nil error for valid jpg, got %v", err)
		}
	})

	t.Run("Valid JPEG", func(t *testing.T) {
		err := validator.ValidateImageBytes(jpegMagic, "photo.jpeg")
		if err != nil {
			t.Fatalf("expected nil error for valid jpeg, got %v", err)
		}
	})

	t.Run("Valid WebP", func(t *testing.T) {
		err := validator.ValidateImageBytes(webpMagic, "banner.webp")
		if err != nil {
			t.Fatalf("expected nil error for valid webp, got %v", err)
		}
	})

	t.Run("Exceeds 2MB limit", func(t *testing.T) {
		largeData := make([]byte, validator.MaxImageSize+1)
		copy(largeData, pngMagic)
		err := validator.ValidateImageBytes(largeData, "huge.png")
		if err != validator.ErrImageTooLarge {
			t.Fatalf("expected ErrImageTooLarge, got %v", err)
		}
	})

	t.Run("Invalid Extension (e.g. .gif)", func(t *testing.T) {
		err := validator.ValidateImageBytes(pngMagic, "animated.gif")
		if err == nil {
			t.Fatal("expected error for gif extension, got nil")
		}
	})

	t.Run("Invalid Extension (e.g. .pdf)", func(t *testing.T) {
		err := validator.ValidateImageBytes(pngMagic, "document.pdf")
		if err == nil {
			t.Fatal("expected error for pdf extension, got nil")
		}
	})

	t.Run("Spoofed extension (text file named .png)", func(t *testing.T) {
		fakeContent := []byte("<html><body>Hello World</body></html>")
		err := validator.ValidateImageBytes(fakeContent, "malicious.png")
		if err == nil {
			t.Fatal("expected error for spoofed png content, got nil")
		}
	})

	t.Run("Empty file", func(t *testing.T) {
		err := validator.ValidateImageBytes([]byte{}, "empty.png")
		if err != validator.ErrEmptyFile {
			t.Fatalf("expected ErrEmptyFile, got %v", err)
		}
	})
}

func TestValidateImageFile(t *testing.T) {
	t.Run("Valid PNG FileHeader", func(t *testing.T) {
		fh := createMultipartHeader("test.png", pngMagic)
		err := validator.ValidateImageFile(fh)
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}
	})

	t.Run("Valid WebP FileHeader", func(t *testing.T) {
		fh := createMultipartHeader("image.webp", webpMagic)
		err := validator.ValidateImageFile(fh)
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}
	})

	t.Run("FileHeader exceeding 2MB", func(t *testing.T) {
		largeData := make([]byte, validator.MaxImageSize+100)
		copy(largeData, jpegMagic)
		fh := createMultipartHeader("big.jpg", largeData)
		err := validator.ValidateImageFile(fh)
		if err != validator.ErrImageTooLarge {
			t.Fatalf("expected ErrImageTooLarge, got %v", err)
		}
	})

	t.Run("Nil FileHeader", func(t *testing.T) {
		err := validator.ValidateImageFile(nil)
		if err != validator.ErrEmptyFile {
			t.Fatalf("expected ErrEmptyFile, got %v", err)
		}
	})
}
