# go-http-server


## TIRED 😣

* go version go1.22.1 linux/amd64


## Encountered Problems

* fmt.Fprintf() and w.Write() do not flush content, response body may be empty

* http.DetectContentType() unable to detect flac 

* err ==  no work, use Errors.Is()

* http.Error() delete content-type filed, so it detect content-type automatically