package auth

import (
	"context"
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
	keys      jwk.Set
	userModel models.UserModel
}

func InitOIDCAuthenticator(jwks_url string, aud []string, iss []string) {
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

}

func (oidc OIDCAuthenticator) VerifyJwt(token []byte) (*jwt.Token, error) {

	jwt, err := jwt.Parse(token, jwt.WithKeySet(oidc.keys))
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

func (oidc OIDCAuthenticator) AuthenticationMiddleWare(w http.ResponseWriter, r *http.Request, pr httprouter.Params) httprouter.Handle {

	bearer := strings.Split(r.Header.Get("Authorization"), "Bearer ")[1]

	tokenbytes := []byte(bearer)

	jsonwebtoken, err := oidc.VerifyJwt(tokenbytes)

	if err != nil {

	}

	useremail := (*jsonwebtoken).Subject()

	user, _ := oidc.userModel.GetUserByIDWithCreate(useremail)

	ctx := context.WithValue(r.Context(), "Authentication", user.Email)

	r.WithContext(ctx)

	return func(w http.ResponseWriter, r *http.Request, pr httprouter.Params) {

	}

}
