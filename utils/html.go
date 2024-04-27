package utils

import (
    "fmt"
    "net/http"
    "net/url"
    "os"
    "strconv"
    "strings"
    "path/filepath"
)


func Index(path string, files []os.DirEntry, showHidden bool) string {
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
    if showHidden {
	    sb.WriteString("<a href=\"?showHidden=false\"><button>Show Hidden Files</button></a> on <p></p>")
    }else {
		sb.WriteString("<a href=\"?showHidden=true\"><button>Show Hidden Files</button></a> off <p></p>")
    }
    sb.WriteString("<form method=\"post\" enctype=\"multipart/form-data\">\n")
    sb.WriteString("<input type=\"file\" name=\"file\" required=\"required\" multiple>")
    sb.WriteString("&gt;&gt;")
    sb.WriteString("<button type=\"submit\">Upload</button>")
    sb.WriteString("</form>")
    sb.WriteString("<hr>\n")

    // parent directory
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
            sb.WriteString(url.PathEscape(f.Name()))
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
            sb.WriteString(url.PathEscape(f.Name()))
    		sb.WriteString("\">")
    		sb.WriteString(f.Name())
    		sb.WriteString("</a>")
    		sb.WriteString("</li>\n")   
    	}
    }

    sb.WriteString("</ul>\n")
    sb.WriteString("</body>\n</html>")

    return sb.String()
}


func  ErrorHtml(w http.ResponseWriter, content string, code int){
	html :=fmt.Sprintf(
		`<!DOCTYPE html>
		<html>
		<head>
		<meta name="Content-Type" content="text/html; charset=utf-8">
		<title>%d</title>
		</head>
		<body>
		<center><h2>%s</h2></center> 
		</body>
		</html>`, code, content)

	w.Header().Set("Content-Length", strconv.Itoa(len([]byte(html))) )
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(code)
	w.Write([]byte(html))
}
