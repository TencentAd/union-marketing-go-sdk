package ams

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"git.code.oa.com/tme-server-component/kg_growth_open/api/sdk"
	"git.code.oa.com/tme-server-component/kg_growth_open/pkg/define"
	"git.code.oa.com/tme-server-component/kg_growth_open/pkg/sdk/config"
	"github.com/antihax/optional"
	log "github.com/sirupsen/logrus"
	amsAds "github.com/tencentad/marketing-api-go-sdk/pkg/ads"
	"github.com/tencentad/marketing-api-go-sdk/pkg/api"
	amsConfig "github.com/tencentad/marketing-api-go-sdk/pkg/config"
)

// HTTPServer 登录授权服务
type AuthService struct {
	config       *config.Config
	amsSDKClient *amsAds.SDKClient
}

// NewAuthService
func NewAuthService(config *config.Config) *AuthService {
	s := &AuthService{
		config: config,
	}

	if err := s.Init(); err != nil {
		log.Errorf("failed to init AuthService, err: %v", err)
	}

	return s
}

func (s *AuthService) Init() error {
	s.amsSDKClient = amsAds.Init(&amsConfig.SDKConfig{})
	return nil
}

// ServeAMS 根据Authorization Code获取账户信息
func (s *AuthService) ServeAuth(w http.ResponseWriter, req *http.Request) {
	authConf := s.config.Auth
	if authConf == nil {
		serveErrorResponse(w, fmt.Errorf("auth no ams config"))
		return
	}

	query := req.URL.Query()
	authorizationCode := query.Get("authorization_code")
	if authorizationCode == "" {
		serveErrorResponse(w, fmt.Errorf("'authorization_code' parameter not exist"))
		return
	}

	amsResp, _, err := s.amsSDKClient.Oauth().Token(
		context.Background(), authConf.ClientID, authConf.ClientSecret, "authorization_code",
		&api.OauthTokenOpts{
			AuthorizationCode: optional.NewString(authorizationCode),
			RedirectUri: optional.NewString("http://dev.ug.com/ams"),
		})
	if err != nil {
		serveErrorResponse(w, err)
		return
	}

	if amsResp.AuthorizerInfo == nil {
		serveErrorResponse(w, fmt.Errorf("no authorizer info returned"))
		return
	}
	info := amsResp.AuthorizerInfo

	// convert response
	resp := &sdk.AuthResponse{
		Code:    0,
		Message: define.Success,
		Data: &sdk.AuthAccountOutput{
			AccountUin:            info.AccountUin,
			AccountId:             info.AccountId,
			ScopeList:             info.ScopeList,
			WechatAccountId:       info.WechatAccountId,
			AccountRoleType:       AccountRoleTypeMapping[info.AccountRoleType],
			AccountType:           AccountTypeMapping[info.AccountType],
			RoleType:              RoleTypeMapping[info.RoleType],
			AccessToken:           amsResp.AccessToken,
			RefreshToken:          amsResp.RefreshToken,
			AccessTokenExpiresIn:  amsResp.AccessTokenExpiresIn,
			RefreshTokenExpiresIn: amsResp.RefreshTokenExpiresIn,
		},
	}

	serverResponse(w, resp)
}

func serveErrorResponse(w http.ResponseWriter, err error) {
	resp := &sdk.AuthResponse{
		Code:    -1,
		Message: err.Error(),
	}

	serverResponse(w, resp)
}

func serverResponse(w http.ResponseWriter, resp *sdk.AuthResponse) {
	data, _ := json.Marshal(resp)
	_, _ = w.Write(data)
}


