package controller_test

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
)

type CorsTest struct {
	Name            string
	Allowed_Origins []string
	Request_Origin  string
	Expect          bool
}

var CorsTests []CorsTest = []CorsTest{
	{
		Name:            "Test Without Atserisks But Later * And Correct",
		Allowed_Origins: []string{"Compress", "*"},
		Request_Origin:  "Compressed",
		Expect:          true,
	},
	{
		Name:            "Test Without Atserisks Incorrect",
		Allowed_Origins: []string{"Compress"},
		Request_Origin:  "Compressed",
		Expect:          false,
	},
	{
		Name:            "Test Allow Subdomains",
		Allowed_Origins: []string{"https://*.compressiblelowcalculator.com", "http://*.compressibleflowcalculator.com"},
		Request_Origin:  "https://auth.compressibleflowcalculator.com",
		Expect:          true,
	},
}

func CorsMiddleware(next httprouter.Handle, re *regexp.Regexp) func(http.ResponseWriter, *http.Request, httprouter.Params) {

	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		if r.Header.Get("origin") == "" || r.Header.Get("origin") == r.Header.Get("host") {
			next(w, r, p)
			return
		}

		if re.MatchString(r.Header.Get("origin")) {
			(w).Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
			next(w, r, p)
			return
		}

		next(w, r, p)
	}
}

func TestCorsMiddleware(t *testing.T) {
	//* == [a-bA-B0-9]+

	expr, err := regexp.Compile("(https://.*\\.compressibleflowcalculator\\.com)|(http://.*\\.compressibleflowcalculator\\.com)")

	if err != nil {
		t.Error("Failed To Compile Regex")
	}

	match := expr.FindString("https://auth.compressibleflowcalculator.com")
	if match == "" {
		t.Error("No Matching Strings")
	}

	matches := expr.Match([]byte("https://auth.compressibleflowcalculator.com"))

	if !matches {
		t.Errorf(" DID NOTMATCH")
	}

	for _, test := range CorsTests {

		var stringLiteral = "^" + "(" + test.Allowed_Origins[0] + ")"
		for _, v := range test.Allowed_Origins[1:] {
			stringLiteral = stringLiteral + "|" + "(" + v + ")"
		}

		stringLiteral = strings.Replace(stringLiteral, ".", "\\.", -1)
		stringLiteral = strings.Replace(stringLiteral, "*", ".*", -1)

		stringLiteral = stringLiteral + "$"

		regex, err := regexp.Compile(stringLiteral)
		if err != nil {

		}

		req := httptest.NewRequest(http.MethodGet, "/upper?word=abc", nil)
		w := httptest.NewRecorder()

		req.Header.Add("Origin", test.Request_Origin)

		t.Run(test.Name, func(t *testing.T) {
			CorsMiddleware(func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

				allowed_origin := w.Header().Get("Access-Control-Allow-Origin") == test.Request_Origin

				if allowed_origin != test.Expect {
					t.Errorf("Expexcted %v, Got %v but got %v", test.Expect, test.Request_Origin, w.Header().Get("Access-Control-Allow-Origin"))
				}

			}, regex)(w, req, httprouter.Params{})
		})
	}
}
