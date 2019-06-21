package app

import (
	"encoding/json"
	"os"
)

// Config stores App configuration.
type Config struct {
	Server *ServerConfig `json:"server"`
	Tor    *TorConfig    `json:"tor,omitempty"`
}

// ServerConfig stores App server configuration.
type ServerConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

// TorConfig stores tor configuration.
type TorConfig struct {
	Controller     *TorControllerConfig `json:"controller"`
	PrivateKeyFile string               `json:"private_key_file,omitempty"`
}

// TorControllerConfig stores tor controller configuration.
type TorControllerConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password,omitempty"`
}

// DefaultConfig returns a new Config filled with default values.
func DefaultConfig() *Config {
	return &Config{
		Server: &ServerConfig{
			Host: "127.0.0.1",
			Port: 0,
		},
		Tor: &TorConfig{
			Controller: &TorControllerConfig{
				Host: "127.0.0.1",
				Port: 9051,
			},
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
