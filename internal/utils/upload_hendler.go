package utils

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

type (
	UploadResult struct {
		Filename  string
		Path      string
		Size      int64
		MimeType  string
		Extension string
	}

	uploadFile struct {
		pathDirSave  string
		UploadResult []UploadResult
		uploadDMY    time.Time
	}
)

func NewUploadFile() *uploadFile {
	return &uploadFile{
		uploadDMY: time.Now(),
	}
}

func (u *uploadFile) SetPath(pathDirSave string) *uploadFile {
	u.pathDirSave = pathDirSave
	return u
}

func (u *uploadFile) GetMultiparts(form *multipart.Form) (*uploadFile, error) {
	files := form.File
	for _, fileHeaders := range files {
		for _, fileHeader := range fileHeaders {
			file, err := fileHeader.Open()
			if err != nil {
				return nil, err
			}
			defer file.Close()

			extension := filepath.Ext(fileHeader.Filename)
			encodedFilename := u.md5Filename(fileHeader.Filename)

			uploadDir := u.pathDirSave
			if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
				err := os.MkdirAll(uploadDir, os.ModePerm)
				if err != nil {
					return nil, err
				}
			}

			path := filepath.Join(uploadDir, encodedFilename+extension)
			dst, err := os.Create(path)
			if err != nil {
				return nil, err
			}
			defer dst.Close()

			if _, err = io.Copy(dst, file); err != nil {
				return nil, err
			}

			u.UploadResult = append(u.UploadResult, UploadResult{
				Filename:  encodedFilename,
				Path:      path,
				Size:      fileHeader.Size,
				MimeType:  fileHeader.Header.Get("Content-Type"),
				Extension: extension,
			})
		}
	}

	return u, nil
}

func (u *uploadFile) GetResultUploader() []UploadResult {
	return u.UploadResult
}

func (u *uploadFile) md5Filename(filename string) string {
	hash := md5.New()

	currentTime := time.Now()
	dateString := currentTime.Format("20060102")
	timeString := currentTime.Format("150405")

	hash.Write([]byte(filename + timeString + dateString))
	return "f" + timeString + hex.EncodeToString(hash.Sum(nil)) + "_" + dateString
}
