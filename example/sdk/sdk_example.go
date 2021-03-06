package main

import (
	"flag"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/tencentad/union-marketing-go-sdk/api/httpx"
	"github.com/tencentad/union-marketing-go-sdk/api/manager"
	"github.com/tencentad/union-marketing-go-sdk/api/sdk"
	"github.com/tencentad/union-marketing-go-sdk/pkg/config"
	"github.com/tencentad/union-marketing-go-sdk/pkg/define"
	"github.com/tencentad/union-marketing-go-sdk/pkg/sdk/account"
	"github.com/tencentad/union-marketing-go-sdk/pkg/sdk/account/mysql"
	sdkConfig "github.com/tencentad/union-marketing-go-sdk/pkg/sdk/config"
	"github.com/tencentad/union-marketing-go-sdk/pkg/sdk/orm"
)

// Config 配置
type Config struct {
	AMS         *sdkConfig.Config `json:"ams"`
	OceanEngine *sdkConfig.Config `json:"ocean_engine"`
	HTTP        *HTTPConfig       `json:"http"`
	DB          *orm.Option       `json:"db"`
}

var conf Config

// HTTPConfig http配置
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
	return nil
}

// serveAuthCallback 提供http接口，在用户授权后获取token信息
func serveAuthCallback(pattern string, impl sdk.MarketingSDK, redirectUrl string) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, req *http.Request) {
		authAccountList, err := impl.ProcessAuthCallback(&sdk.ProcessAuthCallbackInput{
			AuthCallback: req,
		})
		if err != nil {
			httpx.ServeErrorResponse(w, err)
			return
		}
		resp := &httpx.Response{
			Code:    0,
			Message: define.Success,
			Data:    authAccountList,
		}
		httpx.ServerResponse(w, resp)
	})
}

// serveCall 提供http接口，在用户授权后获取token信息
func serveCall(pattern string) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, req *http.Request) {
		query := req.URL.Query()
		method := query["method"][0]
		input := query["input"][0]
		platform := query["platform"][0]

		fmt.Println("method:", method)
		fmt.Println("input:", input)
		fmt.Println("platform:", platform)

		response, err := manager.Call(sdk.MarketingPlatformType(platform), method, input)
		if err != nil {
			httpx.ServeErrorResponse(w, err)
			return
		}
		resp := &httpx.Response{
			Code:    0,
			Message: define.Success,
			Data:    response,
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

	manager.Register(sdk.AMS, conf.AMS)
	amsImpl, _ := manager.GetImpl(sdk.AMS)

	output, err := amsImpl.GenerateAuthURI(&sdk.GenerateAuthURIInput{State: "1"})
	if err != nil {
		log.Errorf("failed to generate auth uri, err: %v", err)
	} else {
		log.Info(output.AuthURI)
	}
	serveAuthCallback("/dashboard/advertiser/callback", amsImpl, conf.AMS.Auth.RedirectUri)
	serveCall("/call")

	// OceanEngine
	manager.Register(sdk.OceanEngine, conf.OceanEngine)
	oceanEgineImpl, err := manager.GetImpl(sdk.OceanEngine)
	if err != nil {
		log.Errorf("failed to get platfrom service, platfrom = %s err: %v", sdk.OceanEngine, err)
	}

	oceanEngineUri, err := oceanEgineImpl.GenerateAuthURI(&sdk.GenerateAuthURIInput{})
	if err != nil {
		log.Errorf("failed to generate auth uri, err: %v", err)
	} else {
		log.Info(oceanEngineUri.AuthURI)
	}

	serveAuthCallback("/ocean_engine", oceanEgineImpl, oceanEgineImpl.GetConfig().Auth.RedirectUri)

	if err := account.Init(mysql.NewTokenStorage(), mysql.NewRefreshLock()); err != nil {
		log.Errorf("failed to init account, err: %v", err)
	}

	if err := http.ListenAndServe(conf.HTTP.ServeAddress, nil); err != nil {
		log.Fatalf("While serving http request: %v", err)
	}

}
