package router

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

//Middleware ...
type Middleware func(http.HandlerFunc) http.HandlerFunc

//Routing ...
type Routing struct {
	Router           *mux.Router
	routerWithPrefix *mux.Router
	prefixes         []string
	middleware       []Middleware
	isGroup          bool
}

//NewRoute ...
func NewRoute(r *mux.Router) *Routing {
	return &Routing{Router: r}
}

//Get ...
func (r *Routing) Get(uri string, action http.HandlerFunc, name ...string) {
	r.compose(uri, action, http.MethodGet, name...)
}

//Head ...
func (r *Routing) Head(uri string, action http.HandlerFunc, name ...string) {
	r.compose(uri, action, http.MethodHead, name...)
}

// Post ...
func (r *Routing) Post(uri string, action http.HandlerFunc, name ...string) {
	r.compose(uri, action, http.MethodPost, name...)
}

// Put ...
func (r *Routing) Put(uri string, action http.HandlerFunc, name ...string) {
	r.compose(uri, action, http.MethodPut, name...)
}

// Patch ...
func (r *Routing) Patch(uri string, action http.HandlerFunc, name ...string) {
	r.compose(uri, action, http.MethodPatch, name...)
}

// Delete ...
func (r *Routing) Delete(uri string, action http.HandlerFunc, name ...string) {
	r.compose(uri, action, http.MethodDelete, name...)
}

//Option ...
func (r *Routing) Option(uri string, action http.HandlerFunc, name ...string) {
	r.compose(uri, action, http.MethodOptions, name...)
}

//HandleFilesystem ...
func (r *Routing) HandleFilesystem(uri string, handler http.Handler)  {
	r.Router.PathPrefix(uri).Handler(handler).Methods("GET")
}

//Middleware provide a convenient mechanism for filtering HTTP requests entering your application.
func (r *Routing) Middleware(middleware ...Middleware) *Routing {
	r.middleware = middleware
	return r
}

func (r *Routing) compose(uri string, action http.HandlerFunc, method string, name ...string) {

	//Grouping middleware
	wrapped := action
	for i := len(r.middleware) - 1; i >= 0; i-- {
		wrapped = r.middleware[i](wrapped)
	}

	url := uri
	if len(url) > 0 {
		url = fmt.Sprintf("/%s", uri)
	}

	var route = r.Router.HandleFunc(url, wrapped).Methods(method)
	if r.routerWithPrefix != nil {
		route = r.routerWithPrefix.HandleFunc(url, wrapped).Methods(method)
	}

	//If a least a name was provided, we get first one from the list and set as the Router name.
	if len(name) > 0 {
		route.Name(name[0])
	}

	//After the route has been configured, if we're outside of group Router then we'll clean up the middleware
	//then when the user try to add a new Router and wan use a middleware then they have to specify again.
	if r.routerWithPrefix == nil && !r.isGroup {
		r.middleware = nil
	}
}

// Prefix method may be used to prefix each route in the group with a given URI.
// For example, you may want to prefix all route URIs within the group with admin:
func (r *Routing) Prefix(prefix string, f func()) {

	defer func() {
		r.routerWithPrefix = nil
		if len(r.prefixes) > 0 {
			r.prefixes = r.prefixes[:len(r.prefixes)-1]
		}
	}()

	if len(prefix) == 0 {
		panic("Prefix(): the prefix can't be empty")
	}

	r.prefixes = append(r.prefixes, prefix)

	var mergePrefix = strings.Join(r.prefixes, "/")

	r.routerWithPrefix = r.Router.PathPrefix(fmt.Sprintf("/%s", mergePrefix)).Subrouter().StrictSlash(true)
	f()

}

// Group allow you to share route attributes, such as middleware or namespaces
// across a large number of routes without needing to define those
// attributes on each individual route. Shared attributes are
// specified in an slice format as the first parameter
// to the Route::group method.
func (r *Routing) Group(f func()) {
	defer func() {
		r.isGroup = false
	}()

	r.isGroup = true

	f()
}
