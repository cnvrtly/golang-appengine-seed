package http

import (
	"net/http"
	"github.com/cnvrtly/adaptr"
	"github.com/julienschmidt/httprouter"
	"github.com/cnvrtly/authorizzr-client"
	"cnv_xconfig"
	"github.com/cnvrtly/dstore"
	"fmt"
	"context"
)

func xConfigPathPOST(lifecycleAdapters *RequestLifecycleAdapters) httprouter.Handle {
	return adaptr.WrapHandleFuncAdapters(
		func(w http.ResponseWriter, r *http.Request) {
			wrksObj := adaptr.GetCtxValue(r, adaptr.CtxWorkspaceIdentObjKey).(*authorizzr_client.WorkspaceIdentObject)
			if wrksObj.Workspace == "" {
				http.Error(w, "workspace not in token", http.StatusBadRequest)
				return
			}

			namespaceId := string(wrksObj.Value)
			xConfPath := &cnv_xconfig.XConfigPathEntityPrivate{}
			xConfPath.FindableEnt = &dstore.FindableEnt{}
			bodyVal := adaptr.GetCtxValue(r, adaptr.CtxRequestBodyByteArrKey)

			if bodyVal == nil || len(bodyVal.([]byte)) == 0 {
				http.Error(w, "no request body value", http.StatusBadRequest)
				return
			}
			err := xConfPath.Unmarshal(bodyVal.([]byte))
			if err != nil {
				http.Error(w, fmt.Sprintf("can not parse json to xConfigPath err=%s json=%s", err, bodyVal), http.StatusBadRequest)
				return
			}

			err = xConfigPathService.Save(r.Context(), namespaceId, xConfPath)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			resJSON, err := xConfPath.MarshallPublic()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Write([]byte(resJSON))
		},
		[]adaptr.Adapter{
			adaptr.AuthPermitAll(nil),
			//adaptr.WriteResponse(`backend wrks :)`),
		},
		lifecycleAdapters.PreHandler,
		lifecycleAdapters.PostHandler,
	)
}

func xConfigPathGET(lifecycleAdapters *RequestLifecycleAdapters) httprouter.Handle {
	getOne := func(ctx context.Context, namespaceId string, xConfPathStr string) ([]byte, error) {
		xConfPath, err := xConfigPathService.Get(ctx, namespaceId, xConfPathStr)
		if err != nil {
			return nil, err
		}

		return xConfPath.MarshallPrivate()
	}

	getList := func(ctx context.Context, namespaceId string) ([]byte, error) {
		resList,err:=xConfigPathService.GetList(ctx, namespaceId)
		if err != nil {
			return nil, err
		}

		return xConfigPathService.MarshallList(resList)
	}
	return adaptr.WrapHandleFuncAdapters(
		func(w http.ResponseWriter, r *http.Request) {

			wrksObj := adaptr.GetCtxValue(r, adaptr.CtxWorkspaceIdentObjKey).(*authorizzr_client.WorkspaceIdentObject)
			if wrksObj.Workspace == "" {
				http.Error(w, "workspace not in token", http.StatusBadRequest)
				return
			}

			namespaceId := string(wrksObj.Value)
			xConfPath := &cnv_xconfig.XConfigPathEntityPrivate{}
			xConfPath.FindableEnt = &dstore.FindableEnt{}
			ctxIdParam := adaptr.GetCtxValue(r, adaptr.CtxRequestIdParamKey)
			var res []byte
			var err error
			if ctxIdParam != nil {
				var idParam string
				idParam = ctxIdParam.(string)
				res, err = getOne(r.Context(), namespaceId, idParam)

			} else {
				res, err = getList(r.Context(), namespaceId)
			}

			if err != nil && err != dstore.ErrorNotFound {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			if res != nil && err != dstore.ErrorNotFound {
				w.Write(res)
				return
			}
			http.Error(w, "not found", http.StatusNotFound)
			return
		},

		[]adaptr.Adapter{
			adaptr.AuthPermitAll(nil),
			//adaptr.WriteResponse(`backend wrks :)`),
		},
		lifecycleAdapters.PreHandler,
		lifecycleAdapters.PostHandler,
	)
}

func xConfigPathDELETE(lifecycleAdapters *RequestLifecycleAdapters) httprouter.Handle {
	return adaptr.WrapHandleFuncAdapters(
		func(w http.ResponseWriter, r *http.Request) {

			wrksObj := adaptr.GetCtxValue(r, adaptr.CtxWorkspaceIdentObjKey).(*authorizzr_client.WorkspaceIdentObject)
			if wrksObj.Workspace == "" {
				http.Error(w, "workspace not in token", http.StatusBadRequest)
				return
			}

			namespaceId := string(wrksObj.Value)
			xConfPath := &cnv_xconfig.XConfigPathEntityPrivate{}
			xConfPath.FindableEnt = &dstore.FindableEnt{}
			ctxIdParam := adaptr.GetCtxValue(r, adaptr.CtxRequestIdParamKey)

			if ctxIdParam != nil {
				var idParam string
				idParam = ctxIdParam.(string)
				err:=xConfigPathService.Delete(r.Context(), namespaceId, idParam)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusNoContent)
				return
			}
			http.Error(w, "no id", http.StatusBadRequest)
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
