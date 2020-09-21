package main

import "net/http"

type HealthcheckHandler struct {
}

func (h HealthcheckHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(200)
}
