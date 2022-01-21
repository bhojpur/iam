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

// GetCerts
// @Title GetCerts
// @Tag Cert API
// @Description get certs
// @Param   owner     query    string  true        "The owner of certs"
// @Success 200 {array} object.Cert The Response object
// @router /get-certs [get]
func (c *ApiController) GetCerts() {
	owner := c.Input().Get("owner")
	limit := c.Input().Get("pageSize")
	page := c.Input().Get("p")
	field := c.Input().Get("field")
	value := c.Input().Get("value")
	sortField := c.Input().Get("sortField")
	sortOrder := c.Input().Get("sortOrder")
	if limit == "" || page == "" {
		c.Data["json"] = object.GetMaskedCerts(object.GetCerts(owner))
		c.ServeJSON()
	} else {
		limit := utils.ParseInt(limit)
		paginator := pagination.SetPaginator(c.Ctx, limit, int64(object.GetCertCount(owner, field, value)))
		certs := object.GetMaskedCerts(object.GetPaginationCerts(owner, paginator.Offset(), limit, field, value, sortField, sortOrder))
		c.ResponseOk(certs, paginator.Nums())
	}
}

// @Title GetCert
// @Tag Cert API
// @Description get cert
// @Param   id    query    string  true        "The id of the cert"
// @Success 200 {object} object.Cert The Response object
// @router /get-cert [get]
func (c *ApiController) GetCert() {
	id := c.Input().Get("id")

	c.Data["json"] = object.GetMaskedCert(object.GetCert(id))
	c.ServeJSON()
}

// @Title UpdateCert
// @Tag Cert API
// @Description update cert
// @Param   id    query    string  true        "The id of the cert"
// @Param   body    body   object.Cert  true        "The details of the cert"
// @Success 200 {object} controllers.Response The Response object
// @router /update-cert [post]
func (c *ApiController) UpdateCert() {
	id := c.Input().Get("id")

	var cert object.Cert
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &cert)
	if err != nil {
		panic(err)
	}

	c.Data["json"] = wrapActionResponse(object.UpdateCert(id, &cert))
	c.ServeJSON()
}

// @Title AddCert
// @Tag Cert API
// @Description add cert
// @Param   body    body   object.Cert  true        "The details of the cert"
// @Success 200 {object} controllers.Response The Response object
// @router /add-cert [post]
func (c *ApiController) AddCert() {
	var cert object.Cert
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &cert)
	if err != nil {
		panic(err)
	}

	c.Data["json"] = wrapActionResponse(object.AddCert(&cert))
	c.ServeJSON()
}

// @Title DeleteCert
// @Tag Cert API
// @Description delete cert
// @Param   body    body   object.Cert  true        "The details of the cert"
// @Success 200 {object} controllers.Response The Response object
// @router /delete-cert [post]
func (c *ApiController) DeleteCert() {
	var cert object.Cert
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &cert)
	if err != nil {
		panic(err)
	}

	c.Data["json"] = wrapActionResponse(object.DeleteCert(&cert))
	c.ServeJSON()
}
