#ifndef __EXEC_TASK_H__
#define __EXEC_TASK_H__

#include "sppincl.h"
#include "rta_read_svr.h"


class CGetLoginTimeTask : public IMtTask
{
	public:
		CGetLoginTimeTask(string & key) {m_key = key;}
		~CGetLoginTimeTask() {};
		virtual int Process();

		//功能函数

	public:
		string m_key;
        string m_strLoginTime;
};

class CGetRTAInfoCacheTask : public IMtTask
{
	public:
		CGetRTAInfoCacheTask(string & key) {m_key = key;}
		~CGetRTAInfoCacheTask() {};
		virtual int Process();

		//功能函数

	public:
		string m_key;
        RTAInfoCache m_stRTAInfoCache;
};

#endif
