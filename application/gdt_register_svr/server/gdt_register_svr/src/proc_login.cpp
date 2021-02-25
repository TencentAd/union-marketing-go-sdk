//create by create_spp_server.sh
/**
 * 名称: 全民K歌
 * (c) Copyright 2019 derickwu@tencent.com. All Rights Reserved.
 *
 * change log:
 * Author: derickwu
 *   Date: 2019-06-26
 *	 Desc:
 */
#include "server_config.h"
#include "klib.h"
#include "gdt_register.h"
#include "exec_task.h"
using namespace proto_user_profile;
using namespace proto_gdt_register;
using namespace ns_common_kafka;

#define g_conf (*CSingleton<CServerConf>::instance())

int FuncStoreResultForSecondKeep(char *data, size_t size, size_t nmemb, std::string *result)
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

bool CheckLoginCallback(const RegisterInfo& registerInfo, const unsigned uLoginTimeSec, unsigned int timeDiv)
{
    //不同渠道的时间格式&单位不一致，可以用ad_company_type来区分，要统一转换为秒
    //只有拉活渠道会调用这个函数，当拉活渠道的点击时间存储到了专属字段，则取专属字段的值
    //否则，取原来字段的值
    unsigned int uClickTimeSec = atoll(registerInfo.click_time.c_str());//广告点击时间
    if (registerInfo.ad_company_type_active != 0)
    {
        uClickTimeSec = atoll(registerInfo.click_time_active.c_str());
    }
    return (uLoginTimeSec > uClickTimeSec) &&  (uLoginTimeSec-uClickTimeSec <= timeDiv);
}

int GdtRegister::ProcLoginReport(proto_gdt_register::LoginReportReq &stReq, NS_API_BASE::RSP_BASE &stRsp)
{
    int iRet = 0;

    WNS::WnsQua stQua;
    WNS::WnsQuaHelper::instance()->ParseQua(stReq.strQua, stQua);
    map<string, string> mapDeviceInfo;
    ns_kg::CWebappCommon::ParseParams(stReq.strDeviceInfo, mapDeviceInfo);

    // uid 不合法直接返回
    if (stReq.uUid < 1000000)
    {
        DEBUG("req:%s invalid", klib::jce2str(stReq).c_str());
        return iRet;
    }

    // [step 1] 判断是否新用户
    bool bNewUser = false;
    iRet = JudgeNewUser(stReq.uUid, bNewUser);
    if (iRet)
    {
        ERROR("JudgeNewUser err, uid:%u, deviceinfo:%s, iRet:%d", stReq.uUid, stReq.strDeviceInfo.c_str(), iRet);
        return iRet;
    }
    DEBUG("JudgeNewUser succ, uid:%u, deviceinfo:%s, bNewUser:%d", stReq.uUid, stReq.strDeviceInfo.c_str(), (int)bNewUser);

    // [step 2] 渠道拉新撞库
    if (bNewUser)
    {
        iRet = HandleNewUser(stReq, mapDeviceInfo, stQua);
        if (iRet)
        {
            ERROR("HandleNewUser err, uid:%u, deviceinfo:%s, iRet:%d", stReq.uUid, stReq.strDeviceInfo.c_str(), iRet);
        }
        else
        {
            DEBUG("HandleNewUser succ, uid:%u, deviceinfo:%s", stReq.uUid, stReq.strDeviceInfo.c_str());
        }
    }

    // [step 3] 渠道拉活和次留撞库
    else
    {
        // [step 4] 处理次留模型
	    iRet = HandleSecondaryModel(stReq, mapDeviceInfo, stQua);
	    if (iRet)
	    {
		    ERROR("HandleSecondaryModel err, iRet:%d, req:%s", iRet, klib::jce2str(stReq).c_str());
		    return iRet;
	    }
	    DEBUG("HandleSecondaryModel succ, req:%s", klib::jce2str(stReq).c_str());

        iRet = HandleActiveUser(stReq, mapDeviceInfo, stQua);
        if (iRet)
        {
            ERROR("HandleActiveUser err, uid:%u, deviceinfo:%s, iRet:%d", stReq.uUid, stReq.strDeviceInfo.c_str(), iRet);
        }
        else
        {
            DEBUG("HandleActiveUser succ, uid:%u, deviceinfo:%s", stReq.uUid, stReq.strDeviceInfo.c_str());
        }
    }

    return iRet;
}

int GdtRegister::GetAllDeviceId(map<string, string> &mapDeviceInfo, const WNS::WnsQua &stQua, proto_gdt_register::LoginReportReq &stReq,
        string &strMd5Imei, string &strMd5Oaid, string &strMd5AndroidId, string &strMd5Idfa, string &strMd5IpUa)
{
    // 安卓设备，获取 imei, oaid, androidid
    // deviceinfo中的imei是原值，oaid是md5之后的值，androidid是原值
    // 注意要与渠道点击信息的存储key格式保持一致
    if (stQua.isAndroid())
    {
        string strImei = mapDeviceInfo["i"];
        if (strImei != "" && strImei != "n/a" && strImei != "N/A" && strImei != "000000000000000")
        {
            std::transform(strImei.begin(), strImei.end(), strImei.begin(), ::tolower);
            strMd5Imei = md5Lower(strImei);
        }
        if (mapDeviceInfo["oaid"] != "")
        {
            strMd5Oaid = "ext_oaid_" + mapDeviceInfo["oaid"];
        }
        if (mapDeviceInfo["android"] != "")
        {
            strMd5AndroidId = "ext_androidid_" + md5Lower(mapDeviceInfo["android"]);
        }
        strMd5IpUa = GenerateKeyIpUa(stReq.uUid, stReq.strUserIp, mapDeviceInfo, stQua);

        // 如果安卓设备既拿不到imei又拿不到oaid以及androidid，则返回失败
        if (strMd5Imei == "" && strMd5Oaid == "" && strMd5AndroidId == "" && strMd5IpUa == "")
        {
            return MUSIC_CODE_GDT_REGISTER_EMPTY_DEVICE_INFO;
        }
    }
    // ios设备，获取 idfa原值
    // deviceinfo中的idfa是原值
    else
    {
        if (mapDeviceInfo["idfa"] == "")
        {
            // 如果deviceinfo中不存在idfa，则使用req.strIdfa覆盖
            mapDeviceInfo["idfa"] = stReq.strIdfa;
        }

        string strIdfa = mapDeviceInfo["idfa"];
        if (strIdfa != "")
        {
            std::transform(strIdfa.begin(), strIdfa.end(), strIdfa.begin(), ::toupper);
            strMd5Idfa = md5Lower(strIdfa);
        }
        strMd5IpUa = GenerateKeyIpUa(stReq.uUid, stReq.strUserIp, mapDeviceInfo, stQua);

        // 如果ios设备既拿不到idfa又拿不到ipua，则返回失败
        if (strMd5Idfa == "" && strMd5IpUa == "")
        {
            return MUSIC_CODE_GDT_REGISTER_EMPTY_DEVICE_INFO;
        }
    }

    return 0;
}

string GdtRegister::GenerateKeyIpUa(const unsigned uUid, const string& strClientIP, map<string, string>& mapDeviceInfo, const WNS::WnsQua& qua)
{
    DEBUG("uid:%u generate muid by ip&ua: (ip:%s o:%s os:%s)", uUid, strClientIP.c_str(), mapDeviceInfo["o"].c_str(), mapDeviceInfo["os"].c_str());
    if (strClientIP == "")
    {
        return "";
    }
    string strKey = "";
    if (qua.isAndroid())
    {
        string strOsVer = mapDeviceInfo["o"];
        if (strOsVer.size() == 0)
        {
            return "";
        }
        strKey = strClientIP + "android" + strOsVer;
    }
    else
    {
        string strOsVer = mapDeviceInfo["os"];
        vector<string> arrOsVer;
        ns_kg::CWebappCommon::SpliteStr(strOsVer, arrOsVer, '/');
        if (arrOsVer.size() != 2 || arrOsVer[1].size() == 0)
        {
            DEBUG("uid:%u, ios osver error device=%s", uUid, strOsVer.c_str());
            return "";
        }
        strKey = strClientIP + "iphone" + arrOsVer[1];
    }
    DEBUG("uid:%u, ipua key=%s", uUid, strKey.c_str());
    char res[33] = {0};//如果是32会出core
    NS_KG_MD5::md5_string(strKey.c_str(), strKey.length(), res);
    string strMD5Key = string(res);
    std::transform(strMD5Key.begin(), strMD5Key.end(), strMD5Key.begin(), ::tolower);
    return strMD5Key;
}

int GdtRegister::FindClickInfo(map<string, string> &mapDeviceInfo, const WNS::WnsQua &stQua, LoginReportReq &stReq,
        vector<string> &vctDeviceId, vector<string> &vctData)
{
    int iRet = 0;

    // [step 1] 对于安卓设备，imei>oaid>androidid，逐层匹配
    // 对于ios设备，idfa>ip+ua，逐层匹配，先直接使用上报的idfa字段覆盖掉mapDeviceInfo中的idfa字段
    // 全部是小写之后的md5值
    string strMd5Imei, strMd5Oaid, strMd5AndroidId, strMd5Idfa, strMd5IpUa;
    int iReason = GetAllDeviceId(mapDeviceInfo, stQua, stReq, strMd5Imei, strMd5Oaid, strMd5AndroidId, strMd5Idfa, strMd5IpUa);
    if (iReason)
    {
        DEBUG("GetAllDeviceId err, uid:%u, deviceinfo:%s, iReason:%d", stReq.uUid, stReq.strDeviceInfo.c_str(), iReason);
        return iReason;
    }
    DEBUG("GetAllDeviceId succ, uid:%u, deviceinfo:%s, imei_md5:%s, oaid_md5:%s, androidid_md5:%s, idfa_md5:%s, ipua_md5:%s",
            stReq.uUid, stReq.strDeviceInfo.c_str(), strMd5Imei.c_str(), strMd5Oaid.c_str(), strMd5AndroidId.c_str(),
            strMd5Idfa.c_str(), strMd5IpUa.c_str());

    // [step 2] 安卓设备逐层匹配 imei>oaid>androidid>ip+ua
    // ios设备逐层匹配，idfa>ip+ua
    IMtTaskList stTaskList;
    if (stQua.isAndroid())
    {
        AddTaskListAndroid(strMd5Imei, strMd5Oaid, strMd5AndroidId, strMd5IpUa, stTaskList);
    }
    else
    {
        AddTaskListIos(strMd5Idfa, strMd5IpUa, stTaskList);
    }

    // [step 3] task模式执行查询点击信息并获取结果
    mt_exec_all_task(stTaskList);
    for (unsigned i = 0; i < stTaskList.size(); ++i)
    {
        CGetClickInfoTask *pGetClickInfoTask = dynamic_cast<CGetClickInfoTask*>(stTaskList[i]);
        if (pGetClickInfoTask->GetResult() == 0)
        {
            vctDeviceId.push_back(pGetClickInfoTask->m_strDeviceId);
            vctData.push_back(pGetClickInfoTask->m_strData);
        }
    }

    // [step 4] 释放task
    for (unsigned i = 0; i < stTaskList.size(); ++i)
    {
        delete stTaskList[i];
    }

    if (vctDeviceId.size() == 0 || vctData.size() == 0 || vctDeviceId.size() != vctData.size())
    {
        return MUSIC_CODE_GDT_REGISTER_EMPTY_IP_UA_KEY;
    }
    return 0;
}

int GdtRegister::FindValidRegisterInfo(map<string, string> &mapDeviceInfo, const WNS::WnsQua &stQua, LoginReportReq &stReq, 
    vector<string>& vctDeviceId, vector<string>& vctData, string& strDeviceId, RegisterInfo& stRegisterInfo)
{
    int iRet = 0;
    int iReason = 0;
    for(unsigned int i = 0;i < vctDeviceId.size(); ++i)
    {
        RegisterInfo tmpRegisterInfo;
        int iLen = vctData[i].length();
        iRet = tmpRegisterInfo.Decode((uint8_t*)vctData[i].c_str(), &iLen, NULL);
        if (iRet != 0)
        {
            ERROR("tmpRegisterInfo.Decode error, iRet:%d, uid:%u, deviceinfo:%s, deviceid:%s",
                    iRet, stReq.uUid, stReq.strDeviceInfo.c_str(), vctDeviceId[i].c_str());
            continue;
        }
        if(g_conf["ipua_enable." + tmpRegisterInfo.ad_company_type] == "1" && g_conf["invalid_deviceid." + vctDeviceId[i]] == "1")
        {
            DEBUG("uid:%u deviceid:%s in black list, filtered.", stReq.uUid, vctDeviceId[i].c_str());
            continue;
        }
        strDeviceId = vctDeviceId[i];
        stRegisterInfo = tmpRegisterInfo;
        return 0;
    }
    return -13200;
}

void GdtRegister::AddTaskListAndroid(const string &strMd5Imei, const string &strMd5Oaid, const string &strMd5AndroidId, const string &strIpUa, IMtTaskList &stTaskList)
{
    DoAddTaskList(strMd5Imei, stTaskList);
    DoAddTaskList(strMd5Oaid, stTaskList);
    DoAddTaskList(strMd5AndroidId, stTaskList);
    DoAddTaskList(strIpUa, stTaskList);
}

void GdtRegister::AddTaskListIos(const string &strMd5Idfa, const string &strMd5IpUa, IMtTaskList &stTaskList)
{
    DoAddTaskList(strMd5Idfa, stTaskList);
    DoAddTaskList(strMd5IpUa, stTaskList);
}

void GdtRegister::DoAddTaskList(const string &strDeviceId, IMtTaskList &stTaskList)
{
    if (strDeviceId != "")
    {
        CGetClickInfoTask *pGetClickInfoTask = new CGetClickInfoTask(strDeviceId);
        stTaskList.push_back(pGetClickInfoTask);
    }
}

int GdtRegister::HandleNewUser(LoginReportReq &stReq, map<string, string> &mapDeviceInfo, const WNS::WnsQua &stQua)
{
    int iRet = 0;

    // [step 1] 使用设备id尝试从ckv里取出渠道点击信息
    // 逐层设备类型进行匹配，最终匹配上的key为strDeviceId，对应的点击信息为strData
    // 返回的错误码即为reason
    vector<string> vctDeviceId;
    vector<string> vctData;
    iRet = FindClickInfo(mapDeviceInfo, stQua, stReq, vctDeviceId, vctData);
    if (iRet)
    {
        DEBUG("FindClickInfo err, uid:%u, deviceinfo:%s, iRet:%d", stReq.uUid, stReq.strDeviceInfo.c_str(), iRet);
        return 0;
    }
    DEBUG("FindClickInfo succ, uid:%u, deviceinfo:%s, deviceid size:%u", stReq.uUid, stReq.strDeviceInfo.c_str(), vctDeviceId.size());

    // [step 2] 从撞库成功的信息中，提取有效的
    int iReason = 0;
    string strDeviceId = "";
    RegisterInfo register_info;
    iRet = FindValidRegisterInfo(mapDeviceInfo, stQua, stReq, vctDeviceId, vctData, strDeviceId, register_info);
    if(iRet != 0 || strDeviceId == "" || strDeviceId.size() == 0)
    {
        DEBUG("uid:%u hasn`t any valid device id, req:%s.", stReq.uUid, klib::jce2str(stReq).c_str());
        return 0;
    }
    DEBUG("uid:%u, deviceinfo:%s, deviceid:%s, register_info:%s",
        stReq.uUid, stReq.strDeviceInfo.c_str(), strDeviceId.c_str(), klib::jce2str(register_info).c_str());

    // [step 3] 将DAU严格口径登录流水写入kafka
    iRet = SetLoginTime(strDeviceId, stReq.uLoginTime);
    if (iRet)
    {
        ERROR("SetLoginTime err, uid:%u, deviceinfo:%s, deviceid:%s, iRet:%d",
                stReq.uUid, stReq.strDeviceInfo.c_str(), strDeviceId.c_str(), iRet);
    }
    else
    {
        DEBUG("SetLoginTime succ, uid:%u, deviceinfo:%s, deviceid:%s",
                stReq.uUid, stReq.strDeviceInfo.c_str(), strDeviceId.c_str());
    }

    // [step 4] 一个uid只回调和上报一次
    bool bReported = false;
    int iCas = -1;
    iRet = ReportBefore(stReq.uUid, bReported, iCas);
    if (iRet)
    {
        ERROR("ReportBefore err, uid:%u, deviceinfo:%s, deviceid:%s, iRet:%d",
                stReq.uUid, stReq.strDeviceInfo.c_str(), strDeviceId.c_str(), iRet);
    }
    else if (bReported)
    {
        DEBUG("ReportBefore succ, uid:%u, deviceinfo:%s, deviceid:%s, bReported=true, no more report",
                stReq.uUid, stReq.strDeviceInfo.c_str(), strDeviceId.c_str());
        return iRet;
    }
    else
    {
        DEBUG("ReportBefore succ, uid:%u, deviceinfo:%s, deviceid:%s, bReported=false, need to report",
                stReq.uUid, stReq.strDeviceInfo.c_str(), strDeviceId.c_str());
    }

    // [step 5] key验证存在后，根据不同的广告商类型通知对应的广告商，告知广告商该用户已转化成功
    // 如果是中台AMS传过来的广点通点击信息or七麦科技渠道不做注册撞库，仅仅需要做上报
    if (register_info.ad_company_type == AD_COMPANY_GDT)
    {
        iRet = SetReportRecord(stReq.uUid, iCas);
        if (iRet == -13104)
        {
            DEBUG("SetReportRecord err, cas err, uid:%u, deviceinfo:%s, deviceid:%s, has been reported",
                    stReq.uUid, stReq.strDeviceInfo.c_str(), strDeviceId.c_str());
            return 0;
        }

        iRet = SendCallBackToGDT(register_info, stReq.strQua, stReq.strDeviceInfo, stQua);
    }
    else if (register_info.ad_company_type == AD_COMPANY_AMS)
    {
        iRet = SetReportRecord(stReq.uUid, iCas);
        if (iRet == -13104)
        {
            DEBUG("SetReportRecord err, cas err, uid:%u, deviceinfo:%s, deviceid:%s, has been reported",
                    stReq.uUid, stReq.strDeviceInfo.c_str(), strDeviceId.c_str());
            return 0;
        }

        string strType = "ACTIVATE_APP";
        iRet = SendCallBackToAMSSelf(register_info, stQua, stReq.uLoginTime, strType);
    }
    else
    {
        ERROR("ad_company_type_active:%u ad_company_type:%u invalid, uid:%u, deviceinfo:%s, deviceid:%s",
                register_info.ad_company_type_active, register_info.ad_company_type,
                stReq.uUid, stReq.strDeviceInfo.c_str(), strDeviceId.c_str());
        return iRet;
    }

    if (iRet != 0)
    {
        ERROR("send callback error, iRet:%d, uid:%u, deviceinfo:%s, strDeviceId:%s, register_info:%s",
                iRet, stReq.uUid, stReq.strDeviceInfo.c_str(), strDeviceId.c_str(), klib::jce2str(register_info).c_str());
        return iRet;
    }

    // [step 7] 回调成功，回填激活信息到ckv中(次留模型需要)
    if (register_info.ad_company_type == AD_COMPANY_AMS)
    {
        register_info.uUid = stReq.uUid;
        register_info.status_bit |= AD_STATUS_REGISTED;
        register_info.register_time = stReq.uLoginTime;

        // 用于次留模型匹配，使用uid作为key，另外存一份
        EncodeAndSetSecondary(register_info);
    }

    return iRet;
}

int GdtRegister::ReportBefore(const unsigned uUid, bool &bReported, int &iCas)
{
    int iRet = 0;

    string strKey = REPORT_RECORD_KEY_PREFIX + klib::ToStr(uUid);
    string strReportInfo;
    proto_gdt_register::ReportInfo stReportInfo;

    iRet = m_tmem_obj.Get(strKey, strReportInfo, iCas, 0, GDT_REGISTER_TMEM_MOD_ID);
    if (iRet == -13106 || iRet == -13200)
    {
        bReported = false;
        DEBUG("iRet:%d, key:%s does not exist in report info ckv, bReported:%d", iRet, strKey.c_str(), (int)bReported);
        return 0;
    }
    else if (0 != iRet)
    {
        ERROR("m_tmem_obj Get err, key:%s, iRet:%d, errMsg:%s", strKey.c_str(), iRet, m_tmem_obj.GetLastError());
        return iRet;
    }

    bReported = true;
    DEBUG("key:%s exists in report info ckv, bReported:%d", strKey.c_str(), (int)bReported);
    return iRet;
}

int GdtRegister::SetReportRecord(const unsigned uUid, int iCas)
{
    int iRet = 0;

    string strKey = REPORT_RECORD_KEY_PREFIX + klib::ToStr(uUid);
    string strReportInfo;

    proto_gdt_register::ReportInfo stReportInfo;
    stReportInfo.bReported = true;
    int iLen = m_bufsize;
    iRet = stReportInfo.Encode((uint8_t*)m_buf, &iLen, NULL);
    if (iRet != 0)
    {
        ERROR("stReportInfo.Encode error, key:%s, iRet:%d", strKey.c_str(), iRet);
        return iRet;
    }
    strReportInfo.assign(m_buf, iLen);

    unsigned int expire = (unsigned)atoll(g_conf["report_info.expire_time_sec"].c_str());
    iRet = m_tmem_obj.Set(strKey, strReportInfo, iCas, 0, GDT_REGISTER_TMEM_MOD_ID, expire);
    DEBUG("m_tmem_obj.Set succ key:%s report_info expire_time_sec:%u", strKey.c_str(), expire);
    if (iRet != 0)
    {
        ERROR("m_tmem_obj.Set error, iRet:%d, key:%s, error:%s", iRet, strKey.c_str(), m_tmem_obj.GetLastError());
        return iRet;
    }
    DEBUG("m_tmem_obj.Set succ, key:%s, stReportInfo:%s", strKey.c_str(), klib::jce2str(stReportInfo).c_str());
    return iRet;
}

int GdtRegister::JudgeNewUser(const unsigned uUid, bool &bNewUser)
{
    int iRet = 0;
    return iRet;
}

int GdtRegister::HandleActiveUser(LoginReportReq &stReq, map<string, string> &mapDeviceInfo, const WNS::WnsQua &stQua)
{
    int iRet = 0;

    // [step 1] 使用设备id尝试从ckv里取出渠道点击信息
    // 逐层设备类型进行匹配，最终匹配上的key为strDeviceId，对应的点击信息为strData
    // 返回的错误码即为reason
    vector<string> vctDeviceId;
    vector<string> vctData;
    iRet = FindClickInfo(mapDeviceInfo, stQua, stReq, vctDeviceId, vctData);
    if (iRet)
    {
        DEBUG("FindClickInfo err, uid:%u, deviceinfo:%s, iRet:%d", stReq.uUid, stReq.strDeviceInfo.c_str(), iRet);
        return 0;
    }
    DEBUG("FindClickInfo succ, uid:%u, deviceinfo:%s, deviceid size:%u", stReq.uUid, stReq.strDeviceInfo.c_str(), vctDeviceId.size());

    // [step 2] 过滤无效的撞库信息,如idfa中全0的设备
    int iReason = 0;
    string strDeviceId = "";
    RegisterInfo register_info;
    iRet = FindValidRegisterInfo(mapDeviceInfo, stQua, stReq, vctDeviceId, vctData, strDeviceId, register_info);
    if(iRet != 0 || strDeviceId == "" || strDeviceId.size() == 0)
    {
        DEBUG("uid:%u hasn`t any valid device id, req:%s.", stReq.uUid, klib::jce2str(stReq).c_str());
        return 0;
    }
    DEBUG("uid:%u, deviceinfo:%s, deviceid:%s, register_info:%s",
        stReq.uUid, stReq.strDeviceInfo.c_str(), strDeviceId.c_str(), klib::jce2str(register_info).c_str());

    // [step 3] 判断该设备当天是否首登
    bool bTodayFirstLogIn = false;
    iRet = JugdeFirstLoginToday(strDeviceId, bTodayFirstLogIn);
    if (iRet)
    {
        ERROR("JugdeFirstLoginToday err, uid:%u, deviceinfo:%s, deviceid:%s, iRet:%d, bTodayFirstLogIn:%d",
                stReq.uUid, stReq.strDeviceInfo.c_str(), strDeviceId.c_str(), iRet, (int)bTodayFirstLogIn);
    }
    else
    {
        DEBUG("JugdeFirstLoginToday succ, uid:%u, deviceinfo:%s, deviceid:%s, bTodayFirstLogIn:%d",
                stReq.uUid, stReq.strDeviceInfo.c_str(), strDeviceId.c_str(), (int)bTodayFirstLogIn);
    }

    // [step 4] 将DAU严格口径登录流水写入kafka
    iRet = SetLoginTime(strDeviceId, stReq.uLoginTime);
    if (iRet)
    {
        ERROR("SetLoginTime err, uid:%u, deviceinfo:%s, deviceid:%s, iRet:%d",
                stReq.uUid, stReq.strDeviceInfo.c_str(), strDeviceId.c_str(), iRet);
    }
    else
    {
        DEBUG("SetLoginTime succ, uid:%u, deviceinfo:%s, deviceid:%s",
                stReq.uUid, stReq.strDeviceInfo.c_str(), strDeviceId.c_str());
    }

    // [step 5] 一个uid只回调和上报一次
    bool bReported = false;
    int iCas = -1;
    iRet = ReportBefore(stReq.uUid, bReported, iCas);
    if (iRet)
    {
        ERROR("ReportBefore err, uid:%u, deviceinfo:%s, deviceid:%s, iRet:%d",
                stReq.uUid, stReq.strDeviceInfo.c_str(), strDeviceId.c_str(), iRet);
    }
    else if (bReported)
    {
        DEBUG("ReportBefore succ, uid:%u, deviceinfo:%s, deviceid:%s, bReported=true, no more report",
                stReq.uUid, stReq.strDeviceInfo.c_str(), strDeviceId.c_str());
        return iRet;
    }
    else
    {
        DEBUG("ReportBefore succ, uid:%u, deviceinfo:%s, deviceid:%s, bReported=false, need to report",
                stReq.uUid, stReq.strDeviceInfo.c_str(), strDeviceId.c_str());
    }

    // [step 6] 特殊逻辑：补充AMS拉新回调（为了让模型更精确，这里对于AMS，在用户登录的时候也做拉新回调，因为AMS在注册的时候，拉新回调过少）
    else if (atoi(g_conf["ams_returner.enable"].c_str()) == 1 && (register_info.ad_company_type == AD_COMPANY_GDT || register_info.ad_company_type == AD_COMPANY_AMS))
    {
        // 判断用户是否是回流
        bool bReturner = false;
        iRet = JudgeReturner(stReq.uUid, strDeviceId, register_info.strAdAccountId, bReturner);
        if (0 == iRet && bReturner)
        {
            iRet = SetReportRecord(stReq.uUid, iCas);
            if (iRet == -13104)
            {
                DEBUG("SetReportRecord err, cas err, uid:%u, deviceinfo:%s, deviceid:%s, has been reported",
                        stReq.uUid, stReq.strDeviceInfo.c_str(), strDeviceId.c_str());
                return 0;
            }

            DEBUG("JudgeReturner succ, uid=%u deviceid=%s deviceinfo=%s is a returner, do SendCallBackToGDT",
                    stReq.uUid, strDeviceId.c_str(), stReq.strDeviceInfo.c_str());

            // 自己接的AMS拉新，走新链路回调
            if (register_info.ad_company_type == AD_COMPANY_AMS)
            {
                string strType = "ACTIVATE_APP";
                iRet = SendCallBackToAMSSelf(register_info, stQua, stReq.uLoginTime, strType);
            }
            // 中台接的AMS
            else
            {
                iRet = SendCallBackToGDT(register_info, stReq.strQua, stReq.strDeviceInfo, stQua);
            }

        }
        else if (iRet)
        {
            ERROR("uid=%u deviceid=%s deviceinfo=%s JudgeReturner err, iRet:%d",
                    stReq.uUid, strDeviceId.c_str(), stReq.strDeviceInfo.c_str(), iRet);
        }
        else
        {
            DEBUG("uid=%u deviceid=%s deviceinfo=%s JudgeReturner succ, not a returner",
                    stReq.uUid, strDeviceId.c_str(), stReq.strDeviceInfo.c_str());
        }
    }
    else if (register_info.ad_company_type_active == AD_COMPANY_AMS_AC)
    {
        // 当日首登且点击操作发生5min内登陆才算拉活成功
        if (bTodayFirstLogIn && CheckLoginCallback(register_info, stReq.uLoginTime, 300))
        {
            iRet = SetReportRecord(stReq.uUid, iCas);
            if (iRet == -13104)
            {
                DEBUG("SetReportRecord err, cas err, uid:%u, deviceinfo:%s, deviceid:%s, has been reported",
                        stReq.uUid, stReq.strDeviceInfo.c_str(), strDeviceId.c_str());
                return 0;
            }

            iRet = ActiveSendCallBackToGDT(register_info, stQua);
            if (iRet != 0)
            {
                ERROR("send active call back to gdt error, uid:%u, deviceinfo:%s, deviceid:%s, iRet:%d, registerInfo:%s",
                        stReq.uUid, stReq.strDeviceInfo.c_str(), strDeviceId.c_str(), iRet, klib::jce2str(register_info).c_str());
            }
            else
            {
                DEBUG("send active call back to gdt succ, uid:%u, deviceinfo:%s, deviceid:%s, registerInfo:%s",
                    stReq.uUid, stReq.strDeviceInfo.c_str(), strDeviceId.c_str(), klib::jce2str(register_info).c_str());
            }
        }
        else
        {
            if (!bTodayFirstLogIn)
            {
                iReason = MUSIC_CODE_GDT_REGISTER_NOT_FIRST_LOGIN_TODAY;
            }
            else
            {
                iReason = MUSIC_CODE_GDT_REGISTER_NOT_LOGIN_AFTER_CLICK;
            }
            DEBUG("uid:%u, deviceinfo:%s, deviceid:%s login time:%u not in 5min after click time:%s or bTodayFirstLogIn=%d",
                    stReq.uUid, stReq.strDeviceInfo.c_str(), strDeviceId.c_str(), stReq.uLoginTime,
                    (register_info.click_time_active == "" ? register_info.click_time.c_str() : register_info.click_time_active.c_str()),
                    (int)bTodayFirstLogIn);
        }
    }

    return iRet;
}

int GdtRegister::JudgeReturner(const unsigned uUid, const string &strDeviceId,
        const string &strAdAccountId, bool &bReturner)
{
    int iRet = 0;
    return iRet;
}

int GdtRegister::JugdeFirstLoginToday(const string &strDeviceId, bool &bTodayFirstLogIn)
{
    int iRet = 0;
	string strLoginKey = PROFILE_LOGIN_CKV_KEY + strDeviceId;
	int iCas = -1;
	string strLoginTime;

    iRet = m_tmem_obj.Get(strLoginKey, strLoginTime, iCas, 0, GDT_REGISTER_TMEM_MOD_ID);
	if (iRet == -13200 || iRet == -13106)
	{
		DEBUG("no login info for key:%s", strLoginKey.c_str());
        bTodayFirstLogIn = true;
	}
	else if (iRet == 0)
	{
		if (IsSameDay((unsigned)atoll(strLoginTime.c_str()), m_uNowTimestamp))
		{
			DEBUG("login multiple times today for key:%s", strLoginKey.c_str());
            bTodayFirstLogIn = false;
		}
		else
		{
			DEBUG("login time:%s not today, update it, key:%s", strLoginTime.c_str(), strLoginKey.c_str());
            bTodayFirstLogIn = true;
		}
	}
	else
	{
		ERROR("m_tmem_obj Get err, iRet:%d, errMsg:%s", iRet, m_tmem_obj.GetLastError());
		return iRet;
	}
	DEBUG("m_tmem_obj Get succ, strLoginTime:%s, bTodayFirstLogIn:%d", strLoginTime.c_str(), (int)bTodayFirstLogIn);

    return 0;
}

int GdtRegister::SetLoginTime(const string &strDeviceId, const unsigned uLoginTime)
{
    int iRet = 0;

    // [step 1] 生成DauLoginMsg
    DauLoginMsg stDauLoginMsg;
    stDauLoginMsg.strDeviceId = strDeviceId;
    stDauLoginMsg.uLoginTime = uLoginTime;

    char szBuff[102400];
    int iLen = sizeof(szBuff);
    iRet = stDauLoginMsg.Encode((uint8_t*)szBuff, &iLen, NULL);
    if (iRet)
    {
        ERROR("stDauLoginMsg encode err, stDauLoginMsg:%s, iRet:%d", klib::jce2str(stDauLoginMsg).c_str(), iRet);
        return iRet;
    }
    DEBUG("stDauLoginMsg encode succ, stDauLoginMsg:%s", klib::jce2str(stDauLoginMsg).c_str());

    // [step 2] 写入kafka
    string strTmpMsg(szBuff, iLen);
    KafkaMsg stKafkaMsg;
    stKafkaMsg.uMsgType = DAU_LOGIN_MSG_TYPE;
    stKafkaMsg.strMsg = taf::TC_Base64::encode(strTmpMsg);

    struct timeval tvLog;
    API_Log_StartTimer(tvLog);
    iRet = g_kafkaProducer.produce_msg(stKafkaMsg.writeToJsonString());
    API_Log_StopAndWrite(tvLog, LM_ERROR, 0, GDT_REGISTER_MOD_ID, 0, LS_OP_INSERT, "", 0, iRet, iRet==0?LS_RET_SUCC:LS_RET_FAIL);
    if (iRet)
    {
        ERROR("g_kafkaProducer produce_msg err, stDauLoginMsg:%s, iRet:%d", klib::jce2str(stDauLoginMsg).c_str(), iRet);
        return iRet;
    }
    DEBUG("g_kafkaProducer produce_msg succ, stDauLoginMsg:%s", klib::jce2str(stDauLoginMsg).c_str());

    return iRet;
}

int GdtRegister::HandleSecondaryModel(LoginReportReq &stReq, map<string, string> &mapDeviceInfo, const WNS::WnsQua &stQua)
{
	int iRet = 0;

    // [step 1] 使用uid从ckv取出对应的注册流水
    string strData;
    string strKey = SECONDARY_INFO_KEY + klib::ToStr(stReq.uUid);
    int iCas = -1;
    RegisterInfo register_info;
    int iReason = 0;
    iRet = m_tmem_obj.Get(strKey, strData, iCas, 0, GDT_REGISTER_TMEM_MOD_ID);
    if (iRet == -13200 || iRet == -13106)
    {
        DEBUG("no key for req:%s, iRet:%d", klib::jce2str(stReq).c_str(), iRet);
        return 0;
    }
    else if (iRet != 0)
    {
        ERROR("Get SECONDARY_INFO error, iRet:%d, errMsg:%s, req:%s", iRet, m_tmem_obj.GetLastError(), klib::jce2str(stReq).c_str());
        return iRet;
    }

    // [step 2] 解析jce数据
    int iLen = strData.length();
    iRet = register_info.Decode((uint8_t*)strData.c_str(), &iLen, NULL);
    if (iRet != 0)
    {
        ERROR("decode error, uid:%u, iRet:%d, deviceinfo:%s", stReq.uUid, iRet, stReq.strDeviceInfo.c_str());
        return iRet;
    }

    // [step 3] 已经发生过注册，且没有发生过次留回调的用
    DEBUG("uid:%u, deviceinfo:%s, register_info:%s", stReq.uUid, stReq.strDeviceInfo.c_str(), klib::jce2str(register_info).c_str());
    if ((register_info.ad_company_type == AD_COMPANY_AMS) &&
        (register_info.status_bit & AD_STATUS_REGISTED) == 1 &&
        (register_info.status_bit & AD_STATUS_SECONDDAY) == 0)
    {
        // 归零时间戳: 得到注册时间的北京时间0点的时间戳(如果用户是7月8日注册，那么次留时间为7月9日00:00:00~23:59:59)
        unsigned int unNowTime = stReq.uLoginTime;
        unsigned int unRegisterTime = register_info.register_time - ((register_info.register_time + 8*3600) % 86400);
        DEBUG("uid:%u, unNowTime=%u unRegisterTime=%u", stReq.uUid, unNowTime, unRegisterTime);
        if (unNowTime > unRegisterTime &&
            unNowTime - unRegisterTime > 86400 &&
            unNowTime - unRegisterTime < 172800)
        {
            ERROR("uid:%u, time satisfy unNowTime=%u, unRegisterTime=%u", stReq.uUid, unNowTime, register_info.register_time);
            register_info.status_bit |= AD_STATUS_SECONDDAY;
            EncodeAndSetSecondary(register_info);
            if (register_info.ad_company_type == AD_COMPANY_AMS)
            {
                string strType = "START_APP";
                iRet = SendCallBackToAMSSelf(register_info, stQua, stReq.uLoginTime, strType);
            }
        }
        else
        {
            iReason = MUSIC_CODE_GDT_REGISTER_INVALID_LOGIN_TIME;
            DEBUG("time not satisfy, unNowTime= %u, unRegisterTime=%u, uid=%u, deviceinfo=%s",
                                unNowTime, register_info.register_time, stReq.uUid, stReq.strDeviceInfo.c_str());
        }
    }
    else
    {
        if (register_info.ad_company_type != AD_COMPANY_AMS)
        {
            iReason = MUSIC_CODE_GDT_REGISTER_INVALID_COMPANY_ID;
        }
        else
        {
            iReason = MUSIC_CODE_GDT_REGISTER_INVALID_STATUS_BIT;
        }
    }

	return iRet;
}

