package appHttp

import (
	"net/http"
	xhttp "cnv_xconfig/http"
	"github.com/julienschmidt/httprouter"
)

func init()  {
	router := httprouter.New()
	router.HandleOPTIONS = true
	http.Handle("/", xhttp.SetupAPIRoutes(router))
}
