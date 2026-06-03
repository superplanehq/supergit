package config

import (
	"os"
	"strconv"
	"strings"
)

const (
	DefaultRoot           = "/var/lib/supergit/repos"
	DefaultPort           = "8080"
	DefaultBranch         = "main"
	DefaultMaxFileBytes   = 10 * 1024 * 1024
	DefaultMaxCommitBytes = 25 * 1024 * 1024
)

type Config struct {
	Root           string
	Port           string
	PublicURL      string
	DefaultBranch  string
	MaxFileBytes   int64
	MaxCommitBytes int64
	ReservedPaths  []string
}

func Load() Config {
	return Config{
		Root:           loadString("SUPERGIT_ROOT", DefaultRoot),
		Port:           loadString("SUPERGIT_PORT", DefaultPort),
		PublicURL:      strings.TrimRight(loadString("SUPERGIT_PUBLIC_URL", ""), "/"),
		DefaultBranch:  loadString("SUPERGIT_DEFAULT_BRANCH", DefaultBranch),
		MaxFileBytes:   loadInt64("SUPERGIT_MAX_FILE_BYTES", DefaultMaxFileBytes),
		MaxCommitBytes: loadInt64("SUPERGIT_MAX_COMMIT_BYTES", DefaultMaxCommitBytes),
		ReservedPaths:  loadStringList("SUPERGIT_RESERVED_PATHS"),
	}
}

func loadString(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func loadInt64(key string, fallback int64) int64 {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil || parsed <= 0 {
		return fallback
	}

	return parsed
}

func loadStringList(key string) []string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return nil
	}

	parts := strings.Split(value, ",")
	values := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			values = append(values, part)
		}
	}

	return values
}
