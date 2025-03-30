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
    context, err := httpRes.NewContext(w, r)
    if err != nil {
        log.Println(err)
        return 
    }

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

    if configs.Certificate != "" && configs.Certificate_Key != ""{
        fmt.Println("Listening on https://127.0.0.1:" + strconv.Itoa(configs.Port), configs.Root)
        err := http.ListenAndServeTLS(":" + strconv.Itoa(port), configs.Certificate, configs.Certificate_Key, httpHandler)
        if err != nil {
            log.Println(err)
            return 
        }
    } else {
        fmt.Println("Listening on http://127.0.0.1:" + strconv.Itoa(configs.Port), configs.Root)
        err := http.ListenAndServe(":" + strconv.Itoa(port), httpHandler)
        if err != nil {
            log.Println(err)
            return 
        }
    }

}


func main() {
    start(configs.Port)
}
