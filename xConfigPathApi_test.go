package http

import (
	"testing"
	"fmt"
	"net/http/httptest"
	"strings"
	"net/http"
	"google.golang.org/appengine/aetest"
	"context"
)

func TestXConfigPathApiHandlers(t *testing.T) {

	ctx, close, err := aetest.NewContext()

	if err != nil {
		panic(err.Error())
	}
	defer close()
	tests := []func(t *testing.T){
		testCtxWrap(ctx, tXConfigPathPOST),
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
	if err != nil {
		t.Errorf("req err=%s", err.Error())
	}
	r = r.WithContext(ctx)

	fmt.Println(record)
	xConfigPathPOST(apiLifecycleAdaptrsPOST(true))(record, r, nil)
	if record.Code != 200 {
		t.Errorf("post req not OK resp=%v", record.Body.String())
	}
	if record.Body.String() != "okkk" {
		t.Errorf("post response not OK body=%v", record.Body.String())
	}

	//t.Errorf("STATUS=%v body=%v", record.Code, record.Body)

	/*err:=xcT.Unmarshal(
		[]byte(jsonXCPath),
	)
	if err != nil {
		t.Errorf("Unmarshal err=%s", err)
	}

	err= xConfigPathService.save(ctx, nsId, xcT)
	if err != nil {
		t.Errorf("Save err=%s", err)
	}

	xcpL, err := xConfigPathService.get(ctx, nsId, xConfigPathStr)
	if err != nil {
	}

	if val, err:=marshallPrivate(xcpL); err!=nil || strings.Compare( string(val), jsonXCPath)!=0 {
	//if len(xcpL.Plugins)< 1 || !xcpL.IsExactPath{
		t.Errorf("datasto save not valid Plugins not saved err=%s val=%v", err, string(val))
	}

	err = xConfigPathService.delete(ctx, nsId, xConfigPathStr)
	if err != nil {
		t.Errorf("Delete err=%s", err)
	}

	_, err = xConfigPathService.get(ctx, nsId, xConfigPathStr)
	if err != dstore.ErrorNotFound {
		t.Errorf("should be deleted err=%s", err)
	}
*/
}
