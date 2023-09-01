package auth

import (
	"encoding/base64"
	"errors"
	"strings"
)

// We Need to Pass it An Sign-Up , Sign-IN Authenticator Essentially
// EIther LDAP or POSTGRES or KERBEROS + LDAP IN THIS INSTANCE
type SimpleSessions struct {
	auth SignUpOnAuthentication
}

type UserInfo struct {
}

// For Now, Just Return Uername
func (s SimpleSessions) AuthenticateContext(header string) (string, error) {
	strings.Split(header, "Bearer: ")
	token := strings.Split(header, "Bearer: ")[1]

	dest, err := base64.StdEncoding.DecodeString(token)
	if err != nil {

	}

	auth_token := string(dest)

	username := strings.Split(auth_token, ":")[0]
	password := strings.Split(auth_token, ":")[1]

	//Call s.auth method to authenticate
	ok := s.auth.Login(username, password)

	if ok {
		return username, nil
	}

	return "", errors.New("Invalid AUthentication")

}
