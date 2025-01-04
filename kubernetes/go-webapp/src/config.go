package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type ConfigLoadResult struct {
	Config *Config
	Errors []error
}

func loadConfigPathWrapper(filepath string) (loadResult *ConfigLoadResult) {
	loadResult = &ConfigLoadResult{}
	configFile, err := os.Open(filepath)
	if err != nil {
		loadResult.Errors = append(loadResult.Errors, fmt.Errorf("Error opening file: %v", err))
		return loadResult
	}
	loadResult = loadConfigFile(configFile)
	return loadResult
}

func loadConfigPath(filepath string) (config *Config, loadErr error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		loadErr = fmt.Errorf("error reading from file %s\n\terror: %w", filepath, err)
	}
	config = &Config{}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		loadErr = errors.Join(loadErr, fmt.Errorf("error unmarshalling config from file: %s: %w", filepath, err))
	}
	return config, loadErr
}

func loadConfigFile(file *os.File) (loadResult *ConfigLoadResult) {
	loadResult = &ConfigLoadResult{}
	data, err := io.ReadAll(file)
	if err != nil {
		loadResult.Errors = append(loadResult.Errors, fmt.Errorf("error reading from file %s\n\terror: %w", file.Name(), err))
	}
	loadResult.Config = &Config{}
	yamlErr := yaml.Unmarshal(data, loadResult.Config)
	if yamlErr != nil {
		loadResult.Errors = append(loadResult.Errors, fmt.Errorf("error loading yaml config from file: %s: %w", file.Name(), yamlErr))
		jsonErr := json.Unmarshal(data, loadResult.Config)
		if jsonErr != nil {
			loadResult.Errors = append(loadResult.Errors, fmt.Errorf("error loading json config from file: %s: %w", file.Name(), jsonErr))
		}
	}
	return loadResult
}
