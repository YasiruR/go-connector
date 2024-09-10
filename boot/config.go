package boot

import (
	"fmt"
	"github.com/YasiruR/connector/domain/boot"
	"github.com/YasiruR/connector/domain/pkg"
	"gopkg.in/yaml.v3"
	"net"
	"os"
)

const configFile = `config.yaml`

func loadConfig(log pkg.Log) boot.Config {
	var c boot.Config
	file, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatal(fmt.Sprintf("%s not found - %s", configFile, err))
	}

	if err = yaml.Unmarshal(file, &c); err != nil {
		log.Fatal(fmt.Sprintf("unmarshalling %s failed - %s", configFile, err))
	}

	pwd, err := os.Getwd()
	if err != nil {
		log.Error(fmt.Sprintf("could not fetch current working directory - %s", err))
	}

	ipAddr := `http://` + getOutboundIP().String()
	c.Servers.IP = ipAddr

	log.Info("loaded configuration values", `file: `+pwd+configFile, "ip: "+ipAddr)
	return c
}

func getOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
