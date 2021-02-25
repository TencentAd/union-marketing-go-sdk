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
#include "rta_read_svr.h"
#include "exec_task.h"

#define g_conf (*CSingleton<CServerConf>::instance())
using namespace proto_gdt_register;

bool RtaReadSvr::g_bInit = false;

void SplitToUint(const string &str, vector<unsigned int> &vct)
{
    vector<string> tmpVct;
    klib::StrSplit(tmpVct, str, ",");
    for (vector<string>::iterator it = tmpVct.begin(); it != tmpVct.end(); it++)
    {
        if (*it != "")
        {
            vct.push_back((unsigned)atoi(it->c_str()));
        }
    }
}

// 获取某一时区当前时刻对应的零点时间戳
static unsigned int GetDayStartTs(unsigned int uNow, int iTimeZone)
{
    int iSurplus = (uNow - (24 - iTimeZone) * 3600) % (24 * 3600);
    return uNow - iSurplus;
}

// 微线程初始化
int RtaReadSvr::InitService()
{
    return 0;
}

// 进程初始化
int RtaReadSvr::Init(const char* conf_file)
{
    if (g_bInit)
    {
        return 0;
    }

    if (NULL == conf_file)
    {
        ERROR("load_conf error! [params is null]\n");
        return -1;
    }

    if (g_conf.ParseFile(conf_file) != 0)
    {
        ERROR("g_conf.ParseFile fail, file:%s", conf_file);
        return -2;
    }

    API_Logapi_Init("rta_read_svr",  RTA_READ_SVR_MOD_ID);
    g_bInit = true;
    return 0;
}

RtaReadSvr::RtaReadSvr()
{
    unsigned int tmem_bid = atoll(g_conf["tmem_bid"].c_str());
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

    m_vctDataReportRta.clear();
    m_uNowTimestamp = time(NULL);
}

RtaReadSvr::~RtaReadSvr()
{
}

// 打解包函数
template<typename REQUEST, typename RESPONSE, typename CB_FN>
int RtaReadSvr::ProcessRequest(char* body_buf, int body_len, char* res_buf, int& res_len, CB_FN fn)
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
int RtaReadSvr::ProcessRequestNoRsp(char* body_buf, int body_len, CB_FN fn)
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

int RtaReadSvr::DealRequest(short iCmd, short iSubCmd, char* pBodyBuf, int iBodyLen, char* pResBuf, int& iLeftSize)
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
     *不需要回包的可以用NS_API_BASE::RSP_BASE
     */
    switch (iSubCmd)
    {
        case CMD_GDT_REGISTER_SVR_RTA_REQUEST:
            iRet = ProcessRequest<RTAReq, RTARsp, tr1::function<int(RTAReq&, RTARsp&)> >
            (pBodyBuf, iBodyLen, pResBuf, iLeftSize, tr1::bind(&RtaReadSvr::DoRTARequest, this, tr1::placeholders::_1, tr1::placeholders::_2));
            break;
        default:
            iLeftSize = 0;
            iRet = MUSIC_CODE_COMM_ParamInvalid;
            ERROR("unknown cmd:%d", iSubCmd);
            break;
    }
    return iRet;
}

int RtaReadSvr::PostProcess(short iCmd, short iSubCmd, char* pBodyBuf, int iBodyLen)
{
    int iRet = 0;
    switch (iSubCmd)
    {
      // 需要分两阶段处理的加
      case CMD_GDT_REGISTER_SVR_RTA_REQUEST:
          iRet = ProcessRequestNoRsp<RTAReq, tr1::function<int(RTAReq&)> >
                (pBodyBuf, iBodyLen, tr1::bind(&RtaReadSvr::PostRtaRequest, this, tr1::placeholders::_1));
          break;
      default:
          //just return
          break;
    }
    return 0;
}

// 处理RTA请求
int RtaReadSvr::DoRTARequest(proto_gdt_register::RTAReq &req, proto_gdt_register::RTARsp &rsp)
{
    int iRet = 0;

    // [step 1] 判断是否需要竞价
    // 先转成小写字母
    string strDeviceId = req.strDeviceId;
    std::transform(strDeviceId.begin(), strDeviceId.end(), strDeviceId.begin(), ::tolower);

    // RTA请求对实时性要求极高，故采用task模式并发拉取信息
    //[step 1] 创建task list
    IMtTaskList stTaskList;

    // [step 2] 开始插入任务
    CGetLoginTimeTask * getLoginTimeTask = new CGetLoginTimeTask(strDeviceId);
    stTaskList.push_back(getLoginTimeTask);

    CGetRTAInfoCacheTask * getRTACacheTask = new CGetRTAInfoCacheTask(strDeviceId);
    stTaskList.push_back(getRTACacheTask);

    // [step 3] task模式执行
    iRet = mt_exec_all_task(stTaskList);
    if (iRet)
    {
        ERROR("mt_exec_all_task err, key:%s, ret:%d", strDeviceId.c_str(), iRet);
    }
    else
    {
        DEBUG("mt_exec_all_task succ ! key:%s", strDeviceId.c_str());
    }

    // [step 4] 取出数据
    int getLoginTimeRet = getLoginTimeTask->GetResult();
    string strLoginTime = getLoginTimeTask->m_strLoginTime;

    int getRTACacheRet = getRTACacheTask->GetResult();
    RTAInfoCache stRTAInfoCache = getRTACacheTask->m_stRTAInfoCache;

    // [step 5] 释放task list
    for (unsigned int j = 0; j < stTaskList.size(); j++)
    {
        if (NULL != stTaskList[j])
        {
            delete stTaskList[j];
        }
    }

    // [step 6] 使用策略
    if (req.vctStrategyGroupId.size() != 0)
    {
        iRet = DoRTAStrategy(req, getLoginTimeRet, strLoginTime, getRTACacheRet, stRTAInfoCache, rsp.mapGroupId2Res);
        if (iRet != 0)
        {
            ERROR("DoRTAStrategy error, iRet: %d, req: %s.", iRet, klib::jce2str(req).c_str());
            return iRet;
        }
        GetRTAStatus(rsp.mapGroupId2Res, rsp.uStatus);
        DEBUG("DoRTAStrategy succ, req:%s, rsp:%s", klib::jce2str(req).c_str(), klib::jce2str(rsp).c_str());
        return 0;
    }
    else
    {
        ERROR("invalid req:%s, empty vctStrategyGroupId", klib::jce2str(req).c_str());
        return MUSIC_CODE_COMM_ParamInvalid;
    }
    return 0;
}

int RtaReadSvr::PostRtaRequest(proto_gdt_register::RTAReq &req)
{
	int iRet = 0;

	// 数据上报
	iRet = DoRTAReport(req.uPlatformType, req.uDeviceType, req.uCompanyId, req.strDeviceId, req.strReqId);
	if (iRet)
	{
		ERROR("DoRTAReport err, iRet:%d, req:%s", iRet, klib::jce2str(req).c_str());
		return iRet;
	}
	else
	{
		DEBUG("DoRTAReport succ, req:%s", klib::jce2str(req).c_str());
	}

	return 0;
}

int RtaReadSvr::DoRTAReport(const unsigned uPlatformType, const unsigned uDeviceType, const unsigned uCompanyId,
        const string &strMd5Id, const string &strReqId)
{
	int iRet = 0;
    //TODO  业务数据上报
	return iRet;
}

int RtaReadSvr::DoRTAStrategy(RTAReq &req, int getLoginTimeRet, const string& strLoginTime,
        int getRTACacheRet, RTAInfoCache& stRTAInfoCache, map<unsigned, bool> &mapGroupId2Res)
{
    int iRet = 0;
    for (vector<unsigned int>::iterator it = req.vctStrategyGroupId.begin(); it != req.vctStrategyGroupId.end(); it++)
    {
        unsigned uStatus = 0, uStrategyId = 0;

        // 数据上报
        DataReportRta stDataReportRta;
        stDataReportRta.uStrategyGroupId = *it;

        // 逐策略组判断是否需要竞价
        // 只要该设备命中策略组中任何一种过滤策略，则返回不竞价
        iRet = DoRTASingleStrategy(*it, req, getLoginTimeRet, strLoginTime, getRTACacheRet, stRTAInfoCache, uStatus, uStrategyId);
        if (iRet != 0)
        {
            // 出错则默认竞价
            stDataReportRta.uStrategyId = uStrategyId;
            stDataReportRta.uStatus = 1;
            m_vctDataReportRta.push_back(stDataReportRta);
            ERROR("DoRTASingleStrategy error, iRet:%d, strategy group id:%u, strategy id:%u, req:%s.", iRet, *it, uStrategyId, klib::jce2str(req).c_str());
            return iRet;
        }
        DEBUG("DoRTASingleStrategy succ, strategy group id:%u, strategy id:%u, req:%s, uStatus:%u",
                *it, uStrategyId, klib::jce2str(req).c_str(), uStatus);

        // 保存该策略组的结果，数据上报用
        mapGroupId2Res[*it] = (uStatus == 0 ?  false : true);
        stDataReportRta.uStrategyId = uStrategyId;
        stDataReportRta.uStatus = uStatus;
        m_vctDataReportRta.push_back(stDataReportRta);
    }
    return iRet;
}

int RtaReadSvr::DoRTASingleStrategy(const unsigned uStrategyGroupId, proto_gdt_register::RTAReq &req,
        const int getLoginTimeRet, const string& strLoginTime, const int getRTACacheRet, RTAInfoCache& stRTAInfoCache,
        unsigned int &uStatus, unsigned int &uStrategyId)
{
    int iRet = 0;
    const string &strDeviceId = req.strDeviceId;
    const unsigned uPlatformType = req.uPlatformType;
    const unsigned uCompanyId = req.uCompanyId;

    // [step 1] 根据策略组id从配置文件中获取该策略组相关的配置，包括策略组中的所有策略id以及该策略组对应的渠道
    proto_gdt_register::StrategyGroup groupStragegy;
    iRet = GetStrategyGroupById(uStrategyGroupId, groupStragegy);
    if (iRet)
    {
        ERROR("GetStrategyGroupById err, strategy group id:%u, req:%s, iRet:%d", uStrategyGroupId, klib::jce2str(req).c_str(), iRet);
        return iRet;
    }
    DEBUG("strategy group id:%u, req:%s, StrategyGroup:%s.", uStrategyGroupId, klib::jce2str(req).c_str(), klib::jce2str(groupStragegy).c_str());

    // [step 2] 渠道过滤
    bool bMatchCompany = false;
    for (vector<unsigned int>::iterator it = groupStragegy.vctCompanyId.begin(); it != groupStragegy.vctCompanyId.end(); it++)
    {
        if (*it == uCompanyId)
        {
            bMatchCompany = true;
            break;
        }
    }
    if (!bMatchCompany)
    {
        ERROR("uCompanyId not match, req:%s, uStrategyGroupId:%u, StrategyGroup:%s",
                klib::jce2str(req).c_str(), uStrategyGroupId, klib::jce2str(groupStragegy).c_str());
        return MUSIC_CODE_GDT_REGISTER_INVALID_COMPANY_ID;
    }

    // [step 3] 策略过滤
    for (vector<unsigned int>::iterator it = groupStragegy.vctStrategy.begin(); it != groupStragegy.vctStrategy.end(); it++)
    {
        uStrategyId = *it;
        switch (uStrategyId)
        {
            case 101:
            {
                iRet = FilterTodayActive(req.strDeviceId, getLoginTimeRet, strLoginTime, uStatus);
                if (iRet != 0)
                {
                    ERROR("FilterTodayActive error, req:%s, strategy group id:%u, strategy id:%u, iRet:%d",
                            klib::jce2str(req).c_str(), uStrategyGroupId, uStrategyId, iRet);
                    return iRet;
                }
                DEBUG("FilterTodayActive succ, req:%s, strategy group id:%u, strategy id:%u, uStatus:%u",
                        klib::jce2str(req).c_str(), uStrategyGroupId, uStrategyId, uStatus);
                break;
            }
            case 103:
            {
                iRet = FilterClickFreq(uCompanyId, strDeviceId, getRTACacheRet, stRTAInfoCache, uStatus);
                if (iRet != 0)
                {
                    ERROR("FilterClickFreq error, req:%s, strategy group id:%u, strategy id:%u, iRet:%d",
                            klib::jce2str(req).c_str(), uStrategyGroupId, uStrategyId, iRet);
                    return iRet;
                }
                DEBUG("FilterClickFreq succ, req:%s, strategy group id:%u, strategy id:%u, uStatus:%u",
                        klib::jce2str(req).c_str(), uStrategyGroupId, uStrategyId, uStatus);
                break;
            }
            case 104: case 106: case 107: case 108:
            {
                iRet = FilterPrediction(strDeviceId, getRTACacheRet, stRTAInfoCache, uStrategyId, uStatus);
                if (iRet != 0)
                {
                    ERROR("FilterPrediction error, req:%s, strategy group id:%u, strategy id:%u, iRet:%d",
                            klib::jce2str(req).c_str(), uStrategyGroupId, uStrategyId, iRet);
                    return iRet;
                }
                DEBUG("FilterPrediction succ, req:%s, strategy group id:%u, strategy id:%u, uStatus:%u",
                        klib::jce2str(req).c_str(), uStrategyGroupId, uStrategyId, uStatus);
                break;
            }
            default:
            {
                iRet = MUSIC_CODE_GDT_REGISTER_INVALID_STRATRGY_ID;
                ERROR("req:%s, strategy group id:%u, unknown uStrategyId:%u, iRet:%d.",
                        klib::jce2str(req).c_str(), uStrategyGroupId, uStrategyId, iRet);
                return iRet;
            }
        }

        // 只要该设备命中策略组中任何一种过滤策略，则不竞价直接返回
        // 否则继续判断策略组中的下一种过滤策略是否满足
        if (uStatus == 0)
        {
            DEBUG("req[%s] reaches the strategyGroup[%u] strategy[%u], don`t need bid.",
                    klib::jce2str(req).c_str(), uStrategyGroupId, uStrategyId);
            break;
        }
    }

    return iRet;
}

// 目前，先从配置文件读
// 后续开发”红石配置平台“后，改为从ckv读
int RtaReadSvr::GetStrategyGroupById(const unsigned uStrategyGroupId, proto_gdt_register::StrategyGroup& group)
{
    int iRet = 0;

    // 如果该策略组已经无效，直接返回非法策略组
    if (g_conf["strategy_group_" + klib::ToStr(uStrategyGroupId) + ".valid"] != "1")
    {
        DEBUG("uStrategyGroupId:%u is invalid", uStrategyGroupId);
        return MUSIC_CODE_GDT_REGISTER_INVALID_STRATRGY_GROUP;
    }

    // 获取该策略组的相关信息
    group.uStrategyGroupId = uStrategyGroupId;
    SplitToUint(g_conf["strategy_group_" + klib::ToStr(uStrategyGroupId) + ".sence"], group.vctSence);
    SplitToUint(g_conf["strategy_group_" + klib::ToStr(uStrategyGroupId) + ".company"], group.vctCompanyId);
    SplitToUint(g_conf["strategy_group_" + klib::ToStr(uStrategyGroupId) + ".strategy"], group.vctStrategy);

    return iRet;
}

int RtaReadSvr::FilterTodayActive(const string &strDeviceId, const int getLoginTimeRet, const string &strLoginTime, unsigned int &uStatus)
{
    int iRet = 0;
    if (-13106 == getLoginTimeRet || -13200 == getLoginTimeRet)
    {
        DEBUG("iRet:%d, strDeviceId:%s does not exist in profile login ckv, need to bid", getLoginTimeRet, strDeviceId.c_str());
		uStatus = 1;
    }
    else if (0 != getLoginTimeRet)
    {
        ERROR("m_tmem_obj Get err, strDeviceId:%s, iRet:%d", strDeviceId.c_str(), getLoginTimeRet);
		return getLoginTimeRet;
    }
    else
    {
        DEBUG("m_tmem_obj Get succ, strDeviceId:%s, strLoginTime:%s", strDeviceId.c_str(), strLoginTime.c_str());

        long long loginTimeStart = GetDayStartTs((unsigned)atoll(strLoginTime.c_str()), 8);
        long long nowTimeStart = GetDayStartTs(m_uNowTimestamp, 8);
        long long deltaTime = nowTimeStart - loginTimeStart;
        if (deltaTime == 0)  // 自然天的开始时间相同
        {
            DEBUG("strDeviceId:%s's login time is today, does not to bid", strDeviceId.c_str());
            uStatus = 0;
        }
        else
        {
            DEBUG("strDeviceId:%s's login time is not today, need to bid, set status = 1", strDeviceId.c_str());
            uStatus = 1;
        }
    }

    return iRet;
}

int RtaReadSvr::FilterClickFreq(unsigned int uCompanyId, const string& strDeviceId, int getRTACacheRet, RTAInfoCache& stRTAInfoCache, unsigned int& uStatus)
{
    if (0 == getRTACacheRet)
    {
        unsigned int lastClickDay = GetDayStartTs(stRTAInfoCache.uLastClickTime, 8);
        unsigned int nowTimeDay = GetDayStartTs(m_uNowTimestamp, 8);

        // 当日已经点击过广告，不需要再竞价
        if (lastClickDay == nowTimeDay)
        {
            DEBUG("strDeviceId:%s has clicked today, needn`t to bid, set status = 0", strDeviceId.c_str());
            uStatus = 0;
        }
        else
        {
            DEBUG("strDeviceId:%s lastClickDay:%u, nowTimeDay:%u, need to bid.", strDeviceId.c_str(), nowTimeDay, lastClickDay);
            uStatus = 1;
        }
    }
    else
    {
        DEBUG("strDeviceId:%s not click before, need to bid.", strDeviceId.c_str());
        uStatus = 1;
    }
    return 0;
}

int RtaReadSvr::FilterPrediction(const string& strDeviceId, int getRTACacheRet, RTAInfoCache& stRTAInfoCache, unsigned int uStrategyId, unsigned int& uStatus)
{
    int iRet = 0;
    if (getRTACacheRet == -13106 || getRTACacheRet == -13200)
    {
        DEBUG("getRTACacheRet:%d, strDeviceId:%s does not exist in prediction strategy ckv, need to bid", getRTACacheRet, strDeviceId.c_str());
		uStatus = 1;
    }
    else if (0 != getRTACacheRet)
    {
        ERROR("m_tmem_obj Get err, strDeviceId:%s, getRTACacheRet:%d", strDeviceId.c_str(), getRTACacheRet);
		return getRTACacheRet;
    }
    else
    {
        map<unsigned int, unsigned int>::iterator it = stRTAInfoCache.mapUpdateTime.find(uStrategyId);
        if (it == stRTAInfoCache.mapUpdateTime.end())
        {
            DEBUG("strDeviceId:%s does not exist in stRTAInfoCache map, uStrategyId:%u, need to bid.", strDeviceId.c_str(), uStrategyId);
            uStatus = 1;
        }
        else
        {
            unsigned int uLastUpdateTime = it->second;

            string strConfKey = "strategy_info_" + klib::ToStr(uStrategyId) + ".";
            long long validTime = atoll(g_conf[strConfKey + "valid_time_sec"].c_str());
            if (m_uNowTimestamp > uLastUpdateTime && m_uNowTimestamp - uLastUpdateTime > validTime)
            {
                DEBUG("strDeviceId:%s, strategy id:%u, uLastUpdateTime:%u, uTimeNow:%u, greater then %lld, need bid.",
                        strDeviceId.c_str(), uStrategyId, uLastUpdateTime, m_uNowTimestamp, validTime);
                uStatus = 1;
            }
            else
            {
                DEBUG("strDeviceId:%s, strategy id:%u, uLastUpdateTime:%u, uTimeNow:%u, less then %lld, needn`t bid.",
                        strDeviceId.c_str(), uStrategyId, uLastUpdateTime, m_uNowTimestamp, validTime);
                uStatus = 0;
            }
        }
    }
    return 0;
}

void RtaReadSvr::GetRTAStatus(map<unsigned, bool> &mapGroupId2Res, unsigned &uStatus)
{
    map<unsigned, bool>::iterator it = mapGroupId2Res.begin();

    // 只要有一个策略组竞价，则状态为竞价
    for (; it != mapGroupId2Res.end(); ++it)
    {
        if (it->second == true)
        {
            uStatus = 1;
            return ;
        }
    }
    uStatus = 0;
}

#define QZASYNCMSG_CLASS RtaReadSvr

//这个放在最后面，请不要修改位置
#include "qza_svr_frame.h"

