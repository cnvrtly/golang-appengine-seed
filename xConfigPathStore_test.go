package cnv_xconfig

import (
	"testing"
	"github.com/cnvrtly/cutils"
	"google.golang.org/appengine"
	"github.com/cnvrtly/dstore/gae"
	"github.com/cnvrtly/dstore"
	"strings"
)

func TestNativeTokenService_GetOrCreateToken(t *testing.T) {

	req1, inst := cutils.Testing_InitAppengineRequest(t)
	defer inst.Close()

	ctx := appengine.NewContext(req1)
////////////////////////

	store := &gae.DatastoreStoreService{}
	cache:= &gae.MemcacheStoreService{}
	xConfigPathService:=&XConfigPathService{store:store, cache:cache}

	var nsId = "test.com"

	//
	xConfigPathStr := "/path1"
	//entName := xConfigEntName

	xcT := &XConfigPathEntityPrivate{}

	jsonXCPath := `{"isExactPath":true,"plugins":[{"ident":"creator:PLUGIN123","options":{"one":"valueee"},"privateOptions":null}],"path":"` + xConfigPathStr + `"}`
	err:=xcT.Unmarshal(
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

}
