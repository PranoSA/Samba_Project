package controller

import (
	"net/http"

	"github.com/PranoSA/samba_share_backend/web_server/auth"
	"github.com/julienschmidt/httprouter"
)

type AppRouter struct {
	CSRF_Protection bool
	CORS_Origins    []string
	Authenticator   auth.Authentication
	SignUpOn        auth.SignUpOnAuthentication
	CSRF            auth.CookieAuthentication
}

func AuthMiddleware(w http.ResponseWriter, r *http.Request, pa httprouter.Params) {

}

func NewAppRouter(approutes AppRouter) *httprouter.Router {

	router := httprouter.New()

	//Group Rotes
	router.POST("/group/:groupid", approutes.InviteUsers)

	//router.POST("")

	//Register Routes & Middleware Here

	//Unprotected Routes

	//Protected Routes

	return router
}
