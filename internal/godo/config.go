package godo

import (
	"log"
	"os"
	"path/filepath"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
)

type Config struct {
	TodoFile string
}

func NewConfig() *Config {
	k := koanf.New(".")

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("failed to determine home directory")
	}

	configFilePath := filepath.Join(homeDir, ".todo", "config.yaml")
	if err := k.Load(file.Provider(configFilePath), yaml.Parser()); err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	return &Config{
		TodoFile: k.String("todoFile"),
	}
}
