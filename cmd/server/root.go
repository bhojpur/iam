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
	"flag"
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
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var verbose bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "iamsvr",
	Short: "Bhojpur IAMengine is an identity & access management system for distributed enterprise",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if verbose {
			log.SetLevel(log.DebugLevel)
			log.Debug("verbose logging enabled")
		}
		createDatabase := flag.Bool("createDatabase", false, "true if you need Bhojpur IAM to create database")
		flag.Parse()
		object.InitAdapter(*createDatabase)
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

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "en/disable verbose logging")
}
