package ams

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"git.code.oa.com/tme-server-component/kg_growth_open/api/sdk"
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

	if err := s.init(); err != nil {
		log.Errorf("failed to init AuthService, err: %v", err)
	}

	go s.refresh()

	return s
}

func (s *AuthService) init() error {
	s.amsSDKClient = amsAds.Init(&amsConfig.SDKConfig{
		IsDebug: true,
	})
	return nil
}

// GenerateAuthURI implement Auth
func (s *AuthService) GenerateAuthURI(input *sdk.GenerateAuthURIInput) (*sdk.GenerateAuthURIOutput, error) {
	authConf := s.config.Auth
	if authConf == nil {
		return nil, fmt.Errorf("auth no ams config")
	}
	authUri := fmt.Sprintf("https://developers.e.qq.com/oauth/authorize?client_id=%d&redirect_uri=%s",
		s.config.Auth.ClientID,
		url.QueryEscape(input.RedirectURI),
	)

	return &sdk.GenerateAuthURIOutput{
		AuthURI: authUri,
	}, nil
}

// ProcessAuthCallback implement Auth
func (s *AuthService) ProcessAuthCallback(input *sdk.ProcessAuthCallbackInput) (*sdk.ProcessAuthCallbackOutput, error) {
	authConf := s.config.Auth
	if authConf == nil {
		return nil, fmt.Errorf("auth no ams config")
	}

	authCode, err := s.getAuthCode(input.AuthCallback)
	if err != nil {
		return nil, err
	}

	amsResp, _, err := s.amsSDKClient.Oauth().Token(
		context.Background(), authConf.ClientID, authConf.ClientSecret, "authorization_code",
		&api.OauthTokenOpts{
			AuthorizationCode: optional.NewString(authCode),
			RedirectUri:       optional.NewString(authConf.RedirectUri),
		})
	if err != nil {
		return nil, err
	}

	if amsResp.AuthorizerInfo == nil {
		return nil, fmt.Errorf("no authorizer info returned")
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

	if err = account.Insert(authAccount); err != nil {
		return nil, err
	}
	return authAccount, nil
}

func (s *AuthService) getAuthCode(req *http.Request) (string, error) {
	query := req.URL.Query()
	authCode := query.Get("authorization_code")
	if authCode == "" {
		return "", fmt.Errorf("'authorization_code' parameter not exist")
	}
	return authCode, nil
}

// calcExpireAt 计算失效时间
func calcExpireAt(expireIn int64) time.Time {
	return time.Now().Add(time.Second * time.Duration(expireIn))
}

func (s *AuthService) refresh() {
	authConfig := s.config.Auth
	for {
		time.Sleep(10 * time.Second)
		authAccount := account.GetAllAuthAccount()

		for _, a := range authAccount {
			if needRefreshToken(a) {
				amsResp, _, err := s.amsSDKClient.Oauth().Token(
					context.Background(), authConfig.ClientID, authConfig.ClientSecret, "refresh_token",
					&api.OauthTokenOpts{
						RefreshToken: optional.NewString(a.RefreshToken),
					})
				if err != nil {
					log.Errorf("failed to call refresh token api for account[%d]", a.AccountId)
				} else {
					if err = account.RefreshToken(&sdk.AuthAccount{
						AccountId:           a.AccountId,
						AccessToken:         amsResp.AccessToken,
						AccessTokenExpireAt: calcExpireAt(amsResp.AccessTokenExpiresIn),
					}); err != nil {
						log.Errorf("failed to refresh account[%d] token", a.AccountId)
					}
				}
			}
		}
	}
}

// needRefreshToken 判断是否需要刷新token
// ams无法获取到refresh_token的失效时间，每次刷新时会更新，所以这里只判断access_token的失效时间
func needRefreshToken(account *sdk.AuthAccount) bool {
	now := time.Now()
	return account.AccessTokenExpireAt.Sub(now) <= time.Hour
}
