package engine

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
	"flag"

	"github.com/bhojpur/iam/pkg/authz"
	"github.com/bhojpur/iam/pkg/object"
	"github.com/bhojpur/iam/pkg/proxy"
	"github.com/bhojpur/iam/pkg/router"
	_ "github.com/bhojpur/iam/pkg/router"
	logs "github.com/bhojpur/logger/pkg/engine"
	websvr "github.com/bhojpur/web/pkg/engine"
	"github.com/bhojpur/web/pkg/plugins/cors"
	_ "github.com/bhojpur/web/pkg/session/redis"
)

func main() {
	createDatabase := flag.Bool("createDatabase", false, "true if you need Bhojpur IAM to create database")
	flag.Parse()
	object.InitAdapter(*createDatabase)
	object.InitDb()
	object.InitDefaultStorageProvider()
	object.InitLdapAutoSynchronizer()
	proxy.InitHttpClient()
	authz.InitAuthz()

	go object.RunSyncUsersJob()

	websvr.InsertFilter("*", websvr.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "PUT", "PATCH"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	//websvr.DelStaticPath("/static")
	websvr.SetStaticPath("/static", "web/build/static")
	websvr.BConfig.WebConfig.DirectoryIndex = true
	websvr.SetStaticPath("/swagger", "swagger")
	websvr.SetStaticPath("/files", "files")
	// https://studygolang.com/articles/2303
	websvr.InsertFilter("*", websvr.BeforeRouter, router.StaticFilter)
	websvr.InsertFilter("*", websvr.BeforeRouter, router.AutoSigninFilter)
	websvr.InsertFilter("*", websvr.BeforeRouter, router.AuthzFilter)
	websvr.InsertFilter("*", websvr.BeforeRouter, router.RecordMessage)

	websvr.BConfig.WebConfig.Session.SessionName = "bhojpur_iam_session_id"
	if websvr.AppConfig.String("redisEndpoint") == "" {
		websvr.BConfig.WebConfig.Session.SessionProvider = "file"
		websvr.BConfig.WebConfig.Session.SessionProviderConfig = "./tmp"
	} else {
		websvr.BConfig.WebConfig.Session.SessionProvider = "redis"
		websvr.BConfig.WebConfig.Session.SessionProviderConfig = websvr.AppConfig.String("redisEndpoint")
	}
	websvr.BConfig.WebConfig.Session.SessionCookieLifeTime = 3600 * 24 * 30
	//websvr.BConfig.WebConfig.Session.SessionCookieSameSite = http.SameSiteNoneMode

	err := logs.SetLogger("file", `{"filename":"logs/bhojpur_iam.log","maxdays":99999,"perm":"0770"}`)
	if err != nil {
		panic(err)
	}
	//logs.SetLevel(logs.LevelInformational)
	logs.SetLogFuncCall(false)

	websvr.Run()
}