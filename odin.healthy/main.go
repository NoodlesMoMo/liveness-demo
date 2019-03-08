package main

import (
	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"odin.healthy/healthy"
)

func main() {

	health := healthy.NewHealthy()
	health.AddLiveness("go-routine-max", healthy.GoRoutingCheck(10))
	health.AddLiveness("filesystem-read-write", healthy.FileReadWriteCheck("/tmp/abc.txt"))

	health.AddLiveness("mysql-conn-detect", healthy.MysqlPingCheck("xxx"))
	health.AddLiveness("redis-conn-detect", healthy.RedisPingCheck("xxx"))
	health.AddLiveness("dns-reslove", healthy.DNSResolveCheck("xxx"))
	health.AddLiveness("memory-max", healthy.MemoryMaxCheck(1024*1024*2))

	router := routing.New()

	router.Get("/healthy", health.HealthyEndpoint)
	router.Get("/access", healthy.AccessHandler)
	router.Get("/debug/ok", healthy.DebugOKHandler)
	router.Get("/debug/error", healthy.DebugErrHandler)
	router.Get("/cat", healthy.CatFileHandler)

	fasthttp.ListenAndServe(":8080", router.HandleRequest)
}
