package common

import (
	"context"
	"fmt"
	"git.code.oa.com/tme-server-component/kg_growth_open/api"
	"net/http"
)

// ThirdPartyClient ...
type ThirdPartyClient struct {
	http.RoundTripper
	Config         *api.MarketingSDKConfig
	Client         *APIClient  // 网络请求接口
	Ctx            *context.Context
}

// Init ...
func Init(cfg *api.MarketingSDKConfig) *ThirdPartyClient {
	ctx := context.Background()
	client := NewAPIClient(cfg)
	ThirdPartyClient := &ThirdPartyClient{
		Config:       cfg,
		Ctx:          &ctx,
		Client:       client,
	}
	ThirdPartyClient.Client.Cfg.HTTPClient.Transport = ThirdPartyClient
	return ThirdPartyClient
}

func (tads *ThirdPartyClient) SetIpPort(ip string, port int, schema string) {
	ipPort := fmt.Sprintf("%s:%d", ip, port)
	tads.SetHost(ipPort, schema)
}

// SetHost ...
func (tads *ThirdPartyClient) SetHost(host string, schema string) *ThirdPartyClient {
	modified := false
	if host != "" {
		tads.Client.Cfg.Host = host
		modified = true
	}
	if schema != "" {
		tads.Client.Cfg.Scheme = schema
		modified = true
	}
	if modified {
		tads.Client.Cfg.BasePath = tads.Client.Cfg.Scheme + "://" + tads.Client.Cfg.Host
	}
	return tads
}

// SetHeaders ...
func (tads *ThirdPartyClient) SetHeaders(header http.Header) *ThirdPartyClient {
	tads.Config.GlobalConfig.HttpOption.Header = header
	return tads
}

// SetHeader ...
func (tads *ThirdPartyClient) SetHeader(key string, value string) *ThirdPartyClient {
	if tads.Config.GlobalConfig.HttpOption.Header == nil {
		tads.Config.GlobalConfig.HttpOption.Header = http.Header{}
	}
	tads.Config.GlobalConfig.HttpOption.Header.Set(key, value)
	return tads
}


// SetNameService ...

