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
	"google.golang.org/appengine/log"
)

const apiBasePublic string = "/api/v1"

const authTokenCheckURL = "https://authorizzer.appspot.com/api/v1/check"
const authTokenApiKey = "PrCMA2D3Xrp_h7bmozgDFXF4vxOFum258wFTAcx-ULYFpIo8m0CyHIyGEfnFFkdU5-4yYcyMoG_RaX_DG_EeJQ"

var store dstore.SaverRetriever
var cache dstore.SaverRetriever
var xConfigPathService *cnv_xconfig.XConfigPathService

// set in every handler and called on 1st request
var initAdaptr = adaptr.CallOnce(func(w http.ResponseWriter, r *http.Request) {
	// This function is executed only once.
	store = &gae.DatastoreStoreService{}
	cache = &gae.MemcacheStoreService{}
	xConfigPathService = cnv_xconfig.NewXConfigService(cache, store)
})

func initGAEAdapters(testEnv bool) []adaptr.Adapter {
	//GAE sets its own test ctx
	// init adapter for server init executed only on 1st call but needs to be set on every path
	initGAECtxAdaptrs := []adaptr.Adapter{initAdaptr}
	if !testEnv {
		gaeCtxAdaptr := adaptr.PlatformXCtxAdapter(appengine.NewContext)
		return append(initGAECtxAdaptrs, gaeCtxAdaptr)
	}
	return initGAECtxAdaptrs
}
func getTokenValidatorAdapter(testEnv bool) adaptr.Adapter {
	// tokens not validated in test env
	tokenValidatorAdapter := authorizzr_client.ValidateToken(adaptr.CtxTokenKey, authTokenCheckURL, authTokenApiKey)
	if testEnv {
		tokenValidatorAdapter = func(h http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				h.ServeHTTP(w, r)
			})
		}
	}
	return tokenValidatorAdapter
}

func apiLifecycleAdaptrsPOST(testEnv bool) *RequestLifecycleAdapters {
	//test env does not validate token and sets its own gae ctx
	initGAECtxAdaptrs := initGAEAdapters(testEnv)
	tokenValidatorAdapter := getTokenValidatorAdapter(testEnv)
	authPostAdaptrs := []adaptr.Adapter{adaptr.AuthBouncer(adaptr.CtxRouteAuthorizedKey)}

	preJSONAuthRequestAdaptrs := append(
		initGAECtxAdaptrs,
		adaptr.Json2Ctx(adaptr.CtxRequestJsonStructKey, false),
		adaptr.Tkn2Ctx(adaptr.CtxTokenKey, "", adaptr.CtxRequestJsonStructKey),
		authorizzr_client.UserIdentAndAudience2Ctx(adaptr.CtxTokenKey, adaptr.CtxTokenUserIdentKey, adaptr.CtxTokenAudienceKey),
		tokenValidatorAdapter,
		adaptr.JsonContentType(),
	)
	return &RequestLifecycleAdapters{PreHandler: preJSONAuthRequestAdaptrs, PostHandler: authPostAdaptrs}
}

func getTestLifecycleAdaptrs(testEnv bool) *RequestLifecycleAdapters {
	// test env sets its own gae ctx
	var initGAECtxAdaptrs = initGAEAdapters(testEnv)
	var preJSONAuthRequestAdaptrs = append(
		initGAECtxAdaptrs,
	)
	return &RequestLifecycleAdapters{PreHandler: preJSONAuthRequestAdaptrs, PostHandler: nil}
}

type RequestLifecycleAdapters struct {
	PreHandler  []adaptr.Adapter
	PostHandler []adaptr.Adapter
}

func SetupAPIRoutes(router *httprouter.Router) *httprouter.Router {
	var postReqAdapters = apiLifecycleAdaptrsPOST(false)

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
