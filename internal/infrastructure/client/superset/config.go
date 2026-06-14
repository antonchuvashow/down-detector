package superset

import (
	"fmt"
	"net/url"
	"os"
)

func ConfigFromEnv() (Config, error) {
	baseURL, err := url.Parse(getEnv("SUPERSET_BASE_URL", "http://localhost:8088"))
	if err != nil {
		return Config{}, fmt.Errorf("parse superset base url: %w", err)
	}

	return Config{
		BaseURL:       *baseURL,
		AdminUser:     getEnv("SUPERSET_ADMIN_USER", "admin"),
		AdminPassword: getEnv("SUPERSET_ADMIN_PASSWORD", "admin"),
	}, nil
}

func GuestDescriptorFromEnv() GuestDescriptor {
	return GuestDescriptor{
		Username:  getEnv("SUPERSET_GUEST_USERNAME", "guest"),
		Firstname: getEnv("SUPERSET_GUEST_FIRSTNAME", "guest"),
		Lastname:  getEnv("SUPERSET_GUEST_LASTNAME", "guest"),
	}
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
