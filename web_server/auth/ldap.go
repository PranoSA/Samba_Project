package auth

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"

	"github.com/go-ldap/ldap/v3"
)

/**
 *
 *	THE LDAP AUTHENTICATOR WILL BE INITIALIZED WITH A
 * BASE DN, URL,
 *
 * LDAP WILL IMPLEMENT BOTH SIGN-UP + SIGN-ON, and SESSIONS
 *
 */
type LDAPAuthenticator struct {
	LDAP_URI string
	BASE_DN  string
}

func (la LDAPAuthenticator) DialLDAP() {

	ldapCert := "/path/to/cert.pem"
	ldapKey := "/path/to/key.pem"
	ldapCAchain := "/path/to/ca_chain.pem"

	// Load client cert and key
	cert, err := tls.LoadX509KeyPair(ldapCert, ldapKey)
	if err != nil {
		log.Fatal(err)
	}

	// Load CA chain
	caCert, err := os.ReadFile(ldapCAchain)
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Setup TLS with ldap client cert
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            caCertPool,
		InsecureSkipVerify: true,

		//ClientAuth: ,
	}

	conn, err := ldap.DialURL(la.LDAP_URI)
	if err != nil {
		log.Fatalf(" Failed LDAP Connection %v", err.Error())

	}

	defer conn.Close()

	fmt.Println(tlsConfig.Certificates)

	//conn.StartTLS(tlsConfig)

	conn.Start()

	/*searchRequest := NewSearchRequest(
		la.BASE_DN, // The base dn to search
		ScopeWholeSubtree, NeverDerefAliases, 0, 0, false,
		"(&(objectClass=organizationalPerson))", // The filter to apply
		[]string{"dn", "cn"},                    // A list attributes to retrieve
		nil,
	)*/

}
