package auth

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type AllAllowedAuthenticator struct {
}

func (aaa AllAllowedAuthenticator) AuthenticationMiddleWare(next httprouter.Handle) httprouter.Handle {

	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

		auth_user := r.Header.Get("Authorization")

		newctx := context.WithValue(r.Context(), "Authorization", auth_user)

		if auth_user == "" {

		}

		newr := r.WithContext(newctx)

		next(w, newr, p)

	}
}
