package config

import (
	"fmt"
	"github.com/YasiruR/connector/core/pkg"
	"gopkg.in/yaml.v3"
	"os"
)

const configFile = `config.yaml`

func Load(log pkg.Log) Config {
	var c Config
	file, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatal(fmt.Sprintf("%s not found - %s", configFile, err))
	}

	if err = yaml.Unmarshal(file, &c); err != nil {
		log.Fatal(fmt.Sprintf("unmarshalling %s failed - %s", configFile, err))
	}

	log.Info(fmt.Sprintf("loaded configuration from %s", configFile))
	return c
}
