package main

import (
    "fmt"
    "net/http"
    "strconv"
    "flag"
    "g-tmp/hs-go/httpRes"
)


type Engine struct{}


func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    fmt.Println("+", r.RemoteAddr, r.Method, r.URL.Path)
    
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
        fmt.Println(err)
        return 
    }
}


func main() {
	port := "11111"
	flag.PrintDefaults()
    flag.Parse() 
    if flag.NArg() == 1 {
	 	port = flag.Arg(0)
    }

    n, err := strconv.Atoi(port)
    if err != nil || n <= 0 || n > 65535{
    	fmt.Println("Invalid port! Port value is a number between 0 and 65535")
    	return 
    }

    fmt.Println("Listening on http://127.0.0.1:" + port)
    start(port)
}
