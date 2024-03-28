package main

import (
    "fmt"
    "net/http"
    "os"
    "bufio"
    "io"
    "strconv"
    "strings"
    "path/filepath"
    "flag"
)

var home string


func uploadFile(w http.ResponseWriter, r *http.Request) {
    
    multipartReader, err := r.MultipartReader()
    if err != nil {
    	fmt.Println(err)
    	return 
    }

    path := r.URL.Path

	partr, err := multipartReader.NextPart()
	if err != nil {
		fmt.Println(err)
	}
	defer partr.Close()

	outputFile, err := os.Create(home + r.URL.Path + partr.FileName())
	if err != nil {
		fmt.Println(err)
		return 
	}
	defer outputFile.Close()
	outputWriter := bufio.NewWriter(outputFile)

	io.Copy(outputWriter, partr)

	fmt.Fprintf(w, partr.FileName() + " has Uploaded to "+ path + "\n")
}


func gepo(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.RemoteAddr, r.Method, r.URL.Path)

	if r.Method == "POST"{
		uploadFile(w, r)
		return
	}


	path := r.URL.Path
	file, err := os.Open(home + path)
	if err != nil {
		http.NotFound(w, r)
		return 
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		http.NotFound(w, r)
		return 
	}

	if info.IsDir(){
		respDir(w, r, path)
		return
	}

	respFile(w, r, file)
}


func respFile(w http.ResponseWriter, r *http.Request, file *os.File){
	if  r.Header.Get("Range") != ""{
		getPart(w, r, file)
		return
	}

	info, err := file.Stat()
	if err != nil {
		http.NotFound(w, r)
		return 
	}

	w.Header().Set("Content-Length", strconv.FormatInt(info.Size(), 10)) 
	w.Header().Set("Accept-Ranges", "bytes")
	io.Copy(w, file)

}


// handle range request
func getPart(w http.ResponseWriter, r *http.Request, file *os.File){

	var start, end int64
	fmt.Sscanf(r.Header.Get("Range"), "bytes=%d-%d", &start, &end)

	info, err := file.Stat()
	if err != nil {
		fmt.Println(err.Error())
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
	buf := make([]byte, 512)
	_, err = file.Read(buf)
	file.Seek(0, 0)
	tp := http.DetectContentType(buf)
   	rg := fmt.Sprintf("bytes %d-%d/%d", start, end, info.Size())
	w.Header().Set("Content-Range", rg)
	w.Header().Set("Content-Length", strconv.FormatInt(end-start+1, 10))
	w.Header().Set("Accept-Ranges", "bytes")
	w.Header().Set("Content-Type", tp)

	w.WriteHeader(http.StatusPartialContent)
	
	_, err = file.Seek(start, 0)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = io.CopyN(w, file, end-start+1)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

}


func respDir(w http.ResponseWriter, r *http.Request, path string){
	files, err := os.ReadDir(home + path)
	if err != nil {
        fmt.Println(err)
        return
    }

	var sb strings.Builder
    sb.WriteString("<!DOCTYPE html>\n")
    sb.WriteString("<html>\n<head>\n")
    sb.WriteString("<meta name=\"Content-Type\" content=\"text/html; charset=utf-8\">\n")
    sb.WriteString("<title>")
    sb.WriteString(path)
    sb.WriteString("</title>\n")
    sb.WriteString("<style type=\"text/css\">\n")
    sb.WriteString("li{margin: 10px 0;}")
    sb.WriteString("\n</style>\n</head>\n")
    sb.WriteString("<body>\n")
    sb.WriteString("<h1>Directory listing for ")
    sb.WriteString(path)
    sb.WriteString("</h1>\n")
    sb.WriteString("<form method=\"post\" enctype=\"multipart/form-data\">\n")
    sb.WriteString("<input type=\"file\" name=\"file\" required=\"required\" >")
    sb.WriteString("&gt;&gt;")
    sb.WriteString("<button type=\"submit\">Upload</button>")
    sb.WriteString("</form>")
    sb.WriteString("<hr>\n")
    if path == "/" {
    	sb.WriteString("/")
    } else {
    	sb.WriteString("<a href=\"")
    	p := filepath.Dir(strings.TrimSuffix(path, "/"))
    	if p != "/"{
    		p += "/"
    	}
    	sb.WriteString(p)
    	sb.WriteString("\">")
    	sb.WriteString("Parent Directory</a>")
    }
    sb.WriteString("<ul>\n")

    for _, f := range files {
    	if f.IsDir(){
    		sb.WriteString("<li>")
    		sb.WriteString("<a href=\"")
    		sb.WriteString(f.Name())
    		sb.WriteString("/")
    		sb.WriteString("\">")
    		sb.WriteString("<strong>")
    		sb.WriteString(f.Name())
    		sb.WriteString("/")
    		sb.WriteString("</strong>")
    		sb.WriteString("</a>")
    		sb.WriteString("</li>\n")   
    	}
    }
    for _, f := range files {
    	if !f.IsDir(){
    		sb.WriteString("<li>")
    		sb.WriteString("<a href=\"")
    		sb.WriteString(f.Name())
    		sb.WriteString("\">")
    		sb.WriteString(f.Name())
    		sb.WriteString("</a>")
    		sb.WriteString("</li>\n")   
    	}
    }

    sb.WriteString("</ul>\n")
    sb.WriteString("</body>\n</html>")

    fmt.Fprintf(w, sb.String())
}


func start(port string) {
    http.HandleFunc("/", gepo)
    http.ListenAndServe(":"+port, nil)
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


func init(){
	home, _ = os.UserHomeDir()
}