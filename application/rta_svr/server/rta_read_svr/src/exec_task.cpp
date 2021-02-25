#include "server_config.h"
#include "exec_task.h"
#include "klib.h"

#define g_conf (*CSingleton<CServerConf>::instance())

int CGetLoginTimeTask::Process()
{
    int iRet = 0;

    CTmemApiMt tmem_obj;
    iRet = tmem_obj.fnInitL5(atoll(g_conf["tmem.bid"].c_str()), atoll(g_conf["tmem.mid"].c_str()), atoll(g_conf["tmem.cid"].c_str()), atoll(g_conf["tmem.timeout_ms"].c_str()));
    if (iRet != 0)
    {
        ERROR("m_tmem_obj.fnInitL5 error, key:%s, iRet:%d", m_key.c_str(), iRet);
        return iRet;
    }

    string strKey = PROFILE_LOGIN_CKV_KEY + m_key;
	int iCas = -1;

	iRet = tmem_obj.Get(strKey, m_strLoginTime, iCas, 0, RTA_READ_SVR_TMEM_MOD_ID);
    if (iRet == -13106 || iRet == -13200)
	{
		DEBUG("iRet:%d, key:%s does not exist in profile login ckv", iRet, strKey.c_str());
	}
	else if (iRet)
	{
		ERROR("m_tmem_obj Get err, key:%s, iRet:%d, errMsg:%s", strKey.c_str(), iRet, tmem_obj.GetLastError());
	}
	else
	{
		DEBUG("m_tmem_obj Get succ, key:%s, strLoginTime:%s", strKey.c_str(), m_strLoginTime.c_str());
    }

    return iRet;
}

int CGetRTAInfoCacheTask::Process()
{
    int iRet = 0;

    CTmemApiMt tmem_obj;
    iRet = tmem_obj.fnInitL5(atoll(g_conf["tmem.bid"].c_str()), atoll(g_conf["tmem.mid"].c_str()), atoll(g_conf["tmem.cid"].c_str()), atoll(g_conf["tmem.timeout_ms"].c_str()));
    if (iRet != 0)
    {
        ERROR("m_tmem_obj.fnInitL5 error, key:%s, iRet:%d", m_key.c_str(), iRet);
        return iRet;
    }

    string strRTACacheKey = RTA_INFO_CACHE_KEY + m_key;
    string dataRTA;
    int casRTA = 0;

    iRet = tmem_obj.Get(strRTACacheKey, dataRTA, casRTA, 0, RTA_READ_SVR_TMEM_MOD_ID);
    if (iRet == -13106 || iRet == -13200)
	{
		DEBUG("iRet:%d, key:%s does not exist in rta cache ckv", iRet, strRTACacheKey.c_str());
        return 0;
	}
	else if (0 != iRet)
	{
		ERROR("m_tmem_obj Get err, key:%s, iRet:%d, errMsg:%s", strRTACacheKey.c_str(), iRet, tmem_obj.GetLastError());
        return iRet;
	}

    int iLen = dataRTA.length();
    iRet = m_stRTAInfoCache.Decode((uint8_t*)dataRTA.c_str(), &iLen, NULL);
    if (iRet != 0)
    {
        ERROR("RTAInfoCache.Decode error, iRet:%d, key:%s", iRet, strRTACacheKey.c_str());
    }
    else
    {
        DEBUG("m_tmem_obj Get succ, key:%s, stRTACacheInfo:%s", strRTACacheKey.c_str(), klib::jce2str(m_stRTAInfoCache).c_str());
    }

    return iRet;
}

