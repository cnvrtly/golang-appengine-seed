package appHttp

import (
	"net/http"
	"encoding/json"
	"fmt"
	"strings"
	"google.golang.org/appengine/log"
	"errors"
	"io/ioutil"
	"github.com/cnvrtly/adaptr"
	"google.golang.org/appengine"
)

const ctxAppNamespaceIdKey = ctxAppNamespaceKeyType("appNamespaceID")

type ctxAppNamespaceKeyType string

const ctxApiKeyKey = ctxApiKeyKeyType("ctxApiKeyKey")

type ctxApiKeyKeyType string

const ctxRouteAuthorizedKey = ctxRouteAuthorizedType("routeAuthorized")

type ctxRouteAuthorizedType string

const ctxTokenUserIdentKey = ctxTokenUserIdentType("tokenUserIdent")

type ctxTokenUserIdentType string

const ctxTokenAudienceKey = tokenAudienceCtxType("tokenAudience")

type tokenAudienceCtxType string

const nativeTokenKey = ctxNativeTknType("nativeTkn")

type ctxNativeTknType string

const nativeTokenClaimsKey = ctxNativeTknClaimsType("nativeClaimsTkn")

type ctxNativeTknClaimsType string

const requestJsonStructCtxKey = requestJsonStructType("reqJsonStruct")

type requestJsonStructType string

const httpRouterUrlParamsKey = httpRouterUrlParamsType("httprouterParamsKey")

type httpRouterUrlParamsType string

func JsonOut(w http.ResponseWriter, jsonOutPointer interface{}) {
	res, _ := json.Marshal(jsonOutPointer)
	fmt.Fprint(w, string(res))
}


func JsonContentType() adaptr.Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			h.ServeHTTP(w, r)
		})
	}
}

func Cors(domain string, allowHeaders ... string) adaptr.Adapter {
	allowAll := domain == ""
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if allowAll {
				domain = r.Header.Get("origin")
			}
			w.Header().Set("Access-Control-Allow-Origin", domain)
			for _, hdr := range allowHeaders {
				w.Header().Add("Access-Control-Allow-Headers", hdr)
			}
			w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
			h.ServeHTTP(w, r)
		})
	}
}

func Json2Ctx(ctxKey interface{}, reset bool, requiredProps ... string) adaptr.Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if currCtxVal := adaptr.GetCtxValue(r, ctxKey); !reset && currCtxVal != nil {
				for _, param := range requiredProps {
					if _, ok := currCtxVal.(map[string]interface{})[param]; !ok {
						http.Error(w, fmt.Sprintf("Missing required JSON property name=%v", param), http.StatusBadRequest)
						return
					}
				}
				///return
			} else {

				valueStructPointer := map[string]interface{}{}
				if (r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch) {
					if r.Body == nil {
						http.Error(w, "Please send a request body", http.StatusBadRequest)
						return
					}

					err := json.NewDecoder(r.Body).Decode(&valueStructPointer)
					if err != nil {
						http.Error(w, "error parsing json err="+err.Error(), http.StatusBadRequest)
						return
					}
					for _, param := range requiredProps {
						if _, ok := valueStructPointer[param]; !ok {
							http.Error(w, fmt.Sprintf("Missing required JSON property=%v", param), http.StatusBadRequest)
							return
						}
					}
				}
				if r.Method == http.MethodGet {
					for _, param := range requiredProps {

						paramVal := r.URL.Query().Get(param)
						if paramVal == "" {
							http.Error(w, fmt.Sprintf("Missing required url param=%v", param), http.StatusBadRequest)
							return
						}

						valueStructPointer[param] = paramVal
					}
				}
				/*///if r.Body == nil {
					http.Error(w, "Please send a request body", http.StatusBadRequest)
					return
				}*/

				//if valueStructPointer == nil {
				///valueStructPointer := map[string]interface{}{}
				//}
				if (len(valueStructPointer) > 0) {
					r = adaptr.SetCtxValue(r, ctxKey, valueStructPointer)
				}
			}
			h.ServeHTTP(w, r)
		})
	}
}

func ReqrdParams(reqMethod string, requiredParams ... string) adaptr.Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for _, param := range requiredParams {
				switch reqMethod {
				case http.MethodGet:
					if r.URL.Query().Get(param) == "" {
						http.Error(w, fmt.Sprintf("Missing required url parameter=%v", param), http.StatusBadRequest)
						return
					}
				case http.MethodPost, http.MethodPut:

					if r.FormValue(param) == "" {
						hah, _ := ioutil.ReadAll(r.Body);
						//defer r.Body.Close()
						//r.Body.Read(str)
						http.Error(w, fmt.Sprintf("Missing required body parameter=%v val=%v", param, string(hah)), http.StatusBadRequest)
						return
					}
				}
			}
			h.ServeHTTP(w, r)
		})
	}
}

func ValidateCtxTkn(ctxTokenKey interface{}) adaptr.Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//check token

			var isValid bool = false
			ctxTknVal := adaptr.GetCtxValue(r, ctxTokenKey)
			if ctxTknVal != nil {
				var tknValue string = ctxTknVal.(string)
				if tknValue != "" {

					isValid=true
					var err error

					if isValid {
						h.ServeHTTP(w, r)
						return
					} else if err != nil {
						log.Errorf(r.Context(), "ValidateCtxTkn adaptr.Adapter nativeTS.Validate err=", err)
						http.Error(w, "Token not valid.", http.StatusUnauthorized)
						return
					}
				}
			}
			http.Error(w, "Authorization token not present", http.StatusUnauthorized)
			return
		})
	}

}

func Tkn2Ctx(ctxTokenKey interface{}, tknParameterName string) adaptr.Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tknValue, err := getTokenFromReq(r, tknParameterName)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			r = adaptr.SetCtxValue(r, ctxTokenKey, tknValue)

			h.ServeHTTP(w, r)
		})
	}
}

func AuthPermitAll() adaptr.Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, adaptr.SetCtxValue(r, ctxRouteAuthorizedKey, true))
		})
	}
}

func WriteResponse(writeValue string) adaptr.Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//_,err:=w.Write(writeValue)
			_, err := fmt.Fprintln(w, writeValue)
			if err != nil {
				http.Error(w, "", http.StatusInternalServerError)
				return
			}
			h.ServeHTTP(w, r)
		})
	}
}

func authBouncer() adaptr.Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Context().Value(ctxRouteAuthorizedKey) == true {
				if h != nil {
					h.ServeHTTP(w, r)
				}
			} else {
				http.Error(w, "Not authorized", http.StatusForbidden)
			}
		})
	}
}

func gaeCtx() adaptr.Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r.WithContext(appengine.NewContext(r)))
		})
	}
}

func getApiKeyFromReq(r *http.Request) string {
	apiKeyStr := r.URL.Query().Get("apiKey")
	if apiKeyStr != "" {
		return apiKeyStr
	}

	apiKeyStr = r.URL.Query().Get("apikey")
	if apiKeyStr != "" {
		return apiKeyStr
	}

	apiKeyStr = r.FormValue("apiKey")
	if apiKeyStr != "" {
		return apiKeyStr
	}

	apiKeyStr = r.FormValue("apikey")
	if apiKeyStr != "" {
		return apiKeyStr
	}

	return ""
}

func getTokenFromReq(r *http.Request, tknParameterName string) (string, error) {
	if tknParameterName != "" {
		var tknParValue string
		ctxJsonStruct := adaptr.GetCtxValue(r, requestJsonStructCtxKey).(map[string]interface{})
		if (ctxJsonStruct != nil) {
			if v, ok := ctxJsonStruct[tknParameterName]; ok {
				tknParValue = v.(string)
			}
		}

		if (tknParValue == "") {
			tknParValue = r.FormValue(tknParameterName)
			if tknParValue == "" {
				return "", fmt.Errorf("no token value in parameter=%v", tknParameterName)
			}
		}

		return tknParValue, nil
	}

	authHeaderVal := r.Header.Get("Authorization")
	if authHeaderVal == "" {
		return "", errors.New("No Authorization header value")
	}

	bearerStr := "Bearer"
	if last := strings.LastIndex(authHeaderVal, bearerStr); last > -1 {
		tknValue := strings.TrimSpace(authHeaderVal[last+len(bearerStr):])
		if tknValue != "" {
			return tknValue, nil
		}
	}

	return "", errors.New("Authorization header parse failed")
}

func choose(ss []string, test func(string) bool) (ret []string) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}
