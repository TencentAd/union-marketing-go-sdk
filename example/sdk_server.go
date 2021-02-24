package main

import (
	"flag"
	"fmt"
	"net/http"

	"git.code.oa.com/tme-server-component/kg_growth_open/api/httpx"
	"git.code.oa.com/tme-server-component/kg_growth_open/api/manager"
	"git.code.oa.com/tme-server-component/kg_growth_open/api/sdk"
	"git.code.oa.com/tme-server-component/kg_growth_open/pkg/config"
	"git.code.oa.com/tme-server-component/kg_growth_open/pkg/define"
	"git.code.oa.com/tme-server-component/kg_growth_open/pkg/sdk/account"
	"git.code.oa.com/tme-server-component/kg_growth_open/pkg/sdk/account/mysql"
	"git.code.oa.com/tme-server-component/kg_growth_open/pkg/sdk/ams"
	sdkConfig "git.code.oa.com/tme-server-component/kg_growth_open/pkg/sdk/config"
	"git.code.oa.com/tme-server-component/kg_growth_open/pkg/sdk/orm"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	AMS  *sdkConfig.Config `json:"ams"`
	HTTP *HTTPConfig       `json:"http"`
	DB   *orm.Option       `json:"db"`
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

	if conf.DB != nil {
		if db := orm.GetDB(conf.DB); db == nil {
			return fmt.Errorf("db not init ok")
		}
	}

	if err := account.Init(mysql.NewTokenStorage()); err != nil {
		return fmt.Errorf("failed to init account")
	}

	return nil
}

// serveAuthCallback 提供http接口，在用户授权后获取token信息
func serveAuthCallback(pattern string, impl sdk.MarketingSDK, redirectUrl string) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, req *http.Request) {
		authAccount, err := impl.ProcessAuthCallback(&sdk.ProcessAuthCallbackInput{
			AuthCallback: req,
			RedirectUri:  redirectUrl,
		})
		if err != nil {
			httpx.ServeErrorResponse(w, err)
			return
		}
		resp := &httpx.Response{
			Code:    0,
			Message: define.Success,
			Data:    authAccount,
		}
		httpx.ServerResponse(w, resp)
	})
}

func main() {
	flag.Parse()
	err := Load(*configPath)
	if err != nil {
		log.Fatalf("failed to load config, err: %v", err)
	}

	amsImpl := ams.NewAMSService(conf.AMS)
	manager.Register("ams", amsImpl)
	serveAuthCallback("/ams", amsImpl, conf.AMS.Auth.RedirectUri)

	if err := http.ListenAndServe(conf.HTTP.ServeAddress, nil); err != nil {
		log.Fatalf("While serving http request: %v", err)
	}
}
