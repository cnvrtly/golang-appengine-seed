package http

import (
	"testing"
	"net/http/httptest"
	"strings"
	"net/http"
	"google.golang.org/appengine/aetest"
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
)

func TestXConfigPathApiHandlers(t *testing.T) {

	ctx, close, err := aetest.NewContext()

	if err != nil {
		panic(err.Error())
	}
	defer close()
	tests := []func(t *testing.T){
		testCtxWrap(ctx, tXConfigPathPOST),
		testCtxWrap(ctx, tXConfigPathGETList),
		testCtxWrap(ctx, tXConfigPathGETItem),
		testCtxWrap(ctx, tXConfigPathDELETE),
	}

	for _, f := range tests {
		t.Run("Logging t.Run test=", f)
	}

}

func testCtxWrap(ctx context.Context, testCtxFunc func(t *testing.T, ctx context.Context)) func(t *testing.T) {
	return func(t *testing.T) {
		testCtxFunc(t, ctx)
	}
}

func tXConfigPathPOST(t *testing.T, ctx context.Context) {

	if ctx == nil {
		var close func()
		var err error
		ctx, close, err = aetest.NewContext()

		if err != nil {
			panic(err.Error())
		}
		defer close()
	}
	/*
	req1, inst := cutils.Testing_InitAppengineRequest(t)
	defer inst.Close()

	ctx := appengine.NewContext(req1)
	fmt.Println(req1, ctx)*/
	////////////////////////

	//var nsId = "test.com"

	//var apiKeyStr string
	var r *http.Request
	xConfigPathStr := "/path1"

	jsonXCPath := `{"isExactPath":true,"plugins":[{"ident":"creator:PLUGIN123","options":{"one":"valueee"},"privateOptions":null}],"path":"` + xConfigPathStr + `"}`
	record := httptest.NewRecorder()

	r, err := http.NewRequest("POST", "/xConfigPath", strings.NewReader(jsonXCPath))
	decoded:= `{"Name":"authorizzer.com-.-5639445604728832._.test1","Permissions":{"":["*"]},"ExtTknCacheId":"facebook§256853481494930§508262126199477","ExtIssuer":"facebook","ExtSubject":"256853481494930","ExtTknAudience":"508262126199477","ExtExpiresAt":1514516400,"aud":"authorizzer.com-.-5639445604728832._.test1","exp":1514514338,"jti":"1514510738§§authorizzer.com","iat":1514510738,"iss":"authorizzer.com","nbf":1514510738,"sub":"authorizzer.com-.-5639445604728832"}`
	tkn:=fmt.Sprintf("eyJhGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.%s.y2o39cJ6V0RM4hjPw0ytVzEH4BQDk1DxELIlOVdYeHA",jwt.EncodeSegment([]byte(decoded)))

	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tkn))
	if err != nil {
		t.Errorf("req err=%s", err.Error())
	}
	r = r.WithContext(ctx)
	xConfigPathPOST(apiLifecycleAdaptrsPOST(true))(record, r, nil)
	if record.Code != 200 {
		t.Errorf("POST ERR=%v", record.Body.String())
		return
	}

	if record.Body.String() != `{"isExactPath":true,"plugins":[{"ident":"creator:PLUGIN123","options":{"one":"valueee"}}],"path":"/path1"}` {
		t.Errorf("response not expected rCode=%v rBody=%s",record.Code, record.Body.String())
		return
	}

	xConfLoaded, err:=xConfigPathService.Get(ctx, "authorizzer.com-.-5639445604728832._.test1", xConfigPathStr)
	xConfLoadedJson, err:=xConfLoaded.MarshallPrivate()
	if err != nil {
		t.Errorf("marshall err=%s", err.Error())
		return
	}
	if string(xConfLoadedJson)!= `{"isExactPath":true,"plugins":[{"ident":"creator:PLUGIN123","options":{"one":"valueee"},"privateOptions":null}],"path":"/path1"}` {
		t.Errorf("saved and loaded not expected=%s", xConfLoadedJson)
	}
	fmt.Println("POST complete")
	//TODO check if auth bouncer is working
	//TODO check authorizzer checks if namespace matches apiKey - looks it does
}

func tXConfigPathGETList(t *testing.T, ctx context.Context) {

	if ctx == nil {
		var close func()
		var err error
		ctx, close, err = aetest.NewContext()

		if err != nil {
			panic(err.Error())
		}
		defer close()
	}
	////////////////////////

	//var nsId = "test.com"

	//var apiKeyStr string
	var r *http.Request

	record := httptest.NewRecorder()

	r, err := http.NewRequest("GET", "/xConfigPath", nil)
	decoded:= `{"Name":"authorizzer.com-.-5639445604728832._.test1","Permissions":{"":["*"]},"ExtTknCacheId":"facebook§256853481494930§508262126199477","ExtIssuer":"facebook","ExtSubject":"256853481494930","ExtTknAudience":"508262126199477","ExtExpiresAt":1514516400,"aud":"authorizzer.com-.-5639445604728832._.test1","exp":1514514338,"jti":"1514510738§§authorizzer.com","iat":1514510738,"iss":"authorizzer.com","nbf":1514510738,"sub":"authorizzer.com-.-5639445604728832"}`
	tkn:=fmt.Sprintf("eyJhGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.%s.y2o39cJ6V0RM4hjPw0ytVzEH4BQDk1DxELIlOVdYeHA",jwt.EncodeSegment([]byte(decoded)))

	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tkn))
	if err != nil {
		t.Errorf("req err=%s", err.Error())
	}
	r = r.WithContext(ctx)

	xConfigPathGET(apiLifecycleAdaptrsPOST(true))(record, r, nil)
	if record.Code != 200 {
		t.Errorf("GET list ERR=%v", record.Body.String())
	}

	if record.Body.String() != `[{"isExactPath":true,"plugins":[{"ident":"creator:PLUGIN123","options":{"one":"valueee"},"privateOptions":null}],"path":"/path1"}]` {
		t.Errorf("response not expected=%s", record.Body.String())
	}
	fmt.Println("GET list complete")
}

func tXConfigPathGETItem(t *testing.T, ctx context.Context) {

	if ctx == nil {
		var close func()
		var err error
		ctx, close, err = aetest.NewContext()

		if err != nil {
			panic(err.Error())
		}
		defer close()
	}
	////////////////////////

	//var nsId = "test.com"

	//var apiKeyStr string
	var r *http.Request
	xConfigPathStr := "/path1"

	record := httptest.NewRecorder()

	r, err := http.NewRequest("GET", "/xConfigPath",nil)
	decoded:= `{"Name":"authorizzer.com-.-5639445604728832._.test1","Permissions":{"":["*"]},"ExtTknCacheId":"facebook§256853481494930§508262126199477","ExtIssuer":"facebook","ExtSubject":"256853481494930","ExtTknAudience":"508262126199477","ExtExpiresAt":1514516400,"aud":"authorizzer.com-.-5639445604728832._.test1","exp":1514514338,"jti":"1514510738§§authorizzer.com","iat":1514510738,"iss":"authorizzer.com","nbf":1514510738,"sub":"authorizzer.com-.-5639445604728832"}`
	tkn:=fmt.Sprintf("eyJhGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.%s.y2o39cJ6V0RM4hjPw0ytVzEH4BQDk1DxELIlOVdYeHA",jwt.EncodeSegment([]byte(decoded)))

	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tkn))
	if err != nil {
		t.Errorf("req err=%s", err.Error())
	}
	r = r.WithContext(ctx)

	xConfigPathGET(apiLifecycleAdaptrsPOST(true))(record, r, httprouter.Params{httprouter.Param{Key:"id", Value:xConfigPathStr}})
	if record.Code != 200 {
		t.Errorf("GET item ERR=%v", record.Body.String())
	}

	if record.Body.String() != `{"isExactPath":true,"plugins":[{"ident":"creator:PLUGIN123","options":{"one":"valueee"},"privateOptions":null}],"path":"/path1"}` {
		t.Errorf("response not expected=%s", record.Body.String())
		return
	}
	fmt.Println("GET item complete")

}

func tXConfigPathDELETE(t *testing.T, ctx context.Context) {

	if ctx == nil {
		var close func()
		var err error
		ctx, close, err = aetest.NewContext()

		if err != nil {
			panic(err.Error())
		}
		defer close()
	}
	////////////////////////

	//var nsId = "test.com"

	//var apiKeyStr string
	var r *http.Request
	xConfigPathStr := "/path1"

	record := httptest.NewRecorder()

	r, err := http.NewRequest("GET", "/xConfigPath",nil)
	decoded:= `{"Name":"authorizzer.com-.-5639445604728832._.test1","Permissions":{"":["*"]},"ExtTknCacheId":"facebook§256853481494930§508262126199477","ExtIssuer":"facebook","ExtSubject":"256853481494930","ExtTknAudience":"508262126199477","ExtExpiresAt":1514516400,"aud":"authorizzer.com-.-5639445604728832._.test1","exp":1514514338,"jti":"1514510738§§authorizzer.com","iat":1514510738,"iss":"authorizzer.com","nbf":1514510738,"sub":"authorizzer.com-.-5639445604728832"}`
	tkn:=fmt.Sprintf("eyJhGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.%s.y2o39cJ6V0RM4hjPw0ytVzEH4BQDk1DxELIlOVdYeHA",jwt.EncodeSegment([]byte(decoded)))

	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tkn))
	if err != nil {
		t.Errorf("req err=%s", err.Error())
	}
	r = r.WithContext(ctx)

	xConfigPathDELETE(apiLifecycleAdaptrsPOST(true))(record, r, httprouter.Params{httprouter.Param{Key:"id", Value:xConfigPathStr}})
	if record.Code != 204 {
		t.Errorf("DELETE list ERR=%v", record.Body.String())
	}


	record = httptest.NewRecorder()
	r, _= http.NewRequest("GET", "/xConfigPath",nil)
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tkn))
	r = r.WithContext(ctx)

	xConfigPathGET(apiLifecycleAdaptrsPOST(true))(record, r, httprouter.Params{httprouter.Param{Key:"id", Value:xConfigPathStr}})
	if record.Code != 404 {
		t.Errorf("GET after DELETE ERR=%v", record.Body.String())
	}

	fmt.Println("DELETE item complete")

}
