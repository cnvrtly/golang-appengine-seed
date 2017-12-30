package http

import (
	"net/http"
	"github.com/cnvrtly/adaptr"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/cnvrtly/authorizzr-client"
	"fmt"
)

func apiSomeJsonHandler(w http.ResponseWriter, r *http.Request) {
	ctxVal := adaptr.GetCtxValue(r, adaptr.CtxRequestJsonStructKey).(map[string]interface{})
	jsonVal:= ctxVal["testVal"].(string)

	retJson, err := json.Marshal(map[string]string{"value": jsonVal})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write(retJson)
}

func xConfigPathGETHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintln(w, `{"valid":`+fmt.Sprint(true)+`}`)
}

func apiEmptyHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintln(w, `{"valid":`+fmt.Sprint(true)+`}`)
}

func xConfigPathPOST(lifecycleAdapters *RequestLifecycleAdapters) httprouter.Handle {
	return adaptr.WrapHandleFuncAdapters(
		func(w http.ResponseWriter, r *http.Request) {

			wrksObj := adaptr.GetCtxValue(r, adaptr.CtxWorkspaceIdentObjKey).(*authorizzr_client.WorkspaceIdentObject)
			w.Write([]byte(fmt.Sprintf("post res=%v", wrksObj.Workspace)))
			if wrksObj.Workspace == "" {
				http.Error(w, "workspace not in token", http.StatusBadRequest)
				return
			}

			namespaceId := wrksObj.Value
			w.Write([]byte(fmt.Sprintf("token namespace=%v", namespaceId)))

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
