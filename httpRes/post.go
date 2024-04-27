package httpRes

import (
    "net/http"
    "os"
    "bufio"
    "net/url"
    "path/filepath"
    "io"
    "log"
    "errors"
    "g-tmp/hs-go/utils"
    "g-tmp/hs-go/configs"
)


func Post(w http.ResponseWriter, r *http.Request){
	path, _ = url.PathUnescape(r.URL.EscapedPath())

	uploadFile(w, r)
}


func uploadFile(w http.ResponseWriter, r *http.Request) {
    
    multipartReader, err := r.MultipartReader()
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
				utils.ErrorHtml(w, "500 Internal Server Error", http.StatusInternalServerError)
				log.Println(err)
				return 
			}
		}
		defer partr.Close()

		abPath := filepath.Join(configs.Root, path, partr.FileName())
		log.Println(abPath)

		// check uploaded files exist or not
		if _, err := os.Stat(abPath); err == nil {
			content +=  "<a style=\"color:orange\">" + partr.FileName() + "</a>" + "<p></p>"
		}else if errors.Is(err, os.ErrNotExist) {
			content +=  "<a style=\"color:green\">" + partr.FileName() + "</a>" + "<p></p>"
		}

		outputFile, err := os.Create(abPath)
		if err != nil {
			log.Println(err)
			utils.ErrorHtml(w, "500 Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer outputFile.Close()

		bWriter := bufio.NewWriter(outputFile)
		io.Copy(bWriter, partr)
	}

	utils.ErrorHtml(w, "Uploaded <p></p>" + content, http.StatusOK)
}