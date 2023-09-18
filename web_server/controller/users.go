package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (ap AppRouter) WhoAmI(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	w.Write([]byte(r.Context().Value("Authorization").(string)))
}
