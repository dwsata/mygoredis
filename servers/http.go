package servers

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/opentechnologysel/mygoredis/caches"
	"io/ioutil"
	"net/http"
)

type HttpServer struct {
	cache *caches.Cache
}

func NewHttpServer(cache *caches.Cache) *HttpServer {
	return &HttpServer{cache: cache}
}

func (hs *HttpServer) Run(address string) error {
	return http.ListenAndServe(address, hs.routerHandler())
}

func (hs *HttpServer) routerHandler() http.Handler {
	router := httprouter.New()
	router.GET("/cache/:key", hs.getHandler)
	router.PUT("/cache/:key", hs.setHandler)
	router.DELETE("/cache/:key", hs.deleteHandler)
	router.GET("/status", hs.statusHandler)
	return router
}

func (hs *HttpServer) getHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	key := params.ByName("key")
	value, ok := hs.cache.Get(key)

	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Write(value)
}

func (hs *HttpServer) setHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	key := params.ByName("key")
	value, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	hs.cache.Set(key, value)
}

func (hs *HttpServer) deleteHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	key := params.ByName("key")

	hs.cache.Delete(key)
}

func (hs *HttpServer) statusHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	status, err := json.Marshal(map[string]interface{}{
		"count": hs.cache.Count(),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(status)
}
