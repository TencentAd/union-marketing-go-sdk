package ocean_engine

import "git.code.oa.com/tme-server-component/kg_growth_open/api/sdk"

type AddImageData struct {
	Code      int           `json:"code"`
	Message   string        `json:"message"`
	Data      *AddImageInfo `json:"data"`
	RequestId string        `json:"request_id"`
}

type AddImageInfo struct {
	ImageId    string `json:"id,omitempty"`
	FileSize   int64  `json:"size,omitempty"`
	Width      int64  `json:"width,omitempty"`
	Height     int64  `json:"height,omitempty"`
	PreviewUrl string `json:"url,omitempty"`
	Format     string `json:"format,omitempty"`
	Signature  string `json:"signature,omitempty"`
	MaterialId int64  `json:"material_id,omitempty"`
}

type GetMaterialData struct {
	Code      int              `json:"code"`
	Message   string           `json:"message"`
	Data      *GetMaterialList `json:"data"`
	RequestId string           `json:"request_id"`
}
type GetMaterialList struct {
	MaterialList []*GetMaterialInfo `json:"list"`
	PageInfo     *sdk.PageConf      `json:"page_info"`
}

type GetMaterialInfo struct {
	Id          string `json:"id,omitempty"`
	FileSize    int64  `json:"size,omitempty"`
	Width       int64  `json:"width,omitempty"`
	Height      int64  `json:"height,omitempty"`
	PreviewUrl  string `json:"url,omitempty"`
	Format      string `json:"format,omitempty"`
	Signature   string `json:"signature,omitempty"`
	MaterialId  int64  `json:"material_id,omitempty"`
	CreatedTime string `json:"create_time,omitempty"` // 素材的上传时间，格式："yyyy-mm-dd HH:MM:SS"
	FileName    string `json:"filename,omitempty"`    // 素材的文件名
	PosterUrl 	string	`json:"poster_url,omitempty"`// 视频首帧截图，仅限同主体进行素材预览查看，若非同主体会返回“素材所属主体与开发者主体不一致无法获取URL”，链接1小时过期
	BitRate 	int64	`json:"bit_rate,omitempty"`// 视频码率，单位bps
	Duration   	float32	`json:"duration,omitempty"`// 视频时长
	Source 	string	`json:"source,omitempty"`// 视频素材来源，详见【附录-素材来源】
	Labels 	string	`json:"labels,omitempty"`// 视频标签
}

type AddVideoData struct {
	Code      int          `json:"code"`
	Message   string       `json:"message"`
	Data      *AddVideoInfo `json:"data"`
	RequestId string       `json:"request_id"`
}

type AddVideoInfo struct {
	VideoId    string `json:"video_id,omitempty"`
	FileSize   int64  `json:"size,omitempty"`
	Width      int64  `json:"width,omitempty"`
	Height     int64  `json:"height,omitempty"`
	VideoUrl   string `json:"video_url,omitempty"`
	Duration   int64  `json:"duration,omitempty"`
	MaterialId int64  `json:"material_id,omitempty"`
}
