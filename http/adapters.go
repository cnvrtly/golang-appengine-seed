package http

import (
	"net/http"
	"github.com/cnvrtly/adaptr"
)

/*const ctxAppNamespaceIdKey = ctxAppNamespaceKeyType("appNamespaceID")

type ctxAppNamespaceKeyType string

const ctxApiKeyKey = ctxApiKeyKeyType("ctxApiKeyKey")

type ctxApiKeyKeyType string*/


/*
const requestJsonStructCtxKey = requestJsonStructType("reqJsonStruct")

type requestJsonStructType string*/

/*const httpRouterUrlParamsKey = httpRouterUrlParamsType("httprouterParamsKey")

type httpRouterUrlParamsType string*/


func ProductionNamespace2Ctx(ctxTokenKey interface{}) adaptr.Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			domainName := adaptr.GetParamFromReqString(r, "ns")
			var err error
			if domainName== "" {
				//domainName=...cont get server name
			}

			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

				ns:=getProductionNamespaceFromDomainName(domainName)

			r = adaptr.SetCtxValue(r, ctxTokenKey, ns)

			h.ServeHTTP(w, r)
		})
	}
}

func getProductionNamespaceFromRequest(r *http.Request, tknParameterName string) {
	/*String domainName
	if(request.getParameter("ns")!=null){
		domainName= request.getParameter("ns")
	}else{
		domainName=request.getServerName()
	}
	return  getProductionNamespaceFromDomainName(domainName)*/

}

func getProductionNamespaceFromDomainName( domainName string ) string{
	return ""
	/*
List<String> subDomainOf=[AppConfig.applicationId+"."+AppConfig.appDomain,AppConfig.appDomain,"appengine.com","localhost"]
String retSubDomain=domainName
for (String subDom in subDomainOf) {
def tlDomainIndex = domainName.indexOf(subDom)
if (tlDomainIndex > -1) {
retSubDomain=domainName.substring(0,tlDomainIndex)

if(retSubDomain.endsWith("-dot-")){
retSubDomain=retSubDomain.substring(0,retSubDomain.size()-5).replaceAll("-dot-",".")
}else if(retSubDomain!=null && retSubDomain.size()>1){
retSubDomain=retSubDomain.substring(0,retSubDomain.size()-1)
}else if(retSubDomain.empty && subDomainOf.any{it==subDom}){
return subDom
}else {
return  null
}

break
}
}
return  retSubDomain*/
}


func choose(ss []string, test func(string) bool) (ret []string) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}

