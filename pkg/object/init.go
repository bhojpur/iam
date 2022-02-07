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
	_ "embed"

	"github.com/bhojpur/iam/pkg/utils"
)

//go:embed token_jwt_key.pem
var tokenJwtPublicKey string

//go:embed token_jwt_key.key
var tokenJwtPrivateKey string

func InitDb() {
	initBuiltInOrganization()
	initBuiltInUser()
	initBuiltInApplication()
	initBuiltInCert()
	initBuiltInLdap()
}

func initBuiltInOrganization() {
	organization := getOrganization("admin", "built-in")
	if organization != nil {
		return
	}

	organization = &Organization{
		Owner:         "admin",
		Name:          "built-in",
		CreatedTime:   utils.GetCurrentTime(),
		DisplayName:   "Built-in Organization",
		WebsiteUrl:    "https://bhojpur.net",
		Favicon:       "https://static.bhojpur.net/favicon.ico",
		PhonePrefix:   "91",
		DefaultAvatar: "https://static.bhojpur.net/image/logo.png",
		PasswordType:  "plain",
	}
	AddOrganization(organization)
}

func initBuiltInUser() {
	user := getUser("built-in", "admin")
	if user != nil {
		return
	}

	user = &User{
		Owner:             "built-in",
		Name:              "admin",
		CreatedTime:       utils.GetCurrentTime(),
		Id:                utils.GenerateId(),
		Type:              "normal-user",
		Password:          "123",
		DisplayName:       "Admin",
		Avatar:            "https://bhojpur.net/image/logo.png",
		Email:             "admin@bhojpur.net",
		Phone:             "12345678910",
		Address:           []string{},
		Affiliation:       "Bhojpur Consulting, Inc.",
		Tag:               "staff",
		Score:             2000,
		Ranking:           1,
		IsAdmin:           true,
		IsGlobalAdmin:     true,
		IsForbidden:       false,
		IsDeleted:         false,
		SignupApplication: "built-in-app",
		CreatedIp:         "127.0.0.1",
		Properties:        make(map[string]string),
	}
	AddUser(user)
}

func initBuiltInApplication() {
	application := getApplication("admin", "app-built-in")
	if application != nil {
		return
	}

	application = &Application{
		Owner:          "admin",
		Name:           "app-built-in",
		CreatedTime:    utils.GetCurrentTime(),
		DisplayName:    "Bhojpur IAM",
		Logo:           "https://static.bhojpur.net/image/logo.png",
		HomepageUrl:    "https://www.bhojpur.net",
		Organization:   "built-in",
		Cert:           "cert-built-in",
		EnablePassword: true,
		EnableSignUp:   true,
		Providers:      []*ProviderItem{},
		SignupItems: []*SignupItem{
			{Name: "ID", Visible: false, Required: true, Prompted: false, Rule: "Random"},
			{Name: "Username", Visible: true, Required: true, Prompted: false, Rule: "None"},
			{Name: "Display name", Visible: true, Required: true, Prompted: false, Rule: "None"},
			{Name: "Password", Visible: true, Required: true, Prompted: false, Rule: "None"},
			{Name: "Confirm password", Visible: true, Required: true, Prompted: false, Rule: "None"},
			{Name: "Email", Visible: true, Required: true, Prompted: false, Rule: "None"},
			{Name: "Phone", Visible: true, Required: true, Prompted: false, Rule: "None"},
			{Name: "Agreement", Visible: true, Required: true, Prompted: false, Rule: "None"},
		},
		RedirectUris:  []string{},
		ExpireInHours: 168,
	}
	AddApplication(application)
}

func initBuiltInCert() {
	cert := getCert("admin", "cert-built-in")
	if cert != nil {
		return
	}

	cert = &Cert{
		Owner:           "admin",
		Name:            "cert-built-in",
		CreatedTime:     utils.GetCurrentTime(),
		DisplayName:     "Built-in Cert",
		Scope:           "JWT",
		Type:            "x509",
		CryptoAlgorithm: "RSA",
		BitSize:         4096,
		ExpireInYears:   20,
		PublicKey:       tokenJwtPublicKey,
		PrivateKey:      tokenJwtPrivateKey,
	}
	AddCert(cert)
}

func initBuiltInLdap() {
	ldap := GetLdap("ldap-built-in")
	if ldap != nil {
		return
	}

	ldap = &Ldap{
		Id:         "ldap-built-in",
		Owner:      "built-in",
		ServerName: "BuildIn LDAP Server",
		Host:       "bhojpure.net",
		Port:       389,
		Admin:      "cn=buildin,dc=bhojpur,dc=net",
		Passwd:     "123",
		BaseDn:     "ou=BuildIn,dc=bhojpur,dc=net",
		AutoSync:   0,
		LastSync:   "",
	}
	AddLdap(ldap)
}
