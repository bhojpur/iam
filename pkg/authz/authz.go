package authz

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
	"github.com/bhojpur/iam/conf"
	plcsvr "github.com/bhojpur/policy/pkg/engine"
	model "github.com/bhojpur/policy/pkg/model"
	ormadapter "github.com/bhojpur/policy/pkg/persist/orm-adapter"
	stringadapter "github.com/bhojpur/policy/pkg/persist/string-adapter"
	websvr "github.com/bhojpur/web/pkg/engine"
)

var Enforcer *plcsvr.Enforcer

func InitAuthz() {
	var err error

	driverName, err := websvr.AppConfig.String("driverName")
	tableNamePrefix, err := websvr.AppConfig.String("tableNamePrefix")
	dbName, err := websvr.AppConfig.String("dbName")
	a, err := ormadapter.NewAdapterWithTableName(driverName, conf.GetBhojpurConfDataSourceName()+dbName, "iam_rule", tableNamePrefix, true)
	if err != nil {
		panic(err)
	}

	modelText := `
[request_definition]
r = subOwner, subName, method, urlPath, objOwner, objName

[policy_definition]
p = subOwner, subName, method, urlPath, objOwner, objName

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = (r.subOwner == p.subOwner || p.subOwner == "*") && \
    (r.subName == p.subName || p.subName == "*" || r.subName != "anonymous" && p.subName == "!anonymous") && \
    (r.method == p.method || p.method == "*") && \
    (r.urlPath == p.urlPath || p.urlPath == "*") && \
    (r.objOwner == p.objOwner || p.objOwner == "*") && \
    (r.objName == p.objName || p.objName == "*") || \
    (r.urlPath == "/api/update-user" && r.subOwner == r.objOwner && r.subName == r.objName)
`

	m, err := model.NewModelFromString(modelText)
	if err != nil {
		panic(err)
	}

	Enforcer, err = plcsvr.NewEnforcer(m, a)
	if err != nil {
		panic(err)
	}

	Enforcer.ClearPolicy()

	//if len(Enforcer.GetPolicy()) == 0 {
	if true {
		ruleText := `
p, built-in, *, *, *, *, *
p, app, *, *, *, *, *
p, *, *, POST, /api/signup, *, *
p, *, *, POST, /api/get-email-and-phone, *, *
p, *, *, POST, /api/login, *, *
p, *, *, GET, /api/get-app-login, *, *
p, *, *, POST, /api/logout, *, *
p, *, *, GET, /api/get-account, *, *
p, *, *, POST, /api/login/oauth/access_token, *, *
p, *, *, POST, /api/login/oauth/refresh_token, *, *
p, *, *, GET, /api/get-application, *, *
p, *, *, GET, /api/get-users, *, *
p, *, *, GET, /api/get-user, *, *
p, *, *, GET, /api/get-organizations, *, *
p, *, *, GET, /api/get-user-application, *, *
p, *, *, GET, /api/get-default-providers, *, *
p, *, *, GET, /api/get-resources, *, *
p, *, *, POST, /api/upload-avatar, *, *
p, *, *, POST, /api/unlink, *, *
p, *, *, POST, /api/set-password, *, *
p, *, *, POST, /api/send-verification-code, *, *
p, *, *, GET, /api/get-human-check, *, *
p, *, *, POST, /api/reset-email-or-phone, *, *
p, *, *, POST, /api/upload-resource, *, *
p, *, *, GET, /.well-known/openid-configuration, *, *
p, *, *, *, /api/certs, *, *
p, *, *, GET, /api/get-saml-login, *, *
p, *, *, POST, /api/acs, *, *
`

		sa := stringadapter.NewAdapter(ruleText)
		// load all rules from string adapter to enforcer's memory
		err := sa.LoadPolicy(Enforcer.GetModel())
		if err != nil {
			panic(err)
		}

		// save all rules from enforcer's memory to ORM adapter (DB)
		// same as:
		// a.SavePolicy(Enforcer.GetModel())
		err = Enforcer.SavePolicy()
		if err != nil {
			panic(err)
		}
	}
}

func IsAllowed(subOwner string, subName string, method string, urlPath string, objOwner string, objName string) bool {
	res, err := Enforcer.Enforce(subOwner, subName, method, urlPath, objOwner, objName)
	if err != nil {
		panic(err)
	}

	return res
}
