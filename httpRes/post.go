package httpRes

import (
    "os"
    "bufio"
    "path/filepath"
    "io"
    "log"
    "errors"
    "g-tmp/hs-go/configs"
)



func Post(context *Context){

	uploadFile(context)
}


func uploadFile(context *Context) {
    
    multipartReader, err := context.R.MultipartReader()
    if err != nil {
    	log.Println(err)
    	return 
    }

    var content string

	for {
		partr, err := multipartReader.NextPart()
		if err != nil  {
			if err == io.EOF{
				break
			}else {
				context.HtmlR(500, "500 Internal Server Error")
				log.Println(err)
				return 
			}
		}
		defer partr.Close()

		absPath := filepath.Join(configs.Root, context.Path, partr.FileName())
		log.Println(absPath)

		// check uploaded files exist or not
		if _, err := os.Stat(absPath); err == nil {
			content +=  "<a style=\"color:orange\">" + partr.FileName() + "</a>" + "<p></p>"
		}else if errors.Is(err, os.ErrNotExist) {
			content +=  "<a style=\"color:green\">" + partr.FileName() + "</a>" + "<p></p>"
		}
		
		outputFile, err := os.Create(absPath)
		if err != nil {
			log.Println(err)
			context.HtmlR(500, "500 Internal Server Error")
			return
		}
		defer outputFile.Close()

		bWriter := bufio.NewWriter(outputFile)
		io.Copy(bWriter, partr)
	}

	context.HtmlR(200, "Uploaded <p></p>" + content)
}