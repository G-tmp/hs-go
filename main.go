package main

import (
    "fmt"
    "net/http"
    "net/url"
    "strconv"
    "flag"
    "log"
    "path/filepath"
    "g-tmp/hs-go/httpRes"
    "g-tmp/hs-go/configs"
)


type Engine struct{}


func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    path, err := url.PathUnescape(r.URL.EscapedPath())
    if err != nil {
        log.Println(err)
    }
    fmt.Println("+", r.RemoteAddr, r.Method, path)

    switch r.Method {
    case "GET", "HEAD":
        httpRes.Get(w, r)
    case "POST":
        httpRes.Post(w, r)
    default:
        http.Error(w, "unsupported method" , http.StatusMethodNotAllowed)
    }
}


func start(port string) {
    engine := new(Engine)
    err := http.ListenAndServe(":" + port, engine)
    if err != nil {
        log.Println(err)
        return 
    }
}


func main() {
	port := "11111"
	flag.PrintDefaults()
    flag.Parse()
    if flag.NArg() == 1 {
        port = flag.Arg(0)
    }else if flag.NArg() == 2 {
        port = flag.Arg(0)
        root, err := filepath.Abs(flag.Arg(1))
        if err != nil {
            log.Println(err)
            return 
        }

        configs.Root = root
    }

    n, err := strconv.Atoi(port)
    if err != nil || n <= 0 || n > 65535{
    	log.Println("Invalid port! Port value is a number between 0 and 65535")
    	return 
    }

    configs.Port = port
    fmt.Println("Listening on http://127.0.0.1:" + port, configs.Root)
    start(configs.Port)
}
