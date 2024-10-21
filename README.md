# go-http-server

 A tiny http server map local files

 ```
-d  mapping as server root directory
-p  port
-c  certificate for tls
-k  certificate key for tls
 ```



## Dev

* build on linux/amd64

* determine files mime type ```https://github.com/gabriel-vasile/mimetype```

* natural sort ```https://github.com/maruel/natural```

* MultipartReader support large size and multi-part upload 



## Solved Problems

* http.DetectContentType() unable to detect flac 

* if err ==  no work, use Errors.Is() instead

* http.Error() delete content-type filed, so browser determine content-type automatically

* URL en/decode ```url.PathUnescape(r.URL.EscapedPath()) & url.PathEscape()```

* gabriel-vasile/mimetype unable determine ```css js``` and return text/plain
