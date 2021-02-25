package main

import (
	"context"
	"git.code.oa.com/going/going/cat/qzs"
	"git.code.oa.com/tme/kg_golang_proj/jce_go/proto_gdt_register"
	"git.code.oa.com/tme/kg_golang_proj/pb_go/ams"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"net/http"
	"strconv"
)

func handleAMSRta(w http.ResponseWriter, r *http.Request) {
	qzsCtx := qzs.NewContext(context.Background())
	qzsCtx.Debug("handleAMSRta starts")

	// [step 1] 获取body：protobuf序列化后的二进制
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		qzsCtx.Error("ioutil.ReadAll body err: %v", err)
		return
	}
	qzsCtx.Debug("method:%v, url:%s body:%v", r.Method, r.URL.String(), body)

	// 过滤无效请求
	if r.Method != "POST" || body == nil {
		qzsCtx.Error("invalid request, skip it")
		return
	}

	// [step 2] Unmarshal 反序列化pb
	var unmarshaledReq tencent_ad_rta.RtaRequest
	err = proto.Unmarshal(body, &unmarshaledReq)
	if err != nil {
		qzsCtx.Error("proto.Unmarshal err: %v, body:%v", err, body)
		return
	}
	qzsCtx.Debug("proto.Unmarshal succ, body:%v, unmarshaledReq:%+v", body, unmarshaledReq)

	var ReqId string
	if unmarshaledReq.Id != nil {
		ReqId = *unmarshaledReq.Id
	}

	var unmarshaledRsp tencent_ad_rta.RtaResponse
	unmarshaledRsp.RequestId = new(string)
	*unmarshaledRsp.RequestId = ReqId

	// 非法请求不参竞
	if unmarshaledReq.Device == nil || unmarshaledReq.Device.Os == nil {
		qzsCtx.Error("req id:%+v, empty os", ReqId)
		unmarshaledRsp.Code = new(uint32)
		*unmarshaledRsp.Code = 1
		rsp, _ := proto.Marshal(&unmarshaledRsp)
		w.Header().Set("Content-Type", "application/x-protobuf;charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write(rsp)
		return
	}

	// 非法请求不参竞
	if *unmarshaledReq.Device.Os != tencent_ad_rta.RtaRequest_OS_IOS && *unmarshaledReq.Device.Os != tencent_ad_rta.RtaRequest_OS_ANDROID {
		qzsCtx.Error("req id:%v, invalid os type:%+v", ReqId, *unmarshaledReq.Device.Os)
		unmarshaledRsp.Code = new(uint32)
		*unmarshaledRsp.Code = 1
		rsp, _ := proto.Marshal(&unmarshaledRsp)
		w.Header().Set("Content-Type", "application/x-protobuf;charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write(rsp)
		return
	}

	// [step 3] 判断设备号和平台类型
	var uPlatform uint32
	if *unmarshaledReq.Device.Os == tencent_ad_rta.RtaRequest_OS_IOS {
		uPlatform = 2
	} else {
		uPlatform = 1
	}

	// 非法请求不参竞
	if unmarshaledReq.Device.CachedDeviceidType == nil {
		qzsCtx.Error("reqid:%v, CachedDeviceidType is nil", ReqId)
		unmarshaledRsp.Code = new(uint32)
		*unmarshaledRsp.Code = 1
		rsp, _ := proto.Marshal(&unmarshaledRsp)
		w.Header().Set("Content-Type", "application/x-protobuf;charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write(rsp)
	}

	var uDeviceType uint32
	var strMd5Id string
	if *unmarshaledReq.Device.CachedDeviceidType == tencent_ad_rta.RtaRequest_Device_IDFA_MD5 {
		uDeviceType = proto_gdt_register.IDFA_MD5

		// 非法请求不参竞
		if unmarshaledReq.Device.IdfaMd5Sum == nil {
			qzsCtx.Error("reqid:%v, CachedDeviceidType is idfa_md5 but empty value", ReqId)
			unmarshaledRsp.Code = new(uint32)
			*unmarshaledRsp.Code = 1
			rsp, _ := proto.Marshal(&unmarshaledRsp)
			w.Header().Set("Content-Type", "application/x-protobuf;charset=UTF-8")
			w.WriteHeader(http.StatusOK)
			w.Write(rsp)
			return
		}
		strMd5Id = *unmarshaledReq.Device.IdfaMd5Sum
	} else if *unmarshaledReq.Device.CachedDeviceidType == tencent_ad_rta.RtaRequest_Device_IMEI_MD5 {
		uDeviceType = proto_gdt_register.IMEI_MD5

		// 非法请求不参竞
		if unmarshaledReq.Device.ImeiMd5Sum == nil {
			qzsCtx.Error("reqid:%v, CachedDeviceidType is imei_md5 but empty value", ReqId)
			unmarshaledRsp.Code = new(uint32)
			*unmarshaledRsp.Code = 1
			rsp, _ := proto.Marshal(&unmarshaledRsp)
			w.Header().Set("Content-Type", "application/x-protobuf;charset=UTF-8")
			w.WriteHeader(http.StatusOK)
			w.Write(rsp)
			return
		}
		strMd5Id = *unmarshaledReq.Device.ImeiMd5Sum
	} else if *unmarshaledReq.Device.CachedDeviceidType == tencent_ad_rta.RtaRequest_Device_OAID_MD5 {
		uDeviceType = proto_gdt_register.OAID_MD5

		// 非法请求不参竞
		if unmarshaledReq.Device.OaidMd5Sum == nil {
			qzsCtx.Error("reqid:%v, CachedDeviceidType is oaid_md5 but empty value", ReqId)
			unmarshaledRsp.Code = new(uint32)
			*unmarshaledRsp.Code = 1
			rsp, _ := proto.Marshal(&unmarshaledRsp)
			w.Header().Set("Content-Type", "application/x-protobuf;charset=UTF-8")
			w.WriteHeader(http.StatusOK)
			w.Write(rsp)
			return
		}
		strMd5Id = *unmarshaledReq.Device.OaidMd5Sum
	} else if *unmarshaledReq.Device.CachedDeviceidType == tencent_ad_rta.RtaRequest_Device_ANDROIDID_MD5 {
		uDeviceType = proto_gdt_register.ANDROID_ID_MD5

		// 非法请求不参竞
		if unmarshaledReq.Device.AndroidIdMd5Sum == nil {
			qzsCtx.Error("reqid:%v, CachedDeviceidType is android_md5 but empty value", ReqId)
			unmarshaledRsp.Code = new(uint32)
			*unmarshaledRsp.Code = 1
			rsp, _ := proto.Marshal(&unmarshaledRsp)
			w.Header().Set("Content-Type", "application/x-protobuf;charset=UTF-8")
			w.WriteHeader(http.StatusOK)
			w.Write(rsp)
			return
		}
		strMd5Id = *unmarshaledReq.Device.AndroidIdMd5Sum
	} else {
		// 非法请求不参竞
		qzsCtx.Error("reqid:%v, no required CachedDeviceidType", ReqId)
		unmarshaledRsp.Code = new(uint32)
		*unmarshaledRsp.Code = 1
		rsp, _ := proto.Marshal(&unmarshaledRsp)
		w.Header().Set("Content-Type", "application/x-protobuf;charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write(rsp)
		return
	}
	qzsCtx.Debug("ReqId:%v\n get md5str:%v\n", ReqId, strMd5Id)

	// [step 4] 调用server判断是否需要继续投放广告
	var uBid uint32
	mapRes := make(map[uint32]bool)
	uStrategyGroupId := conf.AMS.StrategyGroupId
	iRet := DoCheckBid(qzsCtx, ReqId, strMd5Id, uPlatform, uDeviceType, uStrategyGroupId, proto_gdt_register.AD_COMPANY_AMS, mapRes, &uBid)
	if iRet != 0 {
		uBid = 1
		for _, StrategyId := range uStrategyGroupId {
			mapRes[StrategyId] = true
		}
		qzsCtx.Error("DoCheckBid err, reqId:%v\n iRet:%v\n uBid:%v\n mapRes:%+v", ReqId, iRet, uBid, mapRes)
	} else {
		qzsCtx.Debug("DoCheckBid succ, ReqId:%v\n uBid:%v\n mapRes:%+v", ReqId, uBid, mapRes)
	}

	// [step 5] 填充unmarshaledRsp
	// code=1 所有账户均不参竞
	// OutTargetId 表示需要参竞的策略组
	// code=0 & OutTargetId=nil 所有账户都参竞
	// code=0 & OutTargetId!=nil OutTargetId绑定的账户参竞
	unmarshaledRsp.Code = new(uint32)
	*unmarshaledRsp.Code = 1
	for uStrategyId, bBid := range mapRes {
		if bBid == true {
			unmarshaledRsp.OutTargetId = append(unmarshaledRsp.OutTargetId, string(strconv.Itoa(int(uStrategyId))))
		}
	}

	if len(unmarshaledRsp.OutTargetId) > 0 {
		*unmarshaledRsp.Code = 0
	}

	// [step 6] 序列化rsp并回包
	rsp, err := proto.Marshal(&unmarshaledRsp)
	if err != nil {
		qzsCtx.Error("reqid:%v\n proto.Marshal rsp err:", ReqId, err)
		return
	} else {
		qzsCtx.Debug("proto.Marshal succ, reqid:%v\n unmarshaledRsp:%+v\n rsp:%v", ReqId, unmarshaledRsp, rsp)
	}

	w.Header().Set("Content-Type", "application/x-protobuf;charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(rsp)
}

func GetDeviceId(ctx *qzs.Context, ReqId string, uPlatform uint32, Req tencent_ad_rta.RtaRequest) (strMd5Id string) {
	// [step 1] 先把可疑设备数组转成map，方便查询
	mapDoubtfulIds := make(map[tencent_ad_rta.RtaRequest_Device_DeviceIdTag]int)
	for i, v := range Req.Device.DoubtfulIdsList {
		mapDoubtfulIds[v] = i
	}

	// [step 2] 按照优先级获取设备id
	// 如果是ios且idfa不可疑，则取idfa
	// 如果是android，按照不可疑设备的imei>oaid>android
	if uPlatform == 2 {
		// 判断idfa是否为空或者是否可疑
		if Req.Device.IdfaMd5Sum == nil {
			ctx.Error("ios platform but empty idfa, reqid:%v", ReqId)
			return strMd5Id
		}
		if _, ok := mapDoubtfulIds[tencent_ad_rta.RtaRequest_Device_IDFA_MD5_DOUBTFUL]; ok {
			ctx.Error("ios platform but idfa:%s is doubtul, reqid:%v", *Req.Device.IdfaMd5Sum, ReqId)
			return strMd5Id
		} else {
			ctx.Debug("ios platform with trustful idfa:%v, reqid:%v", *Req.Device.IdfaMd5Sum, ReqId)
			strMd5Id = *Req.Device.IdfaMd5Sum
			return strMd5Id
		}
	} else {
		// 判断imei是否为空以及是否可疑
		if Req.Device.ImeiMd5Sum == nil {
			ctx.Error("android platform but empty imei, reqid:%v, try oaid", ReqId)
		} else if _, ok := mapDoubtfulIds[tencent_ad_rta.RtaRequest_Device_IMEI_MD5_DOUBTFUL]; ok {
			ctx.Error("android platform but imei:%s is doubtful, reqid:%v, try oaid", *Req.Device.ImeiMd5Sum, ReqId)
		} else {
			ctx.Debug("android platform with trustful imei:%v, reqid:%v", *Req.Device.ImeiMd5Sum, ReqId)
			strMd5Id = *Req.Device.ImeiMd5Sum
			return strMd5Id
		}

		// 判断oaid是否为空以及是否可疑
		if Req.Device.OaidMd5Sum == nil {
			ctx.Error("android platform but empty oaid, reqid:%v, try android id", ReqId)
		} else if _, ok := mapDoubtfulIds[tencent_ad_rta.RtaRequest_Device_OAID_MD5_DOUBTFUL]; ok {
			ctx.Error("android platform but oaid:%v is doubtful, reqid:%v, try android id", *Req.Device.OaidMd5Sum, ReqId)
		} else {
			ctx.Debug("android platform with trustful oaid:%v, reqid:%v", *Req.Device.OaidMd5Sum, ReqId)
			strMd5Id = *Req.Device.OaidMd5Sum
			return strMd5Id
		}

		// 判断android id是否为空以及是否可疑
		if Req.Device.AndroidIdMd5Sum == nil {
			ctx.Error("android platform but empty android id, reqid:%v", ReqId)
			return strMd5Id
		} else if _, ok := mapDoubtfulIds[tencent_ad_rta.RtaRequest_Device_ANDROIDID_MD5_DOUBTFUL]; ok {
			ctx.Error("android platform but android id:%v is doubtful, reqid:%v", *Req.Device.AndroidIdMd5Sum, ReqId)
		} else {
			ctx.Debug("android platform with trustful android id:%v, reqid:%v", *Req.Device.AndroidIdMd5Sum, ReqId)
			strMd5Id = *Req.Device.AndroidIdMd5Sum
		}
	}

	return strMd5Id
}

func DoCheckBid(qzsCtx *qzs.Context, strReqId string, strMd5Id string, uPlatform uint32, uDeviceType uint32,
	StrategyGroupId []uint32, uCompanyid uint32, mapRes map[uint32]bool, pBid *uint32) (iRet int) {
	var req proto_gdt_register.RTAReq
	var rsp proto_gdt_register.RTARsp

	req.StrDeviceId = strMd5Id
	req.UPlatformType = uPlatform
	req.UCompanyId = uCompanyid
	req.StrReqId = strReqId
	req.UDeviceType = uDeviceType

	if len(StrategyGroupId) == 0 {
		req.VctStrategyGroupId = conf.AMS.StrategyGroupId
	} else {
		for _, val := range StrategyGroupId {
			req.VctStrategyGroupId = append(req.VctStrategyGroupId, val)
		}
	}

	callDesc := qzh.CallDesc{
		CmdId:       proto_gdt_register.MAIN_CMD_GDT_REGISTER,
		SubCmdId:    proto_gdt_register.CMD_GDT_REGISTER_SVR_RTA_REQUEST,
		AppProtocol: "qza",
	}
	var authInfo qzh.AuthInfo

	qc := qzh.NewQzClient(callDesc, authInfo, "gdt_register_svr", &req, &rsp)

	// 分渠道进行模调上报
	qc.InterfaceID = conf.InterfaceID.InterfaceID_AMS

	beginTime := time.Now()
	err := qc.Do(qzsCtx)
	if err != nil {
		err := err.(qzh.QzError)
		iRet = -1
		qzsCtx.Error("qc Do error, err:%+v, \n=============req:%+v \n", err, req)
	} else {
		qzsCtx.Debug("qc Do success\n===========req:%+v \n rsp:%+v\n", req, rsp)
	}
	endTime := time.Now()
	cost := endTime.Sub(beginTime)

	// 总数单独报一份
	if iRet == 0 {
		mic.Report(proto_gdt_register.RTA_WRITE_SVR_MOD_ID, proto_gdt_register.IF_TOTAL_REGISTER_SVR_RTA_REQUEST, "255.255.255.255", 0, 0, cost)
	} else {
		mic.Report(proto_gdt_register.RTA_WRITE_SVR_MOD_ID, proto_gdt_register.IF_TOTAL_REGISTER_SVR_RTA_REQUEST, "255.255.255.255", iRet, 1, cost)
	}

	for k, v := range rsp.MapGroupId2Res {
		mapRes[k] = v
	}
	*pBid = rsp.UStatus
	return iRet
}
