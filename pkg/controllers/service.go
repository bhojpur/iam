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
	"fmt"

	"github.com/bhojpur/iam/pkg/object"
	"github.com/bhojpur/iam/pkg/utils"
)

// SendEmail
// @Title SendEmail
// @Tag Service API
// @Description This API is not for Bhojpur IAM frontend to call, it is for Bhjojpur IAM SDKs.
// @Param   clientId    query    string  true        "The clientId of the application"
// @Param   clientSecret    query    string  true    "The clientSecret of the application"
// @Param   body    body   emailForm    true         "Details of the email request"
// @Success 200 {object}  Response object
// @router /api/send-email [post]
func (c *ApiController) SendEmail() {
	provider, _, ok := c.GetProviderFromContext("Email")
	if !ok {
		return
	}

	var emailForm struct {
		Title     string   `json:"title"`
		Content   string   `json:"content"`
		Sender    string   `json:"sender"`
		Receivers []string `json:"receivers"`
	}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &emailForm)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	if utils.IsStrsEmpty(emailForm.Title, emailForm.Content, emailForm.Sender) {
		c.ResponseError(fmt.Sprintf("Empty parameters for emailForm: %v", emailForm))
		return
	}

	invalidReceivers := []string{}
	for _, receiver := range emailForm.Receivers {
		if !utils.IsEmailValid(receiver) {
			invalidReceivers = append(invalidReceivers, receiver)
		}
	}

	if len(invalidReceivers) != 0 {
		c.ResponseError(fmt.Sprintf("Invalid Email receivers: %s", invalidReceivers))
		return
	}

	for _, receiver := range emailForm.Receivers {
		err = object.SendEmail(provider, emailForm.Title, emailForm.Content, receiver, emailForm.Sender)
		if err != nil {
			c.ResponseError(err.Error())
			return
		}
	}

	c.ResponseOk()
}

// SendSms
// @Title SendSms
// @Tag Service API
// @Description This API is not for Bhojpur IAM frontend to call, it is for Bhojpur IAM SDKs.
// @Param   clientId    query    string  true        "The clientId of the application"
// @Param   clientSecret    query    string  true    "The clientSecret of the application"
// @Param   body    body   smsForm    true           "Details of the sms request"
// @Success 200 {object}  Response object
// @router /api/send-sms [post]
func (c *ApiController) SendSms() {
	provider, _, ok := c.GetProviderFromContext("SMS")
	if !ok {
		return
	}

	var smsForm struct {
		Content   string   `json:"content"`
		Receivers []string `json:"receivers"`
		OrgId     string   `json:"organizationId"` // e.g. "admin/built-in"
	}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &smsForm)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	org := object.GetOrganization(smsForm.OrgId)
	var invalidReceivers []string
	for idx, receiver := range smsForm.Receivers {
		if !utils.IsPhoneCnValid(receiver) {
			invalidReceivers = append(invalidReceivers, receiver)
		} else {
			smsForm.Receivers[idx] = fmt.Sprintf("+%s%s", org.PhonePrefix, receiver)
		}
	}

	if len(invalidReceivers) != 0 {
		c.ResponseError(fmt.Sprintf("Invalid phone receivers: %s", invalidReceivers))
		return
	}

	err = object.SendSms(provider, smsForm.Content, smsForm.Receivers...)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	c.ResponseOk()
}
