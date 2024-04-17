package main

import (
    "fmt"
    "net/http"
    "net/url"
    "strconv"
    "log"
    "g-tmp/hs-go/httpRes"
    "g-tmp/hs-go/configs"
)


type Handler struct{}


func (handler *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    path, err := url.PathUnescape(r.URL.EscapedPath())
    if err != nil {
        log.Println(err)
    }
    fmt.Println("+", r.RemoteAddr, r.Method, path, r.URL.RawQuery)

    switch r.Method {
    case "GET", "HEAD":
        httpRes.Get(w, r)
    case "POST":
        httpRes.Post(w, r)
    default:
        http.Error(w, "unsupported method" , http.StatusMethodNotAllowed)
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
