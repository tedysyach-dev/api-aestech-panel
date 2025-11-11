package utils

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Mengubah datetime unix millisecond ke dalam bentuk string
func FormatTimestamp(ms int64, timeType string) string {
	if ms == 0 {
		return ""
	}
	if timeType == "date" {
		return time.UnixMilli(ms).Format("2006-01-02")
	} else {
		return time.UnixMilli(ms).Format("2006-01-02 15:04:05")
	}
}

// Mengubah datetime stirng kedalam bentuk unix millisecond
func ParseTimestamp(s string, timeType string) (int64, error) {
	if s == "" {
		return 0, nil
	}

	var t time.Time
	var err error
	if timeType == "date" {
		t, err = time.Parse("2006-01-02", s)
	} else {
		t, err = time.Parse("2006-01-02 15:04:05", s)
	}

	if err != nil {
		return 0, err
	}
	return t.UnixMilli(), nil
}

// GenerateHash menghasilkan hash dari password menggunakan bcrypt.
// Fungsi ini mengembalikan string hash dan error (jika terjadi error).
func GenerateHash(password string) (string, error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashBytes), nil
}

// CompareHash membandingkan password plaintext dengan hashed password.
// Fungsi ini mengembalikan true jika password cocok, dan false jika tidak.
func CompareHash(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
