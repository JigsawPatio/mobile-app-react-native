package mobileappreactnative

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// LoadEnv loads environment variables from a .env file.
func LoadEnv() error {
	err := godotenv.Load()
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("error loading .env file: %w", err)
	}
	return nil
}

// GenerateRandomString returns a securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomString(n int) (string, error) {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}
	return string(bytes), nil
}

// GenerateSecureToken generates a random token using base64 encoding.
func GenerateSecureToken(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// GetEnvOrDefault retrieves an environment variable or returns a default value if the variable is not set.
func GetEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// FileExists checks if a file exists at the given path.
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

// CreateDirectoryIfNotExists creates a directory if it doesn't exist.
func CreateDirectoryIfNotExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, os.ModeDir|0755) // Read, write, and execute for owner; read and execute for group and others.
		if err != nil {
			return fmt.Errorf("failed to create directory %s: %w", path, err)
		}
	}
	return nil
}

// GetExecutableDir returns the directory of the currently running executable.
func GetExecutableDir() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %w", err)
	}
	return filepath.Dir(ex), nil
}

// LogRequest logs the details of an HTTP request.
func LogRequest(r *http.Request) {
	log.Printf("Request: Method=%s, URL=%s, RemoteAddr=%s, UserAgent=%s",
		r.Method, r.URL.String(), r.RemoteAddr, r.UserAgent())
}

// SanitizeFilename removes potentially harmful characters from a filename.
func SanitizeFilename(filename string) string {
	// Replace spaces with underscores
	filename = strings.ReplaceAll(filename, " ", "_")

	// Remove characters that are not alphanumeric, underscores, or periods
	var result strings.Builder
	for _, r := range filename {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' || r == '.' {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// Retry executes a function multiple times with a delay between retries.
func Retry(attempts int, sleep time.Duration, f func() error) (err error) {
	for i := 0; i < attempts; i++ {
		err = f()
		if err == nil {
			return nil
		}

		log.Printf("attempt %d failed: %v", i+1, err)
		time.Sleep(sleep)
	}
	return fmt.Errorf("after %d attempts, the last error was: %s", attempts, err)
}