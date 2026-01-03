package main

import (
    "strconv"
    "log/slog"
    "time"
    "fmt"
    
    "g-tmp/hs-go/configs"
    "gup"
)


func logger(next func(*gup.Context) error) gup.HandlerFunc {
    return func(c *gup.Context) {
        t := time.Now()
        err := next(c)
        
        if c.StatusCode >= 500 {
            slog.Error("", "addr", c.R.RemoteAddr, "method", c.Method, "path", c.Path, "query", c.R.URL.RawQuery, "state", c.StatusCode, "time", time.Since(t))
            fmt.Printf("%+v\n", err)
        } else if c.StatusCode >= 400 {
            slog.Warn("", "addr", c.R.RemoteAddr, "method", c.Method, "path", c.Path, "query", c.R.URL.RawQuery, "state", c.StatusCode, "time", time.Since(t), "err", err)
        }else {
            slog.Info("", "addr", c.R.RemoteAddr, "method", c.Method, "path", c.Path, "query", c.R.URL.RawQuery, "state", c.StatusCode, "time", time.Since(t))
        }
    }
}

func main(){
    gp := gup.New()
    gp.Get(logger(get))
    gp.Post(logger(post))

    if configs.Certificate != "" && configs.Certificate_Key != ""{
        slog.Info("Launch HTTPS Server", "addr","https://127.0.0.1:" + strconv.Itoa(configs.Port), "root", configs.Root)
        err := gp.RunTLS(":" + strconv.Itoa(configs.Port), configs.Certificate, configs.Certificate_Key)
        if err != nil {
            slog.Error(err.Error())
        }
    }else {
        slog.Info("Launch HTTP Server", "addr","http://127.0.0.1:" + strconv.Itoa(configs.Port), "root", configs.Root)
        err := gp.Run(":" + strconv.Itoa(configs.Port))
        if err != nil {
            slog.Error(err.Error())
        }
    }

}
