package api

import "net/http"

func setCommonHeaders(w *http.ResponseWriter)(http.ResponseWriter){
	(*w).Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	(*w).Header().Set("Pragma", "no-cache")
	(*w).Header().Set("Expires", "0")
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
	return *w
}

func Index(w http.ResponseWriter, req *http.Request) {
	str := "{\"page\":\"index\"}"
	setCommonHeaders(&w)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(str))
}