package oceanengine

import (
	api2 "git.code.oa.com/tme-server-component/kg_growth_open/api"
	tconfig "github.com/tencentad/marketing-api-go-sdk/pkg/config"
)

// OceanEngineReport
type OceanEngineReport struct {
}

// OceanEngineReport constructor
func NewOceanEngineReport(msdkConfig *api2.MarketingSDKConfig) *OceanEngineReport {
	tSdkConfig := tconfig.SDKConfig{
		AccessToken:   msdkConfig.AccessToken,
		IsDebug:       msdkConfig.IsDebug,
		DebugFile:     msdkConfig.DebugFile,
		SkipMonitor:   msdkConfig.SkipMonitor,
		IsStrictMode:  msdkConfig.IsStrictMode,
	}

	tClient.UseProduction()
	return &OceanEngineReport{
		mSdkConfig: sdkConfig,
	}
}

//func (a *DailyReportsApiService) Get(ctx context.Context, accountId int64, level string, dateRange ReportDateRange, localVarOptionals *DailyReportsGetOpts) (DailyReportsGetResponseData, http.Header, error) {
//	var (
//		localVarHttpMethod  = strings.ToUpper("Get")
//		localVarPostBody    interface{}
//		localVarFileName    string
//		localVarFileBytes   []byte
//		localVarFileKey     string
//		localVarReturnValue DailyReportsGetResponseData
//		localVarResponse    DailyReportsGetResponse
//	)
//
//	// create path and map variables
//	localVarPath := a.client.Cfg.BasePath + "/daily_reports/get"
//
//	localVarHeaderParams := make(map[string]string)
//	localVarQueryParams := url.Values{}
//	localVarFormParams := url.Values{}
//
//	localVarQueryParams.Add("account_id", parameterToString(accountId, ""))
//	localVarQueryParams.Add("level", parameterToString(level, ""))
//	localVarQueryParams.Add("date_range", parameterToString(dateRange, ""))
//	if localVarOptionals != nil && localVarOptionals.Filtering.IsSet() {
//		localVarQueryParams.Add("filtering", parameterToString(localVarOptionals.Filtering.Value(), "multi"))
//	}
//	if localVarOptionals != nil && localVarOptionals.GroupBy.IsSet() {
//		localVarQueryParams.Add("group_by", parameterToString(localVarOptionals.GroupBy.Value(), "multi"))
//	}
//	if localVarOptionals != nil && localVarOptionals.OrderBy.IsSet() {
//		localVarQueryParams.Add("order_by", parameterToString(localVarOptionals.OrderBy.Value(), "multi"))
//	}
//	if localVarOptionals != nil && localVarOptionals.Page.IsSet() {
//		localVarQueryParams.Add("page", parameterToString(localVarOptionals.Page.Value(), ""))
//	}
//	if localVarOptionals != nil && localVarOptionals.PageSize.IsSet() {
//		localVarQueryParams.Add("page_size", parameterToString(localVarOptionals.PageSize.Value(), ""))
//	}
//	if localVarOptionals != nil && localVarOptionals.TimeLine.IsSet() {
//		localVarQueryParams.Add("time_line", parameterToString(localVarOptionals.TimeLine.Value(), ""))
//	}
//	if localVarOptionals != nil && localVarOptionals.Fields.IsSet() {
//		localVarQueryParams.Add("fields", parameterToString(localVarOptionals.Fields.Value(), "multi"))
//	}
//	// to determine the Content-Type header
//	localVarHttpContentTypes := []string{"text/plain"}
//
//	// set Content-Type header
//	localVarHttpContentType := selectHeaderContentType(localVarHttpContentTypes)
//	if localVarHttpContentType != "" {
//		localVarHeaderParams["Content-Type"] = localVarHttpContentType
//	}
//
//	// to determine the Accept header
//	localVarHttpHeaderAccepts := []string{"application/json"}
//
//	// set Accept header
//	localVarHttpHeaderAccept := selectHeaderAccept(localVarHttpHeaderAccepts)
//	if localVarHttpHeaderAccept != "" {
//		localVarHeaderParams["Accept"] = localVarHttpHeaderAccept
//	}
//	r, err := a.client.prepareRequest(ctx, localVarPath, localVarHttpMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFileName, localVarFileBytes, localVarFileKey)
//	if err != nil {
//		return localVarReturnValue, nil, err
//	}
//
//	localVarHttpResponse, err := a.client.callAPI(r)
//	if err != nil || localVarHttpResponse == nil {
//		return localVarReturnValue, nil, err
//	}
//
//	localVarBody, err := ioutil.ReadAll(localVarHttpResponse.Body)
//	defer localVarHttpResponse.Body.Close()
//	if err != nil {
//		return localVarReturnValue, nil, err
//	}
//
//	if localVarHttpResponse.StatusCode < 300 {
//		// If we succeed, return the data, otherwise pass on to decode error.
//		err = a.client.decode(&localVarResponse, localVarBody, localVarHttpResponse.Header.Get("Content-Type"))
//		if err == nil {
//			if localVarResponse.Code != 0 {
//				var localVarResponseErrors []ApiErrorStruct
//				if localVarResponse.Errors != nil {
//					localVarResponseErrors = *localVarResponse.Errors
//				}
//				err = errors.NewError(localVarResponse.Code, localVarResponse.Message, localVarResponse.MessageCn, localVarResponseErrors)
//				return localVarReturnValue, localVarHttpResponse.Header, err
//			}
//			return *localVarResponse.Data, localVarHttpResponse.Header, err
//		} else {
//			return localVarReturnValue, localVarHttpResponse.Header, err
//		}
//	}
//
//	if localVarHttpResponse.StatusCode >= 300 {
//		newErr := GenericSwaggerError{
//			body:  localVarBody,
//			error: localVarHttpResponse.Status,
//		}
//
//		if localVarHttpResponse.StatusCode == 200 {
//			var v DailyReportsGetResponse
//			err = a.client.decode(&v, localVarBody, localVarHttpResponse.Header.Get("Content-Type"))
//			if err != nil {
//				newErr.error = err.Error()
//				return localVarReturnValue, localVarHttpResponse.Header, newErr
//			}
//			newErr.model = v
//			return localVarReturnValue, localVarHttpResponse.Header, newErr
//		}
//
//		return localVarReturnValue, localVarHttpResponse.Header, newErr
//	}
//
//	return localVarReturnValue, localVarHttpResponse.Header, nil
//}
