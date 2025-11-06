package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type ServerConfig struct {
	URL       string   `json:"url"`
	Transport string   `json:"transport"`
	Command   string   `json:"command,omitempty"`
	Args      []string `json:"args,omitempty"`
}

type Config struct {
	DefaultServer string                  `json:"default_server,omitempty"`
	Servers       map[string]ServerConfig `json:"servers"`
}

// Load configuration from file
func Load(path string) (*Config, error) {
	if path == "" {
		path = getDefaultConfigPath()
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{Servers: make(map[string]ServerConfig)}, nil
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	if cfg.Servers == nil {
		cfg.Servers = make(map[string]ServerConfig)
	}

	return &cfg, nil
}

// Save configuration to file
func (c *Config) Save(path string) error {
	if path == "" {
		path = getDefaultConfigPath()
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetServer returns server configuration by name
func (c *Config) GetServer(name string) (ServerConfig, error) {
	if name == "" {
		name = c.DefaultServer
	}

	if name == "" {
		return ServerConfig{}, fmt.Errorf("no server specified and no default server configured")
	}

	server, ok := c.Servers[name]
	if !ok {
		return ServerConfig{}, fmt.Errorf("server '%s' not found in configuration", name)
	}

	return server, nil
}

// AddServer adds or updates a server in the configuration
func (c *Config) AddServer(name string, server ServerConfig) {
	if c.Servers == nil {
		c.Servers = make(map[string]ServerConfig)
	}
	c.Servers[name] = server
}

// getDefaultConfigPath returns the default configuration file path
func getDefaultConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".mcp-config.json"
	}
	return filepath.Join(home, ".mcp-config.json")
}
