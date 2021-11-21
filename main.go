package main

import (
	"flag"

	"github.com/opentechnologysel/mygoredis/caches"
	"github.com/opentechnologysel/mygoredis/servers"
)

func main() {
	address := flag.String("address", ":5837", "the address used to listen,such as 127.0.0.1:5837.")
	options := caches.DefaultOptions()
	flag.Int64Var(&options.MaxEntrySize, "maxEntrySize", options.MaxEntrySize, "The max memory size that entries can use. The unit is GB.")
	flag.IntVar(&options.MaxGCCount, "maxGcCount", options.MaxGCCount, "The max count of entries that gc will clean.")
	flag.Int64Var(&options.GCDuration, "gcDuration", options.GCDuration, "The duration between two gc tasks. The unit is Minute.")

	flag.Parse()
	cache := caches.NewCacheWith(*options)
	cache.AutoGc()
	httpServer := servers.NewHttpServer(cache)
	err := httpServer.Run(*address)
	if err != nil {
		panic(err)
	}
}
