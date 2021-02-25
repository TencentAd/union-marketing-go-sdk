#ifndef __GDT_REGISTER_H__
#define __GDT_REGISTER_H__
//create by create_spp_server.sh, svr_tpl@v1.0
//必须包含spp的头文件
#include <openssl/sha.h>
#include <sstream>
#include <openssl/hmac.h>
#include "sppincl.h"
#include "api_base_new.h"
#include "qza_protocol.h"
#include "conf_file.h"
#include "dcapi_cpp.h"
#include "value.h"
#include "writer.h"
#include "reader.h"
#include "curl/curl.h"
#include "logapi.h"
#include "server_config.h"
#include "string_deal.h"
#include "comm_define.h"
#include "qza_sync_msg.h"
#include "music_errcode.h"

#include "curl_wrap.h"
#include "klib.h"
#include "tc_base64.h"
#include "proto_svr_gdt_register.h"
#include "func_lib.h"
#include "WnsQuaHelper.h"
#include "tmem_api.h"
#include "webapp_common.h"
#include "kg_md5.h"
#include "kafka_producer.h"

using namespace QZ_BILL;
using namespace proto_gdt_register;
using namespace NS_PROFILE;

class GdtRegister: public ns_kg::CQzaSyncMsg
{
public:
    GdtRegister();

    virtual ~GdtRegister();

    static int Init(const char* etc);
    static int InitKafkaProducer();

    static int HandleRoute(unsigned flow, void* arg1, void* arg2);

    virtual int DealRequest( short iCmd, short iSubCmd, char* pBodyBuf, int iBodyLen, char* pResBuf, int& iLeftSize);
	virtual int PostProcess( short iCmd, short iSubCmd, char* pBodyBuf, int iBodyLen);

public:

	//处理cgi点击
    int ProcReportClick(ReportClickReq& req, NS_API_BASE::RSP_BASE& rsp);

    //根据点击请求更新register info
    void UpdateRegisterInfo(ReportClickReq &req, string &strDeviceId, bool bUpdateMuid=false);

    //根据点击请求填充需要存储的信息
    void FillInRegisterInfo(ReportClickReq &req, RegisterInfo &register_info);

    //根据点击请求生成muid
    void FillInMuid(ReportClickReq &req, string &muid);

    //读取存储中的register info
    void GetRegisterInfo(string &strDeviceId, proto_gdt_register::RegisterInfo &register_info, int &iCas);

	//严格口径dau表上报
	int ProcLoginReport(proto_gdt_register::LoginReportReq &stReq, NS_API_BASE::RSP_BASE &stRsp);

    //渠道拉新处理
    int HandleNewUser(proto_gdt_register::LoginReportReq &stReq, map<string, string> &mapDeviceInfo, const WNS::WnsQua &stQua);

    //渠道拉活处理
    int HandleActiveUser(proto_gdt_register::LoginReportReq &stReq, map<string, string> &mapDeviceInfo, const WNS::WnsQua &stQua);

    //判断是否回流用户
    int JudgeReturner(const unsigned uUid, const string &strDeviceId, const string &strAccountId, bool &bReturner);

    //判断是否新用户
    int JudgeNewUser(const unsigned uUid, bool &bNewUser);

    //判断用户是否上报过
    int ReportBefore(const unsigned uUid, bool &bReported, int &iCas);

    //设置上报状态位
    int SetReportRecord(const unsigned uUid, int iCas);

    //对于安卓设备，imei>oaid>androidid，逐层匹配
    //对于ios设备，idfa>ip+ua，逐层匹配
    int FindClickInfo(map<string, string> &mapDeviceInfo, const WNS::WnsQua &stQua, proto_gdt_register::LoginReportReq &stReq,
            vector<string> &vctDeviceId, vector<string> &vctData);

    //针对特定渠道，过滤掉一些无效的deviceInfo，如果 0000-00000 对应的registerInfo
    int FindValidRegisterInfo(map<string, string> &mapDeviceInfo, const WNS::WnsQua &stQua, LoginReportReq &stReq,
        vector<string>& vctDeviceId, vector<string>& vctData, string& strDeviceId, RegisterInfo& stRegisterInfo);

    //从deviceinfo中解析出所有设备信息，均为md5之后的小写
    int GetAllDeviceId(map<string, string> &mapDeviceInfo, const WNS::WnsQua &stQua, proto_gdt_register::LoginReportReq &stReq,
            string &strMd5Imei, string &strMd5Oaid, string &strMd5AndroidId, string &strMd5Idfa, string &strMd5IpUa);

    //将安卓设备的imei oaid以及androidid信息都放入task模式查询
    void AddTaskListAndroid(const string &strMd5Imei, const string &strMd5Oaid, const string &strMd5AndroidId, const string &strIpUa, IMtTaskList &stTaskList);

    //将ios设备的idfa和ipua信息都放入task模式查询
    void AddTaskListIos(const string &strMd5Idfa, const string &strMd5IpUa, IMtTaskList &stTaskList);

    //将设备id以task模式去查是否有对应的渠道点击信息
    void DoAddTaskList(const string &strDeviceId, IMtTaskList &stTaskList);

    //打解包函数
    template<typename REQUEST, typename RESPONSE, typename CB_FN>
    int ProcessRequest(char* body_buf, int body_len, char* res_buf, int& res_len, CB_FN fn);
    template<typename REQUEST, typename RESPONSE, typename CB_FN>
    int ProcessRequestNoResponse(char* body_buf, int body_len, CB_FN fn);

    //deviceinfo中的ipua --> muid
    string GenerateKeyIpUa(const unsigned uUid, const string& strClientIP, map<string, string>& mapDeviceInfo, const WNS::WnsQua& qua);
    //监测链接中的ipua --> muid
    string GenerateMuidByIpUa(const ReportClickReq& req);

    //从UA里提取系统版本号
    string GetSysVerFromUserAgent(int nSystem, size_t nSystemIdx, const string& strUserAgent);

    //回调给广点通
    int SendCallBackToGDT(const RegisterInfo& register_info, const string &strQua, const string &strDeviceInfo, const WNS::WnsQua& qua);
    int ActiveSendCallBackToGDT(const RegisterInfo& register_info, const WNS::WnsQua& qua);

    int report_to_gdt(const RegisterInfo& register_info, string& post, string& nonce, const NS_PROFILE::UserInfoCreateReq& req);
    int refresh_token(const RegisterInfo& register_info);
    int GetGdtToken(unsigned int advertiseId, string& strToken);
    // uHeaderTyoe = 1 : Content-Type: application/json
    // uHeaderType = 2 : application/x-www-formurlencoded
    // uReqType = 1 : GET
    // uReqType = 2 : POST
    int CommonCurl(const string& url, const string& strPostJson, string& strResult, unsigned uHeaderType = 1, unsigned uReqType = 1);

    //-----------------适配广点通的AMS投放----------------------//
    void ParseParamsForAMS(RegisterInfo& register_info, const string& strCallBackParam);
    int SendCallBackToAMS(const RegisterInfo& register_info, const string &strQua, const string &strDeviceInfo, const WNS::WnsQua& qua);

    // AMS拉新回调
    int SendCallBackToAMSSelf(RegisterInfo &stRegisterInfo, const WNS::WnsQua &stQua, const unsigned uLoginTime, const string &strType);

    int GetCurlTime(CURL* curl, double& dTotalTime, double& dDNSTime, double& dConnectTime, double& dTranferTime);

    int EncodeAndSet(string& key, RegisterInfo& register_info, int iCas = -1);
    int EncodeAndSet(string& key, GdtToken& gdtToken, unsigned int expireSec);
    int EncodeAndSetSecondary(RegisterInfo& register_info);

    // 判断该设备是否当天首登
    int JugdeFirstLoginToday(const string &strDeviceId, bool &bTodayFirstLogIn);

    // 将设备登陆时间写入kafka
    int SetLoginTime(const string &strDeviceId, const unsigned uLoginTime);

    //---------------次留模型处理-------------------
    int HandleSecondaryModel(proto_gdt_register::LoginReportReq &stReq, map<string, string> &mapDeviceInfo, const WNS::WnsQua &stQua);

    string md5Lower(const string &str);

public:
    static bool g_bInit;
	static string g_access_token;
    static string g_access_token_bproxy;
    static unsigned int g_token_time;
    static unsigned int g_token_time_bproxy;
    static ns_common_kafka::CKafkaProducer g_kafkaProducer;             // kafka producer
    static unsigned int g_uKafkaSetAcksAll;                             // kafka写入一致性级别开关

public:
    char* m_buf;
    unsigned int m_bufsize;

    CTmemApiMt m_tmem_obj;
    unsigned int m_tmem_expire;
    unsigned int m_curl_mid;
    unsigned int m_curl_cid;
    unsigned int m_curl_timeout_ms;
    unsigned int m_curl_bdtimeout;
    unsigned int m_ams_mid;
    unsigned int m_ams_cid;

    unsigned int m_token_timeout;
    unsigned int m_uNowTimestamp;

    unsigned int m_uFilterIllegalIdfa;
};

#endif

