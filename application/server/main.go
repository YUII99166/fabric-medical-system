package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"application/blockchain"
	"application/db"
	"application/pkg/cron"
	"application/routers"
)

func main() {
	timeLocal, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		log.Printf("时区设置失败 %s", err)
	}
	time.Local = timeLocal

	// 初始化数据库
	if err := db.Init(); err != nil {
		log.Printf("数据库初始化失败 %s", err)
		return
	}
	defer db.Close()

	blockchain.Init()
	go cron.Init()

	endPoint := fmt.Sprintf("0.0.0.0:%d", 8888)
	server := &http.Server{
		Addr:    endPoint,
		Handler: routers.InitRouter(),
	}
	log.Printf("[info] start http server listening %s", endPoint)
	if err := server.ListenAndServe(); err != nil {
		log.Printf("start http server failed %s", err)
	}
}
