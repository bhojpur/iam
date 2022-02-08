package cmd

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
	"fmt"
	"os"

	"github.com/bhojpur/iam/pkg/authz"
	"github.com/bhojpur/iam/pkg/object"
	"github.com/bhojpur/iam/pkg/proxy"
	"github.com/bhojpur/iam/pkg/router"
	_ "github.com/bhojpur/iam/pkg/router"
	logsvr "github.com/bhojpur/logger/pkg/engine"
	_ "github.com/bhojpur/session/pkg/provider/redis"
	websvr "github.com/bhojpur/web/pkg/engine"
	"github.com/spf13/cobra"
)

var (
	createDatabase bool
)

var webCmdOpts struct {
	Host           string
	CreateDatabase bool
}

// webCmd represents the web command
var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Manage web services API backend of Bhojpur IAM",
	Run: func(cmd *cobra.Command, args []string) {
		object.InitAdapter(createDatabase)
		object.InitDb()
		object.InitDefaultStorageProvider()
		object.InitLdapAutoSynchronizer()
		proxy.InitHttpClient()
		authz.InitAuthz()

		go object.RunSyncUsersJob()

		//websvr.DelStaticPath("/static")
		websvr.SetStaticPath("/static", "pkg/webui/build/static")
		websvr.BConfig.WebConfig.DirectoryIndex = true
		websvr.SetStaticPath("/swagger", "swagger")
		websvr.SetStaticPath("/files", "files")
		websvr.InsertFilter("*", websvr.BeforeRouter, router.StaticFilter)
		websvr.InsertFilter("*", websvr.BeforeRouter, router.AutoSigninFilter)
		websvr.InsertFilter("*", websvr.BeforeRouter, router.AuthzFilter)
		websvr.InsertFilter("*", websvr.BeforeRouter, router.RecordMessage)

		websvr.BConfig.WebConfig.Session.SessionName = "bhojpur_session_id"
		redisEndpoint, err := websvr.AppConfig.String("redisEndpoint")
		if err != nil {
			fmt.Errorf("redisEndpoint", err)
		}
		if redisEndpoint == "" {
			websvr.BConfig.WebConfig.Session.SessionProvider = "file"
			websvr.BConfig.WebConfig.Session.SessionProviderConfig = "./tmp"
		} else {
			websvr.BConfig.WebConfig.Session.SessionProvider = "redis"
			websvr.BConfig.WebConfig.Session.SessionProviderConfig = redisEndpoint
		}
		websvr.BConfig.WebConfig.Session.SessionCookieLifeTime = 3600 * 24 * 30
		//websvr.BConfig.WebConfig.Session.SessionCookieSameSite = http.SameSiteNoneMode

		err = logsvr.SetLogger("file", `{"filename":"logs/bhojpur_iam.log","maxdays":99999,"perm":"0770"}`)
		if err != nil {
			panic(err)
		}
		//logsvr.SetLevel(logs.LevelInformational)
		logsvr.SetLogFuncCall(false)

		websvr.Run()
	},
}

func init() {
	rootCmd.AddCommand(webCmd)

	iamHost := os.Getenv("IAM_HOST")
	if iamHost == "" {
		iamHost = "localhost:8000"
	}
	createDB := os.Getenv("IAM_CREATE_DATABASE")
	if createDB == "" {
		createDatabase = false
	}
	webCmd.PersistentFlags().StringVar(&webCmdOpts.Host, "host", iamHost, "[host address] Bhojpur IAM host address (defaults to IAM_HOST env var)")
	webCmd.PersistentFlags().BoolVar(&webCmdOpts.CreateDatabase, "createDatabase", createDatabase, "true if you need Bhojpur IAM to create database")
}
