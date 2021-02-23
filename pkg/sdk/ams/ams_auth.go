package ams

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"git.code.oa.com/tme-server-component/kg_growth_open/api/sdk"
	"git.code.oa.com/tme-server-component/kg_growth_open/pkg/define"
	"git.code.oa.com/tme-server-component/kg_growth_open/pkg/sdk/account"
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

	go s.refresh()

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
			RedirectUri:       optional.NewString(authConf.RedirectUri),
		})
	if err != nil {
		serveErrorResponse(w, err)
		return
	}

	if amsResp.AuthorizerInfo == nil {
		serveErrorResponse(w, fmt.Errorf("no authorizer info returned"))
		return
	}

	// convert response
	info := amsResp.AuthorizerInfo
	authAccount := &sdk.AuthAccount{
		AccountUin:           info.AccountUin,
		AccountId:            info.AccountId,
		ScopeList:            *info.ScopeList,
		WechatAccountId:      info.WechatAccountId,
		AccountRoleType:      AccountRoleTypeMapping[info.AccountRoleType],
		AccountType:          AccountTypeMapping[info.AccountType],
		RoleType:             RoleTypeMapping[info.RoleType],
		AccessToken:          amsResp.AccessToken,
		RefreshToken:         amsResp.RefreshToken,
		AccessTokenExpireAt:  calcExpireAt(amsResp.AccessTokenExpiresIn),
		RefreshTokenExpireAt: calcExpireAt(amsResp.RefreshTokenExpiresIn),
	}

	resp := &sdk.AuthResponse{
		Code:    0,
		Message: define.Success,
		Data:    authAccount,
	}

	serverResponse(w, resp)
}

// calcExpireAt 计算失效时间
func calcExpireAt(expireIn int64) time.Time {
	return time.Now().Add(time.Second * time.Duration(expireIn))
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

func (s *AuthService) refresh() {
	authConfig := s.config.Auth
	for {
		time.Sleep(10 * time.Second)
		now := time.Now()
		authAccount := account.ManagerSingleton.GetAllAuthAccount()

		for _, a := range authAccount {
			if a.RefreshTokenExpireAt.Sub(now) <= time.Hour || a.AccessTokenExpireAt.Sub(now) <= time.Hour {
				amsResp, _, err := s.amsSDKClient.Oauth().Token(
					context.Background(), authConfig.ClientID, authConfig.ClientSecret, "refresh_token",
					&api.OauthTokenOpts{
						RefreshToken: optional.NewString(a.RefreshToken),
					})
				if err != nil {
					log.Errorf("failed to call refresh token api for account[%d]", a.AccountId)
				} else {
					if err = account.ManagerSingleton.RefreshToken(&sdk.AuthAccount{
						RefreshToken:         amsResp.RefreshToken,
						AccessToken:          amsResp.AccessToken,
						RefreshTokenExpireAt: calcExpireAt(amsResp.RefreshTokenExpiresIn),
						AccessTokenExpireAt:  calcExpireAt(amsResp.AccessTokenExpiresIn),
					}); err != nil {
						log.Errorf("failed to refresh account[%d] token", a.AccountId)
					}
				}
			}
		}
	}
}
