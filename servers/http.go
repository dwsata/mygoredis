package servers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/opentechnologyself/mygoredis/caches"
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

func wrapUriWithVersion(uri string) string {
	return path.Join("/", APIVersion, uri)
}
func (hs *HttpServer) routerHandler() http.Handler {
	router := httprouter.New()
	router.GET(wrapUriWithVersion("/cache/:key"), hs.getHandler)
	router.PUT(wrapUriWithVersion("/cache/:key"), hs.setHandler)
	router.DELETE(wrapUriWithVersion("/cache/:key"), hs.deleteHandler)
	router.GET(wrapUriWithVersion("/status"), hs.statusHandler)
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
	ttl, err := ttlOf(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = hs.cache.SetWithTTL(key, value, ttl)
	if err != nil {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		w.Write([]byte("Error: " + err.Error()))
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func ttlOf(r *http.Request) (int64, error) {
	ttls, ok := r.Header["Ttl"]
	if !ok || len(ttls) < 1 {
		return caches.NeverDie, nil
	}
	return strconv.ParseInt(ttls[0], 10, 64)
}
func (hs *HttpServer) deleteHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	key := params.ByName("key")

	hs.cache.Delete(key)
}

func (hs *HttpServer) statusHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	status, err := json.Marshal(hs.cache.Status())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(status)
}
