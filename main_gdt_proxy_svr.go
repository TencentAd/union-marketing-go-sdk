package main

import (
	"net/http"
	"context"
	"git.code.oa.com/going/going/config"
	"git.code.oa.com/going/going/cat/qzs"
)

//=================配置文件===================
var conf = struct {
    InterfaceID struct {
        InterfaceIDVivo     int
        InterfaceIDKShou    int
        InterfaceIDAMS      int
    }
    MarketingAPI struct {
        AppId               string
        Secret              string
        DebugMode           string
    }
}{}

func main() {
	config.Parse(&conf)
	qzsCtx := qzs.NewContext(context.Background())

	http.HandleFunc("/vivo/gdt", handleVivoGdt);	    // 处理vivo http请求
    http.HandleFunc("/kshou/gdt", handleKShouGdt);      // 处理快手 http请求
    http.HandleFunc("/ams/gdt", handleAMSGdt);          // 处理ams http请求
    http.HandleFunc("/toutiao/gdt", handleTouTiaoGdt);  // 处理头条 http请求

    // ------------------------------ 头条 Marketing API 相关接口 ------------------------------------- //
    http.HandleFunc("/toutiao/oauth2/callback/", handleTouTiaoOauth); // 头条获取auth_code回调接口

	qzsCtx.Debug("Start http server ...")
	err := http.ListenAndServe(":5180", nil)
	if err != nil {
		qzsCtx.Error("init http", err)
	}
}

