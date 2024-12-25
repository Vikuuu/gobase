package main

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const configFileName = "gobase.yml"

type Config struct {
	Database   string `yaml:"database"`
	SchemaData string `yaml:"schema_data"`
	Migration  string `yaml:"migration"`
}

func ReadConfig() (Config, error) {
	filePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	cfg := Config{}
	err = decoder.Decode(&cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func getConfigFilePath() (string, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	configFilePath := filepath.Join(workingDir, configFileName)
	return configFilePath, nil
}
