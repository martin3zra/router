package router_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/martin3zra/router"
)

var route *router.Routing

func TestMain(m *testing.M) {
	route = router.NewRoute(mux.NewRouter().StrictSlash(true))
	exitValue := m.Run()
	os.Exit(exitValue)
}

func TestHandlerMustBeCall(t *testing.T) {

	route.Prefix("admin", func() {
		route.Get("dashboard", dashboardHandler)
	})

	cases := []struct {
		method      string
		path        string
		status      int
		responseKey string
		name        string
	}{
		{
			method:      "GET",
			path:        "/admin/dashboard",
			status:      http.StatusOK,
			responseKey: "alive",
			name:        "it returns StatusOk for a valid path",
		},
		{
			method:      "GET",
			path:        "/admin/not-found-dashboard",
			status:      http.StatusNotFound,
			responseKey: "alive",
			name:        "it returns StatusNotFound for a invalid path",
		},
	}

	for _, item := range cases {
		t.Run(item.name, func(t *testing.T) {

			req, err := http.NewRequest(item.method, item.path, nil)
			if err != nil {
				t.Errorf("this is the error: %v", err)
			}

			response := httptest.NewRecorder()
			route.Router.ServeHTTP(response, req)

			if response.Code != item.status {
				t.Errorf("handler returned wrong status code: got %v want %v",
					response.Code, item.status)
			}
		})
	}
}

func TestRouting_Group(t *testing.T) {

	route.Prefix("admin", func() {
		route.Get("dashboard", dashboardHandler)

		route.Prefix("users", func() {
			route.Get("me", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				io.WriteString(w, `{"me": "Alfredo"}`)
			})
		})

		route.Prefix("posts", func() {
			route.Get("", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				io.WriteString(w, `{"posts": "news"}`)
			})
		})
	})

	route.Get("guest", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"guest": "jane"}`)
	})

	cases := []struct {
		method      string
		path        string
		status      int
		responseKey string
		name        string
	}{
		{
			method:      "GET",
			path:        "/admin/dashboard",
			status:      http.StatusOK,
			responseKey: "alive",
			name:        "dashboard it returns StatusOk for a valid path",
		},
		{
			method:      "GET",
			path:        "/admin/users/me",
			status:      http.StatusOK,
			responseKey: "me",
			name:        "me it returns StatusOk for a valid path with prefix",
		},
		{
			method:      "GET",
			path:        "/admin/posts",
			status:      http.StatusOK,
			responseKey: "posts",
			name:        "posts it returns StatusOk for a valid path with prefix",
		},
		{
			method:      "GET",
			path:        "/guest",
			status:      http.StatusOK,
			responseKey: "guest",
			name:        "guest it returns StatusOk for a valid path with prefix",
		},
	}

	for _, item := range cases {
		t.Run(item.name+item.path, func(t *testing.T) {

			req, err := http.NewRequest(item.method, item.path, nil)
			if err != nil {
				t.Errorf("this is the error: %v", err)
			}

			response := httptest.NewRecorder()
			route.Router.ServeHTTP(response, req)

			if response.Code != item.status {
				t.Errorf("handler returned wrong status code: got %v want %v",
					response.Code, item.status)
			}
		})
	}
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	// A very simple health check.
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{"alive": true}`)
}
