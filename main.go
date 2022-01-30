package main

import (
	"net/http"

	"github.com/solenovex/web-tutor/router"
	"github.com/solenovex/web-tutor/service"
)

func main() {
	//路由设置
	router.MapRouter()
	//mysql设置
	service.DBInit()
	//redis配置
	service.RedisInit()
	//日志配置
	service.LogInit()
	//配置布隆过滤器
	service.BloomInit()
	http.ListenAndServe(service.GetConfiguration().Server.Host, nil)

}
