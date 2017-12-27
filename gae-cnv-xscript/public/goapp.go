package appHttp

import (
	"net/http"
	http3 "cnv_xconfig/http"
	"github.com/julienschmidt/httprouter"
)

func init()  {
	router := httprouter.New()
	router.HandleOPTIONS = true
	http.Handle("/", http3.GetRouterPUBLIC(router))
}
