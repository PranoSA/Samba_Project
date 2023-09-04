package auth

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

/**
 *
 * HERE WE WILL DEFINE 2 AUTH INTERFACES THAT WILL BE USED BY THE APPLICATION FOR A) ACCOUNT MANAGEMENT, B)AUTH/SESSION MANAGEMENT
 *
 *
 * 1. If Configuring Authentication with LDAP, Postgres,  Keberos + LDAP, or DynamoDB, you will
 * create a "SIgnUpOnAuthentication" that will be used for certain routes (/signup, /register, /account/delete, login)
 * Otherwise Login and Account Management will be dedicated to an outside entity (Keycloak or other OIDC/SSO Provider)
 *
 * 2. For All Authentication, you will define a "Authentication" interface that will be used to authenticate Session state
 * (After Login by a user), which will be used by middleware for protected routes to create an identity context
 *
 */

/**
 *
 *  The Interfaces Defined in This FIle ARe
 *
 *
 * Authentication -> All Methods of Authenticating Requests Need to To Do This
 * The RedisSessionManager, BasicAuthSession, and OIDCProvider will Implement This
 *
 * UserManagementAuthentication -> Some Authenticators / Session Managers Will Allow Signing Up and Creating Users
 * These Managers Require an AuthenticationStore that Can Login, Signup, and Edit Users
 *
 * The Route Controller Will Open Up These Routes in the case that session-based or simple auth based is used, regardless
 * of the Backin Store
 *
 * CookieAuthentication -> Any Authenticator That Uses Cookies Will Implement This And Respond To This
 *
 * The purpose of this is to be extensible and as insanity check to make sure routes aren't reached that shouldn't be
 *
 * All Authenticators Should Implement These Interfaces, but many should return errors
 *
 */

type AuthInterface interface {
	Authentication
	UserManagementAuthentication
	CookieAuthentication
}

type Authentication interface {
	AuthenticationMiddleWare(next httprouter.Handle) httprouter.Handle
}

type UserManagementAuthentication interface {
	Login(w http.ResponseWriter, r *http.Request, ap httprouter.Params)
	Signup(w http.ResponseWriter, r *http.Request, ap httprouter.Params)
	Logout(w http.ResponseWriter, r *http.Request, ap httprouter.Params)
}

type CookieAuthentication interface {
	CSRF(w http.ResponseWriter, r *http.Request, ap httprouter.Params)
}

/**
 *
 * These Are Interfaces to Define the Backing Stores
 *
 * If you want to initialize any simple auth or session-based controlling, you must define a backing s
 * store for user that implements these methods and pass it into the initialization of those interfaces
 *
 */
type SignUpOnAuthentication interface {
	Login(username string, password string) bool
	Signup(username string, password string) bool
}
