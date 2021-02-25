#PATH_PRJ_ROOT是项目的根目录相对于当前目录的相对路径，必须保证正确
PATH_PRJ_ROOT:=$(if $(PATH_PRJ_ROOT), $(PATH_PRJ_ROOT), ../../..)

#如果自定义需要编译的目录，或者是需要保证各个目录的编译顺序，请定义MAKE_DIRS变量
#MAKE_DIRS=
include $(PATH_PRJ_ROOT)/makeinclude/makefile.dircomm
