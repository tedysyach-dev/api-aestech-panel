package utils

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type UploadConfig struct {
	UploadDir     string
	AllowedExts   []string
	MaxFileSizeMB int64
	BaseURL       string
}

type UploadResult struct {
	FileName string
	Path     string
	Size     int64
	URL      string
}

type fileResult struct {
	result UploadResult
	err    error
}

func UploadFiles(r *http.Request, files []*multipart.FileHeader, config UploadConfig) ([]UploadResult, error) {
	if len(files) == 0 {
		return nil, errors.New("tidak ada file yang dikirim")
	}

	results := make([]UploadResult, 0, len(files))
	resultCh := make(chan fileResult, len(files))

	// Buat folder upload berdasarkan tanggal
	dateFolder := time.Now().Format("2006/01/02")
	uploadPath := filepath.Join(config.UploadDir, dateFolder)
	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		return nil, err
	}

	// Tentukan scheme dan BaseURL jika kosong
	scheme := "http"
	if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	if config.BaseURL == "" {
		config.BaseURL = fmt.Sprintf("%s://%s/uploads", scheme, r.Host)
	}

	// Mulai upload tiap file secara paralel
	for _, fileHeader := range files {
		go func(fh *multipart.FileHeader) {
			ext := strings.ToLower(filepath.Ext(fh.Filename))
			if !isAllowedExt(ext, config.AllowedExts) {
				resultCh <- fileResult{err: fmt.Errorf("file %s memiliki ekstensi tidak diizinkan", fh.Filename)}
				return
			}

			if fh.Size > config.MaxFileSizeMB*1024*1024 {
				resultCh <- fileResult{err: fmt.Errorf("file %s melebihi batas %d MB", fh.Filename, config.MaxFileSizeMB)}
				return
			}

			file, err := fh.Open()
			if err != nil {
				resultCh <- fileResult{err: err}
				return
			}
			defer file.Close()

			newFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
			fullPath := filepath.Join(uploadPath, newFileName)

			out, err := os.Create(fullPath)
			if err != nil {
				resultCh <- fileResult{err: err}
				return
			}

			size, err := io.Copy(out, file)
			out.Close()
			if err != nil {
				resultCh <- fileResult{err: err}
				return
			}

			publicURL := fmt.Sprintf("%s/%s/%s", strings.TrimRight(config.BaseURL, "/"), dateFolder, newFileName)

			resultCh <- fileResult{
				result: UploadResult{
					FileName: newFileName,
					Path:     fullPath,
					Size:     size,
					URL:      publicURL,
				},
			}
		}(fileHeader)
	}

	// Kumpulkan hasil
	for i := 0; i < len(files); i++ {
		res := <-resultCh
		if res.err != nil {
			return nil, res.err // hentikan jika ada error satu saja
		}
		results = append(results, res.result)
	}

	close(resultCh)

	if len(results) == 0 {
		return nil, errors.New("tidak ada file yang berhasil diupload")
	}

	return results, nil
}

func isAllowedExt(ext string, allowed []string) bool {
	ext = strings.ToLower(ext)
	for _, a := range allowed {
		if ext == strings.ToLower(a) {
			return true
		}
	}
	return false
}
