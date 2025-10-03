package httpRes

import (
    "os"
    "bufio"
    "path/filepath"
    "io"
    "log/slog"
    "errors"
    "g-tmp/hs-go/configs"
)



func Post(context *Context){

	uploadFile(context)
}


func uploadFile(context *Context) {
    
    multipartReader, err := context.R.MultipartReader()
    if err != nil {
		slog.Error(err.Error(), "addr", context.R.RemoteAddr, "method", context.Method, "path", context.Path, "query", context.R.URL.RawQuery)
    	context.HtmlR(500, "500 Internal Server Error")
    	return 
    }

    var content string

	for {
		partr, err := multipartReader.NextPart()
		if err != nil  {
			if err == io.EOF{
				break
			}else {
				slog.Error(err.Error(), "addr", context.R.RemoteAddr, "method", context.Method, "path", context.Path, "query", context.R.URL.RawQuery)
				context.HtmlR(500, "500 Internal Server Error")
				return 
			}
		}
		defer partr.Close()

		absPath := filepath.Join(configs.Root, context.Path, partr.FileName())

		// check uploaded files exist in system or not
		if _, err := os.Stat(absPath); err == nil {
			content +=  "<a style=\"color:orange\">" + partr.FileName() + "</a>" + "<p></p>"
		}else if errors.Is(err, os.ErrNotExist) {
			content +=  "<a style=\"color:green\">" + partr.FileName() + "</a>" + "<p></p>"
		}
		
		outputFile, err := os.Create(absPath)
		if err != nil {
			slog.Error(err.Error(), "addr", context.R.RemoteAddr, "method", context.Method, "path", context.Path, "query", context.R.URL.RawQuery)
			context.HtmlR(500, "500 Internal Server Error")
			return
		}
		defer outputFile.Close()

		bWriter := bufio.NewWriter(outputFile)
		io.Copy(bWriter, partr)
		slog.Info("", "addr", context.R.RemoteAddr, "method", context.Method, "path", context.Path, "file", partr.FileName())
	}

	context.HtmlR(200, "Uploaded <p></p>" + content)
}