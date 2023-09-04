package controller

import (
	"net/http"

	"github.com/PranoSA/samba_share_backend/web_server/auth"
	"github.com/PranoSA/samba_share_backend/web_server/models"
	"github.com/julienschmidt/httprouter"
)

type AppRouter struct {
	CORS_Origins  []string
	Authenticator auth.Authentication
	Models        models.Models
}

func AuthMiddleware(w http.ResponseWriter, r *http.Request, pa httprouter.Params) {

}

func NewAppRouter(approutes AppRouter) *httprouter.Router {

	router := httprouter.New()

	middleware := approutes.Authenticator.AuthenticationMiddleWare

	// Inherited Routes

	if usermanagement, ok := approutes.Authenticator.(auth.UserManagementAuthentication); ok {
		router.POST("/signup", usermanagement.Signup)
		router.POST("/login", usermanagement.Login)
		router.GET("/logout", usermanagement.Logout)
	}

	if sessionmanager, ok := approutes.Authenticator.(auth.CookieAuthentication); ok {
		router.POST("/csrf", sessionmanager.CSRF)
	}

	//Group & Share Rotes
	router.DELETE("/group", middleware(approutes.DeleteShare))

	router.POST("/group", middleware(approutes.CreateShare))

	router.POST("/group/:groupid", middleware(approutes.InviteUsers))

	router.POST("/invite/:inviteid", middleware(approutes.AcceptInvite))

	// Space Routes

	// To Be Implemented -> Compute Routes

	return router
}
