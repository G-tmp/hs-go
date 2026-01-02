package configs

import (
    "flag"
    "log/slog"
    "os"
    "path/filepath"
    "time"

    "github.com/lmittmann/tint"
)

var (
    Root string     // root directory
    Port int        // http port
    Certificate string      
    Certificate_Key string
)


func init(){
    home, err := os.UserHomeDir()
    if err != nil {
        slog.Error(err.Error())
        os.Exit(1)
    }
    
    flag.IntVar(&Port, "p", 11111, "listening port, between 0 and 65535")
    flag.StringVar(&Root, "d", home ,"root directory")
    flag.StringVar(&Certificate, "c", "" ,"certificate file")
    flag.StringVar(&Certificate_Key, "k", "" ,"certificate key file")
    flag.Parse()

    Root, err = filepath.Abs(Root)
    if err != nil {
        slog.Error(err.Error())
        os.Exit(1)
    }

    slog.SetDefault(slog.New(
        tint.NewHandler(os.Stdout, &tint.Options{
            Level:      slog.LevelInfo,
            TimeFormat: time.DateTime,
            AddSource:  false,
            ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
                // log ignore string type with empty value
                if a.Value.Kind() == slog.KindString && a.Value.String() == "" {
                    return slog.Attr{}
                }
                // log ignore any type with nil value
                if a.Value.Kind() == slog.KindAny && a.Value.Any() == nil {
                    return slog.Attr{}
                }
                
                return a
            },
        }),
    ))
}