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

	"git.code.oa.com/tme-server-component/kg_growth_open/api/sdk"
	"git.code.oa.com/tme-server-component/kg_growth_open/pkg/sdk/account"
	"git.code.oa.com/tme-server-component/kg_growth_open/pkg/sdk/config"
	"git.code.oa.com/tme-server-component/kg_growth_open/pkg/sdk/http_tools"
)

// HTTPServer 登录授权服务
type AuthService struct {
	config     *config.Config
	httpClient *http_tools.HttpClient
}

// NewAuthService
func NewAuthService(config *config.Config) *AuthService {
	s := &AuthService{
		config:     config,
		httpClient: http_tools.Init(config.HttpConfig),
	}

	account.RegisterGetTokenRefreshTime(sdk.OceanEngine, s.GetTokenRefreshTime)
	account.RegisterRefreshToken(sdk.OceanEngine, s.RefreshToken)

	return s
}

// GenerateAuthURI implement Auth
func (s *AuthService) GenerateAuthURI(input *sdk.GenerateAuthURIInput) (*sdk.GenerateAuthURIOutput, error) {
	authConf := s.config.Auth
	if authConf == nil {
		return nil, fmt.Errorf("auth no ocean engine config")
	}
	authUri := fmt.Sprintf("https://ad.oceanengine.com/openapi/audit/oauth.html?app_id=%d&redirect_uri=%s",
		s.config.Auth.ClientID,
		s.config.Auth.RedirectUri,
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

	method := http_tools.POST
	// create path and map variables
	path := s.httpClient.Config.BasePath + "/oauth2/access_token/"
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

	request, err := s.httpClient.PrepareRequest(context.Background(), path, method, postParams, header,
		nil, nil, "", nil, "")
	if err != nil {
		return nil, err
	}

	authResponse := &AuthReponse{}
	respErr := s.httpClient.DoProcess(request, authResponse)
	return authResponse, respErr
}

type AccountRole string

const (
	Advertise        AccountRole = "ADVERTISER"
	CustomerAdmin    AccountRole = "CUSTOMER_ADMIN"
	Agent            AccountRole = "AGENT"
	ChildAgent       AccountRole = "CHILD_AGENT"
	CustomerOperator AccountRole = "CUSTOMER_OPERATOR"
)

type TokenAdvertiserInfo struct {
	AdvertiserId   int64       `json:"advertiser_id"`
	AdvertiserName string      `json:"advertiser_name"`
	AdvertiserRole int64       `json:"advertiser_role"`
	IsValid        bool        `json:"is_valid"`
	AccountRole    AccountRole `json:"account_role"`
}

type TokenAdvertiserData struct {
	AdvertiserList []*TokenAdvertiserInfo `json:"list"`
}

type TokenAdvertiser struct {
	Code      int                  `json:"code"`
	Message   string               `json:"message"`
	Data      *TokenAdvertiserData `json:"data"`
	RequestId string               `json:"request_id"`
}

// GetAdvertiserListByToken 根据token获取广告主列表
func (s *AuthService) GetAdvertiserListByToken(accessToken string) ([]int64, error) {
	authConf := s.config.Auth
	if authConf == nil {
		return nil, fmt.Errorf("auth no ocean engine config")
	}

	method := http_tools.GET
	// create path and map variables
	path := s.httpClient.Config.BasePath + "/oauth2/advertiser/get/"

	queryParams := url.Values{}
	queryParams["access_token"] = []string{accessToken}
	queryParams["app_id"] = []string{strconv.FormatInt(authConf.ClientID, 10)}
	queryParams["secret"] = []string{authConf.ClientSecret}

	header := make(map[string]string)
	header["Content-Type"] = "application/json"
	header["Accept"] = "application/json"
	header["X-Debug-Mode"] = strconv.FormatInt(1, 10)

	request, err := s.httpClient.PrepareRequest(context.Background(), path, method, nil, header,
		queryParams, nil, "", nil, "")
	if err != nil {
		return nil, err
	}

	tokenAdvertiser := &TokenAdvertiser{}
	respErr := s.httpClient.DoProcess(request, tokenAdvertiser)
	if respErr != nil {
		return nil, respErr
	}
	if len(tokenAdvertiser.Data.AdvertiserList) == 0 {
		return nil, fmt.Errorf("GetAdvertiserListByToken get no account")
	}

	adList := tokenAdvertiser.Data.AdvertiserList
	var res []int64
	for i := 0; i < len(adList); i++ {
		ad := adList[i]
		if !ad.IsValid {
			continue
		}
		switch ad.AccountRole {
		case Advertise:
			res = append(res, ad.AdvertiserId)
		case CustomerOperator:
		case CustomerAdmin:
			// 账号管家
			customerAdvertiserList, err := s.GetAdvertiserListByCustomerID(accessToken, ad.AdvertiserId)
			if err != nil {
				return nil, err
			}
			if customerAdvertiserList.Code != 0 {
				return nil, fmt.Errorf("response : code = %d, message = %s, request_id = %s ", customerAdvertiserList.Code,
					customerAdvertiserList.Message,
					customerAdvertiserList.RequestId)
			}
			adList := customerAdvertiserList.Data.AdList
			for i := 0; i < len(adList); i++ {
				res = append(res, adList[i].AdvertiserId)
			}
			break
		case Agent:
		case ChildAgent:
			// 代理商
			agentAdvertiserList, err := s.GetAdvertiserListByAgentID(accessToken, ad.AdvertiserId, "0", "65535")
			if err != nil {
				return nil, err
			}
			if agentAdvertiserList.Code != 0 {
				return nil, fmt.Errorf("response : code = %d, message = %s, request_id = %s ", agentAdvertiserList.Code,
					agentAdvertiserList.Message,
					agentAdvertiserList.RequestId)
			}
			adList := agentAdvertiserList.Data.AdvertiserIDList
			for i := 0; i < len(adList); i++ {
				res = append(res, adList[i])
			}
			break
		}
	}
	return res, nil
}

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
	adList, err := s.GetAdvertiserListByToken(authResponse.Data.AccessToken)
	if err != nil {
		return nil, err
	}

	resList := make([]*sdk.AuthAccount, 0, len(adList))
	for i := 0; i < len(adList); i++ {
		accID := strconv.FormatInt(adList[i], 10)
		authAccount := &sdk.AuthAccount{
			ID:                   formatAuthAccountID(accID),
			AccountID:            accID,
			AccessToken:          authResponse.Data.AccessToken,
			AccessTokenExpireAt:  calcExpireAt(authResponse.Data.ExpiresIn),
			RefreshToken:         authResponse.Data.RefreshToken,
			RefreshTokenExpireAt: calcExpireAt(authResponse.Data.RefreshTokenExpiresIn),
			Platform:             sdk.OceanEngine,
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

type AgentGetAdvertiserData struct {
	AdvertiserIDList []int64       `json:"list"`
	PageInfo         *sdk.PageConf `json:"page_info"`
}

type AgentGetAdvertiserResponse struct {
	Code      int                     `json:"app_id"`
	Message   string                  `json:"message"`
	Data      *AgentGetAdvertiserData `json:"data"`
	RequestId string                  `json:"request_id"`
}

func (s *AuthService) GetAdvertiserListByAgentID(accessToken string, advertiserID int64, page string,
	pageSize string) (*AgentGetAdvertiserResponse, error) {
	authConf := s.config.Auth
	if authConf == nil {
		return nil, fmt.Errorf("auth no ocean engine config")
	}

	method := http_tools.GET
	// create path and map variables
	path := s.httpClient.Config.BasePath + "/2/agent/advertiser/select/"

	header := make(map[string]string)
	header["Access-Token"] = accessToken
	header["Content-Type"] = "application/json"
	header["Accept"] = "application/json"

	queryParams := url.Values{}
	queryParams["advertiser_id"] = []string{strconv.FormatInt(advertiserID, 10)}
	queryParams["page"] = []string{page}
	queryParams["page_size"] = []string{pageSize}

	request, err := s.httpClient.PrepareRequest(context.Background(), path, method, nil, header,
		queryParams, nil, "", nil, "")
	if err != nil {
		return nil, err
	}

	agentGetAdvertiser := &AgentGetAdvertiserResponse{}
	respErr := s.httpClient.DoProcess(request, agentGetAdvertiser)
	if respErr != nil {
		return nil, respErr
	}
	return agentGetAdvertiser, nil
}

type CustomerGetAdvertiserInfo struct {
	AdvertiserId   int64  `json:"advertiser_id"`
	AdvertiserName string `json:"advertiser_name"`
}

type CustomerGetAdvertiserData struct {
	AdList []*CustomerGetAdvertiserInfo `json:"list"`
}

type CustomerGetAdvertiserResponse struct {
	Code      int                        `json:"app_id"`
	Message   string                     `json:"message"`
	Data      *CustomerGetAdvertiserData `json:"data"`
	RequestId string                     `json:"request_id"`
}

func (s *AuthService) GetAdvertiserListByCustomerID(accessToken string, advertiserID int64) (*CustomerGetAdvertiserResponse, error) {
	authConf := s.config.Auth
	if authConf == nil {
		return nil, fmt.Errorf("auth no ocean engine config")
	}

	method := http_tools.GET
	// create path and map variables
	path := s.httpClient.Config.BasePath + "/2/majordomo/advertiser/select/"

	header := make(map[string]string)
	header["Content-Type"] = "application/json"
	header["Accept"] = "application/json"
	header["Access-Token"] = accessToken

	queryParams := url.Values{}
	queryParams["advertiser_id"] = []string{strconv.FormatInt(advertiserID, 10)}

	request, err := s.httpClient.PrepareRequest(context.Background(), path, method, nil, header,
		queryParams, nil, "", nil, "")
	if err != nil {
		return nil, err
	}

	customerGetAdvertiser := &CustomerGetAdvertiserResponse{}
	respErr := s.httpClient.DoProcess(request, customerGetAdvertiser)
	if respErr != nil {
		return nil, respErr
	}
	return customerGetAdvertiser, nil
}
