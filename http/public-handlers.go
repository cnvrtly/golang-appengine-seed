package http

import (
	"net/http"
	"github.com/cnvrtly/adaptr"
	"encoding/json"
)

func someJsonHandler(w http.ResponseWriter, r *http.Request) {
	ctxVal := adaptr.GetCtxValue(r, adaptr.CtxRequestJsonStructKey).(map[string]interface{})
	jsonVal:= ctxVal["testVal"].(string)

	retJson, err := json.Marshal(map[string]string{"value": jsonVal})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write(retJson)
}

func emptyHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintln(w, `{"valid":`+fmt.Sprint(true)+`}`)
}
