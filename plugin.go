package cnv_xconfig

type PluginTemplate struct {
	Ident string `json:"ident" datastore:"Ident, noindex"`
	Title string `json:"title" datastore:"Title, noindex"`
	Script string `json:"script" datastore:"Script, noindex"`

}