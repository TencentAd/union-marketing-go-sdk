//create by create_spp_server.sh
/**
 * 名称: 全民K歌
 * (c) Copyright 2020 wrenli@tencent.com. All Rights Reserved.
 *
 * change log:
 * Author: wrenli
 *   Date: 2020-11-26
 *	 Desc:
 */
#include "server_config.h"
#include "klib.h"
#include "gdt_register.h"

#define g_conf (*CSingleton<CServerConf>::instance())

// 处理广告商的点击请求
int GdtRegister::ProcReportClick(ReportClickReq& req, NS_API_BASE::RSP_BASE& rsp)
{
    DEBUG("req:%s", klib::jce2str(req).c_str());

    // [step 1] 获取req里面所有的设备字段，md5加密后的小写
    // [step 1.1] muid 是md5小写加密后的imei或者idfa
    string strMuid;
    FillInMuid(req, strMuid);

    if (strMuid != "")
    {
        SetClickTime(strMuid, req.click_time);

        // 先读取存储中已有的register info，再根据req填充register info，最后写入存储
        UpdateRegisterInfo(req, strMuid, true);
        DEBUG("UpdateRegisterInfo key: %s", strMuid.c_str());
    }

    // [step 1.2] oaid md5加密后的小写
    string strMd5Oaid = (req.strOaidMd5 != "" ? req.strOaidMd5 : md5Lower(req.strOaid));
    if (strMd5Oaid != "")
    {
        strMd5Oaid = "ext_oaid_" + strMd5Oaid;
        SetClickTime(strMd5Oaid, req.click_time);

        // 先读取存储中已有的register info，再根据req填充register info，最后写入存储
        UpdateRegisterInfo(req, strMd5Oaid);
        DEBUG("UpdateRegisterInfo key: %s", strMd5Oaid.c_str());
    }

    // [step 1.3] android id md5加密后的小写
    if (req.strMAndroidId != "")
    {
        string strMd5AndroidId = "ext_androidid_" + req.strMAndroidId;
        SetClickTime(strMd5AndroidId, req.click_time);

        // 先读取存储中已有的register info，再根据req填充register info，最后写入存储
        UpdateRegisterInfo(req, strMd5AndroidId);
        DEBUG("UpdateRegisterInfo key: %s", strMd5AndroidId.c_str());
    }

    // [step 1.4] ipua模糊归因有个开关，只对特定的渠道生效
    if (g_conf["ipua_enable."+klib::ToStr(req.ad_company_type)] == "1")
    {
        string strMd5IpUa = GenerateMuidByIpUa(req);
        if (strMd5IpUa != "")
        {
            SetClickTime(strMd5IpUa, req.click_time);

            // 先读取存储中已有的register info，再根据req填充register info，最后写入存储
            UpdateRegisterInfo(req, strMd5IpUa);
            DEBUG("UpdateRegisterInfo key: %s", strMd5IpUa.c_str());
        }
    }
    DEBUG("muid:%s md5oaid:%s md5androidid:%s.", strMuid.c_str(), strMd5Oaid.c_str(), req.strMAndroidId.c_str());

    return 0;
}

void GdtRegister::UpdateRegisterInfo(ReportClickReq &req, string &strDeviceId, bool bUpdateMuid/*=false*/)
{
    RegisterInfo register_info;
    int iRet = 0;

    for(unsigned i = 0; i < (unsigned)atoi(g_conf["global.retry_times"].c_str()); ++i)
    {
        // [step 1] 先读取存储中已有的register info
        int iCas = -1;
        GetRegisterInfo(strDeviceId, register_info, iCas);

        if (bUpdateMuid)
        {
            register_info.muid = strDeviceId;
        }

        // [step 2] 再根据req填充register info
        FillInRegisterInfo(req, register_info);

        // [step 3] 写回存储
        iRet = EncodeAndSet(strDeviceId, register_info, iCas);
        if (iRet == 0)
        {
            DEBUG("EncodeAndSet register_info succ, deviceid:%s, req:%s, register_info:%s",
                    strDeviceId.c_str(), klib::jce2str(req).c_str(), klib::jce2str(register_info).c_str());
            break;
        }
        else if (g_conf["invalid_deviceid."+strDeviceId] == "1")
        {
            DEBUG("EncodeAndSet register_info err, but invalid deviceid:%s, req:%s, register_info:%s, no more retry",
                    strDeviceId.c_str(), klib::jce2str(req).c_str(),klib::jce2str(register_info).c_str());
            break;
        }
        ERROR("EncodeAndSet register_info err, deviceid:%s, req:%s, register_info:%s, iRet:%d",
                strDeviceId.c_str(), klib::jce2str(req).c_str(), klib::jce2str(register_info).c_str(), iRet);
    }
}

void GdtRegister::GetRegisterInfo(string &strDeviceId, RegisterInfo &register_info, int &iCas)
{
    string strData;
    int iRet = m_tmem_obj.Get(strDeviceId, strData, iCas, 0, GDT_REGISTER_TMEM_MOD_ID);
    if (iRet && iRet != -13200 && iRet != -13106)
    {
        ERROR("deviceid:%s, Get RegisterInfo err, iRet:%d, errMsg:%s", strDeviceId.c_str(), iRet, m_tmem_obj.GetLastError());
    }
    else
    {
        DEBUG("deviceid:%s, Get RegisterInfo succ, iRet:%d", strDeviceId.c_str(), iRet);
    }

    int iLen = strData.length();
    iRet = register_info.Decode((uint8_t*)strData.c_str(), &iLen, NULL);
    if (iRet)
    {
        ERROR("deviceid:%s, Decode register_info err, iRet:%d", strDeviceId.c_str(), iRet);
    }
    else
    {
        DEBUG("deviceid:%s Decode register_info succ, RegisterInfo:%s", strDeviceId.c_str(), klib::jce2str(register_info).c_str());
    }
}

void GdtRegister::FillInRegisterInfo(ReportClickReq &req, RegisterInfo &register_info)
{
    // 拉新拉活渠道各存一份数据，避免互相覆盖
    if (req.ad_company_type == AD_COMPANY_AMS_AC)
    {
        register_info.click_time_active = req.click_time;
        register_info.ad_company_type_active = req.ad_company_type;
        register_info.channel_id_active = req.channel_id;
        register_info.callbackparam_active = req.callbackparam;
        register_info.uAdvertiserIdActive = req.advertiser_id;

        // 广告账户相关
        register_info.strAdAccountIdActive = req.strAdAccountId;
        register_info.strAdGroupIdActive = req.strAdGroupId;
        register_info.strAdCreativeIdActive = req.strAdCreativeId;
        register_info.strAdPlanIdActive = req.strAdPlanId;
    }
    else
    {
        register_info.click_time = req.click_time;
        register_info.ad_company_type = req.ad_company_type;
        register_info.channel_id = req.channel_id;

        register_info.click_id = req.click_id;
        register_info.advertiser_id = req.advertiser_id;
        register_info.appid = req.appid;

        // 如果是中台AMS传过来的点击信息，仅仅需要做上报，不需要撞库
        if (req.ad_company_type == AD_COMPANY_GDT && req.channel_id.find("_onlydc") != string::npos)
        {
            ParseParamsForAMS(register_info, req.callbackparam);
        }
        else
        {
            register_info.callbackparam = req.callbackparam;
        }

        // 广告账户相关
        register_info.strAdAccountId = req.strAdAccountId;
        register_info.strAdGroupId = req.strAdGroupId;
        register_info.strAdCreativeId = req.strAdCreativeId;
        register_info.strAdPlanId = req.strAdPlanId;
    }

    register_info.app_type = req.app_type;
    register_info.strAdvertiserName = req.strAdvertiserName;

    // 素材上报相关
    register_info.uMaterialType = req.uMaterialType;
    register_info.strKSongMid = req.strKSongMid;
    register_info.strUserKid = req.strUserKid;
    register_info.uFuncMaterialType = req.uFuncMaterialType;

    // 补充归因
    register_info.strMAndroidId = req.strMAndroidId;
    register_info.strOaid = req.strOaid;

    // device信息相关
    // registerinfo尽量不要新增字段了，回传时尽可能使用deviceinfo中的字段
    register_info.strImeiMd5 = req.strImeiMd5;
    register_info.strOaidMd5 = req.strOaidMd5;
    register_info.strIdfaMd5 = req.strIdfaMd5;
    register_info.strUserAgent = req.strUserAgent;
    register_info.strIp = req.strIp;
}

void GdtRegister::FillInMuid(ReportClickReq &req, string &muid)
{
    // register_info.muid 指的是 md5小写加密之后的idfa或者imei
    if (req.ad_company_type == AD_COMPANY_AMS)
    {
        muid = GenerateMuidForAMS(req);
    }
    // AD_COMPANY_AMS_AC
    // AD_COMPANY_GDT
    // 所有走中台过来的渠道请求以及ams拉活的请求都会落入这个分支
    // imei idfa全部放到了muid字段
    else
    {
        muid = req.muid;
    }
}

string GdtRegister::GetSysVerFromUserAgent(int nSystem, size_t nSystemIdx, const string& strUserAgent) 
{
    size_t nVerBegIdx = 0;
    if (nSystem == 0)//Android
    {
        nVerBegIdx = nSystemIdx + 7;
    }
    else if (nSystem == 1)//iPhone OS
    {
        nVerBegIdx = nSystemIdx + 9;
    }

    string strVer = "";
    bool bFind = false;
    for (size_t idx = nVerBegIdx; idx < strUserAgent.length(); ++idx)
    {
        if (isdigit(strUserAgent[idx]) || strUserAgent[idx] == '_' || strUserAgent[idx] == '.')
        {
            bFind = true;
            if (strUserAgent[idx] == '_')
            {
                strVer += '.';
            }
            else
            {
                strVer += strUserAgent[idx];
            }
        }
        else if (bFind == false)//如果没找到过目标符号，往后继续找
        {
            continue;
        }
        else//如果找到过目标符号，且发现当前符号不是目标符号，终止查找
        {
            break;
        }
    }
    if (strVer.length() == 0)
    {
        return "";
    }
    string strSystemVer = (nSystem == 0 ? "android" : "iphone");
    strSystemVer += strVer;
    return strSystemVer;
}

string GdtRegister::GenerateMuidForAMS(const ReportClickReq& req)
{
    string strMd5Muid;
    if (req.os_type == 0)     //android IMEI 已经进行MD5处理过
    {
        if (req.muid.size() == 0)
        {
            return "";
        }
        strMd5Muid = req.muid;
        std::transform(strMd5Muid.begin(), strMd5Muid.end(), strMd5Muid.begin(), ::tolower);
        if (strMd5Muid == "null")//特殊字符过滤
        {
            return "";
        }
    }
    else  //ios
    {
        if (req.idfa.size() == 0)
        {
            return "";
        }
        if(g_conf["ipua_enable."+klib::ToStr(req.ad_company_type)] == "1" && g_conf["invalid_deviceid." + req.idfa] == "1")
        {
            return "";
        }
        strMd5Muid = req.idfa;
        std::transform(strMd5Muid.begin(), strMd5Muid.end(), strMd5Muid.begin(), ::tolower);
        if (strMd5Muid == "null")//特殊字符过滤
        {
            return "";
        }
    }
    return strMd5Muid;
}

//md5后转小写
string GdtRegister::md5Lower(const string &str)
{
    if (str == "")
    {
        return "";
    }
    char res[33];
    NS_KG_MD5::md5_string(str.c_str(), str.length(), res);
    string strMd5Muid = string(res);
    std::transform(strMd5Muid.begin(), strMd5Muid.end(), strMd5Muid.begin(), ::tolower);
    return strMd5Muid;
}

string GdtRegister::GenerateMuidByIpUa(const ReportClickReq& req)
{
    const string& strIp = req.strIp;
    const string& strUserAgent = req.strUserAgent;
    if (strIp.size() ==  0 || strUserAgent.size() == 0)
    {
        return "";
    }

    string strSystemVer = "";
    size_t nSystemIdx = string::npos;
    if ((nSystemIdx = strUserAgent.find("Android")) != string::npos ||
        (nSystemIdx = strUserAgent.find("android")) != string::npos)
    {
        strSystemVer = GetSysVerFromUserAgent(0, nSystemIdx, strUserAgent);
    }
    else if ((nSystemIdx = strUserAgent.find("iPhone OS")) != string::npos)
    {
        strSystemVer = GetSysVerFromUserAgent(1, nSystemIdx, strUserAgent);
    }
    if (strSystemVer.size() == 0)
    {
        return "";
    }
    string strMd5Target = strIp + strSystemVer;
    DEBUG("strMd5Target=%s", strMd5Target.c_str());

    char res[33];
    NS_KG_MD5::md5_string(strMd5Target.c_str(), strMd5Target.length(), res);
    string md5str = string(res);
    std::transform(md5str.begin(), md5str.end(), md5str.begin(), ::tolower);
    return md5str;
}

int GdtRegister::EncodeAndSet(string& key, RegisterInfo& register_info, int iCas/*=-1*/)
{
    int iLen = m_bufsize;
    int iRet = register_info.Encode((uint8_t*)m_buf, &iLen, NULL);
    if(iRet != 0)
    {
        ERROR("register_info.Encode error, iRet:%d", iRet);
        return iRet;
    }
    string data;
    data.assign(m_buf, iLen);
    int cas = -1;
    iRet = m_tmem_obj.Set(key, data, cas, 0, GDT_REGISTER_TMEM_MOD_ID, m_tmem_expire);
    DEBUG("m_tmem_obj.Set succ key:%s m_tmem_expire:%d", key.c_str(), m_tmem_expire);
    if(iRet != 0)
    {
        ERROR("m_tmem_obj.Set error, iRet:%d, muid:%s, adcid:%d, error:%s",
                iRet, key.c_str(), register_info.ad_company_type, m_tmem_obj.GetLastError());
        return iRet;
    }
    return 0;
}

int GdtRegister::EncodeAndSet(string& key, GdtToken& gdtToken, unsigned int expireSec)
{
    int iLen = m_bufsize;
    int iRet = gdtToken.Encode((uint8_t*)m_buf, &iLen, NULL);
    if(iRet != 0)
    {
        ERROR("gdtToken.Encode error, iRet:%d", iRet);
        return iRet;
    }
    string data;
    data.assign(m_buf, iLen);
    int cas = -1;
    iRet = m_tmem_obj.Set(key, data, cas, 0, GDT_REGISTER_TMEM_MOD_ID, expireSec);
    DEBUG("m_tmem_obj.Set succ key:%s m_tmem_expire:%d", key.c_str(), expireSec);
    if(iRet != 0)
    {
        ERROR("m_tmem_obj.Set error, iRet:%d, key:%s, error:%s", iRet, key.c_str(), m_tmem_obj.GetLastError());
        return iRet;
    }
    return 0;
}

int GdtRegister::EncodeAndSetSecondary(RegisterInfo& register_info)
{
    char buff[102400];
    int iLen = sizeof(buff);
    int iRet = register_info.Encode((uint8_t*)buff, &iLen, NULL);
    if (iRet)
    {
        ERROR("register_info Encode err, uid:%u, iRet:%d", register_info.uUid, iRet);
        return iRet;
    }
    string data;
    data.assign(buff, iLen);
    int cas = -1;
    string key = SECONDARY_INFO_KEY + klib::ToStr(register_info.uUid);
    unsigned expireSec = (unsigned)atoi(g_conf["secondary_info_ckv.expire_time_second"].c_str());
    iRet = m_tmem_obj.Set(key, data, cas, 0, GDT_REGISTER_TMEM_MOD_ID, expireSec);
    if (iRet)
    {
        ERROR("m_tmem_obj Set secondary info err, uid:%u, iRet:%d", register_info.uUid, iRet);
        return iRet;
    }
    DEBUG("m_tmem_obj Set secondary info succ, register_info:%s", klib::jce2str(register_info).c_str());
    return 0;
}

void GdtRegister::ParseParamsForAMS(RegisterInfo& register_info, const string& strCallBackParam)
{
    vector<string> arrParamsPair;
    map<string, string> mapParams;
    ns_kg::CWebappCommon::SpliteStr(strCallBackParam, arrParamsPair, '|');
    if (arrParamsPair.size() == 0)
    {
        return;
    }

    for (size_t nIdx = 0; nIdx < arrParamsPair.size(); ++nIdx)
    {
        vector<string> arrKeyValue;
        ns_kg::CWebappCommon::SpliteStr(arrParamsPair[nIdx], arrKeyValue, ',');
        if (arrKeyValue.size() != 2)
        {
            continue;
        }
        mapParams[arrKeyValue[0]] = arrKeyValue[1];
    }

    register_info.click_time = mapParams["click_time"];
    register_info.click_id = mapParams["click_id"];
    register_info.callbackparam = mapParams["advertiser_id"];
    register_info.advertiser_id = (unsigned)atoll(mapParams["ad_id"].c_str());
    register_info.appid = (unsigned)atoll(mapParams["adgroup_id"].c_str());
    return;
}

int GdtRegister::SetClickTime(const string &strDeviceId, const string &strClickTime)
{
    int iRet = 0;

    // [step 1] 获取 clicktime 有些渠道是ms 需要统一处理成s
    string strClickTimeS = (strClickTime.size() > 10 ? strClickTime.substr(0, 10) : strClickTime);
    unsigned uClickTime = (unsigned)atoll(strClickTimeS.c_str());

    // [step 2] 生成clickmsg并encode
    ClickMsg stClickMsg;
    stClickMsg.strDeviceId = strDeviceId;
    stClickMsg.uClickTime = uClickTime;

    char szBuff[102400];
    int iLen = sizeof(szBuff);
    iRet = stClickMsg.Encode((uint8_t*)szBuff, &iLen, NULL);
    if (iRet)
    {
        ERROR("stClickMsg:%s encode err, iRet:%d", klib::jce2str(stClickMsg).c_str(), iRet);
        return iRet;
    }

    // [step 3] 写入kafka
    string strTmpMsg(szBuff, iLen);

    KafkaMsg stKafkaMsg;
    stKafkaMsg.uMsgType = CLICK_MSG_TYPE;
    stKafkaMsg.strMsg = taf::TC_Base64::encode(strTmpMsg);

    struct timeval tvLog;
    API_Log_StartTimer(tvLog);
    iRet = g_kafkaProducer.produce_msg(stKafkaMsg.writeToJsonString());
    API_Log_StopAndWrite(tvLog, LM_ERROR, 0, GDT_REGISTER_MOD_ID, 0, LS_OP_INSERT, "", 0, iRet, iRet==0?LS_RET_SUCC:LS_RET_FAIL);
    if (iRet)
    {
        ERROR("g_kafkaProducer produce_msg err, stKafkaMsg:%s, stClickMsg:%s, iRet:%d",
                klib::jce2str(stKafkaMsg).c_str(), klib::jce2str(stClickMsg).c_str(), iRet);
        return iRet;
    }
    DEBUG("g_kafkaProducer produce_msg succ, stKafkaMsg:%s, stClickMsg:%s",
            klib::jce2str(stKafkaMsg).c_str(), klib::jce2str(stClickMsg).c_str());

    return iRet;
}
