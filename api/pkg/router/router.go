package router

import (
	"net/http"
	"regexp"
)

type RouteGroup struct {
	prefix     string
	middleware []func(http.Handler) http.Handler
	mux        *http.ServeMux
}

// NewGroup создает новую группу маршрутов
func NewGroup(mux *http.ServeMux, prefix string) *RouteGroup {
	return &RouteGroup{
		prefix:     prefix,
		mux:        mux,
		middleware: []func(http.Handler) http.Handler{},
	}
}

// Use добавляет middleware для группы
func (g *RouteGroup) Use(middleware ...func(http.Handler) http.Handler) {
	g.middleware = append(g.middleware, middleware...)
}

// Handle регистрирует обработчик с учетом префикса группы
func (g *RouteGroup) Handle(path string, handler http.Handler) {
	fullPath := g.prefix + path

	// Применяем все middleware группы
	wrappedHandler := handler
	for _, mw := range g.middleware {
		wrappedHandler = mw(wrappedHandler)
	}

	g.mux.Handle(fullPath, wrappedHandler)
}

// HandleFunc does same as handle, but for functions only
func (g *RouteGroup) HandleFunc(path string, f http.HandlerFunc) {
	g.Handle(path, f)
}

// PatternToRegex converts url pattern string to regexp
func PatternToRegex(pattern string) *regexp.Regexp {
	regexPattern := regexp.MustCompile(`\{[^/]+}`).ReplaceAllString("^"+pattern+"$", `([^/]+)`)
	return regexp.MustCompile(regexPattern)
}
