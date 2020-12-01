//create by create_spp_server.sh
/**
 * 名称: 全民K歌
 * (c) Copyright 2020 wrenli@tencent.com. All Rights Reserved.
 *
 * change log:
 * Author: wrenli
 *   Date: 2020-12-01
 *	 Desc:
 */
#include "rta_write_svr.h"

#define g_conf (*CSingleton<CServerConf>::instance())
using namespace proto_gdt_register;

bool RtaWriteSvr::g_bInit = false;

// 微线程初始化
int RtaWriteSvr::InitService()
{
    return 0;
}

// 进程初始化
int RtaWriteSvr::Init(const char* conf_file)
{
    if (g_bInit)
    {
        return 0;
    }

    if (NULL == conf_file)
    {
        ERROR(" load_conf error! [params is null]\n");
        return -1;
    }

    if (g_conf.ParseFile(conf_file) != 0)
    {
        ERROR("g_conf.ParseFile fail, file:%s", conf_file);
        return -2;
    }

    API_Logapi_Init("rta_write_svr", RTA_WRITE_SVR_MOD_ID);
    g_bInit = true;
    return 0;
}

RtaWriteSvr::RtaWriteSvr()
{
    unsigned int tmem_bid = atoll(g_conf["tmem.bid"].c_str());
    unsigned int tmem_mid = atoll(g_conf["tmem.mid"].c_str());
    unsigned int tmem_cid = atoll(g_conf["tmem.cid"].c_str());
    int iRet = m_tmem_obj.fnInitL5(tmem_bid, tmem_mid, tmem_cid, atoll(g_conf["tmem.timeout_ms"].c_str()));
    if (iRet != 0)
    {
        ERROR("m_tmem_obj.fnInitL5 error, bid:%u, mid:%u, cid:%u, iRet:%d", tmem_bid, tmem_mid, tmem_cid, iRet);
    }
    else
    {
        DEBUG("m_tmem_obj.fnInitL5 succ, bid:%u, mid:%u, cid:%u", tmem_bid, tmem_mid, tmem_cid);
    }
}

RtaWriteSvr::~RtaWriteSvr()
{
}

// 打解包函数
template<typename REQUEST, typename RESPONSE, typename CB_FN>
int RtaWriteSvr::ProcessRequest(char* body_buf, int body_len, char* res_buf, int& res_len, CB_FN fn)
{
    int proc_ret = 0;
    do
    {
        REQUEST req;
        RESPONSE rsp;
        int ret = req.Decode((uint8_t*)body_buf, &body_len, NULL);
        if (ret != 0)
        {
            ERROR("decode req failed, ret:%d, body_len:%d", ret, body_len);
            proc_ret = MUSIC_CODE_COMM_ServUnpack;
            res_len = 0;
            break;
        }
        ret = fn(req, rsp);
        if (ret != 0)
        {
            ERROR("process failed, ret:%d", ret);
            proc_ret = ret;
        }
        ret = rsp.Encode((uint8_t *)res_buf, &res_len, NULL);
        if (ret != 0)
        {
            ERROR("encode rsp failed,ret:%d, res_len:%d", ret, res_len);
            proc_ret = MUSIC_CODE_COMM_ServPack;
            res_len = 0;
            break;
        }
    }
    while (0);
    return proc_ret;
}

// 打解包函数
template<typename REQUEST, typename RESPONSE, typename CB_FN>
int RtaWriteSvr::ProcessRequestNoRsp(char* body_buf, int body_len, CB_FN fn)
{
    int proc_ret = 0;
    do
    {
        REQUEST req;
        int ret = req.Decode((uint8_t*)body_buf, &body_len, NULL);
        if (ret != 0)
        {
            ERROR("decode req failed, ret:%d, body_len:%d", ret, body_len);
            proc_ret = MUSIC_CODE_COMM_ServUnpack;
            break;
        }
        ret = fn(req);
        if (ret != 0)
        {
            ERROR("process failed, ret:%d", ret);
            proc_ret = ret;
        }
      }
    while (0);
    return proc_ret;
}

int RtaWriteSvr::DealRequest(short iCmd, short iSubCmd, char* pBodyBuf, int iBodyLen, char* pResBuf, int& iLeftSize)
{
    int iRet = 0;   //最终的返回值 将保存在qza头里
    iRet = InitService();
    if (iRet != 0)
    {
        ERROR("InitService failed, ret:%d", iRet);
        iLeftSize = 0;
        return iRet;
    }

    /**
    *注意！！！：返回前一定保证设置iLeftSize为正确回包大小，否则导致回包过大（ProcessRequest里设置了）
    */
    switch (iSubCmd)
    {
        case SUB_CMD_SET_STRATEGY_INFO:
            iRet = ProcessRequest<proto_gdt_register::SetStrategyInfoReq, proto_gdt_register::SetStrategyInfoRsp,
                tr1::function<int(proto_gdt_register::SetStrategyInfoReq&, proto_gdt_register::SetStrategyInfoRsp&)> >
                    (pBodyBuf, iBodyLen, pResBuf, iLeftSize, tr1::bind(&RtaWriteSvr::DoSetStrategyInfo, this, tr1::placeholders::_1, tr1::placeholders::_2));
            break;
        case SUB_CMD_SET_CLICK_INFO:
            iRet = ProcessRequest<proto_gdt_register::SetClickInfoReq, proto_gdt_register::SetClickInfoRsp,
                tr1::function<int(proto_gdt_register::SetClickInfoReq&, proto_gdt_register::SetClickInfoRsp&)> >
                    (pBodyBuf, iBodyLen, pResBuf, iLeftSize, tr1::bind(&RtaWriteSvr::DoSetClickInfo, this, tr1::placeholders::_1, tr1::placeholders::_2));
            break;
        case SUB_CMD_SET_LOGIN_INFO:
            iRet = ProcessRequest<proto_gdt_register::SetDauLoginReq, proto_gdt_register::SetDauLoginRsp,
                 tr1::function<int(proto_gdt_register::SetDauLoginReq&, proto_gdt_register::SetDauLoginRsp&)> >
                     (pBodyBuf, iBodyLen, pResBuf, iLeftSize, tr1::bind(&RtaWriteSvr::DoSetLoginInfo, this, tr1::placeholders::_1, tr1::placeholders::_2));
            break;
        default:
            iLeftSize = 0;
            iRet = MUSIC_CODE_COMM_ParamInvalid;
            ERROR("unknown sub cmd:%d", iSubCmd);
            break;
    }

    return iRet;
}

int RtaWriteSvr::PostProcess( short iCmd, short iSubCmd, char* pBodyBuf, int iBodyLen)
{
    return 0;
}

int RtaWriteSvr::DoSetStrategyInfo(proto_gdt_register::SetStrategyInfoReq& req, proto_gdt_register::SetStrategyInfoRsp& rsp)
{
    int iRet = 0;
    string strImei = req.strImei;
    string strIdfa = req.strIdfa;
    unsigned int uStrategyId = req.uStrategyId;
    if (strImei != "")
    {
        iRet = UpdateRtaStrategyInfo(strImei, uStrategyId, req.uTimestamp);
        if (iRet != 0)
        {
            rsp.status = iRet;
            ERROR("UpdateRtaStrategyInfo err, req:%s, iRet:%d", klib::jce2str(req).c_str(), iRet);
            return iRet;
        }
        DEBUG("UpdateRtaStrategyInfo succ, req:%s", klib::jce2str(req).c_str());
    }
    if (strIdfa != "")
    {
        iRet = UpdateRtaStrategyInfo(strIdfa, uStrategyId, req.uTimestamp);
        if (iRet != 0)
        {
            rsp.status = iRet;
            ERROR("UpdateRtaStrategyInfo err, req:%s, iRet:%d", klib::jce2str(req).c_str(), iRet);
            return iRet;
        }
        DEBUG("UpdateRtaStrategyInfo succ, req:%s", klib::jce2str(req).c_str());
    }
    return iRet;
}

int RtaWriteSvr::UpdateRtaStrategyInfo(const string& strDeviceId, const unsigned uStrategyId, const unsigned uTimestamp)
{
    int iRet = 0;

    string strRTACacheKey = RTA_INFO_CACHE_KEY + strDeviceId;
    string dataRTA;
    int casRTA = 0;
    proto_gdt_register::RTAInfoCache stRTAInfoCache;

    iRet = m_tmem_obj.Get(strRTACacheKey, dataRTA, casRTA, 0, RTA_WRITE_SVR_TMEM_MOD_ID);
    if (iRet == -13106 || iRet == -13200)
	{
		DEBUG("iRet:%d, key:%s does not exist in rta cache ckv", iRet, strRTACacheKey.c_str());
	}
	else if (0 != iRet)
	{
		ERROR("m_tmem_obj Get err, key:%s, iRet:%d, errMsg:%s", strRTACacheKey.c_str(), iRet, m_tmem_obj.GetLastError());
        return iRet;
	}
    else
    {
        int iLen = dataRTA.length();
        iRet = stRTAInfoCache.Decode((uint8_t*)dataRTA.c_str(), &iLen, NULL);
        if (iRet != 0)
        {
            ERROR("RTAInfoCache.Decode error, iRet:%d, key:%s", iRet, strRTACacheKey.c_str());
        }
        else
        {
            DEBUG("m_tmem_obj Get succ, key:%s, stRTACacheInfo:%s", strRTACacheKey.c_str(), klib::jce2str(stRTAInfoCache).c_str());
        }
    }
    if (stRTAInfoCache.mapUpdateTime[uStrategyId] < uTimestamp)
    {
        stRTAInfoCache.mapUpdateTime[uStrategyId] = uTimestamp;
    }

    char szBuff[102400];
    int iLen = sizeof(szBuff);
    iRet = stRTAInfoCache.Encode((uint8_t*)szBuff, &iLen, NULL);
    if (iRet != 0)
    {
        // 避免编码失败，疯狂重试
        ERROR("stRTAInfoCache.Encode error, strDeviceId:%s, iRet:%d", strDeviceId.c_str(), iRet);
        return 0;
    }

    string data;
    data.assign(szBuff, iLen);
    unsigned int expire = (unsigned)atoll(g_conf["strategy_info_"+klib::ToStr(uStrategyId)+".expire_time_sec"].c_str());
    iRet = m_tmem_obj.Set(strRTACacheKey, data, casRTA, 0, RTA_WRITE_SVR_TMEM_MOD_ID, expire);
    DEBUG("m_tmem_obj.Set succ key:%s m_tmem_expire:%d", strRTACacheKey.c_str(), expire);
    if (iRet != 0)
    {
        ERROR("m_tmem_obj.Set error, iRet:%d, strRTACacheKey:%s, error:%s", iRet, strRTACacheKey.c_str(), m_tmem_obj.GetLastError());
        return iRet;
    }

    return iRet;
}

int RtaWriteSvr::DoSetClickInfo(SetClickInfoReq &stReq, SetClickInfoRsp &stRsp)
{
    int iRet = 0;

    iRet = SetClickTime(stReq.stClickMsg.strDeviceId, stReq.stClickMsg.uClickTime);
    if (iRet)
    {
        stRsp.iResult = iRet;
        ERROR("SetClickTime err, req:%s, rsp:%s, iRet:%d", klib::jce2str(stReq).c_str(), klib::jce2str(stRsp).c_str(), iRet);
        return iRet;
    }

    DEBUG("SetClickTime succ, req:%s", klib::jce2str(stReq).c_str());
    return iRet;
}

int RtaWriteSvr::SetClickTime(const string &strDeviceId, const unsigned uClickTime)
{
    int iRet = 0;

    string strRTACacheKey = RTA_INFO_CACHE_KEY + strDeviceId;
    string dataRTA;
    RTAInfoCache stRTAInfoCache;
    int casRTA = 0;

    iRet = m_tmem_obj.Get(strRTACacheKey, dataRTA, casRTA, 0, RTA_WRITE_SVR_TMEM_MOD_ID);
    if (iRet == -13106 || iRet == -13200)
    {
        DEBUG("iRet:%d, key:%s does not exist in rta cache ckv", iRet, strRTACacheKey.c_str());
    }
    else if (0 != iRet)
    {
        ERROR("m_tmem_obj Get err, key:%s, iRet:%d, errMsg:%s", strRTACacheKey.c_str(), iRet, m_tmem_obj.GetLastError());
        return iRet;
    }
    else
    {
        int iLen = dataRTA.length();
        iRet = stRTAInfoCache.Decode((uint8_t*)dataRTA.c_str(), &iLen, NULL);
        if (iRet != 0)
        {
            ERROR("RTAInfoCache.Decode error, iRet:%d, key:%s", iRet, strRTACacheKey.c_str());
        }
        else
        {
            DEBUG("m_tmem_obj Get succ, key:%s, stRTACacheInfo:%s", strRTACacheKey.c_str(), klib::jce2str(stRTAInfoCache).c_str());
        }
    }

    if (IsSameDay(stRTAInfoCache.uLastClickTime, uClickTime))
    {
        DEBUG("key:%s, lastClickTime=%u is the same day with click time=%u", strRTACacheKey.c_str(), stRTAInfoCache.uLastClickTime, uClickTime);
        return 0;
    }
    stRTAInfoCache.uLastClickTime = uClickTime;

    char szBuff[102400];
    int iLen = sizeof(szBuff);
    iRet = stRTAInfoCache.Encode((uint8_t*)szBuff, &iLen, NULL);
    if (iRet != 0)
    {
        // 避免编码失败，疯狂重试
        ERROR("key:%s RTAInfoCache.Encode error, stRTAInfoCache:%s, iRet:%d",
                strRTACacheKey.c_str(), klib::jce2str(stRTAInfoCache).c_str(), iRet);
        return 0;
    }

    string data;
    data.assign(szBuff, iLen);

    iRet = m_tmem_obj.Set(strRTACacheKey, data, casRTA, 0,
            RTA_WRITE_SVR_TMEM_MOD_ID, atoll(g_conf["tmem_rta_ckv.expire_time_second"].c_str()));
    if (0 != iRet)
    {
        ERROR("m_tmem_obj Set err, key:%s, iRet:%d, errMsg:%s", strRTACacheKey.c_str(), iRet, m_tmem_obj.GetLastError());
    }
    else
    {
        DEBUG("m_tmem_obj Set succ, key:%s, value:%s", strRTACacheKey.c_str(), klib::jce2str(stRTAInfoCache).c_str());
    }

    return iRet;
}

int RtaWriteSvr::DoSetLoginInfo(SetDauLoginReq &stReq, SetDauLoginRsp &stRsp)
{
    int iRet = 0;

    iRet = SetLoginInfo(stReq.stDauLoginMsg.strDeviceId, stReq.stDauLoginMsg.uLoginTime);
    if (iRet)
    {
        stRsp.iResult = iRet;
        ERROR("SetLoginInfo err, req:%s, rsp:%s, iRet:%d", klib::jce2str(stReq).c_str(), klib::jce2str(stRsp).c_str(), iRet);
        return iRet;
    }

    DEBUG("SetLoginInfo succ, req:%s", klib::jce2str(stReq).c_str());
    return iRet;
}

int RtaWriteSvr::SetLoginInfo(const string &strDeviceId, const unsigned uLoginTime)
{
    int iRet = 0;
    string strLoginKey = PROFILE_LOGIN_CKV_KEY + strDeviceId;
    int iCas = -1;
    string strLoginTime;

    iRet = m_tmem_obj.Get(strLoginKey, strLoginTime, iCas, 0, RTA_WRITE_SVR_TMEM_MOD_ID);
    if (iRet == -13200 || iRet == -13106)
    {
        DEBUG("no login info for key:%s", strLoginKey.c_str());
        strLoginTime = klib::ToStr(uLoginTime);
    }
    else if (iRet == 0)
    {
        if (IsSameDay((unsigned)atoll(strLoginTime.c_str()), uLoginTime))
        {
            DEBUG("login multiple times today for key:%s", strLoginKey.c_str());
        }
        else
        {
            DEBUG("login time:%s not today, update it, key:%s", strLoginTime.c_str(), strLoginKey.c_str());
            strLoginTime = klib::ToStr(uLoginTime);
        }
    }
    else
    {
        ERROR("m_tmem_obj Get err, iRet:%d, errMsg:%s", iRet, m_tmem_obj.GetLastError());
        return iRet;
    }
    DEBUG("m_tmem_obj Get succ, strLoginTime:%s", strLoginTime.c_str());

    // 写入存储
    iRet = m_tmem_obj.Set(strLoginKey, strLoginTime, iCas, 0,
            RTA_WRITE_SVR_TMEM_MOD_ID, atoll(g_conf["tmem_login_ckv.expire_time_second"].c_str()));
    if (iRet)
    {
        ERROR("m_tmem_obj Set err, key:%s, iRet:%d, errMsg:%s", strLoginKey.c_str(), iRet, m_tmem_obj.GetLastError());
    }
    else
    {
        DEBUG("m_tmem_obj Set succ, key:%s, strLoginTime:%s", strLoginKey.c_str(), strLoginTime.c_str());
    }

    return iRet;
}

bool RtaWriteSvr::IsSameDay(time_t time_1, time_t time_2)
{
    tm tmTime1, tmTime2;
    localtime_r(&time_1, &tmTime1);
    localtime_r(&time_2, &tmTime2);

    if (tmTime1.tm_year == tmTime2.tm_year && tmTime1.tm_mon == tmTime2.tm_mon &&
            tmTime1.tm_mday == tmTime2.tm_mday)
    {
        return true;
    }

    return false;
}

#define QZASYNCMSG_CLASS RtaWriteSvr

//这个放在最后面，请不要修改位置
#include "qza_svr_frame.h"

