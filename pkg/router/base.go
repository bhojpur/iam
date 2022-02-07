package router

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
	"strings"

	"github.com/bhojpur/iam/pkg/object"
	"github.com/bhojpur/iam/pkg/utils"
	ctxsvr "github.com/bhojpur/web/pkg/context"
)

type Response struct {
	Status string      `json:"status"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
	Data2  interface{} `json:"data2"`
}

func responseError(ctx *ctxsvr.Context, error string, data ...interface{}) {
	resp := Response{Status: "error", Msg: error}
	switch len(data) {
	case 2:
		resp.Data2 = data[1]
		fallthrough
	case 1:
		resp.Data = data[0]
	}

	err := ctx.Output.JSON(resp, true, false)
	if err != nil {
		panic(err)
	}
}

func denyRequest(ctx *ctxsvr.Context) {
	responseError(ctx, "Unauthorized operation")
}

func getUsernameByClientIdSecret(ctx *ctxsvr.Context) string {
	clientId, clientSecret, ok := ctx.Request.BasicAuth()
	if !ok {
		clientId = ctx.Input.Query("clientId")
		clientSecret = ctx.Input.Query("clientSecret")
	}

	if clientId == "" || clientSecret == "" {
		return ""
	}

	application := object.GetApplicationByClientId(clientId)
	if application == nil || application.ClientSecret != clientSecret {
		return ""
	}

	return fmt.Sprintf("app/%s", application.Name)
}

func getSessionUser(ctx *ctxsvr.Context) string {
	sessvr := ctx.Input.CruSession
	user := sessvr.Get(nil, "username")
	if user == nil {
		return ""
	}

	return user.(string)
}

func setSessionUser(ctx *ctxsvr.Context, user string) {
	sessvr := ctx.Input.CruSession
	err := sessvr.Set(nil, "username", user)
	if err != nil {
		panic(err)
	}

	sessvr.SessionRelease(nil, ctx.ResponseWriter)
}

func setSessionExpire(ctx *ctxsvr.Context, ExpireTime int64) {
	SessionData := struct{ ExpireTime int64 }{ExpireTime: ExpireTime}
	sessvr := ctx.Input.CruSession
	err := sessvr.Set(nil, "SessionData", utils.StructToJson(SessionData))
	if err != nil {
		panic(err)
	}
	sessvr.SessionRelease(nil, ctx.ResponseWriter)
}

func setSessionOidc(ctx *ctxsvr.Context, scope string, aud string) {
	sessvr := ctx.Input.CruSession
	err := sessvr.Set(nil, "scope", scope)
	if err != nil {
		panic(err)
	}
	err = sessvr.Set(nil, "aud", aud)
	if err != nil {
		panic(err)
	}
	sessvr.SessionRelease(nil, ctx.ResponseWriter)
}

func parseBearerToken(ctx *ctxsvr.Context) string {
	header := ctx.Request.Header.Get("Authorization")
	tokens := strings.Split(header, " ")
	if len(tokens) != 2 {
		return ""
	}

	prefix := tokens[0]
	if prefix != "Bearer" {
		return ""
	}

	return tokens[1]
}
