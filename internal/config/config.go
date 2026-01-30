package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	// Company details (the entity issuing invoices)
	CompanyName    string `json:"company_name"`
	CompanyGSTIN   string `json:"company_gstin"`
	CompanyAddress string `json:"company_address"`
	CompanyState   string `json:"company_state"`

	// SendGrid
	SendGridAPIKey  string `json:"sendgrid_api_key"`
	SendGridTemplID string `json:"sendgrid_template_id"`
	FromEmail       string `json:"from_email"`

	// Operational
	TestMode         bool   `json:"test_mode"`
	OutputDir        string `json:"output_dir"`
	TemplatePath     string `json:"template_path"`
	PaymentTermsDays int    `json:"payment_terms_days"`
}

func Load() (*Config, error) {
	cfg := &Config{
		CompanyName:      "Salter Technologies Private Limited",
		CompanyGSTIN:     "29ABICS0071M1ZY",
		CompanyAddress:   "T-9 Shirping Chirping Woods, Villament103, Tower-9, Haralur Road, Shubh Enclave, Ambalipura, Bengaluru, Bengaluru Urban, Karnataka, 560102",
		CompanyState:     "Karnataka",
		TestMode:         true,
		OutputDir:        "./output",
		TemplatePath:     "templates/invoice.html",
		PaymentTermsDays: 30,
	}

	// Try config file
	configPath := findConfigFile()
	if configPath != "" {
		data, err := os.ReadFile(configPath)
		if err == nil {
			if err := json.Unmarshal(data, cfg); err != nil {
				return nil, fmt.Errorf("parsing config %s: %w", configPath, err)
			}
		}
	}

	// Env overrides
	if v := os.Getenv("SENDGRID_API_KEY"); v != "" {
		cfg.SendGridAPIKey = v
	}
	if v := os.Getenv("SENDGRID_TEMPLATE_ID"); v != "" {
		cfg.SendGridTemplID = v
	}
	if v := os.Getenv("INVOICE_FROM_EMAIL"); v != "" {
		cfg.FromEmail = v
	}
	if v := os.Getenv("INVOICE_OUTPUT_DIR"); v != "" {
		cfg.OutputDir = v
	}
	if os.Getenv("INVOICE_TEST_MODE") == "false" {
		cfg.TestMode = false
	}

	return cfg, nil
}

// CompanyStateCode returns the first 2 chars of the company GSTIN (state code).
func (c *Config) CompanyStateCode() string {
	if len(c.CompanyGSTIN) >= 2 {
		return c.CompanyGSTIN[:2]
	}
	return "29" // Karnataka default
}

func findConfigFile() string {
	candidates := []string{
		"invoice-esign.json",
		filepath.Join(homeDir(), ".config", "invoice-esign", "config.json"),
	}
	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	return ""
}

func homeDir() string {
	home, _ := os.UserHomeDir()
	return home
}
