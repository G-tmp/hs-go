package main

import (
    "fmt"
    "net/http"
    "strconv"
    "flag"
    "g-tmp/hs-go/httpRes"
)


func start(port string) {
    http.HandleFunc("/", httpRes.Gepo)
    err := http.ListenAndServe(":" + port, nil)
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