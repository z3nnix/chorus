package chorus

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Target struct {
	Deps     []string `yaml:"deps"`
	Cmds     []string `yaml:"cmds"`
	Phony    bool     `yaml:"phony"`
	executed bool
}

type Config struct {
	Variables map[string]string `yaml:"variables"`
	Targets   map[string]Target `yaml:"targets"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	config := &Config{}
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	if config.Variables == nil {
		config.Variables = make(map[string]string)
	}
	config.Variables["DATE"] = time.Now().Format("2006-01-02")
	return config, nil
}
