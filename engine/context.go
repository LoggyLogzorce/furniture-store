package engine

import "net/http"

type Context struct {
	Response http.ResponseWriter
	Request  *http.Request
}
