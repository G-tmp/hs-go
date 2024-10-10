package configs

import (
    "flag"
    "log"
    "os"
    "path/filepath"
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
        log.Println(err)
        return 
    }
    
    flag.IntVar(&Port, "p", 11111, "listening port, between 0 and 65535")
    flag.StringVar(&Root, "d", home ,"root directory")
    flag.StringVar(&Certificate, "c", "" ,"certificate file")
    flag.StringVar(&Certificate_Key, "k", "" ,"certificate key file")
    flag.Parse()

    Root, err = filepath.Abs(Root)
    if err != nil {
        log.Println(err)
    }

    log.SetFlags(log.LstdFlags | log.Lshortfile)
}