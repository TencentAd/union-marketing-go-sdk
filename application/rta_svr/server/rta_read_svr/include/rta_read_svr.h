#ifndef __RTA_READ_SVR_H__
#define __RTA_READ_SVR_H__
//create by create_spp_server.sh, svr_tpl@v1.0
//必须包含spp的头文件
#include <tr1/functional>
#include <sstream>
#include "sppincl.h"
#include "api_base_new.h"
#include "qza_protocol.h"
#include "conf_file.h"
#include "dcapi_cpp.h"
#include "logapi.h"
#include "server_config.h"
#include "string_deal.h"
#include "comm_define.h"
#include "qza_sync_msg.h"
#include "music_errcode.h"

#include "klib.h"
#include "proto_svr_gdt_register.h"
#include "tmem_api.h"
#include "app_dcreport.h"

using namespace proto_gdt_register;

class RtaReadSvr: public ns_kg::CQzaSyncMsg
{
public:
    RtaReadSvr();

    virtual ~RtaReadSvr();

    // 初始化类静态成员
    static int Init(const char* etc);

    static int HandleRoute(unsigned flow, void* arg1, void* arg2)
    {
        return 1;
    }

    virtual int DealRequest( short iCmd, short iSubCmd, char* pBodyBuf, int iBodyLen, char* pResBuf, int& iLeftSize);

    // 需要先回报后处理的
    virtual int PostProcess( short iCmd, short iSubCmd, char* pBodyBuf, int iBodyLen);

    // 微线程初始化
    int InitService();

public:

    // 打解包函数
    template<typename REQUEST, typename RESPONSE, typename CB_FN>
    int ProcessRequest(char* body_buf, int body_len, char* res_buf, int& res_len, CB_FN fn);

    // 打解包函数，需要先回报后处理的
    template<typename REQUEST, typename RESPONSE, typename CB_FN>
    int ProcessRequestNoRsp(char* body_buf, int body_len, CB_FN fn);

    // 处理RTA请求
	int DoRTARequest(proto_gdt_register::RTAReq &req, proto_gdt_register::RTARsp &rsp);

    // 处理RTA上报
	int DoRTAReport(const unsigned uPlatformType, const unsigned uDeviceType, const unsigned uCompanyId, const string &strMd5Id, const string &strReqId);

    // post阶段上报
	int PostRtaRequest(proto_gdt_register::RTAReq &req);

    int DoRTAStrategy(proto_gdt_register::RTAReq &req, int getLoginTimeRet, const string& strLoginTime,
            int getRTACacheRet, RTAInfoCache& stRTAInfoCache, map<unsigned, bool> &mapGroupId2Res);

    void GetRTAStatus(map<unsigned, bool> &mapGroupId2Res, unsigned &uStatus);

    int GetStrategyGroupById(unsigned int uStrategyGroupId, proto_gdt_register::StrategyGroup& group);

    int FilterTodayActive(const string& strDeviceId, const int getLoginTimeRet, const string& strLoginTime, unsigned int& uStatus);

    int FilterClickFreq(unsigned int uCompanyId, const string& strTmpKey, int getRTACacheRet, proto_gdt_register::RTAInfoCache& stRTAInfoCache, unsigned int& uStatus);

	int FilterPrediction(const string& strTmpKey, int getRTACacheRet, proto_gdt_register::RTAInfoCache& stRTAInfoCache, unsigned int uStrategyId, unsigned int& uStatus);

    int DoRTASingleStrategy(const unsigned uStrategyGroupId, proto_gdt_register::RTAReq &req, const int getLoginTimeRet,
            const string& strLoginTime, const int getRTACacheRet, RTAInfoCache& stRTAInfoCache,
            unsigned int& uStatus, unsigned int &uStrategyId);

public:
    static bool g_bInit;

    CTmemApiMt m_tmem_obj;
    unsigned int m_uNowTimestamp;

    vector<proto_gdt_register::DataReportRta> m_vctDataReportRta;   // rta 上报信息
};

#endif

