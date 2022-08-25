package godo

import (
	"log"
	"os"
	"path/filepath"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
)

type config struct {
	TodoFile string
}

func newConfig() *config {
	k := koanf.New(".")

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("failed to determine home directory")
	}

	configFilePath := filepath.Join(homeDir, ".todo", "config.yaml")
	if err := k.Load(file.Provider(configFilePath), yaml.Parser()); err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	return &config{
		TodoFile: k.String("todoFile"),
	}
}
