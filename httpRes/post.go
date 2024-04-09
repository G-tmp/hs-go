package httpRes

import (
    // "fmt"
    "net/http"
    "os"
    "bufio"
    "io"
    "log"
    "strings"
    "g-tmp/hs-go/utils"
)


func Post(w http.ResponseWriter, r *http.Request){


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

		outputFile, err := os.Create(Home + r.URL.Path + partr.FileName())
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