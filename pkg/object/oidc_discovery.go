package object

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"

	websvr "github.com/bhojpur/web/pkg/engine"
	"gopkg.in/square/go-jose.v2"
)

type OidcDiscovery struct {
	Issuer                                 string   `json:"issuer"`
	AuthorizationEndpoint                  string   `json:"authorization_endpoint"`
	TokenEndpoint                          string   `json:"token_endpoint"`
	UserinfoEndpoint                       string   `json:"userinfo_endpoint"`
	JwksUri                                string   `json:"jwks_uri"`
	ResponseTypesSupported                 []string `json:"response_types_supported"`
	ResponseModesSupported                 []string `json:"response_modes_supported"`
	GrantTypesSupported                    []string `json:"grant_types_supported"`
	SubjectTypesSupported                  []string `json:"subject_types_supported"`
	IdTokenSigningAlgValuesSupported       []string `json:"id_token_signing_alg_values_supported"`
	ScopesSupported                        []string `json:"scopes_supported"`
	ClaimsSupported                        []string `json:"claims_supported"`
	RequestParameterSupported              bool     `json:"request_parameter_supported"`
	RequestObjectSigningAlgValuesSupported []string `json:"request_object_signing_alg_values_supported"`
}

var oidcDiscovery OidcDiscovery

func init() {
	origin := websvr.AppConfig.String("origin")

	// Examples:
	// https://login.okta.com/.well-known/openid-configuration
	// https://auth0.auth0.com/.well-known/openid-configuration
	// https://accounts.google.com/.well-known/openid-configuration
	// https://access.line.me/.well-known/openid-configuration
	oidcDiscovery = OidcDiscovery{
		Issuer:                                 origin,
		AuthorizationEndpoint:                  fmt.Sprintf("%s/login/oauth/authorize", origin),
		TokenEndpoint:                          fmt.Sprintf("%s/api/login/oauth/access_token", origin),
		UserinfoEndpoint:                       fmt.Sprintf("%s/api/get-account", origin),
		JwksUri:                                fmt.Sprintf("%s/api/certs", origin),
		ResponseTypesSupported:                 []string{"id_token"},
		ResponseModesSupported:                 []string{"login", "code", "link"},
		GrantTypesSupported:                    []string{"password", "authorization_code"},
		SubjectTypesSupported:                  []string{"public"},
		IdTokenSigningAlgValuesSupported:       []string{"RS256"},
		ScopesSupported:                        []string{"openid", "email", "profile", "address", "phone", "offline_access"},
		ClaimsSupported:                        []string{"iss", "ver", "sub", "aud", "iat", "exp", "id", "type", "displayName", "avatar", "permanentAvatar", "email", "phone", "location", "affiliation", "title", "homepage", "bio", "tag", "region", "language", "score", "ranking", "isOnline", "isAdmin", "isGlobalAdmin", "isForbidden", "signupApplication", "ldap"},
		RequestParameterSupported:              true,
		RequestObjectSigningAlgValuesSupported: []string{"HS256", "HS384", "HS512"},
	}
}

func GetOidcDiscovery() OidcDiscovery {
	return oidcDiscovery
}

func GetJsonWebKeySet() (jose.JSONWebKeySet, error) {
	cert := GetDefaultCert()

	//follows the protocol rfc 7517(draft)
	//link here: https://self-issued.info/docs/draft-ietf-jose-json-web-key.html
	//or https://datatracker.ietf.org/doc/html/draft-ietf-jose-json-web-key
	certPemBlock := []byte(cert.PublicKey)
	certDerBlock, _ := pem.Decode(certPemBlock)
	x509Cert, _ := x509.ParseCertificate(certDerBlock.Bytes)

	var jwk jose.JSONWebKey
	jwk.Key = x509Cert.PublicKey
	jwk.Certificates = []*x509.Certificate{x509Cert}
	jwk.KeyID = cert.Name

	var jwks jose.JSONWebKeySet
	jwks.Keys = []jose.JSONWebKey{jwk}
	return jwks, nil
}
