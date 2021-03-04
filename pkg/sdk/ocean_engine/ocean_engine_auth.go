package ocean_engine

//
import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	"git.code.oa.com/tme-server-component/kg_growth_open/api/sdk"
	"git.code.oa.com/tme-server-component/kg_growth_open/pkg/sdk/account"
	"git.code.oa.com/tme-server-component/kg_growth_open/pkg/sdk/config"
	"git.code.oa.com/tme-server-component/kg_growth_open/pkg/sdk/http_tools"
)

// HTTPServer 登录授权服务
type AuthService struct {
	config     *config.Config
	httpClinet *http_tools.HttpClient
}

// NewAuthService
func NewAuthService(config *config.Config) *AuthService {
	s := &AuthService{
		config:     config,
		httpClinet: http_tools.Init(config.HttpConfig),
	}

	if err := s.init(); err != nil {
		log.Errorf("failed to init AuthService, err: %v", err)
	}

	account.RegisterGetTokenRefreshTime(sdk.OceanEngine, s.GetTokenRefreshTime)
	account.RegisterRefreshToken(sdk.OceanEngine, s.RefreshToken)

	return s
}

func (s *AuthService) init() error {
	return nil
}

// GenerateAuthURI implement Auth
func (s *AuthService) GenerateAuthURI(input *sdk.GenerateAuthURIInput) (*sdk.GenerateAuthURIOutput, error) {
	authConf := s.config.Auth
	if authConf == nil {
		return nil, fmt.Errorf("auth no ocean engine config")
	}
	authUri := fmt.Sprintf("https://ad.oceanengine.com/openapi/audit/oauth.html?app_id=%d&redirect_uri=%s&state=%s",
		s.config.Auth.ClientID,
		s.config.Auth.RedirectUri,
		url.QueryEscape(input.State),
	)

	return &sdk.GenerateAuthURIOutput{
		AuthURI: authUri,
	}, nil
}

type PostBody struct {
	AppId        int64  `json:"app_id"`
	Secret       string `json:"secret"`
	GrantType    string `json:"grant_type"`
	AuthCode     string `json:"auth_code,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type Data struct {
	AccessToken           string  `json:"access_token"`
	ExpiresIn             int64   `json:"expires_in"`
	RefreshToken          string  `json:"refresh_token"`
	RefreshTokenExpiresIn int64   `json:"refresh_token_expires_in"`
	AdvertiserIds         []int64 `json:"advertiser_ids"`
}

type AuthReponse struct {
	Code      int    `json:"app_id"`
	Message   string `json:"message"`
	Data      *Data  `json:"data"`
	RequestId string `json:"request_id"`
}

func (s *AuthService) getToken(isRefresh bool, val string) (*AuthReponse, error) {
	authConf := s.config.Auth
	if authConf == nil {
		return nil, fmt.Errorf("auth no ocean engine config")
	}

	method := "POST"
	// create path and map variables
	path := s.httpClinet.Config.BasePath + "/oauth2/access_token/"
	postBody := PostBody{
		AppId:     s.config.Auth.ClientID,
		Secret:    s.config.Auth.ClientSecret,
		GrantType: "auth_code",
	}
	if isRefresh {
		postBody.RefreshToken = val
	} else {
		postBody.AuthCode = val
	}
	postParams, _ := json.Marshal(postBody)

	header := make(map[string]string)
	header["Content-Type"] = "application/json"
	header["Accept"] = "application/json"

	request, err := s.httpClinet.PrepareRequest(context.Background(), path, method, postParams, header,
		nil, nil, "", nil, "")
	if err != nil {
		return nil, err
	}

	authResponse := &AuthReponse{}
	respErr := s.httpClinet.DoProcess(request, authResponse)
	return authResponse, respErr
}

//// GetAdvertiserListByToken 根据token获取广告主列表
//func (s *AuthService) GetAdvertiserListByToken(accessToken string) ([]int64, error) {
//	authConf := s.config.Auth
//	if authConf == nil {
//		return nil, fmt.Errorf("auth no ocean engine config")
//	}
//
//	method := "GET"
//	// create path and map variables
//	path := s.httpClinet.Config.BasePath + "/oauth2/advertiser/get/"
//
//	headerParams := make(map[string]string)
//	headerParams["Content-Type"] = "application/json"
//	headerParams["Accept"] = "application/json"
//	queryParams := url.Values{}
//	queryParams["access_token"] = []string{accessToken}
//	queryParams["app_id"] = []string{strconv.FormatInt(authConf.ClientID, 10)}
//	queryParams["secret"] = []string{authConf.ClientSecret}
//
//	//request, err := s.httpClinet.PrepareRequest(context.Background(), path, method, nil, headerParams,
//	//	nil, nil, "", nil, "")
//	//if err != nil {
//	//	return nil, err
//	//}
//
//	type TokenAdvertiser struct {
//		Code      int    `json:"app_id"`
//		Message   string `json:"message"`
//		Data      *Data  `json:data`
//		RequestId string `json:request_id`
//	}
//	//type Data struct {
//	//	AdvertiserId         int64 `json:"advertiser_id"`
//	//	AdvertiserId         int64 `json:"advertiser_name"`
//	//	AdvertiserId         int64 `json:"advertiser_role"`
//	//	AdvertiserId         int64 `json:"is_valid"`
//	//	AdvertiserId         int64 `json:"account_role"`
//	//
//	//}
//	//
//	//authReponse := &AuthReponse{}
//	//resp_err := s.httpClinet.DoProcess(request, authReponse)
//	//return authReponse, resp_err
//	return nil,nil
//}

// ProcessAuthCallback implement Auth
func (s *AuthService) ProcessAuthCallback(input *sdk.ProcessAuthCallbackInput) (*sdk.ProcessAuthCallbackOutput,
	error) {
	authCode, err := s.getAuthCode(input.AuthCallback)
	if err != nil {
		return nil, err
	}
	authResponse, respErr := s.getToken(false, authCode)
	if respErr != nil {
		return nil, respErr
	}
	if authResponse.Code != 0 {
		return nil, fmt.Errorf("response : code = %d, message = %s, request_id = %s ", authResponse.Code,
			authResponse.Message,
			authResponse.RequestId)
	}

	// 头条需要根据Token获取已授权账户

	resList := make([]*sdk.AuthAccount, 0, len(authReponse.Data.AdvertiserIds))
	for i := 0; i < len(authReponse.Data.AdvertiserIds); i++ {
		accid := strconv.FormatInt(authReponse.Data.AdvertiserIds[i], 10)
		authAccount := &sdk.AuthAccount{
			ID:                   formatAuthAccountID(accID),
			AccountID:            accID,
			AccessToken:          authResponse.Data.AccessToken,
			AccessTokenExpireAt:  calcExpireAt(authResponse.Data.ExpiresIn),
			RefreshToken:         authResponse.Data.RefreshToken,
			RefreshTokenExpireAt: calcExpireAt(authResponse.Data.RefreshTokenExpiresIn),
		}
		if err = account.Insert(authAccount); err != nil {
			return nil, err
		}
		resList = append(resList, authAccount)
	}
	result := &sdk.ProcessAuthCallbackOutput{
		AuthAccountList: resList,
	}
	return result, nil
}

func (s *AuthService) getAuthCode(req *http.Request) (string, error) {
	query := req.URL.Query()
	authCode := query.Get("auth_code")
	if authCode == "" {
		return "", fmt.Errorf("'auth_code' parameter not exist")
	}
	return authCode, nil
}

// calcExpireAt 计算失效时间
func calcExpireAt(expireIn int64) time.Time {
	return time.Now().Add(time.Second * time.Duration(expireIn))
}

// GetTokenRefreshTime  这里只判断access_token的失效时间
func (s *AuthService) GetTokenRefreshTime(account *sdk.AuthAccount) time.Time {
	return account.AccessTokenExpireAt
}

func (s *AuthService) RefreshToken(acc *sdk.AuthAccount) (*sdk.RefreshTokenOutput, error) {
	authResponse, respErr := s.getToken(true, acc.RefreshToken)
	if respErr != nil {
		return nil, respErr
	}
	if authResponse.Code != 0 {
		return nil, fmt.Errorf("response : code = %d, message = %s, request_id = %s ", authResponse.Code,
			authResponse.Message,
			authResponse.RequestId)
	}

	return &sdk.RefreshTokenOutput{
		ID:                   acc.ID,
		AccessToken:          authResponse.Data.AccessToken,
		AccessTokenExpireAt:  calcExpireAt(authResponse.Data.ExpiresIn),
		RefreshToken:         authResponse.Data.RefreshToken,
		RefreshTokenExpireAt: calcExpireAt(authResponse.Data.RefreshTokenExpiresIn),
	}, nil
}

// formatAuthAccountID
func formatAuthAccountID(accountID string) string {
	return fmt.Sprintf("%s:%s", sdk.OceanEngine, accountID)
}
