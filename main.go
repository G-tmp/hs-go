package main

import (
    "strconv"
    "log/slog"
    
    "g-tmp/hs-go/configs"
    "gup"
)


func main(){

    gp := gup.New()
    gp.Get(func(c *gup.Context){
        get(c)
    })
    gp.Post(func(c *gup.Context){
        post(c)
    })
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
