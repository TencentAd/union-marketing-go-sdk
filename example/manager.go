package main

import (
	"flag"
	"net/http"

	"git.code.oa.com/tme-server-component/kg_growth_open/api/manager"
	"git.code.oa.com/tme-server-component/kg_growth_open/pkg/config"
	"git.code.oa.com/tme-server-component/kg_growth_open/pkg/sdk/ams"
	sdkConfig "git.code.oa.com/tme-server-component/kg_growth_open/pkg/sdk/config"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	AMS  *sdkConfig.Config `json:"ams"`
	HTTP *HTTPConfig       `json:"http"`
}

var conf Config

type HTTPConfig struct {
	ServeAddress string `json:"serve_address"`
}

var (
	configPath = flag.String("config_path", "", "")
)

// Load 加载配置
func Load(configFile ...string) error {
	if err := config.Init(configFile...); err != nil {
		return err
	}

	if err := config.Scan(&conf); err != nil {
		return err
	}
	return nil
}

func main() {
	flag.Parse()
	err := Load(*configPath)
	if err != nil {
		log.Fatalf("failed to load config, err: %v", err)
	}

	amsImpl := ams.NewAMSService(conf.AMS)
	manager.Register("ams", amsImpl)

	if err := http.ListenAndServe(conf.HTTP.ServeAddress, nil); err != nil {
		log.Fatalf("While serving http request: %v", err)
	}
}
