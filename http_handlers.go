package main

import (
	"encoding/json"
	"fmt"
	"github.com/ncw/swift"
	"net/http"
	"strconv"
)

func listRest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	items := Cache.Keys()
	fmt.Println(Cache.Projects)

	w.Header().Set("Total-Items", strconv.Itoa(len(items)))
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(items)
}

func getRest(w http.ResponseWriter, r *http.Request) {
	projectName := r.URL.Path[len("/item/"):]
	project, found := Cache.Projects[projectName]
	if !found {
		errorResponse(w, "404 Item not found", http.StatusNotFound)
		return
	}
	obj := project.getLatest()
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", obj.Filename))
	w.WriteHeader(http.StatusOK)
	validateCheckSum := true
	connObj.ObjectGet(SETTINGS.Get("OBJCONTAINER"), obj.Path, w, validateCheckSum, make(swift.Headers))
}

func errorResponse(w http.ResponseWriter, reason string, httpStatus int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(ErrorMsg{Error: "Error", Reason: reason, HTTPStatus: httpStatus})
}

// ErrorMsg Response structs
type ErrorMsg struct {
	Error      string
	Reason     string
	HTTPStatus int
}
