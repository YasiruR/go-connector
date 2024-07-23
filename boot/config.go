package boot

import (
	"fmt"
	"github.com/YasiruR/connector/core"
	"github.com/YasiruR/connector/core/pkg"
	"gopkg.in/yaml.v3"
	"os"
)

const configFile = `config.yaml`

func loadConfig(log pkg.Log) core.Config {
	var c core.Config
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
