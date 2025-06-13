package app

import "net/http"

func InitRouters() *http.ServeMux {
	engine := &http.ServeMux{}

	engine.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello world"))
	})

	return engine
}
