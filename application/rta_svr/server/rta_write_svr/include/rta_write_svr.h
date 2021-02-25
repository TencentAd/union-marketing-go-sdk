#ifndef __RTA_WRITE_SVR_H__
#define __RTA_WRITE_SVR_H__
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

using namespace proto_gdt_register;

class RtaWriteSvr: public ns_kg::CQzaSyncMsg
{
public:
    RtaWriteSvr();

    virtual ~RtaWriteSvr();

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

    static bool IsSameDay(time_t time_1, time_t time_2);

    // 存储渠道点击
    int SetClickTime(const string &strDeviceId);

    // 导入数据侧离线策略号码包
    int DoSetStrategyInfo(proto_gdt_register::SetStrategyInfoReq &req, proto_gdt_register::SetStrategyInfoRsp &rsp);

    // 更新离线策略号码包
    int UpdateRtaStrategyInfo(const string &strDeviceId, const unsigned uStrategyId, const unsigned uTimestamp);

    // 设置渠道点击信息
    int DoSetClickInfo(SetClickInfoReq &stReq, SetClickInfoRsp &stRsp);
    int SetClickTime(const string &strDeviceId, const unsigned uClickTime);

    // 设置严格口径DAU登陆信息
    int DoSetLoginInfo(SetDauLoginReq &stReq, SetDauLoginRsp &stRsp);
    int SetLoginInfo(const string &strDeviceId, const unsigned uLoginTime);

public:
    static bool g_bInit;

    CTmemApiMt m_tmem_obj;
    unsigned int m_uNowTimestamp;
};

#endif

