package ustring

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	// Chuyển đổi byte ngẫu nhiên thành chuỗi base64 để dễ sử dụng
	return base64.URLEncoding.EncodeToString(bytes), nil
}
