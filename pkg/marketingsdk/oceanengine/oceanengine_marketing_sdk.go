package oceanengine

import (
	"context"
	"git.code.oa.com/tme-server-component/kg_growth_open/api"
	uuid "github.com/satori/go.uuid"
	"github.com/tencentad/marketing-api-go-sdk/pkg/config"
	"net/http"
	"strconv"
	"time"
)

type OceanEngineMarketingService struct {
	SdkConfig     *api.MarketingSDKConfig // sdk配置
	oceanEngineClient *APIClient // API请求Client
	reportService *OceanEngineReport   // 报表模块
}

// Name 名称
func (oe *OceanEngineReport) Name() string {
	return "ByteDance-Report"
}

// Init ...
func (oe *OceanEngineReport) init(cfg *config.SDKConfig) *SDKClient {
	version := "1.6.0"
	apiVersion := "v1.1"
	ctx := context.Background()
	ctx = context.WithValue(ctx, api2.ContextAPIKey, apiKey)
	client := api.NewAPIClient(cfg)
	sdkClient := &SDKClient{
		Config:       cfg,
		Ctx:          &ctx,
		Client:       client,
		RoundTripper: http.DefaultTransport,
		Version:      version,
		ApiVersion:   apiVersion,
	}
	sdkClient.Client.Cfg.HTTPClient.Transport = sdkClient
	sdkClient.UseSandbox()
	sdkClient.middlewareList = []Middleware{
		&AuthMiddleware{sdkClient},
		&HeaderMiddleware{sdkClient},
		&DiffHostMiddleware{sdkClient},
		&LogMiddleware{sdkClient},
	}
	return sdkClient
}

// 获取报表接口
func (t *OceanEngineMarketingService) GetReport(reqParam *api.ReportInputParam) (*api.ReportOutput, error) {
	//if reqParam.ReportTimeGranularity == api.ReportTimeDaily {
	//	return t.reportService.getDailyReport(t.SdkConfig, reqParam)
	//} else {
	//	return t.reportService.getHourlyReport(t.SdkConfig, reqParam)
	//}
	return nil, nil
}