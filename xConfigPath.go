package cnv_xconfig

import (
	"github.com/cnvrtly/dstore"
	"strings"
	"context"
	"encoding/json"
	"google.golang.org/appengine/datastore"
)

const xConfigEntName string = "XCPTH"
const cacheExpiresInSec int64 = 60 * 60 * 24

//TODO update xConfigPaths plugins when saving
type PluginTemplate2XConfigs struct {
	PluginTemplate string
	XConfigPaths   []string
}

//TODO when saving new xscript -- for later impl.
type XScriptVersion2XConfigs struct {
	XScriptVersion string
	XConfigPaths   []string
}

type XConfigPathEntityPrivate struct {
	IsExactPath bool             `json:"isExactPath" datastore:"IsExactPath, index"`
	PluginsJson []byte           `json:"-" datastore:"PluginsJson, noindex"`
	Plugins     []*PluginPrivate `json:"plugins" datastore:"-"`
	Path        string           `json:"path" datastore:"-"`
	//TODO xscript script, plugins script
	*dstore.FindableEnt `json:"-" datastore:"-"`
}

// called by datastore when loading - use xConfigPathService.Get for retrieving
func (xce *XConfigPathEntityPrivate) Load(ps []datastore.Property) error {
	if err := datastore.LoadStruct(xce, ps); err != nil {
		return err
	}
	//parse json 2 plugins
	if xce.PluginsJson != nil {
		json.Unmarshal(xce.PluginsJson, &xce.Plugins)
	}
	xce.FindableEnt=&dstore.FindableEnt{}
	// Derive the Sum field.
	//xce.Sum = xce.I + xce.J
	return nil
}

// called by datastore when saving - use xConfigPathService.Save for saving
func (xce *XConfigPathEntityPrivate) Save() ([]datastore.Property, error) {
	if xce.Plugins != nil {
		var err error
		xce.PluginsJson, err = json.Marshal(&xce.Plugins)
		if err != nil {
			return nil, err
		}
	}
	// Save I and J as usual. The code below is equivalent to calling
	// "return datastore.SaveStruct(x)", but is done manually for
	// demonstration purposes.
	return datastore.SaveStruct(xce)

	//return []datastore.Property{
	//		{
	//			Name:  "I",
	//			Value: int64(xce.I),
	//		},
	//		{
	//			Name:  "J",
	//			Value: int64(xce.J),
	//		},
	//	}, nil

}

type xConfigPathEntityPublic struct {
	IsExactPath bool            `json:"isExactPath" datastore:"-"`
	Plugins     []*PluginPublic `json:"plugins" datastore:"-"`
	Path        string          `json:"path" datastore:"-"`
}

type PluginPrivate struct {
	Ident          string                 `json:"ident" datastore:"Ident, noindex"`
	Options        map[string]interface{} `json:"options" datastore:"Options, noindex"`
	PrivateOptions map[string]interface{} `json:"privateOptions" datastore:"PrivateOptions, noindex"`
}

type PluginPublic struct {
	Ident   string                 `json:"ident" datastore:"Ident, noindex"`
	Options map[string]interface{} `json:"options" datastore:"Options, noindex"`
}

func (xce *XConfigPathEntityPrivate) Unmarshal(b []byte) error {

	err := json.Unmarshal(b, xce)
	if err != nil {
		return err
	}
	xce.IsExactPath = isExactPath(xce.Path)
	xce.FindableEnt = &dstore.FindableEnt{}
	xce.FindableEnt.FindBy(xce.Path)
	return nil
}

func (xce *XConfigPathEntityPrivate) MarshallPublic() ([]byte, error) {
	path, _ := xce.FindBy("")
	xceCopy := &xConfigPathEntityPublic{IsExactPath: xce.IsExactPath, Path: path}
	for _, plugin := range xce.Plugins {
		xceCopy.Plugins = append(xceCopy.Plugins, &PluginPublic{Ident: plugin.Ident, Options: plugin.Options})
	}
	b, err := json.Marshal(xceCopy)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (xce *XConfigPathEntityPrivate) MarshallPrivate() ([]byte, error) {
	xce.Path, _ = xce.FindableEnt.FindBy("")
	b, err := json.Marshal(xce)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func NewXConfigPathEntity(path string, plugins []*PluginPrivate) *XConfigPathEntityPrivate {
	xPath := &XConfigPathEntityPrivate{
		IsExactPath: isExactPath(path),
		Plugins:     plugins,
		FindableEnt: &dstore.FindableEnt{},
	}
	xPath.FindableEnt.FindBy(path)
	return xPath
}

func NewXConfigService(cache dstore.SaverRetriever, store interface{}) *XConfigPathService {
	return &XConfigPathService{store: store.(dstore.SaverRetriever), cache: cache}
}

type XConfigPathService struct {
	store dstore.SaverRetriever
	cache dstore.SaverRetriever
}

func isExactPath(path string) bool {
	return !strings.ContainsAny(path, "*,") && !strings.HasPrefix(path, "-")
}

func (xcs *XConfigPathService) Save(ctx context.Context, namespaceId string, xcp *XConfigPathEntityPrivate) (error) {

	ident, err := xcp.FindBy("")
	if err != nil {
		return err
	}
	dsKOpt, _ := xcs.store.CreateKeyOptions(ctx, namespaceId, xConfigEntName, ident, 0)
	_, err = xcs.store.Save(dsKOpt, xcp, nil)
	if err != nil {
		return err
	}

	//keyOpt, _ := xcs.cache.CreateKeyOptions(ctx, namespaceId, xConfigEntName, ident, 0)
	//keyOpt.ExpiresIn(cacheExpiresInSec)
	//xcpJson, err := json.Marshal(xcp)
	//if err != nil {
	//	return err
	//}
	//_, err = xcs.cache.Save(keyOpt, xcpJson, nil)
	//if err != nil {
	//	return err
	//}
	return nil
}

func (xcs *XConfigPathService) Get(ctx context.Context, namespaceId string, path string) (*XConfigPathEntityPrivate, error) {

	xcp := &XConfigPathEntityPrivate{FindableEnt: &dstore.FindableEnt{}}
	xcp.FindBy(path)
	//TODO memcache first
	dsKOpt, err := xcs.store.CreateKeyOptions(ctx, namespaceId, xConfigEntName, path, 0)
	if err != nil {
		return nil, err
	}
	_, err = xcs.store.Load(dsKOpt, xcp)
	if err != nil {
		return nil, err
	}

	return xcp, nil
	//TODO implement caching
	/*transId, err := trans.FindBy("")
	if err != nil {
		return "", err
	}

	keyOpt, _ := memcacheService.CreateKeyOptions(ctx, GLOBAL_NS, xConfigEntName, transId, 0)
	keyOpt.ExpiresIn(cacheExpiresInSec)
	trJson, err := json.Marshal(trans)
	if err != nil {
		return "", err
	}
	memcacheService.Save(keyOpt, trJson, nil)*/
}

func (xcs *XConfigPathService) MarshallList(items []*XConfigPathEntityPrivate) ([]byte, error){
	res :=[]byte("[")
	var jsonVal []byte
	var err error
	for _,v:= range items{
		jsonVal, err=v.MarshallPrivate()
		if err!= nil {
			return nil, err
		}
		res=append(res, jsonVal...)
	}
	res=append(res, []byte("]")...)
	return res, nil
}

func (xcs *XConfigPathService) GetList(ctx context.Context, namespaceId string) ([]*XConfigPathEntityPrivate, error) {

	var resList []*XConfigPathEntityPrivate
	_, err:=xcs.store.GetAll(ctx, namespaceId, xConfigEntName, &resList, nil)
	if err != nil {
		return nil, err
	}

	/*for t := iterator.(*datastore.Iterator); ; {
		var ent *XConfigPathEntityPrivate = &XConfigPathEntityPrivate{FindableEnt:&dstore.FindableEnt{}}
		key, err := t.Next(ent)
		if err == datastore.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		ent.FindBy(key.StringID())
		resList = append(resList, ent)

	}*/

	return resList, nil
}

func (xcs *XConfigPathService) Delete(ctx context.Context, namespaceId string, path string) (error) {

	dsKOpt, err := xcs.store.CreateKeyOptions(ctx, namespaceId, xConfigEntName, path, 0)
	if err != nil {
		return err
	}
	return xcs.store.Delete(dsKOpt)
}
