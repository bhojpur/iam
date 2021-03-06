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

	"github.com/bhojpur/iam/pkg/object"
	"github.com/bhojpur/iam/pkg/utils"
	ctxsvr "github.com/bhojpur/web/pkg/context"
)

func AutoSigninFilter(ctx *ctxsvr.Context) {
	//if getSessionUser(ctx) != "" {
	//	return
	//}

	// GET parameter like "/page?access_token=123" or
	// HTTP Bearer token like "Authorization: Bearer 123"
	accessToken := ctx.Input.Query("accessToken")
	if accessToken == "" {
		accessToken = parseBearerToken(ctx)
	}
	if accessToken != "" {
		token := object.GetTokenByAccessToken(accessToken)
		if token == nil {
			responseError(ctx, "Access token doesn't exist")
			return
		}

		if utils.IsTokenExpired(token.CreatedTime, token.ExpiresIn) {
			responseError(ctx, "Access token has expired")
			return
		}

		userId := fmt.Sprintf("%s/%s", token.Organization, token.User)
		application, _ := object.GetApplicationByUserId(fmt.Sprintf("app/%s", token.Application))
		setSessionUser(ctx, userId)
		setSessionOidc(ctx, token.Scope, application.ClientId)
		return
	}

	// "/page?clientId=123&clientSecret=456"
	userId := getUsernameByClientIdSecret(ctx)
	if userId != "" {
		setSessionUser(ctx, userId)
		return
	}

	// "/page?username=abc&password=123"
	userId = ctx.Input.Query("username")
	password := ctx.Input.Query("password")
	if userId != "" && password != "" {
		owner, name := utils.GetOwnerAndNameFromId(userId)
		_, msg := object.CheckUserPassword(owner, name, password)
		if msg != "" {
			responseError(ctx, msg)
			return
		}

		setSessionUser(ctx, userId)
		return
	}

}
