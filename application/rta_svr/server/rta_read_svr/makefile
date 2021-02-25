#PATH_PRJ_ROOT是项目的根目录相对于当前目录的相对路径，必须保证正确
PATH_PRJ_ROOT:=$(if $(PATH_PRJ_ROOT), $(PATH_PRJ_ROOT), ../../../..)

#这里的规则是所有的文件一起生成一个可执行文件
#需要依赖的库和头文件请在这里定义
INC= -I include -I../comm/include -I $(SRC_SPP_SHARE) $(INC_JSONCPP) $(INC_KG_ENCRYPTION) $(INC_STRING_DEAL) $(INC_TBASE) $(INC_QZAPP_FRAME) $(SPP_INC) $(INC_BASECLASS) $(INC_PROTOCOL) $(INC_CMEM_JCE) $(INC_L5)
LIB= $(LIB_MOD_LOG) $(LIB_QZAPP_FRAME_MT) $(LIB_TBASE) $(LIB_LOGAPI) $(LIB_CMEM_JCE) -lpthread $(LIB_PROTOCOL) $(LIB_L5) $(LIB_BASECLASS)

#需要生成的可执行文件，如果是动态库，一定要用.so后缀名
TARGET=rta_read_svr.so

#如果不是在本目录下，则需要定义以下变量
SRC_DIR=./src/
OBJS_DIR=./src/
BIN_DIR=./bin/

include $(PATH_PRJ_ROOT)/makeinclude/makefile.bincomm
C_FLAGS += -DDEF_USE_API -DDEF_USE_MT
