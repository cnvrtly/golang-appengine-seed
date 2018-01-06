package cnv_xconfig

import (
	"testing"
	"fmt"
)

func TestXConfigPathMarshalling(t *testing.T) {

	if !isExactPath("/path"){
		fmt.Println("ERR isExactPath /path",isExactPath("/path"))
		t.FailNow()
	}

	if isExactPath("*"){
		fmt.Println("ERR isExactPath *",isExactPath("*"))
		t.FailNow()
	}

	xc := newXConfigPathEntity("/path", nil)
	b0, _ := marshallPublic(xc)
	if (string(b0) != `{"isExactPath":true,"plugins":null,"path":"/path"}`) {
		t.Errorf("NOT SAME MarshallNoPrivate got=%s", b0)
	}

	var plugins []*PluginPrivate
	plugins=append(plugins, &PluginPrivate{Ident:"creator:PLUGIN123", Options:map[string]interface{}{"one":"valueee"} })
	xc = newXConfigPathEntity("/path", plugins)
	b1, _ := marshallPrivate(xc)

	if (string(b1) != `{"isExactPath":true,"plugins":[{"ident":"creator:PLUGIN123","options":{"one":"valueee"},"privateOptions":null}],"path":"/path"}`) {
		t.Errorf("NOT SAME MarshallWithPrivate")
	}

	var xcEnt *XConfigPathEntityPrivate = &XConfigPathEntityPrivate{}

	err:=xcEnt.Unmarshal(
		[]byte(`{"isExactPath":true,"plugins":[{"ident":"creator:PLUGIN123","options":{"one":"valueee"},"privateOptions":null}],"path":"/path1"}`),
		)
	if err != nil {
		t.Errorf("Unmarshal err=%s", err)
	}
}
