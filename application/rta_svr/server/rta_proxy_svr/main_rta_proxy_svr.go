package main

import (
	"context"
	"git.code.oa.com/going/going/cat/qzs"
	"git.code.oa.com/going/going/config"
	"net/http"
)

//=================配置文件===================
var conf = struct {
	Comm struct {
		HttpsCompanyList   []int
		MapCompanyBusiness map[string]Business
	}
}{}

type Business struct {
	SecretKey string
}

func main() {
	config.Parse(&conf)
	qzsCtx := qzs.NewContext(context.Background())

	http.HandleFunc("/ams/rta", handleAMSRta) // 处理AMS http请求

	qzsCtx.Debug("Start http server ...")
	err := http.ListenAndServe(":5174", nil)
	if err != nil {
		qzsCtx.Error("init http", err)
	}
}
