package configs

import (
    "flag"
    "log"
    "os"
    "path/filepath"
)


var Root string    // root directory
var Port int    // http port


func init(){
    home, _ := os.UserHomeDir()
    
    flag.IntVar(&Port, "p", 11111, "listening port, between 0 and 65535")
    flag.StringVar(&Root, "d", home ,"root directory")
    flag.Parse()

    var err error
    Root, err = filepath.Abs(Root)
    if err != nil {
        log.Println(err)
        return 
    }

    log.SetFlags(log.LstdFlags | log.Lshortfile)
}