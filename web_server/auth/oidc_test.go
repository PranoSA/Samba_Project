package auth_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PranoSA/samba_share_backend/web_server/auth"
	"github.com/julienschmidt/httprouter"
)

/**
 *
 * Semi-Level Integration Testing
 *
 * Run Keycloak on a Local Service
 *
 *
 * 1. Create Realm Test
 * 2. JWKS URL = http://localhost:8080/realms/test/protocol/openid-connect/certs
 * 3. Issuer = http://localhost:8080/realms/test
 * 4. Generate Client ID samba
 * 5. Client AUthenticate Off
 * 6. Implicit Flow On (Client Will Use PKCE )
 * 7. CORS Origins -> Wherever you are running test client if must
 * 8. Valid Redirect URIs -> Do localhost:8000
 * 9. aud Will Be samba (Same as client ID
 *
 * 10. Now Configure With an Identity Provider, e.g. Google
 *
 * Now, Curl Request To Login
 *
 * curl -X POST http://localhost:8080/realms/test/protocol/openid-connect/token -d "client_id=samba" -d "redirect_uri=http:/localhost:8080/callback" -d "grant_type=password" -d "password=prano" -d "username=prano" -d "response_type=oidc"
 *
 *
 * curl -X POST http://localhost:8080/realms/test/protocol/openid-connect/token -d "client_id=samba" -d "redirect_uri=http:/localhost:8080/callback" -d "grant_type=password" -d "password=prano" -d "username=prano" -d "response_type=oidc" -d "scope=openid" | jq '.id_token'
 */

var token string = "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICIwelJHMDdmdUplZ3FtRTBDVDNRdi1yei1ULW5DM0lab05PcjhIY19SRXRFIn0.eyJleHAiOjE2OTM5NDQzOTMsImlhdCI6MTY5Mzk0NDA5MywiYXV0aF90aW1lIjowLCJqdGkiOiI2NDU5ZTU5Zi1mMTM1LTRjYmYtOGRiOC0yODE2NmQ2NzkzYjUiLCJpc3MiOiJodHRwOi8vbG9jYWxob3N0OjgwODAvcmVhbG1zL3Rlc3QiLCJhdWQiOiJzYW1iYSIsInN1YiI6IjM2NzIwNTY1LTUzNTEtNGJjMC1hMzc2LTgwNTJjYTcwZjRiNSIsInR5cCI6IklEIiwiYXpwIjoic2FtYmEiLCJzZXNzaW9uX3N0YXRlIjoiN2ViNTZkYzUtMjljYS00MjMyLThiNGQtOGI0ZjE2MjczMDU4IiwiYXRfaGFzaCI6InZQQkFGYkR1amlodGdfd3BRTlNPVGciLCJhY3IiOiIxIiwic2lkIjoiN2ViNTZkYzUtMjljYS00MjMyLThiNGQtOGI0ZjE2MjczMDU4IiwiZW1haWxfdmVyaWZpZWQiOmZhbHNlLCJwcmVmZXJyZWRfdXNlcm5hbWUiOiJwcmFubyIsImdpdmVuX25hbWUiOiIiLCJmYW1pbHlfbmFtZSI6IiJ9.MB0bYkqQfVyo-2GbUNyRk1K-Dt0GW1PuzNc6rMvVXdlFnp7KlwmyIitHjqbg0jnzvI7ukZnf5i679srF2Bl-qjGNM2dCxqFmEEXYZk4eyvNno5NsjPvpgjQInK1f2FyFhwXMLvyFsBqtzLgHyiWclTxKrvB61aKGW84ShTBVpN6gSSCUKjSF8i8ud5Nnm5kLL3JBIT1JDFQxmeJM_zBu-hqvVVyD-3gK3Zume3CGppL_U4GtDo6lcbDrO8n0MT5VxdTbTZEHI73o3GHZNZRctZ7OdTPBs5zxNepOtI3vpB-x8d-iAClEohCOxxIo8428cIqTNsRFZ-A_LhURbw-XBA"

func TestOidcMiddleware(t *testing.T) {

	jwks_url := "http://localhost:8080/realms/test/protocol/openid-connect/certs"
	issuer := "http://localhost:8080/realms/test"
	aud := "samba"

	t.Run("oidc", func(t *testing.T) {

		oidc_authenticator, e := auth.InitializeOIDCAuthenticator(jwks_url, issuer, aud)
		if e != nil {
			t.Fatalf("Invalid AUthenticator %v", e)
		}

		//newreq := httptest.NewRequest(http.MethodGet, "/anything", nil)
		neww := httptest.NewRecorder()

		var body []byte = []byte{24}

		//Mocking an HTTP Response
		req := httptest.NewRequest(http.MethodGet, "/anything", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		VerifyContext := httprouter.Handle(func(w http.ResponseWriter, r *http.Request, pa httprouter.Params) {
			email := r.Context().Value("Authentication")

			if email != "prano" {
				t.Errorf("Expected %v, Got %v For Expected Context \n", "prano", email)
				return
			}

			t.Logf("Correctly Got %v for %v", "prano", email)

		})

		oidc_authenticator.AuthenticationMiddleWare(VerifyContext)(neww, req, httprouter.Params{})
	})

}
