package ams

import (
	"fmt"
	"strconv"
	"time"

	"git.code.oa.com/tme-server-component/kg_growth_open/api/sdk"
	"git.code.oa.com/tme-server-component/kg_growth_open/pkg/sdk/account"
	config "git.code.oa.com/tme-server-component/kg_growth_open/pkg/sdk/config"
	"github.com/antihax/optional"
	"github.com/tencentad/marketing-api-go-sdk/pkg/api"
	"github.com/tencentad/marketing-api-go-sdk/pkg/model"
)

// AMSMaterialService AMS物料服务
type MaterialService struct {
	config *config.Config
}

// NewAMSMaterialService 获取物料服务
func NewMaterialService(sConfig *config.Config) *MaterialService {
	return &MaterialService{
		config: sConfig,
	}
}

// AddImage 增加图片上传
func (s *MaterialService) AddImage(input *sdk.ImageAddInput) (*sdk.ImagesAddOutput, error) {
	id := formatAuthAccountID(input.BaseInput.AccountId, input.BaseInput.AMSSystemType)
	authAccount, err := account.GetAuthAccount(id)
	if err != nil {
		return nil, fmt.Errorf("AddImage get AuthAccount Info error accid=%s", input.BaseInput.AccountId)
	}

	tClient := getAMSSdkClient(authAccount)
	imagesAddOpts := &api.ImagesAddOpts{}
	if input.File != nil {
		imagesAddOpts.File = optional.NewInterface(input.File)
	}
	if len(input.BytesAMS) > 0 {
		imagesAddOpts.Bytes = optional.NewString(input.BytesAMS)
	}
	if len(input.DescAMS) > 0 {
		imagesAddOpts.Description = optional.NewString(input.DescAMS)
	}
	accID, err := strconv.ParseInt(input.BaseInput.AccountId, 10, 64)
	if err != nil {
		return nil, err
	}
	response, _, err := tClient.Images().Add(*tClient.Ctx, accID, string(input.UploadType),
		input.Signature, imagesAddOpts)
	if err != nil {
		return nil, err
	}
	output := &sdk.ImagesAddOutput{
		ImageId:     response.ImageId,
		PreviewUrl:  response.PreviewUrl,
		Description: response.Description,
		Width:       response.Width,
		Height:      response.Height,
		FileSize:    response.FileSize,
		Type:        sdk.ImageType(string(response.Type_)),
		Signature:   response.Signature,
	}
	return output, err
}

// getFilter 获取过滤信息
func (s *MaterialService) getFilter(input *sdk.MaterialGetInput, isImage bool) []model.FilteringStruct {
	if input.Filtering == nil {
		return nil
	}
	preStr := "image_"
	if !isImage {
		preStr = "media_"
	}
	// Filtering
	var TFiltering []model.FilteringStruct
	mFiltering := input.Filtering
	// image_id
	if mFiltering.MaterialIds != nil {
		TFiltering = append(TFiltering, model.FilteringStruct{
			Field:    preStr + "id",
			Operator: "IN",
			Values:   &mFiltering.MaterialIds,
		})
	}
	// Width
	if mFiltering.Width > 0 {
		TFiltering = append(TFiltering, model.FilteringStruct{
			Field:    preStr + "width",
			Operator: "EQUALS",
			Values:   &[]string{strconv.FormatInt(mFiltering.Width, 10)},
		})
	}
	// Height
	if mFiltering.Height > 0 {
		TFiltering = append(TFiltering, model.FilteringStruct{
			Field:    preStr + "height",
			Operator: "EQUALS",
			Values:   &[]string{strconv.FormatInt(mFiltering.Height, 10)},
		})
	}
	// CreatedStartTime
	if len(mFiltering.CreatedStartTime) > 0 {
		TFiltering = append(TFiltering, model.FilteringStruct{
			Field:    "created_time",
			Operator: "GREATER_EQUALS",
			Values:   &[]string{mFiltering.CreatedStartTime},
		})
	}
	// CreatedEndTime
	if len(mFiltering.CreatedEndTime) > 0 {
		TFiltering = append(TFiltering, model.FilteringStruct{
			Field:    "created_time",
			Operator: "LESS_EQUALS",
			Values:   &[]string{mFiltering.CreatedEndTime},
		})
	}
	return TFiltering
}

// GetImage 获取图片信息
func (s *MaterialService) GetImage(input *sdk.MaterialGetInput) (*sdk.ImageGetOutput, error) {
	id := formatAuthAccountID(input.BaseInput.AccountId, input.BaseInput.AMSSystemType)
	authAccount, err := account.GetAuthAccount(id)
	if err != nil {
		return nil, fmt.Errorf("GetImage get AuthAccount Info error accid=%s", input.BaseInput.AccountId)
	}
	tClient := getAMSSdkClient(authAccount)
	imagesGetOpts := &api.ImagesGetOpts{
		Fields: optional.NewInterface([]string{"image_id", "width", "height", "file_size", "type", "signature",
			"source_signature", "preview_url", "source_type", "created_time", "last_modified_time"}),
	}

	if tFilter := s.getFilter(input, true); tFilter != nil && len(tFilter) > 0 {
		imagesGetOpts.Filtering = optional.NewInterface(tFilter)
	}
	accID, err := strconv.ParseInt(input.BaseInput.AccountId, 10, 64)
	if err != nil {
		return nil, err
	}
	response, _, err := tClient.Images().Get(*tClient.Ctx, accID, imagesGetOpts)
	if err != nil {
		return nil, err
	}
	imageOutput := &sdk.ImageGetOutput{}
	s.copyImageInfoToOutput(&response, imageOutput)
	return imageOutput, err
}

// copyImageInfoToOutput 拷贝图片信息
func (s *MaterialService) copyImageInfoToOutput(imageResponseData *model.ImagesGetResponseData, imageOutput *sdk.ImageGetOutput) {
	if len(*imageResponseData.List) == 0 {
		return
	}
	rList := make([]*sdk.ImageGetOutputStruct, 0, len(*imageResponseData.List))
	for i := 0; i < len(*imageResponseData.List); i++ {
		imageData := (*imageResponseData.List)[i]
		rList = append(rList, &sdk.ImageGetOutputStruct{
			ImageId:          imageData.ImageId,
			Width:            imageData.Width,
			Height:           imageData.Height,
			FileSize:         imageData.FileSize,
			ImageType:        sdk.ImageType(imageData.Type_),
			Signature:        imageData.Signature,
			Description:      imageData.Description,
			SourceSignature:  imageData.SourceSignature,
			PreviewUrl:       imageData.PreviewUrl,
			SourceType:       sdk.MaterialSourceType(imageData.SourceType),
			CreatedTime:      time.Unix(imageData.CreatedTime, 0).Format("2006-01-02 15:04:05"),
			LastModifiedTime: time.Unix(imageData.LastModifiedTime, 0).Format("2006-01-02 15:04:05"),
		})
	}
	imageOutput.List = rList
	imageOutput.PageInfo = &sdk.PageConf{
		Page:        imageResponseData.PageInfo.Page,
		PageSize:    imageResponseData.PageInfo.PageSize,
		TotalNumber: imageResponseData.PageInfo.TotalNumber,
		TotalPage:   imageResponseData.PageInfo.TotalPage,
	}
}

// GetVideo 获取视频信息
func (s *MaterialService) GetVideo(input *sdk.MaterialGetInput) (*sdk.VideoGetOutput, error) {
	id := formatAuthAccountID(input.BaseInput.AccountId, input.BaseInput.AMSSystemType)
	authAccount, err := account.GetAuthAccount(id)
	if err != nil {
		return nil, fmt.Errorf("GetVideo get AuthAccount Info error accid=%s", input.BaseInput.AccountId)
	}
	tClient := getAMSSdkClient(authAccount)
	videoGetOpts := &api.VideosGetOpts{
		Fields: optional.NewInterface([]string{"video_id", "width", "height", "video_frames", "video_fps",
			"video_codec", "video_bit_rate", "audio_codec", "audio_bit_rate", "file_size", "type", "signature",
			"system_status", "description", "preview_url", "created_time", "last_modified_time",
			"video_profile_name", "audio_sample_rate", "max_keyframe_interval", "min_keyframe_interval",
			"sample_aspect_ratio", "audio_profile_name", "scan_type", "image_duration_millisecond",
			"audio_duration_millisecond", "source_type"}),
	}
	if tFilter := s.getFilter(input, false); tFilter != nil && len(tFilter) > 0 {
		videoGetOpts.Filtering = optional.NewInterface(tFilter)
	}
	if input.Page > 0 {
		videoGetOpts.Page = optional.NewInt64(input.Page)
	}
	if input.PageSize > 0 {
		videoGetOpts.PageSize = optional.NewInt64(input.PageSize)
	}

	accountid, err := strconv.ParseInt(input.BaseInput.AccountId, 10, 64)
	if err != nil {
		return nil, err
	}
	response, _, err := tClient.Videos().Get(*tClient.Ctx, accountid, videoGetOpts)
	if err != nil {
		return nil, err
	}
	videoOutput := &sdk.VideoGetOutput{}
	s.copyVideoInfoToOutput(&response, videoOutput)
	return videoOutput, err
}

// copyVideoInfoToOutput 拷贝视频信息
func (s *MaterialService) copyVideoInfoToOutput(videoResponseData *model.VideosGetResponseData, videoOutput *sdk.VideoGetOutput) {
	if len(*videoResponseData.List) == 0 {
		return
	}
	rList := make([]*sdk.VideoGetOutputStruct, 0, len(*videoResponseData.List))
	for i := 0; i < len(*videoResponseData.List); i++ {
		videoData := (*videoResponseData.List)[i]
		rList = append(rList, &sdk.VideoGetOutputStruct{
			VideoId:                  strconv.FormatInt(videoData.VideoId, 10),
			Width:                    videoData.Width,
			Height:                   videoData.Height,
			VideoFrames:              videoData.VideoFrames,
			VideoFps:                 videoData.VideoFps,
			VideoCodec:               videoData.VideoCodec,
			VideoBitRate:             videoData.VideoBitRate,
			AudioCodec:               videoData.AudioCodec,
			AudioBitRate:             videoData.AudioBitRate,
			FileSize:                 videoData.FileSize,
			VideoType:                sdk.VideoType(videoData.Type_),
			Signature:                videoData.Signature,
			SystemStatus:             sdk.SystemStatus(videoData.SystemStatus),
			Description:              videoData.Description,
			PreviewUrl:               videoData.PreviewUrl,
			KeyFrameImageUrl:         videoData.KeyFrameImageUrl,
			CreatedTime:              time.Unix(videoData.CreatedTime, 0).Format("2006-01-02 15:04:05"),
			LastModifiedTime:         videoData.LastModifiedTime,
			VideoProfileName:         videoData.VideoProfileName,
			AudioSampleRate:          videoData.AudioSampleRate,
			MaxKeyframeInterval:      videoData.MaxKeyframeInterval,
			MinKeyframeInterval:      videoData.MinKeyframeInterval,
			SampleAspectRatio:        videoData.SampleAspectRatio,
			AudioProfileName:         videoData.AudioProfileName,
			ScanType:                 videoData.ScanType,
			ImageDurationMillisecond: videoData.ImageDurationMillisecond,
			AudioDurationMillisecond: videoData.AudioDurationMillisecond,
			SourceType:               sdk.MaterialSourceType(videoData.SourceType),
		})
	}
	videoOutput.List = rList
	videoOutput.PageInfo = &sdk.PageConf{
		Page:        videoResponseData.PageInfo.Page,
		PageSize:    videoResponseData.PageInfo.PageSize,
		TotalNumber: videoResponseData.PageInfo.TotalNumber,
		TotalPage:   videoResponseData.PageInfo.TotalPage,
	}
}

// AddVideo 增加视频上传
func (s *MaterialService) AddVideo(input *sdk.VideoAddInput) (*sdk.VideoAddOutput, error) {
	id := formatAuthAccountID(input.BaseInput.AccountId, input.BaseInput.AMSSystemType)
	authAccount, err := account.GetAuthAccount(id)
	if err != nil {
		return nil, fmt.Errorf("AddVideo get AuthAccount Info error accid=%s", input.BaseInput.AccountId)
	}
	tClient := getAMSSdkClient(authAccount)
	videoAddOpts := &api.VideosAddOpts{}
	if len(input.Desc) > 0 {
		videoAddOpts.Description = optional.NewString(input.Desc)
	}
	accID, err := strconv.ParseInt(input.BaseInput.AccountId, 10, 64)
	if err != nil {
		return nil, err
	}
	response, _, err := tClient.Videos().Add(*tClient.Ctx, accID, input.File, input.Signature, videoAddOpts)
	if err != nil {
		return nil, err
	}
	output := &sdk.VideoAddOutput{
		VideoId: response.VideoId,
	}
	return output, err
}

func (s *MaterialService) BindMaterial(input *sdk.MaterialBindInput) (*sdk.MaterialBindOutput, error) {
	panic("implement me")
}
