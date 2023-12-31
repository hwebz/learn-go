package main

import (
	"fmt"
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<h1>Welcome to my awesome site!</h1>")
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<h1>Contact Page</h1><p>To get in touch, email me at <a href=\"mailto:admin@test.com\">admin@test.com</a>.</p>")
}

func pathHandler(w http.ResponseWriter, r *http.Request) {
	// if r.URL.Path == "/" {
	// 	homeHandler(w, r)
	// } else if r.URL.Path == "/contact" {
	// 	contactHandler(w, r)
	// }

	switch r.URL.Path {
	case "/":
		homeHandler(w, r)
	case "/contact":
		contactHandler(w, r)
	default:
		// TODO: handle the page not found error
		// w.WriteHeader(http.StatusNotFound)
		// fmt.Fprint(w, "Page not found")

		http.Error(w, "Page Not Found", http.StatusNotFound)
	}

	// r.URL.Path vs r.URL.RawPath
}

type Router struct{}

func (router Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		homeHandler(w, r)
	case "/contact":
		contactHandler(w, r)
	default:
		http.Error(w, "Page Not Found", http.StatusNotFound)
	}
}

func main() {
	var router Router
	fmt.Println("Starting the server on 3000...")
	http.ListenAndServe(":3000", router)
}
