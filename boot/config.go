package boot

import (
	"fmt"
	"github.com/YasiruR/go-connector/domain/boot"
	"github.com/YasiruR/go-connector/domain/pkg"
	"gopkg.in/yaml.v3"
	"net"
	"os"
)

func loadConfig(log pkg.Log) boot.Config {
	var c boot.Config
	file, err := os.ReadFile(boot.FilePath)
	if err != nil {
		log.Fatal(fmt.Sprintf("%s not found - %s", boot.FilePath, err))
	}

	if err = yaml.Unmarshal(file, &c); err != nil {
		log.Fatal(fmt.Sprintf("unmarshalling %s failed - %s", boot.FilePath, err))
	}

	pwd, err := os.Getwd()
	if err != nil {
		log.Error(fmt.Sprintf("could not fetch current working directory - %s", err))
	}

	ipAddr := `http://` + outboundIP().String()
	c.Servers.IP = ipAddr

	log.Info("loaded configuration values", `file: `+pwd+boot.FilePath, "ip: "+ipAddr)
	return c
}

func outboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
