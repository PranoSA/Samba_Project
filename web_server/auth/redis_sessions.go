package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/redis/go-redis/v9"
)

type RedisSessionManager struct {
	RDB *redis.ClusterClient
	SUO SignUpOnAuthentication
}

func (r RedisSessionManager) ValidateCookie(string) (string, error) {

	return "", nil
}

/**
 *
 *  The Cookies Will Be Random Bytes Base64 Encoded
 *  The Hashes will
 *
 */
type CookieToken struct {
	Hash       string
	Expiration time.Time
	Csrf       string
	Csrf_exp   time.Time
	Email      string
}

/**
 *
 * Pass In THe SessionID Cookie Here
 */

func (rsm RedisSessionManager) AuthenticationMiddleWare(next httprouter.Handle) httprouter.Handle {

	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

		cookie_base64, err := r.Cookie("SESSIONID")

		ctx := context.WithValue(r.Context(), "Authentication", "")

		r = r.WithContext(ctx)

		if err == http.ErrNoCookie {
			fmt.Println("ErrNoCookie")
			next(w, r, p)
			return
		}

		fmt.Printf("Got Cookie %v \n", cookie_base64.Value)

		if err != nil {
			fmt.Println("ErrNoCookie")
			next(w, r, p)
			return
		}

		var cookie []byte

		cookie, err = base64.StdEncoding.DecodeString(cookie_base64.Value)

		if err != nil {
		}

		cookieKey := sha256.New()

		cookieKey.Write(cookie)

		cookieSum := cookieKey.Sum(nil)

		cookieToken := CookieToken{}

		fmt.Printf("Getting Key : %v \n", "Cookie:"+string(cookieSum))

		rfields, err := rsm.RDB.Get(context.Background(), "Cookie:"+string(cookieSum)).Result()

		if err != nil {
			fmt.Printf("Could Not Find Redis Field")
			next(w, r, p)
		}

		json.Unmarshal([]byte(rfields), &cookieToken)

		csrfHeader := r.Header.Get("X-CSRF-TOKEN")

		//csrfBytes, _ := base64.StdEncoding.DecodeString(csrfHeader)

		fmt.Printf("Getting CSRF Bytes : %v \n", csrfHeader)

		if cookieToken.Csrf != csrfHeader {
			fmt.Println(cookieToken.Csrf)
			fmt.Println("Calling Next CSRF WRONG")
			next(w, r, p)
			return
		}

		if cookieToken.Csrf_exp.After(time.Now()) && cookieToken.Expiration.After(time.Now()) {
			ctx := context.WithValue(r.Context(), "Authentication", cookieToken.Email)
			//fmt.Printf("CSRF Expiration : %v, Cookie Expiration %v, Now : %v \n", cookieToken.Csrf_exp.Unix(), cookieToken.Expiration.Unix)
			fmt.Printf("Here!!!!!!! \n")
			r = r.WithContext(ctx)
			next(w, r, p)
			return

		}
		fmt.Println("Calling Next EXPIRED ")
		next(w, r, p)

	})

}

type CSRFResponse struct {
	CSRF_TOKEN string
	Expiration time.Time
}

func (rsm RedisSessionManager) CSRF(w http.ResponseWriter, r *http.Request, rp httprouter.Params) {

	cookie_base64, err := r.Cookie("SESSIONID")

	if err != nil {

	}

	cookie, err := base64.StdEncoding.DecodeString(cookie_base64.Value)
	if err != nil {
	}

	cookieKey := sha256.New()

	cookieKey.Write(cookie)

	cookieSum := cookieKey.Sum(nil)

	cookieToken := CookieToken{}

	fmt.Printf("Setting Key : %v \n", "Cookie:"+string(cookieSum))

	rfields, err := rsm.RDB.Get(context.Background(), "Cookie:"+string(cookieSum)).Result()

	if err != nil {

	}

	json.Unmarshal([]byte(rfields), &cookieToken)

	var csrf_bytes []byte = make([]byte, 25)

	nbytes, err := rand.Read(csrf_bytes)

	if nbytes != 25 || err != nil {

	}

	csrf_header := base64.StdEncoding.EncodeToString(csrf_bytes)

	cookieToken.Csrf = csrf_header
	cookieToken.Csrf_exp = time.Now().Add(time.Minute * 15)

	fmt.Printf(" Setting CSRF  Token String %v \n", cookieToken.Csrf)

	//csrf_header := base64.StdEncoding.EncodeToString(csrf_bytes)

	json.NewEncoder(w).Encode(&CSRFResponse{
		CSRF_TOKEN: csrf_header,
		Expiration: cookieToken.Csrf_exp,
	})

}

type LoginRequest struct {
	Username string
	Password string
}

func (rsm RedisSessionManager) Login(w http.ResponseWriter, r *http.Request, pr httprouter.Params) {

	var Request LoginRequest = LoginRequest{}

	var bodybytes []byte

	_, e := r.Body.Read(bodybytes)
	if e != nil {

	}

	json.NewDecoder(r.Body).Decode(&Request)

	ok := rsm.SUO.Login(Request.Username, Request.Password)

	if ok == false {
		//Bad Login
		w.WriteHeader(501)
		return
	}

	newCookie := CookieToken{}

	cookieBytes := make([]byte, 36)

	rand.Read(cookieBytes)

	cookieByte := sha256.New()

	cookieByte.Write(cookieBytes)

	cookieKey := cookieByte.Sum(nil)

	clientCookie := base64.StdEncoding.EncodeToString(cookieBytes)

	newCookie.Hash = string(cookieKey)
	newCookie.Expiration = time.Now().Add(time.Hour * 24 * 30)

	csrf_bytes := make([]byte, 25)

	rand.Read(csrf_bytes)

	newCookie.Csrf = string(csrf_bytes)
	newCookie.Csrf_exp = time.Now().Add(15 * time.Minute)
	newCookie.Email = Request.Username

	csrfheader := base64.StdEncoding.EncodeToString([]byte(newCookie.Csrf))

	newCookie.Csrf = csrfheader

	redis_value, err := json.Marshal(newCookie)

	if err != nil {

	}

	fmt.Printf("Setting Key : %v \n", "Cookie:"+string(newCookie.Hash))

	err = rsm.RDB.Set(context.Background(), "Cookie:"+newCookie.Hash, string(redis_value), time.Hour*24*30).Err()

	if err != nil {

	}

	fmt.Printf("Setting CSRF Bytes : %v \n", newCookie.Csrf)

	json.NewEncoder(w).Encode(&CSRFResponse{
		CSRF_TOKEN: csrfheader,
		Expiration: newCookie.Expiration,
	})

	fmt.Printf("Setting Cookie %v \n", clientCookie)

	w.Header().Set("Set-Cookie", clientCookie)

}

/**
 * Revoke Both Cookie and CSRF Token
 *
 */
func (rsm RedisSessionManager) Logout(w http.Response, r *http.Request, params httprouter.Params) {

}
