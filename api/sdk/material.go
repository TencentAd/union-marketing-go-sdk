package sdk

import (
	"os"
)

// ImageAddInput 图片请求参数
type ImageAddInput struct {
	BaseInput           BaseInput  `json:"base_input,omitempty"`  // 账户信息
	UploadType          UploadType `json:"upload_type,omitempty"` // 上传类型
	Signature           string     `json:"signature,omitempty"`   // 图片文件签名，使用图片文件的 md5 值
	File                *os.File   `json:"file,omitempty"`        // 被上传的图片文件，当且仅当 upload_type=UPLOAD_TYPE_FILE 时，该字段可填且必填
	ImageUrlOceanEngine string     `json:"image_url,omitempty"`   // 被上传的图片文件URL，仅头条支持
	BytesAMS            string     `json:"bytes,omitempty"`       // 图片 base64 编码，当且仅当 upload_type=UPLOAD_TYPE_BYTES 时，该字段可填且必填
	DescAMS             string     `json:"desc,omitempty"`        // 图片文件描述
}

// UploadType 上传类型
type UploadType string

const (
	UploadTypeFile  UploadType = "UPLOAD_TYPE_FILE"
	UploadTypeBytes UploadType = "UPLOAD_TYPE_BYTES"
	UploadTypeUrl   UploadType = "UPLOAD_TYPE_URL" // AMS不支持
)

// ImagesAddOutput 上传图片返回结果
type ImagesAddOutput struct {
	ImageId     string    `json:"image_id,omitempty"`
	PreviewUrl  string    `json:"preview_url,omitempty"`
	Description string    `json:"description,omitempty"`
	Width       int64     `json:"width,omitempty"`
	Height      int64     `json:"height,omitempty"`
	FileSize    int64     `json:"file_size,omitempty"`
	Type        ImageType `json:"type,omitempty"`
	Signature   string    `json:"signature,omitempty"`
}

// ImageType : 图片类型
type ImageType string

// List of ImageType
const (
	ImageTypeGif ImageType = "IMAGE_TYPE_GIF" // GIF 类型
	ImageTypeJpg ImageType = "IMAGE_TYPE_JPG" //JPG 类型
	ImageTypePng ImageType = "IMAGE_TYPE_PNG" //PNG 类型
	ImageTypeSwf ImageType = "IMAGE_TYPE_SWF" // SWF 类型
)

// MaterialGetInput 获取物料请求结构
type MaterialGetInput struct {
	BaseInput BaseInput          `json:"base_input,omitempty"` // 账户信息
	Filtering *MaterialFiltering `json:"filtering,omitempty"`  // 过滤信息
	Page      int64              `json:"page,omitempty"`       // 搜索页码，默认值：1 最小值 1，最大值 99999
	PageSize  int64              `json:"page_size,omitempty"`  // 一页显示的数据条数，默认值：10。最小值 1，最大值 500
}

// MaterialFiltering 获取物料过滤参数结构
type MaterialFiltering struct {
	Ids                    []string `json:"ids,omitempty"`                // 图片ids 数量限制：<=100  注意：image_ids、material_ids、signatures只能选择一个进行过滤
	MaterialIdsOceanEngine []string `json:"material_ids,omitempty"`       // 图片ids 数量限制：<=100  注意：image_ids、material_ids、signatures只能选择一个进行过滤,头条使用
	Signatures             []string `json:"signatures,omitempty"`         // 图片ids 数量限制：<=100  注意：image_ids、material_ids、signatures只能选择一个进行过滤
	Width                  int64    `json:"width,omitempty"`              // 图片宽度
	Height                 int64    `json:"height,omitempty"`             // 图片高度
	CreatedStartTime       string   `json:"created_start_time,omitempty"` // 根据视频上传时间进行过滤的起始时间，与end_time搭配使用，格式：yyyy-mm-dd
	CreatedEndTime         string   `json:"created_end_time,omitempty"`   // 根据视频上传时间进行过滤的截止时间，与start_time搭配使用，格式：yyyy-mm-dd
}

// ImageGetOutput 获取图片结构
type ImageGetOutput struct {
	List     []*ImageGetOutputStruct `json:"list,omitempty"`
	PageInfo *PageConf               `json:"page_info,omitempty"`
}

// ImageGetOutput 图片信息
type ImageGetOutputStruct struct {
	ImageId          string             `json:"image_id,omitempty"`           // 图片 id
	Width            int64              `json:"width,omitempty"`              // 图片宽度，单位 px
	Height           int64              `json:"height,omitempty"`             // 图片高度，单位 px
	FileSize         int64              `json:"file_size,omitempty"`          // 图片大小 单位 B(byte)
	ImageType        ImageType          `json:"image_type,omitempty"`         // 图片类型，[枚举详情]
	Signature        string             `json:"signature,omitempty"`          //图片文件签名，使用图片文件的 md5 值，用于检查上传图片文件的完整性
	Description      string             `json:"description,omitempty"`        // 图片文件描述
	SourceSignature  string             `json:"source_signature,omitempty"`   // 图片源文件签名，为图片经过裁剪前源文件的 md5 值，若该文件没有经过裁剪，source_signature 为空
	PreviewUrl       string             `json:"preview_url,omitempty"`        // 预览地址
	SourceType       MaterialSourceType `json:"source_type,omitempty"`        // 图片来源
	CreatedTime      string             `json:"created_time,omitempty"`       // 创建时间（时间戳）
	LastModifiedTime string             `json:"last_modified_time,omitempty"` // 最后修改时间（时间戳）
}

// MaterialSourceType 物料来源
type MaterialSourceType string

// List of SourceType
const (
	// AMS
	MaterialSourceTypeUnsupported   MaterialSourceType = "SOURCE_TYPE_UNSUPPORTED"     // 其他上传方式
	MaterialSourceTypeLocal         MaterialSourceType = "SOURCE_TYPE_LOCAL"           // 通过投放端本地自行上传
	MaterialSourceTypeMuse          MaterialSourceType = "SOURCE_TYPE_MUSE"            // 妙思智能制图工具
	MaterialSourceTypeAPI           MaterialSourceType = "SOURCE_TYPE_API"             // 通过 Marketing API 上传
	MaterialSourceTypeQuickDraw     MaterialSourceType = "SOURCE_TYPE_QUICK_DRAW"      // 快速制图工具
	MaterialSourceTypeVideoMakerXSJ MaterialSourceType = "SOURCE_TYPE_VIDEO_MAKER_XSJ" // 视频截图
	MaterialSourceTypeTCC           MaterialSourceType = "SOURCE_TYPE_TCC"             // 腾讯创意订制平台制作，source_reference_id（素材来源关联 id）为 TCC 订单 id
	// OceanEngine
	ADSite         MaterialSourceType = "AD_SITE"         // ad后台本地上传
	CreativeCenter MaterialSourceType = "CREATIVE_CENTER" //创意中心
	OpenApi        MaterialSourceType = "OPEN_API"        //开放平台
	SUPPLIER       MaterialSourceType = "SUPPLIER"        //即合视频
	VideoCapture   MaterialSourceType = "VIDEO_CAPTURE"   //易拍视频
	AccountPush    MaterialSourceType = "ACCOUNT_PUSH"    //推送视频
	STAR           MaterialSourceType = "STAR"            //星图视频
	CewebrityVideo MaterialSourceType = "CEWEBRITY_VIDEO" //达人视频
	OTHERS         MaterialSourceType = "OTHERS"          //其他来源
)

// VideoGetOutput 视频获取
type VideoGetOutput struct {
	List     []*VideoGetOutputStruct `json:"list,omitempty"`
	PageInfo *PageConf               `json:"page_info,omitempty"`
}

// VideoGetOutputStruct 视频信息结构
type VideoGetOutputStruct struct {
	VideoId                  string             `json:"video_id,omitempty"`
	Width                    int64              `json:"width,omitempty"`
	Height                   int64              `json:"height,omitempty"`
	VideoFrames              int64              `json:"video_frames,omitempty"`
	VideoFps                 float64            `json:"video_fps,omitempty"`
	VideoCodec               string             `json:"video_codec,omitempty"`
	VideoBitRate             int64              `json:"video_bit_rate,omitempty"`
	AudioCodec               string             `json:"audio_codec,omitempty"`
	AudioBitRate             int64              `json:"audio_bit_rate,omitempty"`
	FileSize                 int64              `json:"file_size,omitempty"`
	VideoType                VideoType          `json:"video_type,omitempty"`
	Signature                string             `json:"signature,omitempty"`
	SystemStatus             SystemStatus       `json:"system_status,omitempty"`
	Description              string             `json:"description,omitempty"`
	PreviewUrl               string             `json:"preview_url,omitempty"`
	KeyFrameImageUrl         string             `json:"key_frame_image_url,omitempty"`
	CreatedTime              string             `json:"created_time,omitempty"`
	LastModifiedTime         int64              `json:"last_modified_time,omitempty"`
	VideoProfileName         string             `json:"video_profile_name,omitempty"`
	AudioSampleRate          int64              `json:"audio_sample_rate,omitempty"`
	MaxKeyframeInterval      int64              `json:"max_keyframe_interval,omitempty"`
	MinKeyframeInterval      int64              `json:"min_keyframe_interval,omitempty"`
	SampleAspectRatio        string             `json:"sample_aspect_ratio,omitempty"`
	AudioProfileName         string             `json:"audio_profile_name,omitempty"`
	ScanType                 string             `json:"scan_type,omitempty"`
	ImageDurationMillisecond int64              `json:"image_duration_millisecond,omitempty"`
	AudioDurationMillisecond int64              `json:"audio_duration_millisecond,omitempty"`
	SourceType               MaterialSourceType `json:"source_type,omitempty"`
	ProductCatalogId         string             `json:"product_catalog_id,omitempty"`
	ProductOuterId           string             `json:"product_outer_id,omitempty"`
	SourceReferenceId        string             `json:"source_reference_id,omitempty"`
	OwnerAccountId           string             `json:"owner_account_id,omitempty"`
	VideoDurationMillisecond int64              `json:"video_duration_millisecond,omitempty"` // 视频时长
	MaterialID               int64              `json:"material_id,omitempty"`                // 物料时长
	FileName                 string             `json:"filename,omitempty"`                   // 素材文件名
	VideoLabels              string             `json:"labels,omitempty"`                     // 视频标签
}

// VideoType 视频类型
type VideoType string

const (
	MediaTypeMp4 VideoType = "MEDIA_TYPE_MP4"
	MediaTypeAvi VideoType = "MEDIA_TYPE_AVI"
	MediaTypeMov VideoType = "MEDIA_TYPE_MOV"
	MediaTypeFlv VideoType = "MEDIA_TYPE_FLV"
	VideoTypeMp4 VideoType = "VIDEO_TYPE_MP4"
	VideoTypeAVI VideoType = "VIDEO_TYPE_AVI"
	VideoTypeMOV VideoType = "VIDEO_TYPE_MOV"
)

// SystemStatus 视频转码状态
type SystemStatus string

const (
	VideoStatusValid   SystemStatus = "MEDIA_STATUS_VALID"   // 有效
	VideoStatusPending SystemStatus = "MEDIA_STATUS_PENDING" //待处理
	VideoStatusError   SystemStatus = "MEDIA_STATUS_ERROR"   //异常
)

// VideoAddInput 增加视频请求结构
type VideoAddInput struct {
	BaseInput BaseInput `json:"base_input,omitempty"` // 账户信息
	Signature string    `json:"signature,omitempty"`  // 视频文件签名
	File      *os.File  `json:"file,omitempty"`       // 被上传的视频文件，视频二进制流，支持上传的视频文件类型为：mp4、mov、avi
	Desc      string    `json:"desc,omitempty"`       // 视频文件描述
}

// VideoAddOutput 增加视频返回信息
type VideoAddOutput struct {
	VideoId    int64  `json:"video_id,omitempty"` // 视频id,AMS仅返回ID
	FileSize   int64  `json:"size,omitempty"`
	Width      int64  `json:"width,omitempty"`
	Height     int64  `json:"height,omitempty"`
	VideoUrl   string `json:"video_url,omitempty"`
	Duration   int64  `json:"duration,omitempty"`
	MaterialId int64  `json:"material_id,omitempty"`
}

// MaterialBindInput 物料推送接口
type MaterialBindInput struct {
	BaseInput           BaseInput `json:"base_input,omitempty"`  // 账户信息
	TargetAdvertiserIds []int64   `json:"target_advertiser_ids"` // 待推送的广告主，数量限制：<=50
	VideoIds            []string  `json:"video_ids"`             // 视频ID，数量限制：<=50 注意：跟image_ids必须二选一、组织共享视频不可推送
	ImageIds            []string  `json:"image_ids"`             // 图片ID，数量限制：<=50 注意：跟video_ids必须二选一
}

type MaterialBindOutput struct {
	FailList []*MaterialBindInfo `json:"fail_list"` // 推送失败列表

}
type MaterialBindInfo struct {
	VideoId                string                 `json:"video_id"`             // 推送失败的视频id
	ImageId                string                 `json:"image_id"`             // 推送失败的图片id
	TargetAdvertiserId     int64                  `json:"target_advertiser_id"` // 目标推送广告主id
	MaterialBindFailReason MaterialBindFailReason `json:"fail_reason"`          // 推送失败原因
}

type MaterialBindFailReason string

const (
	VideoBindingExisted MaterialBindFailReason = "VIDEO_BINDING_EXISTED" //(视频已存在)
	ImageBindingExisted MaterialBindFailReason = "IMAGE_BINDING_EXISTED" //(视频已存在)
)
