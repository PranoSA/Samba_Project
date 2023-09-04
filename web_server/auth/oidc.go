package auth

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/PranoSA/samba_share_backend/web_server/models"
	"github.com/julienschmidt/httprouter"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
)

type OIDCAuthenticator struct {
	Keys      jwk.Set
	Issuer    string
	Audience  string
	userModel models.UserModel //Why ???

}

func InitOIDCAuthenticatorFromConfig(c map[interface{}]interface{}) (*OIDCAuthenticator, error) {

	jwks_url, ok := c["JWKS_URL"].(string)
	if !ok {
		return nil, errors.New("Please Define JWKS_URL")
	}

	issuer, ok := c["ISSUER"].(string)
	if !ok {
		return nil, errors.New("Please Define OIDC Issuer")
	}

	audience, ok := c["AUDIENCE"].(string)
	if !ok {
		return nil, errors.New("Please Define Mandatory OIDC Audience ")
	}

	authenticator, err := InitializeOIDCAuthenticator(jwks_url, issuer, audience)
	if err != nil {
		return nil, err
	}

	return authenticator, nil
}

/**
 * Switch to An ARray LAter
 *
 */
func InitializeOIDCAuthenticator(jwks_url string, aud string, iss string) (*OIDCAuthenticator, error) {
	req, err := http.NewRequest("GET", jwks_url, nil)

	client := http.Client{}
	if err != nil {

	}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	//var body []byte
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	set, err := jwk.Parse(body)

	fmt.Println(set)

	return &OIDCAuthenticator{
		Keys:     set,
		Issuer:   iss,
		Audience: aud,
	}, nil

}

func (oidc OIDCAuthenticator) VerifyJwt(token []byte) (*jwt.Token, error) {

	jwt, err := jwt.Parse(token, jwt.WithKeySet(oidc.Keys))
	if err != nil {
		return nil, err
	}

	return &jwt, nil
}

func (oidc OIDCAuthenticator) AuthenticateContext(bearer string) (*jwt.Token, error) {
	token_bytes := []byte(strings.Split(bearer, "Bearer: ")[1])
	jsonwebtoken, err := oidc.VerifyJwt(token_bytes)

	useremail := (*jsonwebtoken).Subject()

	user, new := oidc.userModel.GetUserByIDWithCreate(useremail)
	if new == true {

	}
	if user.Email == "" {

	}

	return jsonwebtoken, err
}

func (oidc OIDCAuthenticator) AuthenticationMiddleWare(next httprouter.Handle) httprouter.Handle {

	return func(w http.ResponseWriter, r *http.Request, pa httprouter.Params) {

		bearer := strings.Split(r.Header.Get("Authorization"), "Bearer ")[1]

		tokenbytes := []byte(bearer)

		jsonwebtoken, err := oidc.VerifyJwt(tokenbytes)

		if err != nil {

		}

		useremail := (*jsonwebtoken).Subject()

		user, _ := oidc.userModel.GetUserByIDWithCreate(useremail)

		ctx := context.WithValue(r.Context(), "Authentication", user.Email)

		r.WithContext(ctx)

		next(w, r, pa)
	}

}
