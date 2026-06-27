package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/tfdriftctl/tfdriftctl/internal/model"
)

// File represents the tfdriftctl.yaml configuration.
type File struct {
	Database   string            `yaml:"database"`
	API        APIConfig         `yaml:"api"`
	Alerting   AlertingConfig    `yaml:"alerting"`
	Workspaces []model.Workspace `yaml:"workspaces"`
}

// APIConfig holds server settings.
type APIConfig struct {
	Addr          string `yaml:"addr"`
	TLSCert       string `yaml:"tls_cert"`
	TLSKey        string `yaml:"tls_key"`
	JWTSecret     string `yaml:"jwt_secret"`
	AdminPassword string `yaml:"admin_password"`
}

// AlertingConfig holds SMTP configuration for email alerts.
type AlertingConfig struct {
	Enabled          bool   `yaml:"enabled"`
	SMTPHost         string `yaml:"smtp_host"`
	SMTPPort         int    `yaml:"smtp_port"`
	Username         string `yaml:"username"`
	Password         string `yaml:"password"`
	From             string `yaml:"from"`
	To               string `yaml:"to"`
	MinimumRiskScore int    `yaml:"minimum_risk_score"`
}

// Load reads configuration from a YAML file.
func Load(path string) (*File, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	var cfg File
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	if cfg.Database == "" {
		cfg.Database = "tfdriftctl.db"
	}
	if cfg.API.Addr == "" {
		cfg.API.Addr = ":8080"
	}
	return &cfg, nil
}
