package auth_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PranoSA/samba_share_backend/web_server/auth"
	"github.com/julienschmidt/httprouter"
	"github.com/redis/go-redis/v9"
)

//Lets Test Redis Sessions ....

/**
 *
 *
 * Cases :
 * 1. No Cookie Set
 * 2. Incorrect Cookie Set
 * 3. Correct Cookie Set, No CSRF
 * 4. Correct Cookie Set with Incorrect CSRF
 * 5. Correct Cookie WIth Correct CSRF
 *
 * 6. Generate CSRF
 * 7. Use Said Generated CSRF With Correct Cookie
 *
 * 8. Generate Cookie With Login
 * 9. Generate CSRF
 */

var RedisClient *redis.ClusterClient = redis.NewClusterClient(&redis.ClusterOptions{
	Addrs:    []string{"localhost:6379"},
	Password: "", // no password set
	//DB:       0,  // use default DB
})

type DummyUserDatastore struct {
}

func (d DummyUserDatastore) Login(username string, password string) bool {

	//Correct Password for padler@uci.edu is passwordforpadler

	fmt.Printf(" Username : %v : Password : %v \n", username, password)

	if username == "padler@uci.edu" {
		if password == "passwordforpadler" {
			return true
		}

		return false
	}
	return true
}

func (d DummyUserDatastore) Signup(username string, password string) bool {
	return true
}

var sessions auth.RedisSessionManager = auth.RedisSessionManager{
	RDB: RedisClient,
	SUO: DummyUserDatastore{},
}

func TestRedisSession(t *testing.T) {

	fmt.Println("Starting Testing Of Login & Context Auth")

	type TestLoginAndAuth struct {
		Login            bool
		Name             string
		Password         string
		Email            string
		Expected_Context string
	}

	var TestLoginsAndAuth []TestLoginAndAuth = []TestLoginAndAuth{
		{
			Login:            true,
			Name:             "Correct Password, padler@uci.edu",
			Password:         "passwordforpadler",
			Email:            "padler@uci.edu",
			Expected_Context: "padler@uci.edu",
		},
		{
			Login:            true,
			Name:             "Incorrect Password, padler@uci.edu",
			Password:         "ppppppp",
			Email:            "padler@uci.edu",
			Expected_Context: "",
		},
		{
			Login:            true,
			Name:             "Other User",
			Password:         "",
			Email:            "pcadler@gmail.com",
			Expected_Context: "pcadler@gmail.com",
		},
		{
			Login:            false,
			Name:             "No Login ",
			Password:         "",
			Email:            "pcadler@gmail.com",
			Expected_Context: "",
		},
	}

	for _, test := range TestLoginsAndAuth {
		t.Run(test.Name, func(t *testing.T) {

			// Run a Login + Later Authentication Here ...
			body, _ := json.Marshal(auth.LoginRequest{
				Username: test.Email,
				Password: test.Password,
			})
			//Mocking an HTTP Response
			req := httptest.NewRequest(http.MethodGet, "/login", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			sessions.Login(w, req, httprouter.Params{})

			if w.Header().Get("Set-Cookie") == "" && test.Expected_Context != "" {
				//Not Actually Authenticated But Expected To
				t.Errorf("Expected To Login with %v : %v, But Did not ", test.Email, test.Password)
			}

			if w.Header().Get("Set-Cookie") != "" && test.Expected_Context == "" {
				t.Logf("Status Code %v", w.Result().StatusCode)
				t.Errorf("Expected Not To Authentiate with %v : %v, But Did. Got Cookie %v", test.Email, test.Password, w.Header().Get("Set-Cookie"))
			}

			var bodyjson auth.CSRFResponse = auth.CSRFResponse{}

			json.NewDecoder(w.Body).Decode(&bodyjson)

			newreq := httptest.NewRequest(http.MethodGet, "/anything", nil)
			neww := httptest.NewRecorder()

			newreq.Header.Set("X-CSRF-Token", bodyjson.CSRF_TOKEN)

			if w.Header().Get("Set-Cookie") != "" || true {
				if test.Login == true {
					newreq.AddCookie(&http.Cookie{
						Name:  "SESSIONID",
						Value: w.Header().Get("Set-Cookie"),
					})
				}

				VerifyContext := httprouter.Handle(func(w http.ResponseWriter, r *http.Request, pa httprouter.Params) {
					email := r.Context().Value("Authentication")

					if email != test.Expected_Context {
						t.Errorf("Expected %v, Got %v For Expected Context \n", test.Expected_Context, email)
						return
					}

					t.Logf("Correctly Got %v for %v", test.Expected_Context, email)

				})

				sessions.AuthenticationMiddleWare(VerifyContext)(neww, newreq, httprouter.Params{})
			}

		})
	}

	//Returns a httprouter.Handle

	//Figure Out Context in httprouter.Handle

}
