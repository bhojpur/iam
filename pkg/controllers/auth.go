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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/bhojpur/iam/pkg/idp"
	"github.com/bhojpur/iam/pkg/object"
	"github.com/bhojpur/iam/pkg/proxy"
	"github.com/bhojpur/iam/pkg/utils"
	websvr "github.com/bhojpur/web/pkg/engine"
)

func codeToResponse(code *object.Code) *Response {
	if code.Code == "" {
		return &Response{Status: "error", Msg: code.Message, Data: code.Code}
	}

	return &Response{Status: "ok", Msg: "", Data: code.Code}
}

// HandleLoggedIn ...
func (c *ApiController) HandleLoggedIn(application *object.Application, user *object.User, form *RequestForm) (resp *Response) {
	userId := user.GetId()
	if form.Type == ResponseTypeLogin {
		c.SetSessionUsername(userId)
		utils.LogInfo(c.Ctx, "API: [%s] signed in", userId)
		resp = &Response{Status: "ok", Msg: "", Data: userId}
	} else if form.Type == ResponseTypeCode {
		clientId := c.Input().Get("clientId")
		responseType := c.Input().Get("responseType")
		redirectUri := c.Input().Get("redirectUri")
		scope := c.Input().Get("scope")
		state := c.Input().Get("state")
		nonce := c.Input().Get("nonce")
		challengeMethod := c.Input().Get("code_challenge_method")
		codeChallenge := c.Input().Get("code_challenge")

		if challengeMethod != "S256" && challengeMethod != "null" {
			c.ResponseError("Challenge method should be S256")
			return
		}
		code := object.GetOAuthCode(userId, clientId, responseType, redirectUri, scope, state, nonce, codeChallenge)
		resp = codeToResponse(code)

		if application.EnableSigninSession || application.HasPromptPage() {
			// The prompt page needs the user to be signed in
			c.SetSessionUsername(userId)
		}
	} else {
		resp = &Response{Status: "error", Msg: fmt.Sprintf("Unknown response type: %s", form.Type)}
	}

	// if user did not check auto signin
	if resp.Status == "ok" && !form.AutoSignin {
		timestamp := time.Now().Unix()
		timestamp += 3600 * 24
		c.SetSessionData(&SessionData{
			ExpireTime: timestamp,
		})
	}

	return resp
}

// GetApplicationLogin ...
// @Title GetApplicationLogin
// @Tag Login API
// @Description get application login
// @Param   clientId    query    string  true        "client id"
// @Param   responseType    query    string  true        "response type"
// @Param   redirectUri    query    string  true        "redirect uri"
// @Param   scope    query    string  true        "scope"
// @Param   state    query    string  true        "state"
// @Success 200 {object} controllers.api_controller.Response The Response object
// @router /update-application [get]
func (c *ApiController) GetApplicationLogin() {
	clientId := c.Input().Get("clientId")
	responseType := c.Input().Get("responseType")
	redirectUri := c.Input().Get("redirectUri")
	scope := c.Input().Get("scope")
	state := c.Input().Get("state")

	msg, application := object.CheckOAuthLogin(clientId, responseType, redirectUri, scope, state)
	if msg != "" {
		c.ResponseError(msg, application)
	} else {
		c.ResponseOk(application)
	}
}

func setHttpClient(idProvider idp.IdProvider, providerType string) {
	if providerType == "GitHub" || providerType == "Google" || providerType == "Facebook" || providerType == "LinkedIn" {
		idProvider.SetHttpClient(proxy.ProxyHttpClient)
	} else {
		idProvider.SetHttpClient(proxy.DefaultHttpClient)
	}
}

// Login ...
// @Title Login
// @Tag Login API
// @Description login
// @Param   oAuthParams     query    string  true        "oAuth parameters"
// @Param   body    body   RequestForm  true        "Login information"
// @Success 200 {object} controllers.api_controller.Response The Response object
// @router /login [post]
func (c *ApiController) Login() {
	resp := &Response{}

	var form RequestForm
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &form)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	if form.Username != "" {
		if form.Type == ResponseTypeLogin {
			if c.GetSessionUsername() != "" {
				c.ResponseError("Please sign out first before signing in", c.GetSessionUsername())
				return
			}
		}

		var user *object.User
		var msg string

		if form.Password == "" {
			var verificationCodeType string
			var checkResult string

			// check result through Email or Phone
			if strings.Contains(form.Username, "@") {
				verificationCodeType = "email"
				checkResult = object.CheckVerificationCode(form.Username, form.Code)
			} else {
				verificationCodeType = "phone"
				if len(form.PhonePrefix) == 0 {
					responseText := fmt.Sprintf("%s%s", verificationCodeType, "No phone prefix")
					c.ResponseError(responseText)
					return
				}
				checkPhone := fmt.Sprintf("+%s%s", form.PhonePrefix, form.Username)
				checkResult = object.CheckVerificationCode(checkPhone, form.Code)
			}
			if len(checkResult) != 0 {
				responseText := fmt.Sprintf("%s%s", verificationCodeType, checkResult)
				c.ResponseError(responseText)
				return
			}

			// disable the verification code
			object.DisableVerificationCode(form.Username)

			user = object.GetUserByFields(form.Organization, form.Username)
			if user == nil {
				c.ResponseError("No such user.")
				return
			}
		} else {
			password := form.Password
			user, msg = object.CheckUserPassword(form.Organization, form.Username, password)
		}

		if msg != "" {
			resp = &Response{Status: "error", Msg: msg}
		} else {
			application := object.GetApplication(fmt.Sprintf("admin/%s", form.Application))
			resp = c.HandleLoggedIn(application, user, &form)

			record := object.NewRecord(c.Ctx)
			record.Organization = application.Organization
			record.User = user.Name
			go object.AddRecord(record)
		}
	} else if form.Provider != "" {
		application := object.GetApplication(fmt.Sprintf("admin/%s", form.Application))
		organization := object.GetOrganization(fmt.Sprintf("%s/%s", "admin", application.Organization))
		provider := object.GetProvider(fmt.Sprintf("admin/%s", form.Provider))
		providerItem := application.GetProviderItem(provider.Name)
		if !providerItem.IsProviderVisible() {
			c.ResponseError(fmt.Sprintf("The provider: %s is not enabled for the application", provider.Name))
			return
		}

		userInfo := &idp.UserInfo{}
		if provider.Category == "SAML" {
			// SAML
			userInfo.Id, err = object.ParseSamlResponse(form.SamlResponse, provider.Type)
			if err != nil {
				c.ResponseError(err.Error())
				return
			}
		} else if provider.Category == "OAuth" {
			// OAuth

			clientId := provider.ClientId
			clientSecret := provider.ClientSecret
			if provider.Type == "WeChat" && strings.Contains(c.Ctx.Request.UserAgent(), "MicroMessenger") {
				clientId = provider.ClientId2
				clientSecret = provider.ClientSecret2
			}

			idProvider := idp.GetIdProvider(provider.Type, clientId, clientSecret, form.RedirectUri)
			if idProvider == nil {
				c.ResponseError(fmt.Sprintf("The provider type: %s is not supported", provider.Type))
				return
			}

			setHttpClient(idProvider, provider.Type)

			if form.State != websvr.AppConfig.String("authState") && form.State != application.Name {
				c.ResponseError(fmt.Sprintf("state expected: \"%s\", but got: \"%s\"", websvr.AppConfig.String("authState"), form.State))
				return
			}

			token, err := idProvider.GetToken(form.Code)
			if err != nil {
				c.ResponseError(err.Error())
				return
			}

			if !token.Valid() {
				c.ResponseError("Invalid token")
				return
			}

			userInfo, err = idProvider.GetUserInfo(token)
			if err != nil {
				c.ResponseError(fmt.Sprintf("Failed to login in: %s", err.Error()))
				return
			}
		}

		if form.Method == "signup" {
			user := &object.User{}
			if provider.Category == "SAML" {
				user = object.GetUser(fmt.Sprintf("%s/%s", application.Organization, userInfo.Id))
			} else if provider.Category == "OAuth" {
				user = object.GetUserByField(application.Organization, provider.Type, userInfo.Id)
				if user == nil {
					user = object.GetUserByField(application.Organization, provider.Type, userInfo.Username)
				}
				if user == nil {
					user = object.GetUserByField(application.Organization, "name", userInfo.Username)
				}
			}

			if user != nil && user.IsDeleted == false {
				// Sign in via OAuth (want to sign up but already have account)

				if user.IsForbidden {
					c.ResponseError("the user is forbidden to sign in, please contact the administrator")
				}

				resp = c.HandleLoggedIn(application, user, &form)

				record := object.NewRecord(c.Ctx)
				record.Organization = application.Organization
				record.User = user.Name
				go object.AddRecord(record)
			} else if provider.Category == "OAuth" {
				// Sign up via OAuth
				if !application.EnableSignUp {
					c.ResponseError(fmt.Sprintf("The account for provider: %s and username: %s (%s) does not exist and is not allowed to sign up as new account, please contact your IT support", provider.Type, userInfo.Username, userInfo.DisplayName))
					return
				}

				if !providerItem.CanSignUp {
					c.ResponseError(fmt.Sprintf("The account for provider: %s and username: %s (%s) does not exist and is not allowed to sign up as new account via %s, please use another way to sign up", provider.Type, userInfo.Username, userInfo.DisplayName, provider.Type))
					return
				}

				properties := map[string]string{}
				properties["no"] = strconv.Itoa(len(object.GetUsers(application.Organization)) + 2)
				user = &object.User{
					Owner:             application.Organization,
					Name:              userInfo.Username,
					CreatedTime:       utils.GetCurrentTime(),
					Id:                utils.GenerateId(),
					Type:              "normal-user",
					DisplayName:       userInfo.DisplayName,
					Avatar:            userInfo.AvatarUrl,
					Address:           []string{},
					Email:             userInfo.Email,
					Score:             getInitScore(),
					IsAdmin:           false,
					IsGlobalAdmin:     false,
					IsForbidden:       false,
					IsDeleted:         false,
					SignupApplication: application.Name,
					Properties:        properties,
				}
				// sync info from 3rd-party if possible
				object.SetUserOAuthProperties(organization, user, provider.Type, userInfo)

				affected := object.AddUser(user)
				if !affected {
					c.ResponseError(fmt.Sprintf("Failed to create user, user information is invalid: %s", utils.StructToJson(user)))
					return
				}

				object.LinkUserAccount(user, provider.Type, userInfo.Id)

				resp = c.HandleLoggedIn(application, user, &form)

				record := object.NewRecord(c.Ctx)
				record.Organization = application.Organization
				record.User = user.Name
				go object.AddRecord(record)
			} else if provider.Category == "SAML" {
				resp = &Response{Status: "error", Msg: "The account does not exist"}
			}
			//resp = &Response{Status: "ok", Msg: "", Data: res}
		} else { // form.Method != "signup"
			userId := c.GetSessionUsername()
			if userId == "" {
				c.ResponseError("The account does not exist", userInfo)
				return
			}

			oldUser := object.GetUserByField(application.Organization, provider.Type, userInfo.Id)
			if oldUser == nil {
				oldUser = object.GetUserByField(application.Organization, provider.Type, userInfo.Username)
			}
			if oldUser != nil {
				c.ResponseError(fmt.Sprintf("The account for provider: %s and username: %s (%s) is already linked to another account: %s (%s)", provider.Type, userInfo.Username, userInfo.DisplayName, oldUser.Name, oldUser.DisplayName))
				return
			}

			user := object.GetUser(userId)

			// sync info from 3rd-party if possible
			object.SetUserOAuthProperties(organization, user, provider.Type, userInfo)

			isLinked := object.LinkUserAccount(user, provider.Type, userInfo.Id)
			if isLinked {
				resp = &Response{Status: "ok", Msg: "", Data: isLinked}
			} else {
				resp = &Response{Status: "error", Msg: "Failed to link user account", Data: isLinked}
			}
		}
	} else {
		if c.GetSessionUsername() != "" {
			// user already signed in to Bhojpur IAM, so let the user click the avatar button to do the quick sign-in
			application := object.GetApplication(fmt.Sprintf("admin/%s", form.Application))
			user := c.getCurrentUser()
			resp = c.HandleLoggedIn(application, user, &form)
		} else {
			c.ResponseError(fmt.Sprintf("unknown authentication type (not password or provider), form = %s", utils.StructToJson(form)))
			return
		}
	}

	c.Data["json"] = resp
	c.ServeJSON()
}

func (c *ApiController) GetSamlLogin() {
	providerId := c.Input().Get("id")
	relayState := c.Input().Get("relayState")
	authURL, method, err := object.GenerateSamlLoginUrl(providerId, relayState)
	if err != nil {
		c.ResponseError(err.Error())
	}
	c.ResponseOk(authURL, method)
}

func (c *ApiController) HandleSamlLogin() {
	relayState := c.Input().Get("RelayState")
	samlResponse := c.Input().Get("SAMLResponse")
	decode, err := base64.StdEncoding.DecodeString(relayState)
	if err != nil {
		c.ResponseError(err.Error())
	}
	slice := strings.Split(string(decode), "&")
	relayState = url.QueryEscape(relayState)
	samlResponse = url.QueryEscape(samlResponse)
	targetUrl := fmt.Sprintf("%s?relayState=%s&samlResponse=%s",
		slice[4], relayState, samlResponse)
	c.Redirect(targetUrl, 303)
}
