package env

import "os"

func Get(key string, fallback string) string {
	if os.Getenv(key) == "" {
		return fallback
	}
	return os.Getenv(key)
}
