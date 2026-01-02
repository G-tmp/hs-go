/*
	Only handle get und post method, no router
*/

package gup

import(
    "net/http"
	"log/slog"
)


type HandlerFunc func(*Context)

type engine struct{}

var getFunc HandlerFunc
var postFunc HandlerFunc

func (engine *engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    context, err := newContext(w, r)
    if err != nil {
        slog.Error(err.Error(), "addr", r.RemoteAddr, "method", r.Method, "path", context.Path, "query", r.URL.RawQuery)
        return 
    }

    switch context.Method {
    case "GET", "HEAD":
    	if getFunc == nil {
    		context.Html(200, "200 OK")
    	} else {
        	getFunc(context)
    	}
    case "POST":
    	if postFunc == nil {
    		context.Html(200, "200 OK")
    	} else {
        	postFunc(context)
    	}
    default:
        context.HtmlR(405, "unsupported method")
    }

}


func New() *engine{
	return new(engine)
}

func (engine *engine) Run(addr string) (err error){
    return http.ListenAndServe(addr, engine)
}

func (engine *engine) RunTLS(addr string, certificate string, key string) (err error){
    return http.ListenAndServeTLS(addr, certificate, key, engine)
}

func (engine *engine) Get(handler HandlerFunc){
	getFunc = handler
}

func (engine *engine) Post(handler HandlerFunc){
	postFunc = handler
}