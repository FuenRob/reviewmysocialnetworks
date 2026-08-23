package config

import (
	"bufio"
	"os"
	"strings"
	"sync"
)

type Config struct {
	mu                   sync.RWMutex
	Port                 string `json:"port"`
	InstagramAppID       string `json:"instagram_app_id"`
	InstagramAppSecret   string `json:"instagram_app_secret"`
	InstagramRedirectURI string `json:"instagram_redirect_uri"`
	FrontendURL          string `json:"frontend_url"`
}

var AppConfig *Config

func init() {
	AppConfig = &Config{
		Port:                 getEnv("PORT", "8080"),
		InstagramAppID:       getEnv("INSTAGRAM_APP_ID", ""),
		InstagramAppSecret:   getEnv("INSTAGRAM_APP_SECRET", ""),
		InstagramRedirectURI: getEnv("INSTAGRAM_REDIRECT_URI", "http://localhost:8080/api/auth/callback"),
		FrontendURL:          getEnv("FRONTEND_URL", "http://localhost:8080"),
	}
}

func LoadEnvFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			val = strings.Trim(val, `"'`)
			if os.Getenv(key) == "" {
				_ = os.Setenv(key, val)
			}
		}
	}

	AppConfig.Update(
		getEnv("INSTAGRAM_APP_ID", AppConfig.InstagramAppID),
		getEnv("INSTAGRAM_APP_SECRET", AppConfig.InstagramAppSecret),
		getEnv("INSTAGRAM_REDIRECT_URI", AppConfig.InstagramRedirectURI),
	)
}

func (c *Config) Update(appID, appSecret, redirectURI string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if appID != "" {
		c.InstagramAppID = appID
		_ = os.Setenv("INSTAGRAM_APP_ID", appID)
	}
	if appSecret != "" {
		c.InstagramAppSecret = appSecret
		_ = os.Setenv("INSTAGRAM_APP_SECRET", appSecret)
	}
	if redirectURI != "" {
		c.InstagramRedirectURI = redirectURI
		_ = os.Setenv("INSTAGRAM_REDIRECT_URI", redirectURI)
	}
}

func (c *Config) Get() (appID, appSecret, redirectURI, port string) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.InstagramAppID, c.InstagramAppSecret, c.InstagramRedirectURI, c.Port
}

func (c *Config) IsConfigured() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.InstagramAppID != "" && c.InstagramAppSecret != ""
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
