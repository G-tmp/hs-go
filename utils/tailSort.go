package utils

import (
    "encoding/binary"
    "path/filepath"
    "strconv"
)


func TailSort(filename string) string {
    ext := filepath.Ext(filename)
    name := filename[:len(filename)-len(ext)]
    // split numeric suffix
    i := len(name) - 1
    for ; i >= 0; i-- {
        if  name[i] < '0' || name[i] > '9' {
            break
        }
    }
    i++
    // string numeric suffix to uint64 bytes
    // empty string is zero, so integers are plus one
    num := name[i:]
    b64 := make([]byte, 8)
    if len(num) > 0 {
        u64, err := strconv.ParseUint(num, 10, 64)
        if err == nil {
            binary.BigEndian.PutUint64(b64, u64+1)
        }
    }
    // prefix + numeric-suffix + ext
    return name[:i] + string(b64) + ext
}
