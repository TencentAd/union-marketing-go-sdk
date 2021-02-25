package main

import (
	"net/http"
	"context"
	"git.code.oa.com/going/going/config"
	"git.code.oa.com/going/going/cat/qzs"
)

func main() {
	config.Parse(&conf)
	qzsCtx := qzs.NewContext(context.Background())

    http.HandleFunc("/ams/gdt", handleAMSGdt);          // 处理ams http请求

	qzsCtx.Debug("Start http server ...")
	err := http.ListenAndServe(":5180", nil)
	if err != nil {
		qzsCtx.Error("init http", err)
	}
}

