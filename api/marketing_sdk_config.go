package api

import (
	"net/http"
)

// MarketingSDKConfig
type MarketingSDKConfig struct {
	Configuration
	GlobalConfig GlobalConfig
}

// GlobalConfig ...
type GlobalConfig struct {
	ServiceName ServiceName
	HttpOption  HttpOption
}

// ServiceName ...
type ServiceName struct {
	Name   string
	Schema string
}

// HttpOption ...
type HttpOption struct {
	Header http.Header
}
