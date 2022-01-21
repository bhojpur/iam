package controllers

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
	"encoding/json"

	"github.com/bhojpur/iam/pkg/object"
	"github.com/bhojpur/iam/pkg/utils"
	pagination "github.com/bhojpur/web/pkg/utils/pagination"
)

// GetTokens
// @Title GetTokens
// @Tag Token API
// @Description get tokens
// @Param   owner     query    string  true        "The owner of tokens"
// @Param   pageSize     query    string  true        "The size of each page"
// @Param   p     query    string  true        "The number of the page"
// @Success 200 {array} object.Token The Response object
// @router /get-tokens [get]
func (c *ApiController) GetTokens() {
	owner := c.Input().Get("owner")
	limit := c.Input().Get("pageSize")
	page := c.Input().Get("p")
	field := c.Input().Get("field")
	value := c.Input().Get("value")
	sortField := c.Input().Get("sortField")
	sortOrder := c.Input().Get("sortOrder")
	if limit == "" || page == "" {
		c.Data["json"] = object.GetTokens(owner)
		c.ServeJSON()
	} else {
		limit := utils.ParseInt(limit)
		paginator := pagination.SetPaginator(c.Ctx, limit, int64(object.GetTokenCount(owner, field, value)))
		tokens := object.GetPaginationTokens(owner, paginator.Offset(), limit, field, value, sortField, sortOrder)
		c.ResponseOk(tokens, paginator.Nums())
	}
}

// GetToken
// @Title GetToken
// @Tag Token API
// @Description get token
// @Param   id     query    string  true        "The id of token"
// @Success 200 {object} object.Token The Response object
// @router /get-token [get]
func (c *ApiController) GetToken() {
	id := c.Input().Get("id")

	c.Data["json"] = object.GetToken(id)
	c.ServeJSON()
}

// UpdateToken
// @Title UpdateToken
// @Tag Token API
// @Description update token
// @Param   id     query    string  true        "The id of token"
// @Param   body    body   object.Token  true        "Details of the token"
// @Success 200 {object} controllers.Response The Response object
// @router /update-token [post]
func (c *ApiController) UpdateToken() {
	id := c.Input().Get("id")

	var token object.Token
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &token)
	if err != nil {
		panic(err)
	}

	c.Data["json"] = wrapActionResponse(object.UpdateToken(id, &token))
	c.ServeJSON()
}

// AddToken
// @Title AddToken
// @Tag Token API
// @Description add token
// @Param   body    body   object.Token  true        "Details of the token"
// @Success 200 {object} controllers.Response The Response object
// @router /add-token [post]
func (c *ApiController) AddToken() {
	var token object.Token
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &token)
	if err != nil {
		panic(err)
	}

	c.Data["json"] = wrapActionResponse(object.AddToken(&token))
	c.ServeJSON()
}

// DeleteToken
// @Tag Token API
// @Title DeleteToken
// @Description delete token
// @Param   body    body   object.Token  true        "Details of the token"
// @Success 200 {object} controllers.Response The Response object
// @router /delete-token [post]
func (c *ApiController) DeleteToken() {
	var token object.Token
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &token)
	if err != nil {
		panic(err)
	}

	c.Data["json"] = wrapActionResponse(object.DeleteToken(&token))
	c.ServeJSON()
}

// GetOAuthCode
// @Title GetOAuthCode
// @Tag Token API
// @Description get OAuth code
// @Param   user_id     query    string  true        "The id of user"
// @Param   client_id     query    string  true        "OAuth client id"
// @Param   response_type     query    string  true        "OAuth response type"
// @Param   redirect_uri     query    string  true        "OAuth redirect URI"
// @Param   scope     query    string  true        "OAuth scope"
// @Param   state     query    string  true        "OAuth state"
// @Success 200 {object} object.TokenWrapper The Response object
// @router /login/oauth/code [post]
func (c *ApiController) GetOAuthCode() {
	userId := c.Input().Get("user_id")
	clientId := c.Input().Get("client_id")
	responseType := c.Input().Get("response_type")
	redirectUri := c.Input().Get("redirect_uri")
	scope := c.Input().Get("scope")
	state := c.Input().Get("state")
	nonce := c.Input().Get("nonce")

	challengeMethod := c.Input().Get("code_challenge_method")
	codeChallenge := c.Input().Get("code_challenge")

	if challengeMethod != "S256" && challengeMethod != "null" {
		c.ResponseError("Challenge method should be S256")
		return
	}

	c.Data["json"] = object.GetOAuthCode(userId, clientId, responseType, redirectUri, scope, state, nonce, codeChallenge)
	c.ServeJSON()
}

// GetOAuthToken
// @Title GetOAuthToken
// @Tag Token API
// @Description get OAuth access token
// @Param   grant_type     query    string  true        "OAuth grant type"
// @Param   client_id     query    string  true        "OAuth client id"
// @Param   client_secret     query    string  true        "OAuth client secret"
// @Param   code     query    string  true        "OAuth code"
// @Success 200 {object} object.TokenWrapper The Response object
// @router /login/oauth/access_token [post]
func (c *ApiController) GetOAuthToken() {
	grantType := c.Input().Get("grant_type")
	clientId := c.Input().Get("client_id")
	clientSecret := c.Input().Get("client_secret")
	code := c.Input().Get("code")
	verifier := c.Input().Get("code_verifier")

	if clientId == "" && clientSecret == "" {
		clientId, clientSecret, _ = c.Ctx.Request.BasicAuth()
	}

	c.Data["json"] = object.GetOAuthToken(grantType, clientId, clientSecret, code, verifier)
	c.ServeJSON()
}

// RefreshToken
// @Title RefreshToken
// @Description refresh OAuth access token
// @Param   grant_type     query    string  true        "OAuth grant type"
// @Param	refresh_token	query	string	true		"OAuth refresh token"
// @Param   scope     query    string  true        "OAuth scope"
// @Param   client_id     query    string  true        "OAuth client id"
// @Param   client_secret     query    string  true        "OAuth client secret"
// @Success 200 {object} object.TokenWrapper The Response object
// @router /login/oauth/refresh_token [post]
func (c *ApiController) RefreshToken() {
	grantType := c.Input().Get("grant_type")
	refreshToken := c.Input().Get("refresh_token")
	scope := c.Input().Get("scope")
	clientId := c.Input().Get("client_id")
	clientSecret := c.Input().Get("client_secret")

	c.Data["json"] = object.RefreshToken(grantType, refreshToken, scope, clientId, clientSecret)
	c.ServeJSON()
}
