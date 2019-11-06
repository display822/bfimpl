package services

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

// ComputeMD5 计算文件的md5码
func ComputeMD5(filePath string) (string, error) {
	var result []byte
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum(result)), nil
}
