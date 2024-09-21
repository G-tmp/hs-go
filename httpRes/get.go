package httpRes

import (
    "fmt"
    "net/http"
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
    "github.com/maruel/natural"
)



func Get(context *Context){

	file, err := os.Open(filepath.Join(configs.Root, context.Path))
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			context.HtmlR(404, "404 Not Found")
		} else if errors.Is(err, fs.ErrPermission) {
			context.HtmlR(403, "403 Forbidden")
		}else {
			log.Println(err)
		}
		return 
	}
	defer file.Close()

	// check parameter, set cookie and redirect
	showHidden := context.Query("showHidden")
	if showHidden != "" {
		var cookie http.Cookie
		if showHidden == "true" {
			cookie = http.Cookie{Name: "showHidden", Value: "true", Path: "/"}
		} else{
			cookie = http.Cookie{Name: "showHidden", Value: "false", Path: "/"}
		}

		context.SetCookie(&cookie)
		context.Redirect(context.Path)
		return
	}


	info, err := file.Stat()
	if err != nil {
		log.Println(err)
			context.HtmlR(404, "404 Not Found")
		return 
	}

	if info.IsDir(){

	    // check cookie
		showHidden, err := context.Cookie("showHidden")
	    if err == nil {
	    	if showHidden.Value == "true" {
				respDir(context, true)
	    	}else {
				respDir(context, false)
	    	}
	    }else if err ==  http.ErrNoCookie {
	    	respDir(context, false)
	    }

	}else {
		respFile(context, file)
	}

}


func respFile(context *Context, file *os.File){
	if  context.GetHeader("Range") != ""{
		partialReq(context, file)
		return
	}

	info, err := file.Stat()
	if err != nil {
		log.Println(err)
		context.HtmlR(404, "404 Not Found")
		return 
	}

	mtype, _ := mimetype.DetectFile(filepath.Join(configs.Root, context.Path))
	t := mtype.String()
	if strings.HasSuffix(context.Path, ".css"){
		t = "text/css; charset=utf-8"
	}

	context.SetHeader("Content-Type", t) 
	context.SetHeader("Content-Length", strconv.FormatInt(info.Size(), 10)) 
	context.SetHeader("Accept-Ranges", "bytes")
	
	io.Copy(context.W, file)

}


// handle range request
func partialReq(context *Context, file *os.File){

	var start, end int64
	fmt.Sscanf(context.GetHeader("Range"), "bytes=%d-%d", &start, &end)

	info, err := file.Stat()
	if err != nil {
		log.Println(err)
		context.HtmlR(404, "404 Not Found")
		return
	}
	if start < 0 ||start >= info.Size() ||end < 0 || end >= info.Size(){
		context.HtmlR(416, fmt.Sprintf("out of index, length:%d",info.Size()))
		return
	}
	if end == 0 {
		end = info.Size() - 1
	}
	
	mtype, _ := mimetype.DetectFile(filepath.Join(configs.Root, context.Path))
   	rg := fmt.Sprintf("bytes %d-%d/%d", start, end, info.Size())
	context.SetHeader("Content-Range", rg)
	context.SetHeader("Content-Length", strconv.FormatInt(end - start + 1, 10))
	context.SetHeader("Accept-Ranges", "bytes")
	context.SetHeader("Content-Type", mtype.String())

	context.Status(206)
	
	_, err = file.Seek(start, 0)
	if err != nil {
		log.Println(err)
		context.HtmlR(500, "500 Internal Server Error")
		return
	}
	
	io.CopyN(context.W, file, end - start + 1)
	
}


func respDir(context *Context, showHidden bool){
	files, err := os.ReadDir(filepath.Join(configs.Root, context.Path))
	if err != nil {
        log.Println(err)
        return
    }

    // tail sorting and ignore case sensitive
	sort.Slice(files, 
        func(i, j int) bool {
        	return natural.Less(strings.ToLower(files[i].Name()), strings.ToLower(files[j].Name()))
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

    index := utils.Index(context.Path, files, showHidden)

    context.Html(200, index)

}
