package controller

import (
	"log"
	"net/http"

	"github.com/PranoSA/samba_share_backend/web_server/auth"
	"github.com/PranoSA/samba_share_backend/web_server/models"
	"github.com/julienschmidt/httprouter"
	"github.com/rabbitmq/amqp091-go"
	httpSwagger "github.com/swaggo/http-swagger"
)

type BadRequestResponse []struct {
	ParameterName string
	Param_Type    string //Path Or Query
	Value_Type    string //int32, int64, string, []bytes or utf-8, base64, etc...
	Message       string //Extra Information Such as Valid Ranges, etc....
}

type AppRouter struct {
	CORS_Origins  []string
	Authenticator auth.Authentication
	Models        models.Models
	Queue         *amqp091.Connection
}

/**
 *
 * Host & Origin Checks
 *
 * Only Need to Set Cors Header if Origin Header Exists
 *
 * Need to add further * Functionality
 *
 */

func (approutes AppRouter) CorsMiddleware(next httprouter.Handle) func(http.ResponseWriter, *http.Request, httprouter.Params) {

	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		if r.Header.Get("origin") == "" || r.Header.Get("origin") == r.Header.Get("host") {
			next(w, r, p)
			return
		}

		if approutes.CORS_Origins[0] == "*" {
			(w).Header().Set("Access Controller-Allow-Origin", "*")
			next(w, r, p)
			return
		}

		for i, v := range approutes.CORS_Origins {
			if r.Header.Get("Origin") == approutes.CORS_Origins[i] {

				(w).Header().Set("Access-Control-Allow-Origin", v)
				next(w, r, p)
				return
			}
		}

		next(w, r, p)
	}
}

func swaggerHandler(res http.ResponseWriter, req *http.Request, p httprouter.Params) {
	httpSwagger.WrapHandler(res, req)
}

func (approutes AppRouter) NewAppRouter() *httprouter.Router {

	router := httprouter.New()

	middleware := approutes.Authenticator.AuthenticationMiddleWare

	// Inherited Routes
	if _, ok := approutes.Authenticator.(auth.Authentication); !ok {
		log.Fatal("approutes.Authenticator Must Implement auth.Authenticator")
	}

	if usermanagement, ok := approutes.Authenticator.(auth.UserManagementAuthentication); ok {
		router.POST("/signup", usermanagement.Signup)
		router.POST("/login", usermanagement.Login)
		router.GET("/logout", usermanagement.Logout)
	}

	if sessionmanager, ok := approutes.Authenticator.(auth.CookieAuthentication); ok {
		router.POST("/csrf", sessionmanager.CSRF)
	}

	approutes.StartCompressPublisher()
	approutes.StartDashPublisher()
	//Group & Share Rotes
	router.DELETE("/group/:shareid", middleware(approutes.DeleteShare)) //Only Owner Can DO THis

	router.POST("/space/:spaceid/group", middleware(approutes.CreateShare))

	router.POST("/group/:shareid", middleware(approutes.InviteUsers)) //Only Owners Can DO This

	router.POST("/invite/:inviteid", middleware(approutes.AcceptInvite)) //Only Users with Invite ID Can Do THis

	// Space Routes

	router.GET("/spaces", middleware(approutes.GetMySpaces))

	router.POST("/spaces", middleware(approutes.CreateSpace))

	router.DELETE("/spaces/:spaceid", middleware(approutes.DeleteSpace))

	router.GET("/doc/:any", swaggerHandler)

	// To Be Implemented -> Compute Routes
	//

	router.POST("/spaces/mpegdash/:filename", middleware(approutes.CreateShare))

	/**
	 *
	 * Compute Route -> Share + File Name on Mp4 File ....
	 * /shareid/mount
	 *
	 * Goes to RabbitMQ Queue
	 *
	 * Golang App Consumed The Queue Basedon Server Routing Key
	 *
	 *
	 *
	 */

	/**
	 *
	 * Auxillary Functionality
	 *
	 */

	// router.POST("/space/:spaceid/publications")

	// router.POST("/space/:spaceid/snapshots")

	return router
}
