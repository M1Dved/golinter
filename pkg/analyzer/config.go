package analyzer

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Rules         map[string]bool `yaml:"rules"`
	ExtraKeywords []string        `yaml:"extra_keywords"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) DisabledRules() map[string]bool {
	disabled := make(map[string]bool)
	if c == nil {
		return disabled
	}
	for rule, enabled := range c.Rules {
		if !enabled {
			disabled[rule] = true
		}
	}
	return disabled
}
