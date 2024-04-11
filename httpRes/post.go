package httpRes

import (
    // "fmt"
    "net/http"
    "os"
    "bufio"
    "net/url"
    "path/filepath"
    "io"
    "log"
    "strings"
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

    var names []string

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
		outputFile, err := os.Create(abPath)
		if err != nil {
			log.Println(err)
			utils.ErrorHtml(w, "500 Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer outputFile.Close()
		outputWriter := bufio.NewWriter(outputFile)
		names = append(names, partr.FileName())

		io.Copy(outputWriter, partr)
	}

	s := strings.Join(names, "<p></p>")
	utils.ErrorHtml(w, "Have Uploaded <p></p>" + s, http.StatusOK)
}