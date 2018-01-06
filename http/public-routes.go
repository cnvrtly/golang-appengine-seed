
package http

import (
	"github.com/julienschmidt/httprouter"
	"github.com/cnvrtly/dstore/gae"
	"github.com/cnvrtly/dstore"
	"github.com/cnvrtly/adaptr"
	"google.golang.org/appengine"
)

var memService dstore.SaverRetriever = &gae.MemcacheStoreService{}
var datastoreService dstore.SaverRetriever = &gae.DatastoreStoreService{}

// const appDomainName string = "messenger.appspot.com"

func publicLifecycleAdaptrsGETTest(testEnv bool) *RequestLifecycleAdapters {
	var initGAECtxAdaptrs =[]adaptr.Adapter{initAdaptr}
	if !testEnv{
		var gaeCtxAdaptr = adaptr.PlatformXCtxAdapter(appengine.NewContext)
		initGAECtxAdaptrs=append(initGAECtxAdaptrs, gaeCtxAdaptr)
	}
	var preJSONAuthRequestAdaptrs = append(
		initGAECtxAdaptrs,
	)
	return &RequestLifecycleAdapters{PreHandler:preJSONAuthRequestAdaptrs, PostHandler:nil}
}


func GetRouterPUBLIC(router *httprouter.Router) *httprouter.Router {

	router.GET(apiBasePublic+"/test", publicTestGET(publicLifecycleAdaptrsGETTest(false)))
	//router.POST(apiBasePublic+"/test/:id", testGET())
	//router.OPTIONS(apiBasePublic+"/test", adaptr.CreateOptionsRouterHandle(adaptr.Cors("")))

	//router.GET(apiBasePublic+"/xconfig", apiHandlers.backendApiAdaptersGET() )

	//http.Handle("/", router)
	return router
}


func publicTestGET(lifecycleAdapters *RequestLifecycleAdapters) httprouter.Handle {
	return adaptr.WrapHandleFuncAdapters(
		emptyHandler,
		[]adaptr.Adapter{
			adaptr.WriteResponse(`backend wrks :)`),
		},
		lifecycleAdapters.PreHandler,
		lifecycleAdapters.PostHandler,
	)
}

	/*//STATIC FILES START
	router.GET("/static/*filepath", handlers.staticFilesHandle())
	*/
/*fileServer := http.FileServer(http.Dir("static"))
	router.GET("/static/*filepath", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		*//*
/*w.Header().Set("Vary", "Accept-Encoding")
		w.Header().Set("Cache-Control", "public, max-age=7776000")*//*
/*
		r.URL.Path = p.ByName("filepath")
		//log.Debugf(r.Context(), "OOOOOOOOOOOOOFFF")
		fileServer.ServeHTTP(w, r)
	})*//*

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
				WriteResponse(`frontend wrks :)`),
			},
			h.frontendApiAdapters...
		))
}

func (h *handlrs) backendApiAdaptersGET() httprouter.Handle {
	return h.newRouteHandle(emptyHandler,
		append(
			[]adaptr.Adapter{
				AuthPermitAll(),
				WriteResponse(`backend wrks :)`),
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
*/
