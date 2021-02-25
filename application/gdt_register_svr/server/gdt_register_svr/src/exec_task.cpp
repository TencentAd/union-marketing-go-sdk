#include "server_config.h"
#include "exec_task.h"
#include "klib.h"

#define g_conf (*CSingleton<CServerConf>::instance())

int CGetClickInfoTask::Process()
{
    int iRet = 0;

    CTmemApiMt tmem_obj;
    iRet = tmem_obj.fnInitL5(atoll(g_conf["tmem.bid"].c_str()), atoll(g_conf["tmem.mid"].c_str()), atoll(g_conf["tmem.cid"].c_str()), atoll(g_conf["tmem.timeout_ms"].c_str()));
    if (iRet != 0)
    {
        ERROR("m_tmem_obj.fnInitL5 error, strDeviceId:%s, iRet:%d", m_strDeviceId.c_str(), iRet);
        return iRet;
    }

    int iCas = -1;
    iRet = tmem_obj.Get(m_strDeviceId, m_strData, iCas, 0, GDT_REGISTER_TMEM_MOD_ID);
    return iRet;
}
