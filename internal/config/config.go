package config

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Config struct {
	mu                   sync.RWMutex
	Port                 string `json:"port"`
	InstagramAppID       string `json:"instagram_app_id"`
	InstagramAppSecret   string `json:"instagram_app_secret"`
	InstagramRedirectURI string `json:"instagram_redirect_uri"`
	TikTokClientKey      string `json:"tiktok_client_key"`
	TikTokClientSecret   string `json:"tiktok_client_secret"`
	TikTokRedirectURI    string `json:"tiktok_redirect_uri"`
	FrontendURL          string `json:"frontend_url"`
	TrustProxy           bool   `json:"trust_proxy"`
}

var AppConfig *Config

func init() {
	AppConfig = &Config{
		Port:                 getEnv("PORT", "8080"),
		InstagramAppID:       getEnv("INSTAGRAM_APP_ID", ""),
		InstagramAppSecret:   getEnv("INSTAGRAM_APP_SECRET", ""),
		InstagramRedirectURI: getEnv("INSTAGRAM_REDIRECT_URI", "http://localhost:8080/api/instagram/auth/callback"),
		TikTokClientKey:      getEnv("TIKTOK_CLIENT_KEY", ""),
		TikTokClientSecret:   getEnv("TIKTOK_CLIENT_SECRET", ""),
		TikTokRedirectURI:    getEnv("TIKTOK_REDIRECT_URI", ""),
		FrontendURL:          getEnv("FRONTEND_URL", "http://localhost:8080"),
		TrustProxy:           getEnvBool("TRUST_PROXY", false),
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
	AppConfig.mu.Lock()
	AppConfig.TikTokClientKey = getEnv("TIKTOK_CLIENT_KEY", AppConfig.TikTokClientKey)
	AppConfig.TikTokClientSecret = getEnv("TIKTOK_CLIENT_SECRET", AppConfig.TikTokClientSecret)
	AppConfig.TikTokRedirectURI = getEnv("TIKTOK_REDIRECT_URI", AppConfig.TikTokRedirectURI)
	AppConfig.Port = getEnv("PORT", AppConfig.Port)
	AppConfig.FrontendURL = getEnv("FRONTEND_URL", AppConfig.FrontendURL)
	AppConfig.TrustProxy = getEnvBool("TRUST_PROXY", AppConfig.TrustProxy)
	AppConfig.mu.Unlock()
}

func (c *Config) Update(appID, appSecret, redirectURI string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if appID != "" {
		c.InstagramAppID = appID
	}
	if appSecret != "" {
		c.InstagramAppSecret = appSecret
	}
	if redirectURI != "" {
		c.InstagramRedirectURI = redirectURI
	}
}

func (c *Config) Get() (appID, appSecret, redirectURI, port string) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.InstagramAppID, c.InstagramAppSecret, c.InstagramRedirectURI, c.Port
}

func (c *Config) GetFrontendURL() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.FrontendURL
}

func (c *Config) GetTikTok() (clientKey, clientSecret, redirectURI string) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.TikTokClientKey, c.TikTokClientSecret, c.TikTokRedirectURI
}

func (c *Config) IsTrustedProxy() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.TrustProxy
}

func (c *Config) Validate() error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	port, err := strconv.Atoi(c.Port)
	if err != nil || port < 1 || port > 65535 {
		return fmt.Errorf("PORT debe ser un número entre 1 y 65535")
	}
	if err := validateHTTPURL("INSTAGRAM_REDIRECT_URI", c.InstagramRedirectURI); err != nil {
		return err
	}
	if c.TikTokRedirectURI != "" {
		if err := validateHTTPURL("TIKTOK_REDIRECT_URI", c.TikTokRedirectURI); err != nil {
			return err
		}
		if !strings.HasPrefix(c.TikTokRedirectURI, "https://") {
			return errors.New("TIKTOK_REDIRECT_URI debe usar https://")
		}
		parsed, _ := url.Parse(c.TikTokRedirectURI)
		if parsed.RawQuery != "" || parsed.Fragment != "" {
			return errors.New("TIKTOK_REDIRECT_URI no puede contener parámetros ni fragmentos")
		}
	}
	return validateHTTPURL("FRONTEND_URL", c.FrontendURL)
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
	if filename := strings.TrimSpace(os.Getenv(key + "_FILE")); filename != "" {
		file, err := os.Open(filename)
		if err == nil {
			defer file.Close()
			value, readErr := io.ReadAll(io.LimitReader(file, 64<<10))
			if readErr == nil {
				if val := strings.TrimSpace(string(value)); val != "" {
					return val
				}
			}
		}
	}
	return defaultVal
}

func getEnvBool(key string, defaultVal bool) bool {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return defaultVal
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return defaultVal
	}
	return parsed
}

func validateHTTPURL(name, value string) error {
	parsed, err := url.Parse(value)
	if err != nil || parsed.Host == "" || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return errors.New(name + " debe ser una URL absoluta http:// o https://")
	}
	return nil
}
