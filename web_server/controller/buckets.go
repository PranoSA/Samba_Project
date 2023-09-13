package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

/**
 *
 * Eventually Add These Features
 *
 */

type DashHTTPRequest struct {
	Share_id    string
	File_name   string
	Resolutions []struct {
		Width   int
		Height  int
		Bitrate int
	}
}

func (ar AppRouter) RequestDash(w http.ResponseWriter, r *http.Request, ap httprouter.Params) {

	email := r.Context().Value("Authorization")

	var request DashHTTPRequest

	json.NewDecoder(r.Body).Decode(&request)

	fmt.Println(email)
}

func (ar AppRouter) CompressShare(w http.ResponseWriter, r *http.Request, ap httprouter.Params) {

}

func (ar AppRouter) GetCompressLinks(w http.ResponseWriter, r *http.Request, ap httprouter.Params) {

}
