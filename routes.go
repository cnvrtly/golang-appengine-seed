package appHttp

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"github.com/cnvrtly/dstore/gae"
	"github.com/cnvrtly/adaptr"
	"github.com/cnvrtly/dstore"
)

const publicApiKey = "apiKY"

var nativeApiKeys []string = []string{publicApiKey}

var memService dstore.SaverRetriever = &gae.MemcacheStoreService{}
var datastoreService dstore.SaverRetriever = &gae.DatastoreStoreService{}

const apiBasePublic string = "" // "/api/v1"
// const appDomainName string = "messenger.appspot.com"

var handlers = newHandlers(gaeCtx())

func GetRouter() *httprouter.Router {

	router := httprouter.New()
	//STATIC FILES START
	router.GET("/static/*filepath", handlers.staticFilesHandle())
	/*fileServer := http.FileServer(http.Dir("static"))
	router.GET("/static/*filepath", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		*//*w.Header().Set("Vary", "Accept-Encoding")
		w.Header().Set("Cache-Control", "public, max-age=7776000")*//*
		r.URL.Path = p.ByName("filepath")
		//log.Debugf(r.Context(), "OOOOOOOOOOOOOFFF")
		fileServer.ServeHTTP(w, r)
	})*/
	//STATIC FILES END

	router.GET(apiBasePublic+"/test", handlers.frontendApiAdaptersGET() )
	router.POST(apiBasePublic+"/test/:id", handlers.frontendApiAdaptersPOST() )
	router.OPTIONS(apiBasePublic+"/test", handlers.createOptionsRouterHandle(Cors("")))

	router.HandleOPTIONS = true

	//http.Handle("/", router)
	return router
}



type handlrs struct {
	initAdapter         adaptr.Adapter
	gaeCtxAdapter       adaptr.Adapter
	frontendApiAdapters []adaptr.Adapter
}

func (h *handlrs) createOptionsRouterHandle(corsAdapter adaptr.Adapter) httprouter.Handle {
	return h.newRouteHandle(emptyHandler, []adaptr.Adapter{
		AuthPermitAll(),
		corsAdapter,
	})
}

func (h *handlrs) staticFilesHandle() httprouter.Handle {

	fileServer := http.FileServer(http.Dir("static"))
	return h.newRouteHandleStatic(func(w http.ResponseWriter, r *http.Request) {
		//r.URL.Path = p.ByName("filepath")
		path := adaptr.GetCtxValue(r, httpRouterUrlParamsKey).(string)
		appPath := "/app/index.html"
		if path == "/app/api-key" {
			path = appPath
		}
		r.URL.Path = path
		fileServer.ServeHTTP(w, r)
	})
}

func (h *handlrs) frontendApiAdaptersGET() httprouter.Handle {
	return h.newRouteHandle(emptyHandler,
		append(
			[]adaptr.Adapter{
				AuthPermitAll(),
				WriteResponse(`wrks :)`),
			},
			h.frontendApiAdapters...
		))
}

func (h *handlrs) frontendApiAdaptersPOST() httprouter.Handle {
	return h.newRouteHandle(someJsonHandler,
		append(
			[]adaptr.Adapter{
				Json2Ctx(requestJsonStructCtxKey, false, "testVal"),
				AuthPermitAll(),
			},
			h.frontendApiAdapters...
		))
}

// ############################### util

func (h *handlrs) newRouteHandleStatic(hFn http.HandlerFunc) httprouter.Handle {
	return adaptr.HttprouterAdaptFn(hFn, httpRouterUrlParamsKey)
}

func (h *handlrs) newRouteHandle(hFn http.HandlerFunc, adapters []adaptr.Adapter) httprouter.Handle {
	//to beginning
	adapters = append([]adaptr.Adapter{h.initAdapter}, adapters...)
	if h.gaeCtxAdapter != nil {
		adapters = append([]adaptr.Adapter{h.gaeCtxAdapter}, adapters...)
	}
	//to end
	adapters = append(adapters, authBouncer())

	return adaptr.HttprouterAdaptFn(hFn, httpRouterUrlParamsKey, adapters...)
}

func newHandlers(gaeCtxAdapter adaptr.Adapter) *handlrs {
	initServerAdaptr := adaptr.CallOnce(func(w http.ResponseWriter, r *http.Request) {
		// This function is executed only once.
	})
	frontendApiAdapters := []adaptr.Adapter{
		//Tkn2Ctx(nativeTokenKey, ""), //
		//Json2Ctx(requestJsonStructCtxKey, false, "apiKey"),
		//ValidateCtxTkn(nativeTokenKey),
		JsonContentType(),
	}
	return &handlrs{initAdapter: initServerAdaptr, frontendApiAdapters: frontendApiAdapters, gaeCtxAdapter: gaeCtxAdapter}
}
