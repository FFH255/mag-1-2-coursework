package proto

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	Handle(c *gin.Context)
	Path() string
	Method() string
	Middleware() []string
}

type Middleware interface {
	Handle(c *gin.Context)
	Name() string
}

type Router struct {
	engine     *gin.Engine
	middleware map[string]gin.HandlerFunc
}

func (r *Router) getHandlers(handler Handler) []gin.HandlerFunc {
	var middlewares []gin.HandlerFunc

	for _, key := range handler.Middleware() {
		middleware, ok := r.middleware[key]
		if !ok {
			panic(fmt.Sprintf("middleware with %s key not found", key))
		}

		middlewares = append(middlewares, middleware)
	}

	middlewares = append(middlewares, handler.Handle)

	return middlewares
}

func (r *Router) Use(middleware ...gin.HandlerFunc) {
	r.engine.Use(middleware...)
}

func (r *Router) Middleware(middlewares ...Middleware) {
	for _, middleware := range middlewares {
		r.middleware[middleware.Name()] = middleware.Handle
	}
}

func (r *Router) Handle(handlers ...Handler) {
	for _, handler := range handlers {
		r.engine.Handle(handler.Method(), handler.Path(), r.getHandlers(handler)...)
	}
}

func NewRouter() *Router {
	engine := gin.Default()

	return &Router{
		engine:     engine,
		middleware: make(map[string]gin.HandlerFunc),
	}
}
