package app

import (
	"encoding/json"
	"os"
)

// Config stores App configuration.
type Config struct {
	Server *ServerConfig `json:"server"`
}

// ServerConfig stores App server configuration.
type ServerConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

// DefaultConfig returns a new Config filled with default values.
func DefaultConfig() *Config {
	return &Config{
		Server: &ServerConfig{
			Host: "127.0.0.1",
			Port: 0,
		},
	}
}

// ReadFile reads Config from a JSON config file.
func (c *Config) ReadFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	d := json.NewDecoder(f)
	return d.Decode(c)
}

// WriteFile writes Config to a JSON config file.
func (c *Config) WriteFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	e := json.NewEncoder(f)
	e.SetIndent("", "  ")
	return e.Encode(c)
}
