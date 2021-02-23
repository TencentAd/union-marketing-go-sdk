package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/antihax/optional"
	amsAds "github.com/tencentad/marketing-api-go-sdk/pkg/ads"
	"github.com/tencentad/marketing-api-go-sdk/pkg/api"
	amsConfig "github.com/tencentad/marketing-api-go-sdk/pkg/config"
)

// HTTPServer 登录授权服务
type HTTPServer struct {
	config       *Config
	amsSDKClient *amsAds.SDKClient
}

func (s *HTTPServer) Init(config *Config) error {
	s.config = config
	s.amsSDKClient = amsAds.Init(&amsConfig.SDKConfig{})
	return nil
}

func (s *HTTPServer) ServeHTTP(w http.ResponseWriter, req *http.Request, ) {
	s.ServeAMS(w, req)
}


type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
