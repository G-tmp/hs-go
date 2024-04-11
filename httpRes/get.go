package httpRes

import (
    "fmt"
    "net/http"
    "net/url"
    "os"
    "io"
    "io/fs"
    "path/filepath"
    "strconv"
    "strings"
    "sort"
    "log"
    "errors"
    "g-tmp/hs-go/utils"
    "g-tmp/hs-go/configs"
    "github.com/gabriel-vasile/mimetype"
)


var path string

func Get(w http.ResponseWriter, r *http.Request){
	path, _ = url.PathUnescape(r.URL.EscapedPath())
	file, err := os.Open(filepath.Join(configs.Root, path))
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			utils.ErrorHtml(w, "404 Not Found", http.StatusNotFound)
		} else if errors.Is(err, fs.ErrPermission) {
			utils.ErrorHtml(w, "403 Forbidden", http.StatusForbidden)
		}
		log.Println(err)
		return 
	}
	defer file.Close()

	// check parameter, set cookie and redirect
	showHidden := r.URL.Query().Get("showHidden")
	if showHidden != "" {
		var cookie http.Cookie
		if showHidden == "true" {
			cookie = http.Cookie{Name: "showHidden", Value: "true", Path: "/"}
		} else{
			cookie = http.Cookie{Name: "showHidden", Value: "false", Path: "/"}
		}

		http.SetCookie(w, &cookie)
		http.Redirect(w, r, path, http.StatusFound)
		return
	}


	info, err := file.Stat()
	if err != nil {
		log.Println(err)
		utils.ErrorHtml(w, "404 Not Found", http.StatusNotFound)
		return 
	}

	if info.IsDir(){

	    // check cookie
		showHidden, err := r.Cookie("showHidden")
	    if err == nil {
	    	if showHidden.Value == "true" {
				respDir(w, r, true)
	    	}else {
				respDir(w, r, false)
	    	}
	    }else if err ==  http.ErrNoCookie {
	    	respDir(w, r, false)
	    }

	}else {
		respFile(w, r, file)
	}

}


func respFile(w http.ResponseWriter, r *http.Request, file *os.File){
	if  r.Header.Get("Range") != ""{
		partialReq(w, r, file)
		return
	}

	info, err := file.Stat()
	if err != nil {
		log.Println(err)
		http.NotFound(w, r)
		return 
	}

	mtype, _ := mimetype.DetectFile(filepath.Join(configs.Root, path))

	w.Header().Set("Content-Type", mtype.String()) 
	w.Header().Set("Content-Length", strconv.FormatInt(info.Size(), 10)) 
	w.Header().Set("Accept-Ranges", "bytes")
	
	io.Copy(w, file)

}


// handle range request
func partialReq(w http.ResponseWriter, r *http.Request, file *os.File){

	var start, end int64
	fmt.Sscanf(r.Header.Get("Range"), "bytes=%d-%d", &start, &end)

	info, err := file.Stat()
	if err != nil {
		log.Println(err)
		http.NotFound(w, r)
		return
	}
	if start < 0 ||start >= info.Size() ||end < 0 || end >= info.Size(){
		w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
		w.Write([]byte(fmt.Sprintf("out of index, length:%d",info.Size())))
		return
	}
	if end == 0 {
		end = info.Size() - 1
	}
	
	mtype, _ := mimetype.DetectFile(filepath.Join(configs.Root, path))
   	rg := fmt.Sprintf("bytes %d-%d/%d", start, end, info.Size())
	w.Header().Set("Content-Range", rg)
	w.Header().Set("Content-Length", strconv.FormatInt(end-start+1, 10))
	w.Header().Set("Accept-Ranges", "bytes")
	w.Header().Set("Content-Type", mtype.String())

	w.WriteHeader(http.StatusPartialContent)
	
	_, err = file.Seek(start, 0)
	if err != nil {
		log.Println(err)
		utils.ErrorHtml(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	
	io.CopyN(w, file, end-start+1)
	
}


func respDir(w http.ResponseWriter, r *http.Request, showHidden bool){
	files, err := os.ReadDir(filepath.Join(configs.Root, path))
	if err != nil {
        log.Println(err)
        return
    }

    // ignore case sensitive
    sort.Slice(files, 
        func(i, j int) bool {
            return strings.ToLower(files[i].Name()) < strings.ToLower(files[j].Name()) 
    })

    // filter hidden files 
    if !showHidden {
    	n := 0
    	for _, f := range files {
    		if f.Name()[0] != '.'{
    			files[n] = f
    			n++
    		}
    	}
    	files = files[0:n]
    }

    index := utils.Index(path, files, showHidden)

    w.Write([]byte(index))
}
