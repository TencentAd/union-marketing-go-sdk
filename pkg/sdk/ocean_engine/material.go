package ocean_engine

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"git.code.oa.com/tme-server-component/kg_growth_open/api/sdk"
	"git.code.oa.com/tme-server-component/kg_growth_open/pkg/sdk/account"
	"git.code.oa.com/tme-server-component/kg_growth_open/pkg/sdk/config"
	"git.code.oa.com/tme-server-component/kg_growth_open/pkg/sdk/http_tools"
)

// MaterialService 物料服务
type MaterialService struct {
	config     *config.Config
	httpClient *http_tools.HttpClient
}

// NewMaterialService 获取物料服务
func NewMaterialService(sConfig *config.Config) *MaterialService {
	return &MaterialService{
		config:     sConfig,
		httpClient: http_tools.Init(sConfig.HttpConfig),
	}
}

// AddImage 增加图片上传
func (s *MaterialService) AddImage(input *sdk.ImageAddInput) (*sdk.ImagesAddOutput, error) {
	switch input.UploadType {
	case sdk.UploadTypeFile:
		return s.AddImageByFile(input)
	case sdk.UploadTypeUrl:
		return s.AddImageByURL(input)
	default:
		return nil, fmt.Errorf("no support upload type = %s", input.UploadType)
	}
}

// AddImageByFile File增加图片上传
func (s *MaterialService) AddImageByFile(input *sdk.ImageAddInput) (*sdk.ImagesAddOutput, error) {
	id := formatAuthAccountID(input.BaseInput.AccountId)
	authAccount, err := account.GetAuthAccount(id)
	if err != nil {
		return nil, fmt.Errorf("AddImage get AuthAccount Info error accid=%s", input.BaseInput.AccountId)
	}

	method := http_tools.POST
	// create path and map variables
	path := s.config.HttpConfig.BasePath + "/2/file/image/ad/"

	header := make(map[string]string)
	header["Accept"] = "application/json"
	header["Access-Token"] = authAccount.AccessToken
	var request *http.Request

	formParams := url.Values{}
	header["Content-Type"] = "multipart/form-data"
	formParams["advertiser_id"] = []string{input.BaseInput.AccountId}
	formParams["upload_type"] = []string{string(input.UploadType)}
	formParams["image_signature"] = []string{input.Signature}
	if input.File == nil {
		return nil, fmt.Errorf("AddImage file is empty")
	}
	localNewFile, localErr := os.Open(input.File.Name())
	if localErr != nil {
		return nil, localErr
	}
	defer localNewFile.Close()
	fileBytes, localErr := ioutil.ReadAll(localNewFile)
	if localErr != nil {
		return nil, localErr
	}
	fileName := input.File.Name()
	fileKey := "image_file"
	request, err = s.httpClient.PrepareRequest(context.Background(), path, method, nil, header,
		nil, formParams, fileName, fileBytes, fileKey)
	if err != nil {
		return nil, err
	}

	response := &AddImageData{}
	respErr := s.httpClient.DoProcess(request, response)
	if respErr != nil {
		return nil, respErr
	}
	if response.Code != 0 {
		return nil, fmt.Errorf("response : code = %d, message = %s, request_id = %s ", response.Code,
			response.Message,
			response.RequestId)
	}

	imageInfo := response.Data
	output := &sdk.ImagesAddOutput{
		ImageId:     imageInfo.ImageId,
		PreviewUrl:  imageInfo.PreviewUrl,
		Description: "",
		Width:       imageInfo.Width,
		Height:      imageInfo.Height,
		FileSize:    imageInfo.FileSize,
		Signature:   imageInfo.Signature,
	}

	return output, nil
}

func (s *MaterialService) getUploadType(uploadType sdk.UploadType) string {
	switch uploadType {
	case sdk.UploadTypeFile:
		return "UPLOAD_BY_FILE"
	case sdk.UploadTypeUrl:
		return "UPLOAD_BY_URL"
	default:
		return ""
	}
}

type AddImageStruct struct {
	AdvertiserID int64  `json:"advertiser_id,omitempty"`   // 广告主ID
	UploadType   string `json:"upload_type,omitempty"`     // 图片上传方式
	Signature    string `json:"image_signature,omitempty"` // 图片的md5值
	ImageUrl     string `json:"image_url,omitempty"`       // 图片url地址
}

// AddImageByFile File增加图片上传
func (s *MaterialService) AddImageByURL(input *sdk.ImageAddInput) (*sdk.ImagesAddOutput, error) {
	id := formatAuthAccountID(input.BaseInput.AccountId)
	authAccount, err := account.GetAuthAccount(id)
	if err != nil {
		return nil, fmt.Errorf("AddImageByURL get AuthAccount Info error accid=%s", input.BaseInput.AccountId)
	}

	method := http_tools.POST
	// create path and map variables
	path := s.config.HttpConfig.BasePath + "/2/file/image/ad/"

	header := make(map[string]string)
	header["Accept"] = "application/json"
	header["Access-Token"] = authAccount.AccessToken
	header["Content-Type"] = "application/json"

	var request *http.Request
	accID, _ := strconv.ParseInt(input.BaseInput.AccountId, 10, 64)
	addImage := &AddImageStruct{
		AdvertiserID: accID,
		UploadType:   s.getUploadType(input.UploadType),
		Signature:    input.Signature,
		ImageUrl:     input.ImageUrlOceanEngine,
	}

	imageJson, _ := json.Marshal(addImage)

	request, err = s.httpClient.PrepareRequest(context.Background(), path, method, imageJson, header,
		nil, nil, "", nil, "")
	if err != nil {
		return nil, err
	}

	response := &AddImageData{}
	respErr := s.httpClient.DoProcess(request, response)
	if respErr != nil {
		return nil, respErr
	}
	if response.Code != 0 {
		return nil, fmt.Errorf("response : code = %d, message = %s, request_id = %s ", response.Code,
			response.Message,
			response.RequestId)
	}

	imageInfo := response.Data
	output := &sdk.ImagesAddOutput{
		ImageId:     imageInfo.ImageId,
		PreviewUrl:  imageInfo.PreviewUrl,
		Description: "",
		Width:       imageInfo.Width,
		Height:      imageInfo.Height,
		FileSize:    imageInfo.FileSize,
		Signature:   imageInfo.Signature,
	}

	return output, nil
}

type MaterialFilterInfo struct {
	ImageIds    []string `json:"image_ids,omitempty"`    // 图片ids 数量限制：<=100  注意：image_ids、material_ids、signatures只能选择一个进行过滤
	MaterialIds []int64  `json:"material_ids,omitempty"` // 图片ids 数量限制：<=100  注意：image_ids、material_ids、signatures只能选择一个进行过滤
	Signatures  []string `json:"signatures,omitempty"`   // 图片ids 数量限制：<=100  注意：image_ids、material_ids、signatures只能选择一个进行过滤
	Width       int64    `json:"width,omitempty"`        // 图片宽度
	Height      int64    `json:"height,omitempty"`       // 图片高度
	StartTime   string   `json:"start_time,omitempty"`   // 根据视频上传时间进行过滤的起始时间，与end_time搭配使用，格式：yyyy-mm-dd
	EndTime     string   `json:"end_time,omitempty"`     // 根据视频上传时间进行过滤的截止时间，与start_time搭配使用，格式：yyyy-mm-dd
}

// getFilter 获取过滤信息
func (s *MaterialService) getFilter(input *sdk.MaterialGetInput) (string, error) {
	if input.Filtering == nil {
		return "", nil
	}

	mFilterInfo := &MaterialFilterInfo{}

	if len(input.Filtering.Ids) > 0 {
		imageIDList := input.Filtering.Ids
		for i := 0; i < len(imageIDList); i++ {
			mFilterInfo.ImageIds = append(mFilterInfo.ImageIds, imageIDList[i])
		}
	}

	if len(input.Filtering.MaterialIds) > 0 {
		materialIDList := input.Filtering.MaterialIds
		for i := 0; i < len(materialIDList); i++ {
			mID, err := strconv.ParseInt(materialIDList[i], 10, 64)
			if err != nil {
				return "", err
			}
			mFilterInfo.MaterialIds = append(mFilterInfo.MaterialIds, mID)
		}
	}

	if len(input.Filtering.Signatures) > 0 {
		sigList := input.Filtering.Signatures
		for i := 0; i < len(sigList); i++ {
			mFilterInfo.Signatures = append(mFilterInfo.Signatures, sigList[i])
		}
	}

	// Width
	if input.Filtering.Width > 0 {
		mFilterInfo.Width = input.Filtering.Width
	}
	// Height
	if input.Filtering.Height > 0 {
		mFilterInfo.Height = input.Filtering.Height
	}
	// CreatedStartTime
	if len(input.Filtering.CreatedStartTime) > 0 {
		mFilterInfo.StartTime = input.Filtering.CreatedStartTime
	}
	// CreatedEndTime
	if len(input.Filtering.CreatedEndTime) > 0 {
		mFilterInfo.EndTime = input.Filtering.CreatedEndTime
	}

	filterJson, _ := json.Marshal(mFilterInfo)
	return string(filterJson), nil
}

// GetImage 获取图片信息
func (s *MaterialService) GetImage(input *sdk.MaterialGetInput) (*sdk.ImageGetOutput, error) {
	id := formatAuthAccountID(input.BaseInput.AccountId)
	authAccount, err := account.GetAuthAccount(id)
	if err != nil {
		return nil, fmt.Errorf("GetImage get AuthAccount Info error accid=%s", input.BaseInput.AccountId)
	}

	method := http_tools.POST
	// create path and map variables
	path := s.config.HttpConfig.BasePath + "/2/file/image/get"

	header := make(map[string]string)
	header["Accept"] = "application/json"
	header["Access-Token"] = authAccount.AccessToken
	header["Content-Type"] = "application/json"

	query := url.Values{}

	var request *http.Request
	query["advertiser_id"] = []string{input.BaseInput.AccountId}

	videoFilter, err := s.getFilter(input)
	if err != nil {
		return nil, err
	}
	if len(videoFilter) > 0 {
		query["filtering"] = []string{videoFilter}
	}

	if input.Page > 0 {
		query["page"] = []string{strconv.FormatInt(input.Page, 10)}
	}
	if input.PageSize > 0 {
		query["page_size"] = []string{strconv.FormatInt(input.PageSize, 10)}
	}

	request, err = s.httpClient.PrepareRequest(context.Background(), path, method, nil, header,
		query, nil, "", nil, "")
	if err != nil {
		return nil, err
	}

	response := &GetMaterialData{}
	respErr := s.httpClient.DoProcess(request, response)
	if respErr != nil {
		return nil, respErr
	}
	if response.Code != 0 {
		return nil, fmt.Errorf("response : code = %d, message = %s, request_id = %s ", response.Code,
			response.Message,
			response.RequestId)
	}

	imageInfo := response.Data
	imageOutput := &sdk.ImageGetOutput{}
	s.copyImageInfoToOutput(imageInfo, imageOutput)
	return imageOutput, err
}

// copyImageToOutput 拷贝物料信息
func (s *MaterialService) copyImageInfoToOutput(imageData *GetMaterialList, imageOutput *sdk.ImageGetOutput) {
	if imageData == nil {
		return
	}
	rList := make([]*sdk.ImageGetOutputStruct, 0, len(imageData.MaterialList))
	materialList := imageData.MaterialList
	for i := 0; i < len(materialList); i++ {
		imageData := (materialList)[i]
		rList = append(rList, &sdk.ImageGetOutputStruct{
			ImageId:     imageData.Id,
			Width:       imageData.Width,
			Height:      imageData.Height,
			FileSize:    imageData.FileSize,
			Signature:   imageData.Signature,
			PreviewUrl:  imageData.PreviewUrl,
			CreatedTime: imageData.CreatedTime,
		})
	}
	imageOutput.List = rList
	imageOutput.PageInfo = &sdk.PageConf{
		Page:        imageData.PageInfo.Page,
		PageSize:    imageData.PageInfo.PageSize,
		TotalNumber: imageData.PageInfo.TotalNumber,
		TotalPage:   imageData.PageInfo.TotalPage,
	}
}

// GetVideo 获取视频信息
func (s *MaterialService) GetVideo(input *sdk.MaterialGetInput) (*sdk.VideoGetOutput, error) {
	id := formatAuthAccountID(input.BaseInput.AccountId)
	authAccount, err := account.GetAuthAccount(id)
	if err != nil {
		return nil, fmt.Errorf("GetImage get AuthAccount Info error accid=%s", input.BaseInput.AccountId)
	}

	method := http_tools.POST
	// create path and map variables
	path := s.config.HttpConfig.BasePath + "/2/file/video/get"

	header := make(map[string]string)
	header["Accept"] = "application/json"
	header["Access-Token"] = authAccount.AccessToken
	header["Content-Type"] = "application/json"

	query := url.Values{}

	var request *http.Request
	query["advertiser_id"] = []string{input.BaseInput.AccountId}

	imageFilter, err := s.getFilter(input)
	if err != nil {
		return nil, err
	}
	if len(imageFilter) > 0 {
		query["filtering"] = []string{imageFilter}
	}

	if input.Page > 0 {
		query["page"] = []string{strconv.FormatInt(input.Page, 10)}
	}
	if input.PageSize > 0 {
		query["page_size"] = []string{strconv.FormatInt(input.PageSize, 10)}
	}

	request, err = s.httpClient.PrepareRequest(context.Background(), path, method, nil, header,
		query, nil, "", nil, "")
	if err != nil {
		return nil, err
	}

	response := &GetMaterialData{}
	respErr := s.httpClient.DoProcess(request, response)
	if respErr != nil {
		return nil, respErr
	}
	if response.Code != 0 {
		return nil, fmt.Errorf("response : code = %d, message = %s, request_id = %s ", response.Code,
			response.Message,
			response.RequestId)
	}

	videoInfo := response.Data
	videoOutput := &sdk.VideoGetOutput{}
	s.copyVideoInfoToOutput(videoInfo, videoOutput)
	return videoOutput, err
}

// copyVideoInfoToOutput 拷贝视频物料信息
func (s *MaterialService) copyVideoInfoToOutput(videoData *GetMaterialList, videoOutput *sdk.VideoGetOutput) {
	if videoData == nil {
		return
	}
	rList := make([]*sdk.VideoGetOutputStruct, 0, len(videoData.MaterialList))
	materialList := videoData.MaterialList
	for i := 0; i < len(materialList); i++ {
		videoInfo := (materialList)[i]
		rList = append(rList, &sdk.VideoGetOutputStruct{
			VideoId:                  videoInfo.Id,
			Width:                    videoInfo.Width,
			Height:                   videoInfo.Height,
			VideoCodec:               videoInfo.Format,
			VideoBitRate:             videoInfo.BitRate,
			Signature:                videoInfo.Signature,
			PreviewUrl:               videoInfo.PreviewUrl,
			KeyFrameImageUrl:         videoInfo.PosterUrl,
			CreatedTime:              videoInfo.CreatedTime,
			SourceType:               sdk.MaterialSourceType(videoInfo.Source),
			VideoDurationMillisecond: int64(videoInfo.Duration * 1000),
			MaterialID:               videoInfo.MaterialId,
			FileName:                 videoInfo.FileName,
			VideoLabels:              videoInfo.Labels,
		})
	}
	videoOutput.List = rList
	videoOutput.PageInfo = &sdk.PageConf{
		Page:        videoData.PageInfo.Page,
		PageSize:    videoData.PageInfo.PageSize,
		TotalNumber: videoData.PageInfo.TotalNumber,
		TotalPage:   videoData.PageInfo.TotalPage,
	}
}

// AddVideo 增加视频上传
func (s *MaterialService) AddVideo(input *sdk.VideoAddInput) (*sdk.VideoAddOutput, error) {
	id := formatAuthAccountID(input.BaseInput.AccountId)
	authAccount, err := account.GetAuthAccount(id)
	if err != nil {
		return nil, fmt.Errorf("AddImage get AuthAccount Info error accid=%s", input.BaseInput.AccountId)
	}

	method := http_tools.POST
	// create path and map variables
	path := s.config.HttpConfig.BasePath + "/2/file/video/ad/"

	header := make(map[string]string)
	header["Accept"] = "application/json"
	header["Access-Token"] = authAccount.AccessToken
	var request *http.Request

	formParams := url.Values{}
	header["Content-Type"] = "multipart/form-data"
	formParams["advertiser_id"] = []string{input.BaseInput.AccountId}
	formParams["video_signature"] = []string{input.Signature}
	if input.File == nil {
		return nil, fmt.Errorf("AddVideo file is empty")
	}
	localNewFile, localErr := os.Open(input.File.Name())
	if localErr != nil {
		return nil, localErr
	}
	defer localNewFile.Close()
	fileBytes, localErr := ioutil.ReadAll(localNewFile)
	if localErr != nil {
		return nil, localErr
	}
	fileName := input.File.Name()
	fileKey := "image_file"
	request, err = s.httpClient.PrepareRequest(context.Background(), path, method, nil, header,
		nil, formParams, fileName, fileBytes, fileKey)
	if err != nil {
		return nil, err
	}

	response := &AddVideoData{}
	respErr := s.httpClient.DoProcess(request, response)
	if respErr != nil {
		return nil, respErr
	}
	if response.Code != 0 {
		return nil, fmt.Errorf("response : code = %d, message = %s, request_id = %s ", response.Code,
			response.Message,
			response.RequestId)
	}

	videoInfo := response.Data
	videoID, _ := strconv.ParseInt(videoInfo.VideoId, 10, 64)
	output := &sdk.VideoAddOutput{
		VideoId:    videoID,
		FileSize:   videoInfo.FileSize,
		Width:      videoInfo.Width,
		Height:     videoInfo.Height,
		VideoUrl:   videoInfo.VideoUrl,
		Duration:   videoInfo.Duration,
		MaterialId: videoInfo.MaterialId,
	}
	return output, nil
}

type BindMaterialRequest struct {
	AdvertiserId        int64    `json:"advertiser_id,omitempty"`         // 素材归属广告主
	TargetAdvertiserIds []int64  `json:"target_advertiser_ids,omitempty"` // 待推送的广告主，数量限制：<=50
	VideoIds            []string `json:"video_ids,omitempty"`             // 视频ID，数量限制：<=50 注意：跟image_ids必须二选一、组织共享视频不可推送
	ImageIds            []string `json:"image_ids,omitempty"`             //图片ID，数量限制：<=50
}

type MaterialBindResponse struct {
	Code      int                     `json:"code"`
	Message   string                  `json:"message"`
	Data      *sdk.MaterialBindOutput `json:"data"`
	RequestId string                  `json:"request_id"`
}

// BindMaterial 素材推送
func (s *MaterialService) BindMaterial(input *sdk.MaterialBindInput) (*sdk.MaterialBindOutput, error) {
	id := formatAuthAccountID(input.BaseInput.AccountId)
	authAccount, err := account.GetAuthAccount(id)
	if err != nil {
		return nil, fmt.Errorf("BindMaterial get AuthAccount Info error accid=%s", input.BaseInput.AccountId)
	}

	method := http_tools.POST
	// create path and map variables
	path := s.config.HttpConfig.BasePath + "/2/file/material/bind/"

	header := make(map[string]string)
	header["Accept"] = "application/json"
	header["Content-Type"] = "application/json"
	header["Access-Token"] = authAccount.AccessToken

	var request *http.Request
	accID, _ := strconv.ParseInt(input.BaseInput.AccountId, 10, 64)
	if len(input.TargetAdvertiserIds) == 0 {
		return nil, fmt.Errorf("bindMaterial TargetAdvertiserIds is empty")
	}
	bindMaterial := &BindMaterialRequest{
		AdvertiserId: accID,
		TargetAdvertiserIds: input.TargetAdvertiserIds,
		VideoIds: input.VideoIds,
		ImageIds: input.ImageIds,
	}

	postBody, _ := json.Marshal(bindMaterial)

	request, err = s.httpClient.PrepareRequest(context.Background(), path, method, postBody, header,
		nil, nil, "", nil, "")
	if err != nil {
		return nil, err
	}

	response := &MaterialBindResponse{}
	respErr := s.httpClient.DoProcess(request, response)
	if respErr != nil {
		return nil, respErr
	}
	if response.Code != 0 {
		return nil, fmt.Errorf("response : code = %d, message = %s, request_id = %s ", response.Code,
			response.Message,
			response.RequestId)
	}
	return response.Data, nil
}
