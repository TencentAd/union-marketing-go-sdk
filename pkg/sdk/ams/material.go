package ams

import (
	"strconv"

	"git.code.oa.com/tme-server-component/kg_growth_open/api/sdk"
	sdkconfig "git.code.oa.com/tme-server-component/kg_growth_open/pkg/sdk/config"
	"github.com/antihax/optional"
	"github.com/tencentad/marketing-api-go-sdk/pkg/api"
	"github.com/tencentad/marketing-api-go-sdk/pkg/model"
)

type AMSMaterialService struct {
	config *sdkconfig.Config
}

func NewAMSMaterialService(sConfig *sdkconfig.Config) *AMSMaterialService {
	return &AMSMaterialService{
		config: sConfig,
	}
}

// 增加图片上传
func (t *AMSMaterialService) AddImage(input *sdk.ImageAddInput) (*sdk.ImagesAddOutput, error) {
	tClient := getAMSSdkClient(&input.BaseInput)
	imagesAddOpts := &api.ImagesAddOpts{}
	if input.File != nil {
		imagesAddOpts.File = optional.NewInterface(input.File)
	}
	if len(input.Bytes) > 0 {
		imagesAddOpts.Bytes = optional.NewString(input.Bytes)
	}
	if len(input.Desc) > 0 {
		imagesAddOpts.Description = optional.NewString(input.Desc)
	}

	response, _, err := tClient.Images().Add(*tClient.Ctx, input.BaseInput.AccountId, string(input.UploadType), input.Signature, imagesAddOpts)
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

var TMaterialFilterMax = 4

func (t *AMSMaterialService) getFilter(input *sdk.MaterialGetInput) interface{} {
	if input.Filtering == nil {
		return nil
	}
	// Filtering
	TFiltering := make([]model.FilteringStruct, 0, TMaterialFilterMax)
	// image_id
	imageFiltering := input.Filtering.(*sdk.MaterialFiltering)
	if len(imageFiltering.MaterialIds) > 0 {
		TFiltering = append(TFiltering, model.FilteringStruct{
			Field:    "image_id",
			Operator: "IN",
			Values:   &imageFiltering.MaterialIds,
		})
	}
	// Width
	if imageFiltering.Width > 0 {
		TFiltering = append(TFiltering, model.FilteringStruct{
			Field:    "image_width",
			Operator: "EQUALS",
			Values:   &[]string{strconv.FormatInt(imageFiltering.Width, 10)},
		})
	}

	// Height
	if imageFiltering.Height > 0 {
		TFiltering = append(TFiltering, model.FilteringStruct{
			Field:    "image_height",
			Operator: "EQUALS",
			Values:   &[]string{strconv.FormatInt(imageFiltering.Height, 10)},
		})
	}

	// CreatedStartTime
	if len(imageFiltering.CreatedStartTime) > 0 {
		TFiltering = append(TFiltering, model.FilteringStruct{
			Field:    "created_time",
			Operator: "GREATER_EQUALS",
			Values:   &[]string{imageFiltering.CreatedStartTime},
		})
	}

	// CreatedEndTime
	if len(imageFiltering.CreatedEndTime) > 0 {
		TFiltering = append(TFiltering, model.FilteringStruct{
			Field:    "created_time",
			Operator: "LESS_EQUALS",
			Values:   &[]string{imageFiltering.CreatedEndTime},
		})
	}
	return TFiltering
}

// 获取图片信息
func (t *AMSMaterialService) GetImage(input *sdk.MaterialGetInput) (*sdk.ImageGetOutput, error) {
	tClient := getAMSSdkClient(&input.BaseInput)
	imagesGetOpts := &api.ImagesGetOpts{
		Fields: optional.NewInterface([]string{"image_id", "width", "height", "file_size", "type", "signature", "source_signature", "preview_url", "source_type", "created_time", "last_modified_time"}),
	}
	tFilter := t.getFilter(input)
	if tFilter != nil {
		imagesGetOpts.Filtering = optional.NewInterface(tFilter)
	}

	response, _, err := tClient.Images().Get(*tClient.Ctx, input.BaseInput.AccountId, imagesGetOpts)
	if err != nil {
		return nil, err
	}
	imageOutput := &sdk.ImageGetOutput{}
	t.copyImageInfoToOutput(&response, imageOutput)
	return imageOutput, err
}

func (t *AMSMaterialService) copyImageInfoToOutput(imageResponseData *model.ImagesGetResponseData, imageOutput *sdk.ImageGetOutput) {
	if len(*imageResponseData.List) == 0 {
		return
	}
	rList := make([]sdk.ImageGetOutputStruct, 0, len(*imageResponseData.List))
	for i := 0; i < len(*imageResponseData.List); i++ {
		imageData := (*imageResponseData.List)[i]
		rList = append(rList, sdk.ImageGetOutputStruct{
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
			CreatedTime:      imageData.CreatedTime,
			LastModifiedTime: imageData.LastModifiedTime,
		})
	}
	imageOutput.List = &rList
	imageOutput.PageInfo = &sdk.PageConf{
		Page:        imageResponseData.PageInfo.Page,
		PageSize:    imageResponseData.PageInfo.PageSize,
		TotalNumber: imageResponseData.PageInfo.TotalNumber,
		TotalPage:   imageResponseData.PageInfo.TotalPage,
	}
}

// 获取视频信息
func (t *AMSMaterialService) GetVideo(input *sdk.MaterialGetInput) (*sdk.VideoGetOutput, error) {
	tClient := getAMSSdkClient(&input.BaseInput)
	videoGetOpts := &api.VideosGetOpts{
		Fields: optional.NewInterface([]string{"video_id", "width", "height", "video_frames", "video_fps", "video_codec", "video_bit_rate", "audio_codec", "audio_bit_rate", "file_size", "type", "signature", "system_status", "description", "preview_url", "created_time", "last_modified_time", "video_profile_name", "audio_sample_rate", "max_keyframe_interval", "min_keyframe_interval", "sample_aspect_ratio", "audio_profile_name", "scan_type", "image_duration_millisecond", "audio_duration_millisecond", "source_type"}),
	}
	tFilter := t.getFilter(input)
	if tFilter != nil {
		videoGetOpts.Filtering = optional.NewInterface(tFilter)
	}
	if input.Page > 0 {
		videoGetOpts.Page = optional.NewInt64(input.Page)
	}
	if input.PageSize > 0 {
		videoGetOpts.PageSize = optional.NewInt64(input.PageSize)
	}
	response, _, err := tClient.Videos().Get(*tClient.Ctx, input.BaseInput.AccountId, videoGetOpts)
	if err != nil {
		return nil, err
	}
	videoOutput := &sdk.VideoGetOutput{}
	t.copyVideoInfoToOutput(&response, videoOutput)
	return videoOutput, err
}

func (t *AMSMaterialService) copyVideoInfoToOutput(videoResponseData *model.VideosGetResponseData, videoOutput *sdk.VideoGetOutput) {
	if len(*videoResponseData.List) == 0 {
		return
	}
	rList := make([]sdk.VideoGetOutputStruct, 0, len(*videoResponseData.List))
	for i := 0; i < len(*videoResponseData.List); i++ {
		videoData := (*videoResponseData.List)[i]
		rList = append(rList, sdk.VideoGetOutputStruct{
			VideoId:                  videoData.VideoId,
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
			CreatedTime:              videoData.CreatedTime,
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
	videoOutput.List = &rList
	videoOutput.PageInfo = &sdk.PageConf{
		Page:        videoResponseData.PageInfo.Page,
		PageSize:    videoResponseData.PageInfo.PageSize,
		TotalNumber: videoResponseData.PageInfo.TotalNumber,
		TotalPage:   videoResponseData.PageInfo.TotalPage,
	}
}

// 增加视频上传
func (t *AMSMaterialService) AddVideo(input *sdk.VideoAddInput) (*sdk.VideoAddOutput, error) {
	tClient := getAMSSdkClient(&input.BaseInput)
	videoAddOpts := &api.VideosAddOpts{}
	if len(input.Desc) > 0 {
		videoAddOpts.Description = optional.NewString(input.Desc)
	}
	response, _, err := tClient.Videos().Add(*tClient.Ctx, input.BaseInput.AccountId, input.File, input.Signature, videoAddOpts)
	output := &sdk.VideoAddOutput{
		VideoId: response.VideoId,
	}
	return output,err
}
