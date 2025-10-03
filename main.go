package main

import (
    "net/http"
    "strconv"
    "log/slog"
    "g-tmp/hs-go/httpRes"
    "g-tmp/hs-go/configs"
)


type Handler struct{}


// implement this method to taking over http request handling 
func (handler *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    context, err := httpRes.NewContext(w, r)
    if err != nil {
        slog.Error(err.Error(), "addr", r.RemoteAddr, "method", r.Method, "path", context.Path, "query", r.URL.RawQuery)
        return 
    }

    switch context.Method {
    case "GET", "HEAD":
        httpRes.Get(context)
    case "POST":
        httpRes.Post(context)
    default:
        context.HtmlR(405, "unsupported method")
    }
}


func start(port int) {
    httpHandler := new(Handler)

    if configs.Certificate != "" && configs.Certificate_Key != ""{
        slog.Info("Launch HTTPS Server", "addr","https://127.0.0.1:" + strconv.Itoa(configs.Port), "root", configs.Root)
        err := http.ListenAndServeTLS(":" + strconv.Itoa(port), configs.Certificate, configs.Certificate_Key, httpHandler)
        if err != nil {
            slog.Error(err.Error())
            return 
        }
    } else {
        slog.Info("Launch HTTP Server", "addr","http://127.0.0.1:" + strconv.Itoa(configs.Port), "root", configs.Root)
        err := http.ListenAndServe(":" + strconv.Itoa(port), httpHandler)
        if err != nil {
            slog.Error(err.Error())
            return 
        }
    }


}


func main() {
    start(configs.Port)
}
