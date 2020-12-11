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

