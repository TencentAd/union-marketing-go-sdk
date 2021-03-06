package ams

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/tencentad/union-marketing-go-sdk/api/sdk"
	"github.com/tencentad/union-marketing-go-sdk/pkg/sdk/account"
	"github.com/tencentad/union-marketing-go-sdk/pkg/sdk/config"
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

	account.RegisterGetTokenRefreshTime(sdk.AMS, s.GetTokenRefreshTime)
	account.RegisterRefreshToken(sdk.AMS, s.RefreshToken)

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
		url.QueryEscape(authConf.RedirectUri),
	)

	if len(input.State) > 0 {
		authUri = fmt.Sprintf("%s&state=%s",
			authUri,
			url.QueryEscape(input.State),
		)
	}

	return &sdk.GenerateAuthURIOutput{
		AuthURI: authUri,
	}, nil
}

// ProcessAuthCallback implement Auth
func (s *AuthService) ProcessAuthCallback(input *sdk.ProcessAuthCallbackInput) (*sdk.ProcessAuthCallbackOutput,
	error) {
	authConf := s.config.Auth
	if authConf == nil {
		return nil, fmt.Errorf("auth no ams config")
	}

	authCode, err := s.getAuthCode(input.AuthCallback)
	if err != nil {
		return nil, err
	}

	state, err := s.getState(input.AuthCallback)

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

	var accid string
	var amsSystemType sdk.AMSSystemType
	if info.AccountId > 0 {
		accid = strconv.FormatInt(info.AccountId, 10)
		amsSystemType = sdk.AmsEqq
	} else if len(info.WechatAccountId) > 0 {
		accid = info.WechatAccountId
		amsSystemType = sdk.AmsMp
	} else {
		return nil, fmt.Errorf("invalid accid")
	}

	authAccount := &sdk.AuthAccount{
		ID:                   formatAuthAccountID(accid, amsSystemType),
		Platform:             sdk.AMS,
		AccountUin:           info.AccountUin,
		AccountID:            strconv.FormatInt(info.AccountId, 10),
		ScopeList:            *info.ScopeList,
		WechatAccountID:      info.WechatAccountId,
		AccountRoleType:      AccountRoleTypeMapping[info.AccountRoleType],
		AccountType:          AccountTypeMapping[info.AccountType],
		AMSSystemType:        amsSystemType,
		RoleType:             RoleTypeMapping[info.RoleType],
		AccessToken:          amsResp.AccessToken,
		RefreshToken:         amsResp.RefreshToken,
		AccessTokenExpireAt:  calcExpireAt(amsResp.AccessTokenExpiresIn),
		RefreshTokenExpireAt: calcExpireAt(amsResp.RefreshTokenExpiresIn),
	}

	if err = account.Insert(authAccount); err != nil {
		return nil, err
	}

	resList := make([]*sdk.AuthAccount, 1)
	resList[0] = authAccount

	result := &sdk.ProcessAuthCallbackOutput{
		State:           state,
		AuthAccountList: resList,
	}
	return result, nil
}

func (s *AuthService) getAuthCode(req *http.Request) (string, error) {
	query := req.URL.Query()
	authCode := query.Get("authorization_code")
	if authCode == "" {
		return "", fmt.Errorf("'authorization_code' parameter not exist")
	}
	return authCode, nil
}

func (s *AuthService) getState(req *http.Request) (string, error) {
	query := req.URL.Query()
	state := query.Get("state")
	if state == "" {
		return "", fmt.Errorf("'state' parameter not exist")
	}
	return state, nil
}

// calcExpireAt 计算失效时间
func calcExpireAt(expireIn int64) time.Time {
	return time.Now().Add(time.Second * time.Duration(expireIn))
}

// GetTokenRefreshTime
// ams无法获取到refresh_token的失效时间，每次刷新时会更新，所以这里只判断access_token的失效时间
func (s *AuthService) GetTokenRefreshTime(account *sdk.AuthAccount) time.Time {
	return account.AccessTokenExpireAt
}

func (s *AuthService) RefreshToken(acc *sdk.AuthAccount) (*sdk.RefreshTokenOutput, error) {
	authConfig := s.config.Auth

	amsResp, _, err := s.amsSDKClient.Oauth().Token(
		context.Background(), authConfig.ClientID, authConfig.ClientSecret, "refresh_token",
		&api.OauthTokenOpts{
			RefreshToken: optional.NewString(acc.RefreshToken),
		})
	if err != nil {
		log.Errorf("failed to call refresh token api for account[%s]", acc.AccountID)
		return nil, err
	}

	return &sdk.RefreshTokenOutput{
		ID:                  acc.ID,
		AccessToken:         amsResp.AccessToken,
		AccessTokenExpireAt: calcExpireAt(amsResp.AccessTokenExpiresIn),
	}, nil
}

// formatAuthAccountID
func formatAuthAccountID(accountID string, systemType sdk.AMSSystemType) string {
	return fmt.Sprintf("%s:%s:%s", sdk.AMS, systemType, accountID)
}
