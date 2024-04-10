package configs

import (
    "os"
    "log"
)


var Root string    // root directory

var Port string    // http port


func init(){
    Root, _ = os.UserHomeDir()
    log.SetFlags(log.LstdFlags | log.Lshortfile)
}