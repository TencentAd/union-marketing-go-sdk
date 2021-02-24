package config

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"

	"github.com/awatercolorpen/nitro/config"
	"github.com/awatercolorpen/nitro/config/source/env"
	"github.com/awatercolorpen/nitro/config/source/file"
	"github.com/awatercolorpen/nitro/config/source/flag"
	log "github.com/sirupsen/logrus"
)

var (
	DefaultDelimiter = "_"
	ENVPodIP         = "POD_IP"
)

// Init 配置加载模块初始
func Init(configFile ...string) error {
	if err := config.Load(env.NewSource()); err != nil {
		return err
	}

	if err := config.Load(flag.NewSource(flag.IncludeUnset(true))); err != nil {
		if err.Error() != "error loading source flag: flags not parsed" {
			return err
		}
		fmt.Println(err)
	}

	for _, filePath := range configFile {
		if filePath == "" {
			continue
		}

		if err := config.Load(file.NewSource(file.WithPath(filePath))); err != nil {
			return err
		}
	}

	return nil
}

// Scan scan config from path to v
// default delimiter is "_"
func Scan(v interface{}, path ...string) error {
	key := strings.Join(path, DefaultDelimiter)
	if err := config.Get(path...).Scan(v); err != nil {
		return fmt.Errorf("load config [%v] err: %v", key, err)
	}

	b, _ := json.Marshal(v)
	log.Infof("load config [%v] : %v", key, string(b))
	return nil
}

func currentIpAddress() string {
	addr, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}

	for _, address := range addr {
		if ip, ok := address.(*net.IPNet); ok && !ip.IP.IsLoopback() {
			if ip.IP.To4() != nil {
				return ip.IP.String()
			}
		}
	}

	return ""
}

// GetIP get current host ip
// support get pod ip from a pod, default env is "POD_IP"
func GetIP() string {
	ip := currentIpAddress()
	return config.Get(ENVPodIP).String(ip)
}

// AssignStringIfNotEmpty if source is not empty then assign source to target
func AssignStringIfNotEmpty(source string, target *string) {
	if source != "" {
		*target = source
	}
}
