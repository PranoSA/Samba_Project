package auth_test

import (
	"log"
	"testing"

	"github.com/go-ldap/ldap/v3"
)

func TestLdapConnection(t *testing.T) {

	/*ldap_authenticator := auth.LDAPAuthenticator{
		LDAP_URI: fmt.Sprintf("%v:%v", "localhost", 389),
	}
	*/

	t.Run("LDAP Bind", func(t *testing.T) {

		l, err := ldap.DialURL("ldap://localhost:389")

		if err != nil {
			t.Errorf("Could Not Conect %v", err.Error())
		}

		if err != nil {
			log.Fatal(err)
		}
		defer l.Close()

		//basedn := "dc=ldap,dc=compressibleflowcalculator,dc=com"
		basedn := "uid=prano,ou=people,dc=ldap,dc=compressibleflowcalculator,dc=com"

		//testLogin := fmt.Sprintf("uid=%v,ou=Users,%v", "pranopassword", basedn)

		//err = l.Bind(testLogin, "pranopassword")
		err = l.Bind(basedn, "pranopassword")
		if err != nil {
			t.Errorf("Could Not Bind %v", err)
			log.Fatal(err)
		}

	})

	t.Run("LDAP Search", func(t *testing.T) {

		l, err := ldap.DialURL("ldap://localhost:389")
		if err != nil {
			t.Errorf("Failed To COnnect")
		}

		basedn := "uid=prano,ou=people,dc=ldap,dc=compressibleflowcalculator,dc=com"

		searchRequest := ldap.NewSearchRequest(
			basedn, // The base dn to search
			ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
			"(&(objectClass=posixAccount))",    // The filter to apply
			[]string{"dn", "uid", "uidNumber"}, // A list attributes to retrieve
			nil,
		)

		sr, err := l.Search(searchRequest)

		ent := sr.Entries[0]
		t.Log("#Entries : ", len(sr.Entries))

		ent.GetAttributeValue("dn")

		//dn := searchRequest.Attributes[1]
		t.Errorf("uid is %v \n", ent.GetAttributeValue("uid"))
		t.Errorf("dn is %v \n", ent.DN)

		if ent.GetAttributeValue("dn") != "" {
			t.Errorf("dn is %v \n", ent.GetAttributeValue("dn"))
		}

	})

}
