package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	configDir  = ".aads"
	configFile = "config.yaml"
)

type Config struct {
	ClientID       string `yaml:"client_id"`
	TeamID         string `yaml:"team_id"`
	KeyID          string `yaml:"key_id"`
	OrgID          string `yaml:"org_id"`
	PrivateKeyPath string `yaml:"private_key_path"`
	// DefaultCurrency is used for Money fields when the CLI builds requests from flags.
	// If empty, the CLI attempts to infer it from GET /acls.
	DefaultCurrency string `yaml:"default_currency,omitempty"`
}

func expandHome(path string) string {
	if path == "" {
		return path
	}
	if path == "~" {
		home, _ := os.UserHomeDir()
		return home
	}
	if strings.HasPrefix(path, "~/") {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, strings.TrimPrefix(path, "~/"))
	}
	return path
}

func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("get home dir: %w", err)
	}
	return filepath.Join(home, configDir, configFile), nil
}

func Load() (*Config, error) {
	path, err := configPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		// Allow env-var-only config without requiring a file on disk.
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("read config %s: %w", path, err)
		}
		data = nil
	}

	var cfg Config
	if len(data) > 0 {
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return nil, fmt.Errorf("parse config: %w", err)
		}
	}

	// Env var overrides
	if v := os.Getenv("AADS_CLIENT_ID"); v != "" {
		cfg.ClientID = v
	}
	if v := os.Getenv("AADS_TEAM_ID"); v != "" {
		cfg.TeamID = v
	}
	if v := os.Getenv("AADS_KEY_ID"); v != "" {
		cfg.KeyID = v
	}
	if v := os.Getenv("AADS_ORG_ID"); v != "" {
		cfg.OrgID = v
	}
	if v := os.Getenv("AADS_PRIVATE_KEY_PATH"); v != "" {
		cfg.PrivateKeyPath = v
	}
	if v := os.Getenv("AADS_DEFAULT_CURRENCY"); v != "" {
		cfg.DefaultCurrency = v
	}
	// Alias for convenience.
	if v := os.Getenv("AADS_CURRENCY"); v != "" {
		cfg.DefaultCurrency = v
	}

	cfg.PrivateKeyPath = expandHome(cfg.PrivateKeyPath)

	return &cfg, nil
}

func (c *Config) ValidateAuth() error {
	if c.ClientID == "" {
		return fmt.Errorf("client_id is required")
	}
	if c.TeamID == "" {
		return fmt.Errorf("team_id is required")
	}
	if c.KeyID == "" {
		return fmt.Errorf("key_id is required")
	}
	if c.PrivateKeyPath == "" {
		return fmt.Errorf("private_key_path is required")
	}
	if _, err := os.Stat(c.PrivateKeyPath); err != nil {
		return fmt.Errorf("private key not found at %s: %w", c.PrivateKeyPath, err)
	}
	return nil
}

func (c *Config) Validate() error {
	if err := c.ValidateAuth(); err != nil {
		return err
	}
	if c.OrgID == "" {
		return fmt.Errorf("org_id is required")
	}
	return nil
}

func Save(cfg *Config) error {
	path, err := configPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	return nil
}

func RunInteractiveSetup() error {
	reader := bufio.NewReader(os.Stdin)

	// Try loading existing config for defaults
	existing, _ := Load()
	if existing == nil {
		existing = &Config{}
	}

	cfg := &Config{}

	cfg.ClientID = prompt(reader, "Client ID", existing.ClientID)
	cfg.TeamID = prompt(reader, "Team ID", existing.TeamID)
	cfg.KeyID = prompt(reader, "Key ID", existing.KeyID)
	cfg.OrgID = prompt(reader, "Org ID", existing.OrgID)
	cfg.PrivateKeyPath = prompt(reader, "Private Key Path", existing.PrivateKeyPath)
	cfg.DefaultCurrency = prompt(reader, "Default Currency (optional, e.g. USD)", existing.DefaultCurrency)

	// Expand ~ in path
	if strings.HasPrefix(cfg.PrivateKeyPath, "~/") {
		home, _ := os.UserHomeDir()
		cfg.PrivateKeyPath = filepath.Join(home, cfg.PrivateKeyPath[2:])
	}

	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	if err := Save(cfg); err != nil {
		return err
	}

	path, _ := configPath()
	fmt.Printf("Config saved to %s\n", path)
	return nil
}

func prompt(reader *bufio.Reader, label, defaultVal string) string {
	if defaultVal != "" {
		fmt.Printf("%s [%s]: ", label, defaultVal)
	} else {
		fmt.Printf("%s: ", label)
	}
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		return defaultVal
	}
	return input
}
