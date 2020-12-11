package main
import (
	"net/http"
	"encoding/json"
	"context"
    "net/url"
    "unicode"
    "strings"
	"git.code.oa.com/going/going/cat/qzs"
	"git.code.oa.com/tme/kg_golang_proj/jce_go/proto_gdt_register"
    "strconv"
)

// ===================== rsp 格式 =====================
type AMSClickRsp struct {
	Code            int					`json:"code"`           // 0-成功 -1-失败
	Msg				string				`json:"msg"`            // 说明
}

// ==================== 处理流程 ======================
func handleAMSGdt(w http.ResponseWriter, r *http.Request) {
	qzsCtx := qzs.NewContext(context.Background())
	qzsCtx.Debug("handleAMSGdt starts")

	// [step 1.1] 获取url中的参数
    // https://debuggdtproxy.kg.qq.com/ams/gdt    测试域名
    // https://gdtproxy.kg.qq.com/ams/gdt         外网域名
    // 给AMS的监测链接(Android):
    // https://gdtproxy.kg.qq.com/ams/gdt?ad=30&channelid=xxx&click_id=__CLICK_ID__&os=__DEVICE_OS_TYPE__&promoted_object_id=__PROMOTED_OBJECT_ID__&imei=__MUID__&oaid=__OAID__&androidid=__HASH_ANDROID_ID__&click_time=__CLICK_TIME__&callback=__CALLBACK__&account_id=__ACCOUNT_ID__&account_groupid=__ADGROUP_ID__&creative_id=__AD_ID__&plan_id=__CAMPAIGN_ID__&ip=__IP__&ua=__USER_AGENT__&materialtype=xxx&songmid=xxx&userkid=xxx&adgroup_name=__ADGROUP_NAME__&plan_name=__CAMPAIGN_NAME__
    // 给AMS的监测链接(IOS):
    // https://gdtproxy.kg.qq.com/ams/gdt?ad=30&channelid=xxx&click_id=__CLICK_ID__&os=__DEVICE_OS_TYPE__&promoted_object_id=__PROMOTED_OBJECT_ID__&idfa=__MUID__&click_time=__CLICK_TIME__&callback=__CALLBACK__&account_id=__ACCOUNT_ID__&account_groupid=__ADGROUP_ID__&creative_id=__AD_ID__&plan_id=__CAMPAIGN_ID__&ip=__IP__&ua=__USER_AGENT__&materialtype=xxx&songmid=xxx&userkid=xxx&adgroup_name=__ADGROUP_NAME__&plan_name=__CAMPAIGN_NAME__
    param := r.URL.Query()
    Channel_id := param.Get("channelid")
    Imei := param.Get("imei")                               // md5之后的imei
    Idfa := param.Get("idfa")                               // md5之后的idfa
    Click_time := param.Get("click_time")                   // 点击时间(ms)
    Callback,_ := url.QueryUnescape(param.Get("callback"))  // 回传用到,唯一表示一次请求
    Os := param.Get("os")                                   // ios android
    Oaid := param.Get("oaid")                               // oaid原值
    AndroidId := param.Get("androidid")                     // md5之后的android id
    AccountId := param.Get("account_id")                    // 账户id
    AccountGroupId := param.Get("account_groupid")          // 账户组id
    CreativeId := param.Get("creative_id")                  // 创意id
    PlanId := param.Get("plan_id")                          // 计划id
    IP := param.Get("ip")
    UA := param.Get("ua")
    MaterialType := param.Get("materialtype")               // 素材类型
    SongId := param.Get("songmid")                          // 伴奏id
    UserId := param.Get("userkid")                          // uid
    AdGroupName := param.Get("adgroup_name")                // 广告组名称，用于解析channelid和素材类型
    PlanName := param.Get("plan_name")                      // 计划名称，用于解析伴奏id或者kid
    qzsCtx.Debug("request decode succ, param:%+v, channelid:%v, imei:%v, idfa:%v, click_time:%v, callback:%v, os:%v, oaid:%v, androidid:%v, AccountId:%v, AccountGroupId:%v, creativeId:%v, PlanId:%v, ip:%v, ua:%v, materialtype:%v, songid:%v, userid:%v, AdGroupName:%v, PlanName:%v", param, Channel_id, Imei, Idfa, Click_time, Callback, Os, Oaid, AndroidId, AccountId, AccountGroupId, CreativeId, PlanId, IP, UA, MaterialType, SongId, UserId, AdGroupName, PlanName)

	// [step 2] 填充 ReportClickReq
    var SvrClickReq proto_gdt_register.ReportClickReq
    SvrClickReq.Muid = Imei
    SvrClickReq.Click_time = Click_time
    SvrClickReq.Ad_company_type = proto_gdt_register.AD_COMPANY_AMS
    if Os == "ios" {
        SvrClickReq.Os_type = 1
    } else {
        SvrClickReq.Os_type =0
    }

    SvrClickReq.Channel_id = Channel_id
    tmpValue, _ := strconv.ParseUint(MaterialType, 10, 32)
    SvrClickReq.UMaterialType = uint32(tmpValue)
    // 广告组名称构成原则：appid_渠道号_拉新/拉活_素材一级标签ID_一级标签_二级标签_三级标签_版位名称_出价模式_定向_上线日期_自定义
    // 101097681_1100114954_拉新_101_修音_效果对比_竖版视频_穿山甲_OCPC激活双出价_通投_1201_自定义
    if AdGroupName != "" {
        elements := strings.Split(AdGroupName, "_")
        // 如果无法从检测链接中直接拿到channelid，则从解析后的广告组名称中拿channelid
        if len(elements) >= 2 && SvrClickReq.Channel_id == "" {
            SvrClickReq.Channel_id = elements[1]
        }
        // 如果无法从检测链接中直接拿到material type，则从解析后的广告组名称中拿material type
        if len(elements) >= 4 && SvrClickReq.UMaterialType == 0 {
            tmpValue, _ = strconv.ParseUint(elements[3], 10, 32)
            SvrClickReq.UMaterialType = uint32(tmpValue)
        }
    }

    SvrClickReq.Idfa = Idfa
    SvrClickReq.StrOaid = Oaid
    SvrClickReq.StrMAndroidId = AndroidId
    SvrClickReq.Callbackparam = Callback
    SvrClickReq.StrIp = IP
    SvrClickReq.StrUserAgent = UA

    SvrClickReq.StrKSongMid = SongId
    SvrClickReq.StrUserKid = UserId
    // 计划名称构成原则：songmid/userkid_歌曲名/kol名_素材一级标签_二级标签_三级标签_素材中文名_代理_自定义
    // 000YKeyK3ui0Fn_酒醉的蝴蝶_老歌_土味_竖版视频_街拍小姐姐_腾讯
    // 888888888_一修_kol短视频_教唱_竖版视频_街拍小姐姐_腾讯
    if PlanName != "" {
        elements := strings.Split(PlanName, "_")
        if len(elements) > 1 {
            materialInfo := elements[0]
            bMid := false
            for _,char := range materialInfo {
                // 出现了非数字的字符，则为mid；否则为uid
                if !unicode.IsDigit(char) {
                    bMid = true
                    break
                }
            }
            // 如果无法从检测链接中直接拿到mid，则从解析后的计划名称中拿mid
            if bMid && SvrClickReq.StrKSongMid == "" {
                SvrClickReq.StrKSongMid = materialInfo
            }
            // 如果无法从检测链接中直接拿到kid，则从解析后的计划名称中拿kid
            if !bMid && SvrClickReq.StrUserKid == "" {
                SvrClickReq.StrUserKid = materialInfo
            }
        }
    }

    SvrClickReq.MapExtInfo = make(map[string]string)
    SvrClickReq.StrAdAccountId = AccountId
    SvrClickReq.StrAdGroupId = AccountGroupId
    SvrClickReq.StrAdCreativeId = CreativeId
    SvrClickReq.StrAdPlanId = PlanId

	// [step 3] 先回包后处理
    // 上报点击行为
	go DoReportClick(qzsCtx, SvrClickReq)

	// [step 4] 填充rsp并回包
	var rsp AMSClickRsp
	rsp.Code = 0

	w.Header().Set("Content-Type", "application/json");
    err := json.NewEncoder(w).Encode(&rsp)
	if err != nil {
		qzsCtx.Error("encode rsp err, param:%+v\n, rsp:%+v\n, err:%v", param, rsp, err)
		return
	}
	qzsCtx.Debug("encode rsp succ, param:%+v\n, rsp:%+v\n", param, rsp)
}

