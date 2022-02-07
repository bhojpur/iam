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
	pagination "github.com/bhojpur/web/pkg/pagination"
)

// GetApplications
// @Title GetApplications
// @Tag Application API
// @Description get all applications
// @Param   owner     query    string  true        "The owner of applications."
// @Success 200 {array} object.Application The Response object
// @router /get-applications [get]
func (c *ApiController) GetApplications() {
	userId := c.GetSessionUsername()
	webform, _ := c.Input()
	owner := webform.Get("owner")
	limit := webform.Get("pageSize")
	page := webform.Get("p")
	field := webform.Get("field")
	value := webform.Get("value")
	sortField := webform.Get("sortField")
	sortOrder := webform.Get("sortOrder")
	organization := webform.Get("organization")

	if limit == "" || page == "" {
		var applications []*object.Application
		if organization == "" {
			applications = object.GetApplications(owner)
		} else {
			applications = object.GetApplicationsByOrganizationName(owner, organization)
		}

		c.Data["json"] = object.GetMaskedApplications(applications, userId)
		c.ServeJSON()
	} else {
		limit := utils.ParseInt(limit)
		paginator := pagination.SetPaginator(c.Ctx, limit, int64(object.GetApplicationCount(owner, field, value)))
		applications := object.GetMaskedApplications(object.GetPaginationApplications(owner, paginator.Offset(), limit, field, value, sortField, sortOrder), userId)
		c.ResponseOk(applications, paginator.Nums())
	}
}

// GetApplication
// @Title GetApplication
// @Tag Application API
// @Description get the detail of an application
// @Param   id     query    string  true        "The id of the application."
// @Success 200 {object} object.Application The Response object
// @router /get-application [get]
func (c *ApiController) GetApplication() {
	userId := c.GetSessionUsername()
	webform, _ := c.Input()
	id := webform.Get("id")

	c.Data["json"] = object.GetMaskedApplication(object.GetApplication(id), userId)
	c.ServeJSON()
}

// GetUserApplication
// @Title GetUserApplication
// @Tag Application API
// @Description get the detail of the user's application
// @Param   id     query    string  true        "The id of the user"
// @Success 200 {object} object.Application The Response object
// @router /get-user-application [get]
func (c *ApiController) GetUserApplication() {
	userId := c.GetSessionUsername()
	webform, _ := c.Input()
	id := webform.Get("id")
	user := object.GetUser(id)
	if user == nil {
		c.ResponseError("No such user.")
		return
	}

	c.Data["json"] = object.GetMaskedApplication(object.GetApplicationByUser(user), userId)
	c.ServeJSON()
}

// UpdateApplication
// @Title UpdateApplication
// @Tag Application API
// @Description update an application
// @Param   id     query    string  true        "The id of the application"
// @Param   body    body   object.Application  true        "The details of the application"
// @Success 200 {object} controllers.Response The Response object
// @router /update-application [post]
func (c *ApiController) UpdateApplication() {
	webform, _ := c.Input()
	id := webform.Get("id")

	var application object.Application
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &application)
	if err != nil {
		panic(err)
	}

	c.Data["json"] = wrapActionResponse(object.UpdateApplication(id, &application))
	c.ServeJSON()
}

// AddApplication
// @Title AddApplication
// @Tag Application API
// @Description add an application
// @Param   body    body   object.Application  true        "The details of the application"
// @Success 200 {object} controllers.Response The Response object
// @router /add-application [post]
func (c *ApiController) AddApplication() {
	var application object.Application
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &application)
	if err != nil {
		panic(err)
	}

	c.Data["json"] = wrapActionResponse(object.AddApplication(&application))
	c.ServeJSON()
}

// DeleteApplication
// @Title DeleteApplication
// @Tag Application API
// @Description delete an application
// @Param   body    body   object.Application  true        "The details of the application"
// @Success 200 {object} controllers.Response The Response object
// @router /delete-application [post]
func (c *ApiController) DeleteApplication() {
	var application object.Application
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &application)
	if err != nil {
		panic(err)
	}

	c.Data["json"] = wrapActionResponse(object.DeleteApplication(&application))
	c.ServeJSON()
}
