package config

import "git.code.oa.com/tme-server-component/kg_growth_open/pkg/sdk/http_tools"

// Config SDK配置
type Config struct {
	Auth       *Auth                  `json:"auth"`
	HttpConfig *http_tools.HttpConfig
}

// Auth 授权配置
type Auth struct {
	ClientID     int64  `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectUri  string `json:"redirect_uri"`
}
