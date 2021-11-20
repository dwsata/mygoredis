package main

import (
	"github.com/opentechnologysel/mygoredis/caches"
	"github.com/opentechnologysel/mygoredis/servers"
)

func main() {
	cache := caches.NewCache()
	httpServer := servers.NewHttpServer(cache)
	err := httpServer.Run(":5837")
	if err != nil {
		panic(err)
	}
}
