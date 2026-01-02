package main

import (
    "os"
    "path/filepath"
    "io"
    "errors"

    "g-tmp/hs-go/configs"
    "gup"
)



func post(context *gup.Context) error {

	return uploadFile(context)
}


func uploadFile(context *gup.Context) error {
	defer context.R.Body.Close()

    multipartReader, err := context.R.MultipartReader()
    if err != nil {
    	// Isn't multipart/form-data or a multipart/mixed POST request
    	context.HtmlR(400, "400 Bad Request")
    	return err
    }

    var content string
	var errs []error

	for {
		partr, err := multipartReader.NextPart()
		if err != nil  {
			if err == io.EOF{
				break
			}else {
				context.HtmlR(500, "500 Internal Server Error")
				return err
			}
		}

		if partr.FileName() == "" {
            // Skip non-file part
            io.Copy(io.Discard, partr)
            partr.Close()
            continue
        }

        filename := partr.FileName()
        sysDir := filepath.Join(configs.Root, context.Path)

        tmp, err := os.CreateTemp(sysDir, filename + "_tmp*")
        if err != nil {
        	partr.Close()
        	content += "<a style=\"color:red\">" + "Save " + filename + " failed" + "</a>" + "<p></p>"
        	errs = append(errs, err)
        	continue 
        }

        _, err = io.Copy(tmp, partr)
        if err != nil {
        	partr.Close()
        	tmp.Close()
        	os.Remove(tmp.Name())
        	content += "<a style=\"color:red\">" + "Save " + filename + " failed" + "</a>" + "<p></p>"
        	errs = append(errs, err)
        	continue
        }

        if err = tmp.Sync(); err != nil {
        	partr.Close()
        	tmp.Close()
        	os.Remove(tmp.Name())
        	content += "<a style=\"color:red\">" + "Save " + filename + " failed" + "</a>" + "<p></p>"
        	errs = append(errs, err)
        	continue
        }

        partr.Close()

        if err = tmp.Close(); err != nil {
			partr.Close()
			os.Remove(tmp.Name())
        	content += "<a style=\"color:red\">" + "Save " + filename + " failed" + "</a>" + "<p></p>"
        	errs = append(errs, err)
			continue
		}

		// check uploaded files exist or not
		if _, err := os.Stat(filepath.Join(sysDir, filename)); err == nil {
			content +=  "<a style=\"color:orange\">" + partr.FileName() + "</a>" + "<p></p>"
		}else if errors.Is(err, os.ErrNotExist) {
			content +=  "<a style=\"color:green\">" + partr.FileName() + "</a>" + "<p></p>"
		}

		// rename temp file 
        if err = os.Rename(tmp.Name(), filepath.Join(sysDir, filename)); err != nil {
        	os.Remove(tmp.Name())
			content += "<a style=\"color:red\">" + "Save " + filename + " failed" + "</a>" + "<p></p>"
			errs = append(errs, err)
			continue
        }
	}

	context.HtmlR(200, "Uploaded <p></p>" + content)
	return errors.Join(errs...)
}