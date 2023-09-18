package main

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
)

//curl -X POST http://localhost:8080/realms/test/protocol/openid-connect/token -d "client_id=samba" -d "redirect_uri=http:/localhost:8080/callback" -d "grant_type=password" -d "password=prano" -d "username=prano" -d "response_type=oidc" -d "scope=openid" | jq '.id_token'

/**
 *
 * Localhost Keycloak
 * 
 *     apiUrl := "https://api.com"
    resource := "/user/"
    data := url.Values{}
    data.Set("name", "foo")
    data.Set("surname", "bar")

    u, _ := url.ParseRequestURI(apiUrl)
    u.Path = resource
    urlStr := u.String() // "https://api.com/user/"
 */

var TestRequestsCreateShare []struct{	
}

func TestRest(t *testing.T) {
	realm := "samba"
	auth_url := fmt.Sprintf("http://localhost:8080/realms/%s/protocol/openid-connect/token", realm)


	request_data := url.Values{}

	u, _ := url.ParseRequestURI(req)

	url.QueryEscape(s)

	http.NewRequest(, url, body)

	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/realms/samba/protocol/openid-connect/token")
}
