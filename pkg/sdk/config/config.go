package config

type Config struct {
	Auth *Auth `json:"auth"`
}

type Auth struct {
	ClientID     int64  `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}
