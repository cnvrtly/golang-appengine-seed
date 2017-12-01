package appHttp

import (
	"net/http"
)

func init()  {
	http.Handle("/", GetRouter())
}
