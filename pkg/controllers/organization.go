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

// GetOrganizations ...
// @Title GetOrganizations
// @Tag Organization API
// @Description get organizations
// @Param   owner     query    string  true        "owner"
// @Success 200 {array} object.Organization The Response object
// @router /get-organizations [get]
func (c *ApiController) GetOrganizations() {
	webform, _ := c.Input()
	owner := webform.Get("owner")
	limit := webform.Get("pageSize")
	page := webform.Get("p")
	field := webform.Get("field")
	value := webform.Get("value")
	sortField := webform.Get("sortField")
	sortOrder := webform.Get("sortOrder")
	if limit == "" || page == "" {
		c.Data["json"] = object.GetMaskedOrganizations(object.GetOrganizations(owner))
		c.ServeJSON()
	} else {
		limit := utils.ParseInt(limit)
		paginator := pagination.SetPaginator(c.Ctx, limit, int64(object.GetOrganizationCount(owner, field, value)))
		organizations := object.GetMaskedOrganizations(object.GetPaginationOrganizations(owner, paginator.Offset(), limit, field, value, sortField, sortOrder))
		c.ResponseOk(organizations, paginator.Nums())
	}
}

// GetOrganization ...
// @Title GetOrganization
// @Tag Organization API
// @Description get organization
// @Param   id     query    string  true        "organization id"
// @Success 200 {object} object.Organization The Response object
// @router /get-organization [get]
func (c *ApiController) GetOrganization() {
	webform, _ := c.Input()
	id := webform.Get("id")

	c.Data["json"] = object.GetMaskedOrganization(object.GetOrganization(id))
	c.ServeJSON()
}

// UpdateOrganization ...
// @Title UpdateOrganization
// @Tag Organization API
// @Description update organization
// @Param   id     query    string  true        "The id of the organization"
// @Param   body    body   object.Organization  true        "The details of the organization"
// @Success 200 {object} controllers.Response The Response object
// @router /update-organization [post]
func (c *ApiController) UpdateOrganization() {
	webform, _ := c.Input()
	id := webform.Get("id")

	var organization object.Organization
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &organization)
	if err != nil {
		panic(err)
	}

	c.Data["json"] = wrapActionResponse(object.UpdateOrganization(id, &organization))
	c.ServeJSON()
}

// AddOrganization ...
// @Title AddOrganization
// @Tag Organization API
// @Description add organization
// @Param   body    body   object.Organization  true        "The details of the organization"
// @Success 200 {object} controllers.Response The Response object
// @router /add-organization [post]
func (c *ApiController) AddOrganization() {
	var organization object.Organization
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &organization)
	if err != nil {
		panic(err)
	}

	c.Data["json"] = wrapActionResponse(object.AddOrganization(&organization))
	c.ServeJSON()
}

// DeleteOrganization ...
// @Title DeleteOrganization
// @Tag Organization API
// @Description delete organization
// @Param   body    body   object.Organization  true        "The details of the organization"
// @Success 200 {object} controllers.Response The Response object
// @router /delete-organization [post]
func (c *ApiController) DeleteOrganization() {
	var organization object.Organization
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &organization)
	if err != nil {
		panic(err)
	}

	c.Data["json"] = wrapActionResponse(object.DeleteOrganization(&organization))
	c.ServeJSON()
}
