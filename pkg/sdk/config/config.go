package config

// Config SDK配置
type Config struct {
	Auth *Auth `json:"auth"`
}

// Auth 授权配置
type Auth struct {
	ClientID     int64  `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectUri  string `json:"redirect_uri"`
}

// Server配置
type Server struct {
	BasePath string
	Timeout  int64
	apiVersion string
}
