# go-http-server

 A tiny http(s) server for local files mapping

 ```
-d  local directory as server root directory
-p  port
-c  CA
-k  CA key
 ```



## Dev

* develop on linux/amd64

* determine files mime type ```https://github.com/gabriel-vasile/mimetype```

* natural sort ```https://github.com/maruel/natural```

* MultipartReader support large size and multi-part uploading

* local CA ```https://github.com/FiloSottile/mkcert```



## Solved Problems

* official http.DetectContentType() unable to detect flac 

* if err ==  no work, use Errors.Is() instead

* http.Error() delete content-type filed, so browser determine content-type automatically

* URL en/decode ```url.PathUnescape(r.URL.EscapedPath()) & url.PathEscape()```

* gabriel-vasile/mimetype unable determine ```css js``` so return text/plain manually
