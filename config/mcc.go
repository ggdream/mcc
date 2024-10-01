package config

import "gopkg.in/yaml.v3"

type MCCConfig struct {
	Type    string   `yaml:"type"`
	Steps   []string `yaml:"steps"`
	Apply   string   `yaml:"apply"`
	Command []string `yaml:"command"`
}

func ParseMCCConfig(data []byte) (conf *MCCConfig, err error) {
	err = yaml.Unmarshal(data, &conf)
	return
}
