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

// GetPermissions
// @Title GetPermissions
// @Tag Permission API
// @Description get permissions
// @Param   owner     query    string  true        "The owner of permissions"
// @Success 200 {array} object.Permission The Response object
// @router /get-permissions [get]
func (c *ApiController) GetPermissions() {
	webform, _ := c.Input()
	owner := webform.Get("owner")
	limit := webform.Get("pageSize")
	page := webform.Get("p")
	field := webform.Get("field")
	value := webform.Get("value")
	sortField := webform.Get("sortField")
	sortOrder := webform.Get("sortOrder")
	if limit == "" || page == "" {
		c.Data["json"] = object.GetPermissions(owner)
		c.ServeJSON()
	} else {
		limit := utils.ParseInt(limit)
		paginator := pagination.SetPaginator(c.Ctx, limit, int64(object.GetPermissionCount(owner, field, value)))
		permissions := object.GetPaginationPermissions(owner, paginator.Offset(), limit, field, value, sortField, sortOrder)
		c.ResponseOk(permissions, paginator.Nums())
	}
}

// @Title GetPermission
// @Tag Permission API
// @Description get permission
// @Param   id    query    string  true        "The id of the permission"
// @Success 200 {object} object.Permission The Response object
// @router /get-permission [get]
func (c *ApiController) GetPermission() {
	webform, _ := c.Input()
	id := webform.Get("id")

	c.Data["json"] = object.GetPermission(id)
	c.ServeJSON()
}

// @Title UpdatePermission
// @Tag Permission API
// @Description update permission
// @Param   id    query    string  true        "The id of the permission"
// @Param   body    body   object.Permission  true        "The details of the permission"
// @Success 200 {object} controllers.Response The Response object
// @router /update-permission [post]
func (c *ApiController) UpdatePermission() {
	webform, _ := c.Input()
	id := webform.Get("id")

	var permission object.Permission
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &permission)
	if err != nil {
		panic(err)
	}

	c.Data["json"] = wrapActionResponse(object.UpdatePermission(id, &permission))
	c.ServeJSON()
}

// @Title AddPermission
// @Tag Permission API
// @Description add permission
// @Param   body    body   object.Permission  true        "The details of the permission"
// @Success 200 {object} controllers.Response The Response object
// @router /add-permission [post]
func (c *ApiController) AddPermission() {
	var permission object.Permission
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &permission)
	if err != nil {
		panic(err)
	}

	c.Data["json"] = wrapActionResponse(object.AddPermission(&permission))
	c.ServeJSON()
}

// @Title DeletePermission
// @Tag Permission API
// @Description delete permission
// @Param   body    body   object.Permission  true        "The details of the permission"
// @Success 200 {object} controllers.Response The Response object
// @router /delete-permission [post]
func (c *ApiController) DeletePermission() {
	var permission object.Permission
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &permission)
	if err != nil {
		panic(err)
	}

	c.Data["json"] = wrapActionResponse(object.DeletePermission(&permission))
	c.ServeJSON()
}
