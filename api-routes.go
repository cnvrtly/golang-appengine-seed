package http

import (
	"github.com/julienschmidt/httprouter"
	"github.com/cnvrtly/adaptr"
	"net/http"
	"google.golang.org/appengine"
	"github.com/cnvrtly/authorizzr-client"
	"github.com/cnvrtly/dstore"
	"github.com/cnvrtly/dstore/gae"
	"cnv_xconfig"
)

const apiBasePublic string = "/api/v1"

const authTokenCheckURL= "https://authorizzer.appspot.comapi/v1/check"
const authTokenApiKey = "fdasf"


var store dstore.SaverRetriever
var cache dstore.SaverRetriever
var xConfigPathService *cnv_xconfig.XConfigPathService


var initAdaptr = adaptr.CallOnce(func(w http.ResponseWriter, r *http.Request) {
	// This function is executed only once.
	store= &gae.DatastoreStoreService{}
	cache= &gae.MemcacheStoreService{}
	xConfigPathService=cnv_xconfig.NewXConfigService(cache, store)
})

func apiLifecycleAdaptrsPOST(testEnv bool) *RequestLifecycleAdapters {
	var initGAECtxAdaptrs =[]adaptr.Adapter{initAdaptr}
	if !testEnv{
		var gaeCtxAdaptr = adaptr.PlatformXCtxAdapter(appengine.NewContext)
		initGAECtxAdaptrs=append(initGAECtxAdaptrs, gaeCtxAdaptr)
	}
	var authPostAdaptrs = []adaptr.Adapter{adaptr.AuthBouncer(adaptr.CtxRouteAuthorizedKey)}
	var preJSONAuthRequestAdaptrs = append(
		initGAECtxAdaptrs,
		adaptr.Json2Ctx(adaptr.CtxRequestJsonStructKey, false ),
		adaptr.Tkn2Ctx(adaptr.CtxTokenKey, "", adaptr.CtxRequestJsonStructKey),
		authorizzr_client.ValidateToken(adaptr.CtxTokenKey, authTokenCheckURL, authTokenApiKey),
		adaptr.JsonContentType(),
	)
	return &RequestLifecycleAdapters{PreHandler:preJSONAuthRequestAdaptrs, PostHandler:authPostAdaptrs}
}

func getTestLifecycleAdaptrs(testEnv bool) *RequestLifecycleAdapters {
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

type RequestLifecycleAdapters struct {
	PreHandler []adaptr.Adapter
	PostHandler []adaptr.Adapter
}

func SetupAPIRoutes(router *httprouter.Router) *httprouter.Router {

	var postReqAdapters= apiLifecycleAdaptrsPOST(false)
	
	router.GET(apiBasePublic+"/test", testGET(getTestLifecycleAdaptrs(false)))
	router.POST(apiBasePublic+"/test/:id", testGET(postReqAdapters))
	//router.OPTIONS(apiBasePublic+"/test", adaptr.CreateOptionsRouterHandle(adaptr.Cors("")))
	//router.GET(apiBasePublic+"/xconfig", apiHandlers.backendApiAdaptersGET() )

	router.POST(apiBasePublic+"/xConfigPath", xConfigPathPOST(postReqAdapters))
	//http.Handle("/", router)
	return router
}

/*func newRouteHandleStatic(hFn http.HandlerFunc) httprouter.Handle {
	return adaptr.HttprouterAdaptFn(hFn, httpRouterUrlParamsKey)
}*/

func xConfigPathPOST(lifecycleAdapters *RequestLifecycleAdapters) httprouter.Handle {
	return adaptr.WrapHandleFuncAdapters(
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("ok POST"))
		},
		[]adaptr.Adapter{
			adaptr.AuthPermitAll(nil),
			//adaptr.WriteResponse(`backend wrks :)`),
		},
		lifecycleAdapters.PreHandler,
		lifecycleAdapters.PostHandler,
	)
}
func testGET(lifecycleAdapters *RequestLifecycleAdapters) httprouter.Handle {
	return adaptr.WrapHandleFuncAdapters(
		emptyHandler,
		[]adaptr.Adapter{
			adaptr.WriteResponse(`backend wrks :)`),
		},
		lifecycleAdapters.PreHandler,
		lifecycleAdapters.PostHandler,
	)
}
