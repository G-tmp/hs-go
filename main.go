package main

import (
    "fmt"
    "net/http"
    "strconv"
    "log"
    "g-tmp/hs-go/httpRes"
    "g-tmp/hs-go/configs"
)


type Handler struct{}


func (handler *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    context := httpRes.NewContext(w, r)

    fmt.Println("+", r.RemoteAddr, r.Method, context.Path, r.URL.RawQuery)

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
    err := http.ListenAndServe(":" + strconv.Itoa(port), httpHandler)
    if err != nil {
        log.Println(err)
        return 
    }
}


func main() {
    fmt.Println("Listening on http://127.0.0.1:" + strconv.Itoa(configs.Port), configs.Root)
    start(configs.Port)
}
