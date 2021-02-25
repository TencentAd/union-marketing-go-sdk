//create by create_spp_server.sh
/**
 * 名称: 全民K歌
 * (c) Copyright 2017 harleyhuang@tencent.com. All Rights Reserved.
 *
 * change log:
 * Author: harleyhuang
 *   Date: 2017-06-26
 *	 Desc:
 */
#include "server_config.h"
//#include "app_log.h"
#include "klib.h"
#include <openssl/sha.h>
#include <uuid.h>
#include "gdt_register.h"
#include "exec_task.h"

#define g_conf (*CSingleton<CServerConf>::instance())
//#define g_server_base g_base
using namespace ns_common_kafka;

bool GdtRegister::g_bInit = false;
unsigned int GdtRegister::g_token_time = 0;
string GdtRegister::g_access_token;
unsigned int GdtRegister::g_token_time_bproxy = 0;
string GdtRegister::g_access_token_bproxy;
CKafkaProducer GdtRegister::g_kafkaProducer;
unsigned int GdtRegister::g_uKafkaSetAcksAll = 0;

void c2x(unsigned char c,char* buf)
{
    int i1,i2;

    i1 = c / 16;
    i2 = c % 16;

    if(i1 >= 10)
        buf[0] = 'A' + i1 - 10;
    else
        buf[0] = '0' + i1;

    if(i2 >= 10)
        buf[1] = 'A' + i2 - 10;
    else
        buf[1] = '0' + i2;

    buf[2] = 0;
}

string Bin2Hex(const string& bin)
{
    string hex;
    for(size_t i = 0; i < bin.size(); i ++)
    {
        char buf[3];
        c2x(bin.data()[i], buf);
        hex += buf;
    }
    return hex;
}

int FuncStoreResult(char *data, size_t size, size_t nmemb, std::string *result)
{
    if (result == NULL)
        return 0;

    /*if(size*nmemb > 204800)
    {
        snprintf(g_errmsg, sizeof(g_errmsg),  "(a)huge size, size[%d]", size*nmemb);
    }*/

    for(unsigned int i=0; i<size*nmemb; i++)
        result->push_back(*(data++));

    return size * nmemb;
}

//32位 uuid
string genUuid()
{
    uuid_t uu;
    uuid_generate(uu);
    string tmp((char*)uu, sizeof(uu));
    return klib::ToLower(Bin2Hex(tmp));
}

string GetGdtConfigStr(unsigned int advertiseId, string param)
{
    return g_conf["gdt_" + klib::ToStr(advertiseId) + "." + param];
}

//进程初始化
int GdtRegister::Init(const char* conf_file)
{
    if (g_bInit)
    {
        return 0;
    }

    if (NULL == conf_file)
    {
        g_base->log_.LOG_P_ALL( tbase::tlog::LOG_ERROR, " load_conf error! [params is null]\n");
        return -1;
    }

    if (g_conf.ParseFile(conf_file) != 0)
    {
        g_base->log_.LOG_P_ALL( tbase::tlog::LOG_ERROR, "g_conf.ParseFile fail, file:%s", conf_file);
        return -2;
    }

    int iRet = InitKafkaProducer();
    if (0 != iRet)
    {
        ERROR("InitKafkaProducer fail, iRet:%d", iRet);
        return -4;
    }
    g_uKafkaSetAcksAll = atoi(g_conf["kafka.set_acks_all"].c_str());

    API_Logapi_Init("gdt_register",  GDT_REGISTER_MOD_ID);
    g_bInit = true;
    return 0;
}

int GdtRegister::InitKafkaProducer()
{
    int iRet = 0;

    string kafka_topic = g_conf["kafka.topic"];
    string broker_addr = g_conf["kafka.kafka_broker_addr"];
    ERROR("kafka_topic:%s, broker_addr:%s", kafka_topic.c_str(), broker_addr.c_str());

    KafkaProducerConfig kafka_producer_config;
    string errInfo;
    iRet = kafka_producer_config.set_broker_addr(broker_addr, errInfo);
    if (iRet != 0)
    {
        ERROR("set_broker_addr failed, ret:%d, errInfo:%s", iRet, errInfo.c_str());
        return iRet;
    }
    else
    {
        ERROR("set_broker_addr succ %s", broker_addr.c_str());
    }

    if (1 == g_uKafkaSetAcksAll)
    {
        iRet = kafka_producer_config.set_acks("all");
        if (iRet != 0)
        {
            ERROR("kafka_producer_config.set_acks fail, iRet:%d", iRet);
            return iRet;
        }
    }

    string debug = g_conf["kafka.debug"];
    if (!debug.empty())
    {
        int tmp_ret = kafka_producer_config.set_debug(debug);
    }

    string statistics_interval_ms = g_conf["kafka.statistics_interval_ms"];
    if (!statistics_interval_ms.empty())
    {
        int tmp_ret = kafka_producer_config.set_statistics_interval_ms(statistics_interval_ms);
    }

    iRet = g_kafkaProducer.init(kafka_producer_config, NULL, NULL, NULL, NULL);
    if (iRet != 0)
    {
        ERROR("init g_kafkaProducer failed %d", iRet);
        return iRet;
    }
    else
    {
        ERROR("init g_kafkaProducer succ");
    }

    iRet = g_kafkaProducer.add_topic(kafka_topic);
    if (iRet != 0)
    {
        ERROR("add topic [%s] failed %d", kafka_topic.c_str(), iRet);
        return iRet;
    }
    else
    {
        ERROR("add topic [%s] succ", kafka_topic.c_str());
    }

    return 0;
}

GdtRegister::GdtRegister()
{
    m_DcObj.module = "gdt_register";

    m_bufsize = 1024*100;
    m_buf = new char[m_bufsize];

    unsigned int tmem_mid = atoll(g_conf["tmem.mid"].c_str());
    unsigned int tmem_cid = atoll(g_conf["tmem.cid"].c_str());
    m_tmem_expire = atoll(g_conf["tmem.expire"].c_str());
    int iRet = m_tmem_obj.fnInitL5(atoll(g_conf["tmem.bid"].c_str()), tmem_mid, tmem_cid, atoll(g_conf["tmem.timeout_ms"].c_str()));
    if(iRet != 0)
    {
        ERROR("m_tmem_obj.fnInitL5 error, iRet:%d", iRet);
    }

    m_curl_mid  = atoll(g_conf["curl.mid"].c_str());
    m_curl_cid  = atoll(g_conf["curl.cid"].c_str());
    m_ams_mid  = atoll(g_conf["curl.ams_mid"].c_str());
    m_ams_cid  = atoll(g_conf["curl.ams_cid"].c_str());
    m_curl_timeout_ms = atoll(g_conf["curl.timeout_ms"].c_str());
    m_curl_bdtimeout = atoll(g_conf["curl.bd_timeout_ms"].c_str()) / 1000;
    m_token_timeout = atoll(g_conf["curl.token_timeout"].c_str());
    m_uNowTimestamp = time(NULL);
}

GdtRegister::~GdtRegister()
{
    if(m_buf != NULL)
    {
        delete[] m_buf;
        m_buf = NULL;
    }
}

//打解包函数
template<typename REQUEST, typename RESPONSE, typename CB_FN>
int GdtRegister::ProcessRequest(char* body_buf, int body_len, char* res_buf, int& res_len, CB_FN fn)
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

//打解包函数
template<typename REQUEST, typename RESPONSE, typename CB_FN>
int GdtRegister::ProcessRequestNoResponse(char* body_buf, int body_len, CB_FN fn)
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
            break;
        }
        ret = fn(req, rsp);
        if (ret != 0)
        {
            ERROR("process failed, ret:%d", ret);
            proc_ret = ret;
        }
    }
    while (0);
    return proc_ret;
}

int GdtRegister::HandleRoute(unsigned flow, void* arg1, void* arg2)
{
    int nRoute = 1;
    if (arg1 == NULL)
    {
        return nRoute;
    }
    blob_type* blob = (blob_type*)arg1;
    QZAHEAD* pQzaHead = (QZAHEAD*)blob->data;
    if (pQzaHead == NULL)
    {
        return nRoute;
    }
    short iCmd, iSubCmd;
    pQzaHead->GetCmd(iCmd, iSubCmd);
    if (iCmd == MAIN_CMD_GDT_REGISTER && iSubCmd == CMD_GDT_REGISTER_SVR_LOGIN_REPORT)
    {
        nRoute = 2;
    }
    else
    {
       nRoute = 1;
    }
    DEBUG("icmd=%d iSubCmd=%d nRoute=%d", iCmd, iSubCmd, nRoute);
    return nRoute;
}

int GdtRegister::DealRequest(short iCmd, short iSubCmd, char* pBodyBuf, int iBodyLen, char* pResBuf, int& iLeftSize)
{
    int iRet = 0;   //最终的返回值 将保存在qza头里

    /**
    *注意！！！：返回前一定保证设置iLeftSize为正确回包大小，否则导致回包过大（ProcessRequest里设置了）
    */
    if (iCmd == MAIN_CMD_GDT_REGISTER)
    {
        m_DcObj.setCommandName(proto_gdt_register::etos((GDT_REGISTER_CMD)iSubCmd));
        switch (iSubCmd)
        {
        case CMD_GDT_REGISTER_SVR_REPORT_CLICK:
            iRet = ProcessRequest<ReportClickReq, NS_API_BASE::RSP_BASE,
            tr1::function<int(ReportClickReq&, NS_API_BASE::RSP_BASE&)> >
            (pBodyBuf, iBodyLen, pResBuf, iLeftSize, tr1::bind(&GdtRegister::ProcReportClick, this,  tr1::placeholders::_1, tr1::placeholders::_2));
            break;
        case CMD_GDT_REGISTER_SVR_LOGIN_REPORT:
            iLeftSize = 0;
            iRet = 0;
            break;
		default:
            iLeftSize = 0;
            iRet = MUSIC_CODE_COMM_ParamInvalid;
            ERROR("unknown sub cmd:%d", iSubCmd);
            break;
        }
    }

    return iRet;
}

int GdtRegister::PostProcess( short iCmd, short iSubCmd, char* pBodyBuf, int iBodyLen)
{
    //如果是用户注册，先回包给上游服务，然后再处理和广告商的调用，避免阻塞上游服务
    int iRet = 0;
	if (iCmd == MAIN_CMD_GDT_REGISTER)
	{
		switch(iSubCmd)
		{
            case CMD_GDT_REGISTER_SVR_LOGIN_REPORT:
				iRet = ProcessRequestNoResponse<proto_gdt_register::LoginReportReq, NS_API_BASE::RSP_BASE,
				tr1::function<int(proto_gdt_register::LoginReportReq&, NS_API_BASE::RSP_BASE&)> >
				(pBodyBuf, iBodyLen, tr1::bind(&GdtRegister::ProcLoginReport, this, tr1::placeholders::_1, tr1::placeholders::_2));
				break;
		    default:
				break;
		}
	}
	else
	{
		ERROR("unknown main cmd:%d", iCmd);
		iRet = MUSIC_CODE_COMM_ParamInvalid;
	}

    return iRet;
}

bool GdtRegister::IsSameDay(time_t time_1, time_t time_2)
{
    DEBUG("time_1:%u, time_2:%u", (unsigned)time_1, (unsigned)time_2);

    tm tmTime1, tmTime2;
    localtime_r(&time_1, &tmTime1);
    localtime_r(&time_2, &tmTime2);
    DEBUG("tmTime1:%d year %d month %d day, tmTime2:%d year %d month %d day",
        tmTime1.tm_year, tmTime1.tm_mon, tmTime1.tm_mday, tmTime2.tm_year, tmTime2.tm_mon, tmTime2.tm_mday);

    if (tmTime1.tm_year == tmTime2.tm_year && tmTime1.tm_mon == tmTime2.tm_mon &&
            tmTime1.tm_mday == tmTime2.tm_mday)
    {
        return true;
    }

    return false;
}

int GdtRegister::SendCallBackToGDT(const RegisterInfo& register_info, const string &strQua,
        const string &strDeviceInfo, const WNS::WnsQua& qua)
{
    // 通过中台接的ams会进入这个if分支
    // 回调中台
    if (register_info.ad_company_type == AD_COMPANY_GDT && register_info.channel_id.find("_onlydc") != string::npos)
    {
        return SendCallBackToAMS(register_info, strQua, strDeviceInfo, qua);
    }
    int iRet = 0;
    //1 获取Token
    string strToken;
    //改用token永不过期的方案
    strToken = GetGdtConfigStr(register_info.advertiser_id, "token");
    DEBUG("token: %s.", strToken.c_str());
    //2 构造回调链接
    struct timeval tvLog;
    API_Log_StartTimer(tvLog);
    Json::Value jsonPost;
    Json::Value actions;
    Json::Value action;
    Json::Value userId;

    jsonPost["account_id"] = register_info.advertiser_id;
    if(qua.isIOS()){
        jsonPost["user_action_set_id"] = GetGdtConfigStr(register_info.advertiser_id, "user_action_set_id_ios");
        userId["hash_idfa"] = register_info.muid;
    }else{
        jsonPost["user_action_set_id"] = GetGdtConfigStr(register_info.advertiser_id, "user_action_set_id_android");
        if(register_info.muid.find("ext_") == string::npos){
            userId["hash_imei"] = register_info.muid;
        }
        userId["hash_android_id"] = register_info.strMAndroidId;
        userId["oaid"] = register_info.strOaid;
    }
    time_t tt;
    time(&tt);
    action["action_time"] = (int)tt;
    action["action_type"] = "ACTIVATE_APP";
    action["user_id"] = userId;
    actions.append(action);
    jsonPost["actions"] = actions;

    Json::FastWriter writer;
    string strJsonPost = writer.write(jsonPost);

    string url = "https://api.e.qq.com/v1.1/user_actions/add";
    time_t curTime;
    time(&curTime);
    url.append("?").append("access_token=").append(strToken)
        .append("&").append("timestamp=").append(klib::ToStr(curTime))
        .append("&").append("nonce=").append(genUuid());
    ERROR("url: %s, post json str: %s.", url.c_str(), strJsonPost.c_str());

    //3 curl链接
    string strHttpResult = "";
    iRet = CommonCurl(url, strJsonPost, strHttpResult);
    if(iRet != 0){
        ERROR("curl url[%s] postStr[%s] error.", url.c_str(), strJsonPost.c_str());
        return iRet;
    }

    //4 解析返回
    iRet = MUSIC_CODE_GDT_REGISTER_CALLBACK_GDT_FAIL;
    Json::Reader jReader;
    Json::Value jsonValue;
    if (jReader.parse(strHttpResult, jsonValue)){
        unsigned int code = jsonValue["code"].asUInt();
        if(code == 0){
            iRet = 0;
            ERROR("send call back to gdt success, rsp: %s", strHttpResult.c_str());
        }else{
            ERROR("gdt callback failed, code: %u.", code);
        }
    }else{
        ERROR("json parse error, jsonStr: %s.", jsonValue.asCString());
    }
    API_Log_StopAndWrite(tvLog, LM_ERROR, 0, GDT_REGISTER_MOD_ID, IF_GDT_REGISTER_POST_GDT_ACTION, "", "", 0, iRet, iRet!=0?LS_RET_FAIL:LS_RET_SUCC);
    return iRet;

}

int GdtRegister::ActiveSendCallBackToGDT(const RegisterInfo& register_info, const WNS::WnsQua& qua)
{
    int iRet = 0;
    //1 获取Token
    string strToken = GetGdtConfigStr(register_info.uAdvertiserIdActive, "token");
    DEBUG("token: %s.", strToken.c_str());
    //2 构造回调链接
    struct timeval tvLog;
    API_Log_StartTimer(tvLog);
    Json::Value jsonPost;
    Json::Value actions;
    Json::Value action;
    Json::Value userId;

    jsonPost["account_id"] = register_info.uAdvertiserIdActive;
    if(qua.isIOS()){
        jsonPost["user_action_set_id"] = GetGdtConfigStr(register_info.uAdvertiserIdActive, "user_action_set_id_ios");
        userId["hash_idfa"] = register_info.muid;
    }else{
        jsonPost["user_action_set_id"] = GetGdtConfigStr(register_info.uAdvertiserIdActive, "user_action_set_id_android");
        if(register_info.muid.find("ext_") == string::npos){
            userId["hash_imei"] = register_info.muid;
        }
        userId["hash_android_id"] = register_info.strMAndroidId;
        userId["oaid"] = register_info.strOaid;
    }
    time_t tt;
    time(&tt);
    action["action_time"] = (int)tt;
    action["action_type"] = "VIEW_CONTENT";
    action["user_id"] = userId;
    actions.append(action);
    jsonPost["actions"] = actions;

    Json::FastWriter writer;
    string strJsonPost = writer.write(jsonPost);

    string url = "https://api.e.qq.com/v1.1/user_actions/add";
    time_t curTime;
    time(&curTime);
    url.append("?").append("access_token=").append(strToken)
        .append("&").append("timestamp=").append(klib::ToStr(curTime))
        .append("&").append("nonce=").append(genUuid());
    ERROR("url: %s, post json str: %s.", url.c_str(), strJsonPost.c_str());

    //3 curl链接
    string strHttpResult = "";
    iRet = CommonCurl(url, strJsonPost, strHttpResult);
    if(iRet != 0){
        ERROR("curl url[%s] postStr[%s] error.", url.c_str(), strJsonPost.c_str());
        return iRet;
    }

    //4 解析返回
    iRet = MUSIC_CODE_GDT_REGISTER_CALLBACK_GDT_FAIL;
    Json::Reader jReader;
    Json::Value jsonValue;
    if (jReader.parse(strHttpResult, jsonValue)){
        unsigned int code = jsonValue["code"].asUInt();
        if(code == 0){
            iRet = 0;
            ERROR("send call back to gdt success, rsp: %s", strHttpResult.c_str());
        }else{
            ERROR("gdt callback failed, code: %u.", code);
        }
    }else{
        ERROR("json parse error, jsonStr: %s.", jsonValue.asCString());
    }
    API_Log_StopAndWrite(tvLog, LM_ERROR, 0, GDT_REGISTER_MOD_ID, IF_GDT_REGISTER_POST_GDT_ACTION, "", "", 0, iRet, iRet!=0?LS_RET_FAIL:LS_RET_SUCC);
    return iRet;

}

int GdtRegister::SendCallBackToAMS(const RegisterInfo& register_info, const string &strQua,
        const string &strDeviceInfo, const WNS::WnsQua& qua)
{
    Json::Value post_value;
    post_value["account_id"] = register_info.callbackparam;
    post_value["click_id"] = register_info.click_id;
    post_value["click_time"] = register_info.click_time;
    post_value["action_time"] = (unsigned)time(NULL);
    post_value["action_type"] = "ACTIVATE_APP";

    if(register_info.muid.find("ext_") == string::npos)
    {
        post_value["md5_imei"] = register_info.muid;
        post_value["user_id_type"] = "3";
        post_value["user_id"] = register_info.muid;
    }
    else if(register_info.muid.find("ext_oaid_") != string::npos)
    {
        post_value["user_id_type"] = "1";
        post_value["user_id"] = register_info.strOaid;
    }
    else if(register_info.muid.find("ext_androidid_") != string::npos)
    {
        post_value["user_id_type"] = "5";
        post_value["user_id"] = register_info.strMAndroidId;
    }

    Json::FastWriter writer;
    string strPostJson = writer.write(post_value);

    ERROR("SendCallBackToAMS, post:%s, qua:%s, strDeviceInfo:%s, register_time:%u",
            strPostJson.c_str(), strQua.c_str(), strDeviceInfo.c_str(), (unsigned)time(NULL));
    int iRet = 0;
    CURL *curl = NULL;
    curl = curl_easy_init();
    if (curl == NULL)
    {
        ERROR("curl_easy_init error");
        return -1;
    }
    struct timeval tvLog;
    API_Log_StartTimer(tvLog);
    curl_easy_setopt(curl, CURLOPT_SSL_VERIFYPEER, 0);
    curl_easy_setopt(curl, CURLOPT_SSL_VERIFYHOST, 0);
    curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, FuncStoreResult);
    curl_easy_setopt(curl, CURLOPT_VERBOSE, 0);
    curl_easy_setopt(curl, CURLOPT_USERAGENT, ".curl/7.15.1.(i686-suse-linux).libcurl/7.15.1.OpenSSL/0.9.8a.zlib/1.2.3.libidn/0.6.0" );
    curl_easy_setopt(curl, CURLOPT_POSTFIELDS, strPostJson.c_str());



    QOSREQUEST  stQosReq;
    L5Helper    stL5Client;
    stQosReq._modid = m_ams_mid;
    stQosReq._cmd = m_ams_cid;
    char err_msg[512];
    iRet = stL5Client.GetL5Route(stQosReq, err_msg, sizeof(err_msg));
    if (iRet != 0)
    {
        ERROR("GetL5Route error, iRet:%d, l5:%u:%u, err_msg:%s", iRet, m_ams_mid, m_ams_cid, err_msg);
        API_Log_StopAndWrite(tvLog, LM_ERROR, 0, GDT_REGISTER_MOD_ID, IF_GDT_REGISTER_POST_GDT_ACTION, "", "", 0, iRet, LS_RET_FAIL);
        return iRet;
    }

    snprintf(m_buf, m_bufsize, "http://%s:%d/imei", stQosReq._host_ip.c_str(), stQosReq._host_port);
    string strUrl = m_buf;
    DEBUG("url:%s", strUrl.c_str());
    string http_result;
    struct curl_slist *hs=NULL;
    hs = curl_slist_append(hs, "Content-Type: application/json");
    curl_easy_setopt(curl, CURLOPT_HTTPHEADER, hs);
    curl_easy_setopt(curl, CURLOPT_HEADER, 0);
    curl_easy_setopt(curl, CURLOPT_CONNECTTIMEOUT_MS, 3000);
    curl_easy_setopt(curl, CURLOPT_TIMEOUT, m_curl_bdtimeout);
    curl_easy_setopt(curl, CURLOPT_WRITEDATA, &http_result);
    curl_easy_setopt(curl, CURLOPT_URL, strUrl.c_str());

    //发送请求
    CURLcode res = curl_easy_perform(curl);
    if (CURLE_OPERATION_TIMEDOUT == res)
    {
        stL5Client.UpDataL5Route(stQosReq, ENUM_SOCKRSP_ERR_RECV);
    }
    else
    {
        stL5Client.UpDataL5Route(stQosReq, 0);
    }
    if (res != CURLE_OK)
    {
        ERROR("curl_easy_perform[%d] errmsg=%s, url:%s", res, curl_easy_strerror(res), strUrl.c_str());
        curl_easy_cleanup(curl);
        API_Log_StopAndWrite(tvLog, LM_ERROR, 0, GDT_REGISTER_MOD_ID, IF_GDT_REGISTER_POST_GDT_ACTION, "", "", 0, res, LS_RET_FAIL);
        return -2;
    }
    long http_status;
    curl_easy_getinfo(curl, CURLINFO_RESPONSE_CODE, &http_status);
    DEBUG("http_status:%ld", http_status);
    DEBUG("http_result:%s", http_result.c_str());
    curl_easy_cleanup(curl);
    curl_slist_free_all(hs);

    //处理结果
    Json::Reader jReader;
    Json::Value json_ret;
    iRet = 2;
    if (http_status == 200)
    {
        iRet = 0;
    }
    else
    {
        ERROR("error_code:%s, url:%s", http_result.c_str(), strUrl.c_str());
    }
    API_Log_StopAndWrite(tvLog, LM_ERROR, 0, GDT_REGISTER_MOD_ID, IF_GDT_REGISTER_POST_GDT_ACTION, "", "", 0, iRet, iRet!=0?LS_RET_FAIL:LS_RET_SUCC);
    return iRet;
}

int GdtRegister::refresh_token(const RegisterInfo& register_info)
{
    //定时更新token
    unsigned int now_time = time(NULL);
    if (register_info.advertiser_id == 10833188) // B代理
    {
        unsigned int dist = now_time - g_token_time_bproxy;
        if(dist > m_token_timeout)
        {
            g_token_time_bproxy = now_time;
            DEBUG("g_token_time_bproxy:%u, now_time:%u, m_token_timeout:%u, dist:%u", g_token_time_bproxy, now_time, m_token_timeout, dist);
        }
        else
        {
            DEBUG("use old token:%s", g_access_token.c_str());
            return 0;
        }
    }
    else
    {
        unsigned int dist = now_time - g_token_time;
        if(dist > m_token_timeout)
        {
            g_token_time = now_time;
            DEBUG("g_token_time:%u, now_time:%u, m_token_timeout:%u, dist:%u", g_token_time, now_time, m_token_timeout, dist);
        }
        else
        {
            DEBUG("use old token:%s", g_access_token.c_str());
            return 0;
        }
    }

    int iRet = 0;
	CURL *curl = NULL;
    curl = curl_easy_init();
    if( curl == NULL )
    {
        ERROR("curl_easy_init error");
        return -1;
    }
    struct timeval tvLog;
    API_Log_StartTimer(tvLog);
    curl_easy_setopt(curl, CURLOPT_SSL_VERIFYPEER, 0);
	curl_easy_setopt(curl, CURLOPT_SSL_VERIFYHOST, 0);
    curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, FuncStoreResult);
    curl_easy_setopt(curl, CURLOPT_VERBOSE, 0);
    curl_easy_setopt(curl, CURLOPT_NOSIGNAL, 1 );
    curl_easy_setopt(curl, CURLOPT_USERAGENT, ".curl/7.15.1.(i686-suse-linux).libcurl/7.15.1.OpenSSL/0.9.8a.zlib/1.2.3.libidn/0.6.0" );
	//curl_easy_setopt(curl, CURLOPT_PROTOCOLS, CURLPROTO_HTTP | CURLPROTO_HTTPS);
    //curl_easy_setopt(curl, CURLOPT_POSTFIELDS, post.c_str());
    // curl init end

    QOSREQUEST  _qos_req;
    L5Helper    _l5_client;
    _qos_req._modid = m_curl_mid;
    _qos_req._cmd = m_curl_cid;
    char err_msg[512];
    iRet = _l5_client.GetL5Route(_qos_req,err_msg,sizeof(err_msg));
    if( iRet != 0 )
    {
        ERROR("_l5_client.GetL5Route error, iRet:%d, l5:%u:%u, err_msg:%s", iRet, m_curl_mid, m_curl_cid, err_msg);
        API_Log_StopAndWrite(tvLog, LM_ERROR, 0, GDT_REGISTER_MOD_ID, IF_GDT_REGISTER_REFRESH_TOKEN, "", "", 0, iRet, LS_RET_FAIL);
        return iRet;
    }

    //http://100.107.138.75/oauth/token?client_id=1107986938&client_secret=JnfLPtpaA3QSwrSi&grant_type=refresh_token&refresh_token=858b88c24750006a6acb6474ce570282
    if (register_info.advertiser_id == 10833188) // B代理
    {
        snprintf(m_buf, m_bufsize, "http://%s/oauth/token?client_id=1107986938&client_secret=JnfLPtpaA3QSwrSi&grant_type=refresh_token&refresh_token=4e3892a043f897f48d6406594d5c9148", _qos_req._host_ip.c_str());
    }
    else
    {
        snprintf(m_buf, m_bufsize, "http://%s/oauth/token?client_id=1107986938&client_secret=JnfLPtpaA3QSwrSi&grant_type=refresh_token&refresh_token=858b88c24750006a6acb6474ce570282", _qos_req._host_ip.c_str());
    }

    string url = m_buf;
    DEBUG("url:%s", url.c_str());
    string http_result;
    curl_easy_setopt(curl, CURLOPT_HEADER, 0);
    //curl_easy_setopt(curl, CURLOPT_PORT, 80);
    curl_easy_setopt(curl, CURLOPT_CONNECTTIMEOUT_MS, 3000);
    curl_easy_setopt(curl, CURLOPT_TIMEOUT_MS, m_curl_timeout_ms);
    curl_easy_setopt(curl, CURLOPT_WRITEDATA, &http_result);
    curl_easy_setopt(curl, CURLOPT_URL, url.c_str() ); /*Set URL*/
    CURLcode res = curl_easy_perform(curl);

    if(CURLE_OPERATION_TIMEDOUT == res)
    {
        _l5_client.UpDataL5Route(_qos_req, ENUM_SOCKRSP_ERR_RECV);
    }
    else
    {
        _l5_client.UpDataL5Route(_qos_req, 0);
    }

    if(res != CURLE_OK)
    {
        ERROR("curl_easy_perform[%d] errmsg=%s, url:%s", res, curl_easy_strerror(res), url.c_str());
        curl_easy_cleanup(curl);
        API_Log_StopAndWrite(tvLog, LM_ERROR, 0, GDT_REGISTER_MOD_ID, IF_GDT_REGISTER_REFRESH_TOKEN, "", _qos_req._host_ip, 0, res, LS_RET_FAIL);
        return -2;
    }

    long http_status;
    curl_easy_getinfo(curl, CURLINFO_RESPONSE_CODE, &http_status);
    DEBUG("http_status:%ld", http_status);
    DEBUG("http_result:%s", http_result.c_str());

    curl_easy_cleanup(curl);

    Json::Reader jReader;
    Json::Value json_ret;
    // 0:success, 1: file not found, 2: failed
    iRet = 2;    //先置为出错
    if (jReader.parse(http_result, json_ret))
    {
        if (json_ret.isMember("code"))
        {
        	iRet = json_ret["code"].asInt();
            if(iRet != 0)
            {
                ERROR("error_code iRet:%d, msg:%s, url:%s", iRet, json_ret["message"].asString().c_str(), url.c_str());
            }
            else
            {
                if (register_info.advertiser_id == 10833188) // B代理
                {
                    g_access_token_bproxy= json_ret["data"]["access_token"].asString();
                }
                else
                {
                    g_access_token = json_ret["data"]["access_token"].asString();
                }
            }
        }
        else
        {
            ERROR("no code, http_result:%s", http_result.c_str());
        }
    }
    else
    {
        ERROR("jReader.parse error, http_result:%s", http_result.c_str());
    }
    API_Log_StopAndWrite(tvLog, LM_ERROR, 0, GDT_REGISTER_MOD_ID, IF_GDT_REGISTER_REFRESH_TOKEN, "", _qos_req._host_ip, 0, iRet, iRet!=0?LS_RET_FAIL:LS_RET_SUCC);
    return iRet;
}

int GdtRegister::CommonCurl(const string& url, const string& strPostJson, string& strResult, unsigned uHeaderType, unsigned uReqType)
{
    int iRet = 0;
	CURL *curl = NULL;
    curl = curl_easy_init();
    if (curl == NULL)
    {
        ERROR("curl_easy_init error");
        return MUSIC_CODE_GDT_REGISTER_CURL_FAIL;
    }

    string http_result;
    curl_easy_setopt(curl, CURLOPT_SSL_VERIFYPEER, 0);
	curl_easy_setopt(curl, CURLOPT_SSL_VERIFYHOST, 0);
    curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, FuncStoreResult);
    curl_easy_setopt(curl, CURLOPT_VERBOSE, 0);
    curl_easy_setopt(curl, CURLOPT_NOSIGNAL, 1 );
    curl_easy_setopt(curl, CURLOPT_USERAGENT, ".curl/7.15.1.(i686-suse-linux).libcurl/7.15.1.OpenSSL/0.9.8a.zlib/1.2.3.libidn/0.6.0" );
    if (strPostJson != "")
    {
        curl_easy_setopt(curl, CURLOPT_POSTFIELDS, strPostJson.c_str());
    }
    struct curl_slist *hs = NULL;

    // HEADER
    if (uHeaderType == 1)
    {
        hs = curl_slist_append(hs, "Content-Type: application/json");
    }
    else
    {
        hs = curl_slist_append(hs, "Content-Type: application/x-www-form-urlencoded");
    }

    // GET or POST
    if (uReqType == 2)
    {
        curl_easy_setopt(curl, CURLOPT_CUSTOMREQUEST, "POST");
    }

    curl_easy_setopt(curl, CURLOPT_HTTPHEADER, hs);
    curl_easy_setopt(curl, CURLOPT_HEADER, 0);
    curl_easy_setopt(curl, CURLOPT_CONNECTTIMEOUT_MS, 3000);
    curl_easy_setopt(curl, CURLOPT_TIMEOUT_MS, m_curl_timeout_ms);
    curl_easy_setopt(curl, CURLOPT_WRITEDATA, &http_result);
    curl_easy_setopt(curl, CURLOPT_URL, url.c_str() ); /*Set URL*/
    CURLcode res = curl_easy_perform(curl);
    if (res != CURLE_OK)
    {
        double total, nl, con, trans = 0.0;
        GetCurlTime(curl, total, nl, con, trans);
        ERROR("curl_easy_perform[%d] errmsg=%s, url:%s total=%f namelookup=%f connect=%f trans=%f", res, curl_easy_strerror(res), url.c_str(), total, nl, con, trans);
        curl_easy_cleanup(curl);
        return MUSIC_CODE_GDT_REGISTER_CURL_FAIL;
    }
    strResult = http_result;
    long http_status;
    curl_easy_getinfo(curl, CURLINFO_RESPONSE_CODE, &http_status);
    DEBUG("url:%s, strPostJson:%s, http_status:%ld, http_result:%s",
            url.c_str(), strPostJson.c_str(), http_status, http_result.c_str());
    curl_easy_cleanup(curl);
    curl_slist_free_all(hs);
    if (http_status != 200)
    {
        ERROR("http rsp error, url:%s, strPostJson:%s, status: %ld.", url.c_str(), strPostJson.c_str(), http_status);
        return MUSIC_CODE_GDT_REGISTER_CURL_FAIL;
    }
    return 0;
}

int GdtRegister::GetGdtToken(unsigned int advertiseId, string& strToken){
    int iRet = 0;

    //1 从CKV中获取缓存的 token
    string key = GDT_TOKEN_KEY + klib::ToStr(advertiseId);
    string data = "";
    int cas = 0;
    iRet = m_tmem_obj.Get(key, data, cas, 0, GDT_REGISTER_TMEM_MOD_ID);
    if(iRet == 0){
        GdtToken gdtToken;
        int iLen = data.length();
        iRet = gdtToken.Decode((uint8_t*)data.c_str(), &iLen, NULL);
        if (iRet != 0)
        {
            ERROR("decode error, iRet:%d, key:%s", iRet, key.c_str());
            return iRet;
        }
        DEBUG("token exists in ckv: %s.", gdtToken.strToken.c_str());
        strToken = gdtToken.strToken;
        return 0;
    }
    if(iRet != 0 && iRet != -13200 && iRet != -13106){
        ERROR("ckv get key[%s] error, iRet: %d.", key.c_str(), iRet);
        return iRet;
    }
    DEBUG("Token not exist or expired, refresh it.");

    //2 ckv中不存在token（过期），curl从广点通平台更新token
    struct timeval tvLog;
    API_Log_StartTimer(tvLog);
    string url = "https://api.e.qq.com/oauth/token";
    url.append("?").append("client_id=").append(GetGdtConfigStr(advertiseId, "client_id"))
        .append("&").append("client_secret=").append(GetGdtConfigStr(advertiseId, "client_secret"))
        .append("&").append("grant_type=refresh_token")
        .append("&").append("refresh_token=").append(GetGdtConfigStr(advertiseId, "refresh_token"));
    DEBUG("url: %s.", url.c_str());
    string httpResult;
    iRet = CommonCurl(url, "", httpResult);
    if(iRet != 0){
        ERROR("curl url %s fail, iRet: %d", url.c_str(), iRet);
        API_Log_StopAndWrite(tvLog, LM_ERROR, 0, GDT_REGISTER_MOD_ID, IF_GDT_REGISTER_REFRESH_TOKEN, "", "", 0, iRet, LS_RET_FAIL);
        return iRet;
    }

    //3 解析curl结果
    Json::Reader jReader;
    Json::Value jsonValue;
    iRet = MUSIC_CODE_GDT_REGISTER_REFRESH_TOKEN_FAIL;
    if (jReader.parse(httpResult, jsonValue))
    {
        if (jsonValue.isMember("access_token"))
        {
                strToken = jsonValue["access_token"].asString();
                GdtToken gdtToken;
                gdtToken.strToken = strToken;
                EncodeAndSet(key, gdtToken, 23*60*60);
                iRet = 0;
        }
        else
        {
            ERROR("refresh token error: no access_token return. jsonValue: %s.", httpResult.c_str());
        }
    }
    else
    {
        ERROR("json parse error, jsonValue: %s.", httpResult.c_str());
    }
    API_Log_StopAndWrite(tvLog, LM_ERROR, 0, GDT_REGISTER_MOD_ID, IF_GDT_REGISTER_REFRESH_TOKEN, "", "", 0, iRet, iRet!=0?LS_RET_FAIL:LS_RET_SUCC);
    return iRet;
}

// AMS 拉新回调
// https://developers.e.qq.com/docs/guide/conversion/api
int GdtRegister::SendCallBackToAMSSelf(RegisterInfo &stRegisterInfo, const WNS::WnsQua &stQua,
        const unsigned uLoginTime, const string &strType)
{
    int iRet = 0;

    // [step 1] 填充需要发送的数据
    ActionsData stActionsData;
    stActionsData.action_time = uLoginTime;
    stActionsData.action_type = strType;

    UserIdData &stUserIdData = stActionsData.user_id;
    if (stQua.isAndroid())
    {
        stUserIdData.hash_imei = stRegisterInfo.muid;
        stUserIdData.oaid = stRegisterInfo.strOaid;
        stUserIdData.hash_android_id = stRegisterInfo.strMAndroidId;
    }
    else
    {
        stUserIdData.hash_idfa = stRegisterInfo.muid;
    }

    AMSData stAMSData;
    stAMSData.actions.push_back(stActionsData);

    // [step 2] 回调
    string strJsonPost = stAMSData.writeToJsonString();
    string strUrl = stRegisterInfo.callbackparam;
    string strHttpResult = "";
    ERROR("url:%s, strJsonPost:%s", strUrl.c_str(), strJsonPost.c_str());

    struct timeval tvLog;
    API_Log_StartTimer(tvLog);

    iRet = CommonCurl(strUrl, strJsonPost, strHttpResult, 1, 2);
    if (iRet)
    {
        ERROR("CommonCurl err, strUrl:%s, strJsonPost:%s", strUrl.c_str(), strJsonPost.c_str());
        return iRet;
    }
    DEBUG("CommonCurl succ, strUrl:%s, strJsonPost:%s, strHttpResult:%s",
            strUrl.c_str(), strJsonPost.c_str(), strHttpResult.c_str());

    // [step 4] 解析返回
    Json::Reader jReader;
    Json::Value jsonValue;
    if (jReader.parse(strHttpResult, jsonValue))
    {
        unsigned int code = jsonValue["code"].asUInt();
        if (code)
        {
            iRet = MUSIC_CODE_GDT_REGISTER_CALLBACK_GDT_FAIL;
            ERROR("call ams url:%s err, strJsonPost:%s, code:%u", strUrl.c_str(), strJsonPost.c_str(), code);
        }
        else
        {
            DEBUG("call ams url:%s succ, strJsonPost:%s", strUrl.c_str(), strJsonPost.c_str());
        }
    }
    else
    {
        ERROR("json parse err, strUrl:%s, strJsonPost:%s", strUrl.c_str(), strJsonPost.c_str());
    }

    API_Log_StopAndWrite(tvLog, LM_ERROR, 0, GDT_REGISTER_MOD_ID, IF_GDT_REGISTER_POST_GDT_ACTION, "", "", 0, iRet, iRet!=0?LS_RET_FAIL:LS_RET_SUCC);

    return iRet;
}

int GdtRegister::GetCurlTime(CURL* curl, double& dTotalTime, double& dDNSTime, double& dConnectTime, double& dTranferTime)
{
    curl_easy_getinfo(curl, CURLINFO_TOTAL_TIME, &dTotalTime);
    curl_easy_getinfo(curl, CURLINFO_NAMELOOKUP_TIME, &dDNSTime);
    curl_easy_getinfo(curl, CURLINFO_CONNECT_TIME, &dConnectTime);
    curl_easy_getinfo(curl, CURLINFO_STARTTRANSFER_TIME, &dTranferTime);
    return 1;
}

#define QZASYNCMSG_CLASS GdtRegister

//这个放在最后面，请不要修改位置
#include "qza_svr_frame.h"

