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
	"errors"
	"fmt"
	"strings"

	"github.com/bhojpur/iam/pkg/object"
	"github.com/bhojpur/iam/pkg/utils"
)

func (c *ApiController) getCurrentUser() *object.User {
	var user *object.User
	userId := c.GetSessionUsername()
	if userId == "" {
		user = nil
	} else {
		user = object.GetUser(userId)
	}
	return user
}

// SendVerificationCode ...
// @Title SendVerificationCode
// @Tag Verification API
// @router /send-verification-code [post]
func (c *ApiController) SendVerificationCode() {
	destType := c.Ctx.Request.Form.Get("type")
	dest := c.Ctx.Request.Form.Get("dest")
	orgId := c.Ctx.Request.Form.Get("organizationId")
	checkType := c.Ctx.Request.Form.Get("checkType")
	checkId := c.Ctx.Request.Form.Get("checkId")
	checkKey := c.Ctx.Request.Form.Get("checkKey")
	checkUser := c.Ctx.Request.Form.Get("checkUser")
	remoteAddr := utils.GetIPFromRequest(c.Ctx.Request)

	if len(destType) == 0 || len(dest) == 0 || len(orgId) == 0 || !strings.Contains(orgId, "/") || len(checkType) == 0 || len(checkId) == 0 || len(checkKey) == 0 {
		c.ResponseError("Missing parameter.")
		return
	}

	isHuman := false
	captchaProvider := object.GetDefaultHumanCheckProvider()
	if captchaProvider == nil {
		isHuman = object.VerifyCaptcha(checkId, checkKey)
	}

	if !isHuman {
		c.ResponseError("Turing test failed.")
		return
	}

	user := c.getCurrentUser()
	organization := object.GetOrganization(orgId)
	application := object.GetApplicationByOrganizationName(organization.Name)

	if checkUser == "true" && user == nil &&
		object.GetUserByFields(organization.Name, dest) == nil {
		c.ResponseError("No such user.")
		return
	}

	sendResp := errors.New("Invalid dest type.")
	switch destType {
	case "email":
		if !utils.IsEmailValid(dest) {
			c.ResponseError("Invalid Email address")
			return
		}

		provider := application.GetEmailProvider()
		sendResp = object.SendVerificationCodeToEmail(organization, user, provider, remoteAddr, dest)
	case "phone":
		if !utils.IsPhoneCnValid(dest) {
			c.ResponseError("Invalid phone number")
			return
		}
		org := object.GetOrganization(orgId)
		if org == nil {
			c.ResponseError("Missing parameter.")
			return
		}

		dest = fmt.Sprintf("+%s%s", org.PhonePrefix, dest)
		provider := application.GetSmsProvider()
		sendResp = object.SendVerificationCodeToPhone(organization, user, provider, remoteAddr, dest)
	}

	if sendResp != nil {
		c.Data["json"] = Response{Status: "error", Msg: sendResp.Error()}
	} else {
		c.Data["json"] = Response{Status: "ok"}
	}

	c.ServeJSON()
}

// ResetEmailOrPhone ...
// @Tag Account API
// @Title ResetEmailOrPhone
// @router /api/reset-email-or-phone [post]
func (c *ApiController) ResetEmailOrPhone() {
	userId, ok := c.RequireSignedIn()
	if !ok {
		return
	}

	user := object.GetUser(userId)
	if user == nil {
		c.ResponseError("No such user.")
		return
	}

	destType := c.Ctx.Request.Form.Get("type")
	dest := c.Ctx.Request.Form.Get("dest")
	code := c.Ctx.Request.Form.Get("code")
	if len(dest) == 0 || len(code) == 0 || len(destType) == 0 {
		c.ResponseError("Missing parameter.")
		return
	}

	checkDest := dest
	if destType == "phone" {
		org := object.GetOrganizationByUser(user)
		phonePrefix := "86"
		if org != nil && org.PhonePrefix != "" {
			phonePrefix = org.PhonePrefix
		}
		checkDest = fmt.Sprintf("+%s%s", phonePrefix, dest)
	}
	if ret := object.CheckVerificationCode(checkDest, code); len(ret) != 0 {
		c.ResponseError(ret)
		return
	}

	switch destType {
	case "email":
		user.Email = dest
		object.SetUserField(user, "email", user.Email)
	case "phone":
		user.Phone = dest
		object.SetUserField(user, "phone", user.Phone)
	default:
		c.ResponseError("Unknown type.")
		return
	}

	object.DisableVerificationCode(checkDest)
	c.Data["json"] = Response{Status: "ok"}
	c.ServeJSON()
}
