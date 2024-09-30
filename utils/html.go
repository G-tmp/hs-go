package utils

import (
    "net/url"
    "os"
    "fmt"
    "strings"
    "path/filepath"
)


func Index(path string, files []os.DirEntry, showHidden bool) string {
	var sb strings.Builder

    var tf, oo string
    if showHidden {
        tf = "false"
        oo = "on"
    }else {
        tf = "true"
        oo = "off"
    }

    pp := fmt.Sprintf(
        `<!DOCTYPE html>
        <html>
        <head>
        <meta name="Content-Type" content="text/html; charset=utf-8">
        <title>%s</title>
        <style type="text/css">
            li{margin: 10px 0;}
        </style>
        </head>
        <body>
        <h1>Directory listing for %s</h1>
        <a href="?showHidden=%s"><button>Show Hidden Files</button></a> %s
        <p></p>
        <form method="post" enctype="multipart/form-data">
            <input type="file" name="file", required="required" multiple=""> >> <button type="submit">Upload</button>
        </form>
        <hr>
        `, path, path, tf, oo)
    sb.WriteString(pp)

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
            li := fmt.Sprintf(
                `<li>
                    <a href="%s/"><strong>%s</strong></a>
                </li>
                `, url.PathEscape(f.Name()), f.Name() + "/")
            sb.WriteString(li)
    	}
    }
    for _, f := range files {
    	if !f.IsDir(){
            li := fmt.Sprintf(`
                <li>
                    <a href="%s">%s</a>&emsp; - &emsp;<a href="%s">%s</a>
                </li>
                `, url.PathEscape(f.Name()), f.Name(), url.PathEscape(f.Name()) + "?download", "DL")
            sb.WriteString(li)
    	}
    }

    sb.WriteString("</ul>\n")
    sb.WriteString("<hr>")
    sb.WriteString("</body>\n</html>")

    return sb.String()
}
