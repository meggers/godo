package godo

import (
	"log"
	"os"
	"path/filepath"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"

	util "github.com/meggers/godo/internal"
)

type Config struct {
	TodoFile string
}

func NewConfig() *Config {
	k := koanf.New(".")

	wd, err := os.Getwd()
	util.CheckError(err, "failed to get current working directory")

	devConfigPath := filepath.Join(wd, "configs/config.yaml")
	if err := k.Load(file.Provider(devConfigPath), yaml.Parser()); err != nil {
		log.Printf("failed to pickup dev config: %v", err)
	}

	homeDir, err := os.UserHomeDir()
	util.CheckError(err, "failed to determine home directory")

	configFilePath := filepath.Join(homeDir, ".todo", "config.yaml")
	if err := k.Load(file.Provider(configFilePath), yaml.Parser()); err != nil {
		log.Printf("error loading config, falling back to default: %v", err)
	}

	return &Config{
		TodoFile: k.String("todoFile"),
	}
}
