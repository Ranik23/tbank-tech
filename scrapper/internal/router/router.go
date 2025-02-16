package router

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	ErrMethodUnknown = fmt.Errorf("unknown method name")
)

type Router struct {
	engine *gin.Engine
}

func NewRouter() *Router {
	return &Router{
		engine: gin.Default(),
	}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.engine.ServeHTTP(w, req)
}


func (r *Router) AddHandler(methodName string, path string, handler gin.HandlerFunc) error {
	switch methodName {
	case "POST":
		r.engine.POST(path, handler)
	case "GET":
		r.engine.GET(path, handler)
	default:
		return ErrMethodUnknown
	}
	return nil
}
