package httpRes

import(
	"net/http"
	"net/url"
    "strconv"
    "fmt"
)



type Context struct{
	R			*http.Request
	W 			http.ResponseWriter
	Path 		string
	Method 		string
	StatusCode 	int
	// Cookies 	[]http.Cookie
}


func NewContext(w http.ResponseWriter, r *http.Request) *Context{
	path, _ := url.PathUnescape(r.URL.EscapedPath())

	return &Context{
		R: r,
		W: w,
		Path: path,
		Method: r.Method,
	}
}


func (context *Context) Query(key string) string {
	return context.R.URL.Query().Get(key)
}


func (context *Context) Redirect(location string)  {
	http.Redirect(context.W, context.R, location, http.StatusFound)
}


func (context *Context) Cookie(key string) (*http.Cookie, error) {
	return context.R.Cookie(key) 
}


func (context *Context) SetCookie(cookie *http.Cookie)  {
	http.SetCookie(context.W, cookie)
}


func (context *Context) GetHeader(key string) string {
	return context.R.Header.Get(key)
}


func (context *Context) SetHeader(key string, value string){
	context.W.Header().Set(key, value)
}


func (context *Context) Status(code int){
	context.StatusCode = code
	context.W.WriteHeader(code)
}


func (context *Context) Html(code int, html string){
	context.Status(code)
	context.SetHeader("Content-Length", strconv.Itoa(len([]byte(html))))
	context.SetHeader("Content-Type", "text/html; charset=utf-8")
	context.W.Write([]byte(html))
}


func (context *Context) HtmlR(code int, content string){
	html :=fmt.Sprintf(
		`<!DOCTYPE html>
		<html>
		<head>
		<meta name="Content-Type" content="text/html; charset=utf-8">
		<title>%d</title>
		</head>
		<body>
		<center><h2>%s</h2></center> 
		</body>
		</html>`, code, content)
	context.Html(code, html)
}


func (context *Context) Data(code int, data []byte){
	context.Status(code)
	// io.Copy(context.w, file)
	context.W.Write(data)
}
