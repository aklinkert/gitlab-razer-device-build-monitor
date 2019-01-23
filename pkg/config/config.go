package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

// Config holds the config of the gitlab
type Config struct {
	Groups []string `json:"groups"`
	Repos  []string `json:"repos"`
}

var (
	errFileDoesNotExist = errors.New("given config file does not exist")
)

// Parse takes the given filePath and reads the containing config file into a config struct
func Parse(filePath string) (*Config, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, errFileDoesNotExist
	} else if err != nil {
		return nil, fmt.Errorf("error checking config file: %v", err)
	}

	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %q: %v", filePath, err)
	}

	var cfg Config
	if err := json.Unmarshal(b, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file %q: %v", filePath, err)
	}

	return &cfg, nil
}
