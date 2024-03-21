package tools

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

// MakeDir return path and error
func MakeDir(file multipart.File, dirName string, fileName string) (string, error) {
	err := os.MkdirAll(dirName+"_files", 0666)
	if err != nil {
		return "", err
	}

	fullFileName := fmt.Sprintf("%v", time.Now().Unix())

	path := filepath.Join(dirName+"_files", fullFileName+"."+fileName)

	targetFile, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return "", err
	}
	defer targetFile.Close()

	_, err = io.Copy(targetFile, file)
	if err != nil {
		return "", err
	}

	return path, nil
}
