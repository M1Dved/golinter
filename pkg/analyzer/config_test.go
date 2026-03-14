package analyzer

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	content := `
rules:
  lowercase: true
  english: true
  special: false
  sensitive: true

extra_keywords:
  - ssn
  - credit_card
`

	tmpFile, err := os.CreateTemp("", "golinter-*.yml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	cfg, err := LoadConfig(tmpFile.Name())
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	disabled := cfg.DisabledRules()
	if disabled["lowercase"] {
		t.Error("lowercase should not be disabled")
	}
	if disabled["english"] {
		t.Error("english should not be disabled")
	}
	if !disabled["special"] {
		t.Error("special should be disabled")
	}
	if disabled["sensitive"] {
		t.Error("sensitive should not be disabled")
	}

	if len(cfg.ExtraKeywords) != 2 {
		t.Fatalf("expected 2 extra keywords, got %d", len(cfg.ExtraKeywords))
	}
	if cfg.ExtraKeywords[0] != "ssn" {
		t.Errorf("expected first keyword 'ssn', got %q", cfg.ExtraKeywords[0])
	}
	if cfg.ExtraKeywords[1] != "credit_card" {
		t.Errorf("expected second keyword 'credit_card', got %q", cfg.ExtraKeywords[1])
	}
}

func TestLoadConfigFileNotFound(t *testing.T) {
	_, err := LoadConfig("nonexistent.yml")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestLoadConfigOnlyRules(t *testing.T) {
	content := `
rules:
  sensitive: false
`

	tmpFile, err := os.CreateTemp("", "golinter-*.yml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	tmpFile.WriteString(content)
	tmpFile.Close()

	cfg, err := LoadConfig(tmpFile.Name())
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	disabled := cfg.DisabledRules()
	if !disabled["sensitive"] {
		t.Error("sensitive should be disabled")
	}

	if len(cfg.ExtraKeywords) != 0 {
		t.Error("expected no extra keywords")
	}
}

func TestLoadConfigOnlyKeywords(t *testing.T) {
	content := `
extra_keywords:
  - bank_account
  - ssn
`

	tmpFile, err := os.CreateTemp("", "golinter-*.yml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	tmpFile.WriteString(content)
	tmpFile.Close()

	cfg, err := LoadConfig(tmpFile.Name())
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	disabled := cfg.DisabledRules()
	if len(disabled) != 0 {
		t.Error("no rules should be disabled")
	}

	if len(cfg.ExtraKeywords) != 2 {
		t.Fatalf("expected 2 extra keywords, got %d", len(cfg.ExtraKeywords))
	}
}

func TestDisabledRulesNilConfig(t *testing.T) {
	var cfg *Config
	disabled := cfg.DisabledRules()
	if len(disabled) != 0 {
		t.Error("expected empty map for nil config")
	}
}
