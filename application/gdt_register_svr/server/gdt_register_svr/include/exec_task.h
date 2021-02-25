#ifndef __EXEC_TASK_H__
#define __EXEC_TASK_H__

#include "sppincl.h"
#include "gdt_register.h"

class CGetClickInfoTask : public IMtTask
{
    public:
        CGetClickInfoTask(const string &strDeviceId) : m_strDeviceId(strDeviceId) {}
        ~CGetClickInfoTask() {}
        virtual int Process();

    public:
        string m_strDeviceId;
        string m_strData;
};

#endif
