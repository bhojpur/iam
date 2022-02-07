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

// GetProviders
// @Title GetProviders
// @Tag Provider API
// @Description get providers
// @Param   owner     query    string  true        "The owner of providers"
// @Success 200 {array} object.Provider The Response object
// @router /get-providers [get]
func (c *ApiController) GetProviders() {
	webform, _ := c.Input()
	owner := webform.Get("owner")
	limit := webform.Get("pageSize")
	page := webform.Get("p")
	field := webform.Get("field")
	value := webform.Get("value")
	sortField := webform.Get("sortField")
	sortOrder := webform.Get("sortOrder")
	if limit == "" || page == "" {
		c.Data["json"] = object.GetMaskedProviders(object.GetProviders(owner))
		c.ServeJSON()
	} else {
		limit := utils.ParseInt(limit)
		paginator := pagination.SetPaginator(c.Ctx, limit, int64(object.GetProviderCount(owner, field, value)))
		providers := object.GetMaskedProviders(object.GetPaginationProviders(owner, paginator.Offset(), limit, field, value, sortField, sortOrder))
		c.ResponseOk(providers, paginator.Nums())
	}
}

// @Title GetProvider
// @Tag Provider API
// @Description get provider
// @Param   id    query    string  true        "The id of the provider"
// @Success 200 {object} object.Provider The Response object
// @router /get-provider [get]
func (c *ApiController) GetProvider() {
	webform, _ := c.Input()
	id := webform.Get("id")

	c.Data["json"] = object.GetMaskedProvider(object.GetProvider(id))
	c.ServeJSON()
}

// @Title UpdateProvider
// @Tag Provider API
// @Description update provider
// @Param   id    query    string  true        "The id of the provider"
// @Param   body    body   object.Provider  true        "The details of the provider"
// @Success 200 {object} controllers.Response The Response object
// @router /update-provider [post]
func (c *ApiController) UpdateProvider() {
	webform, _ := c.Input()
	id := webform.Get("id")

	var provider object.Provider
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &provider)
	if err != nil {
		panic(err)
	}

	c.Data["json"] = wrapActionResponse(object.UpdateProvider(id, &provider))
	c.ServeJSON()
}

// @Title AddProvider
// @Tag Provider API
// @Description add provider
// @Param   body    body   object.Provider  true        "The details of the provider"
// @Success 200 {object} controllers.Response The Response object
// @router /add-provider [post]
func (c *ApiController) AddProvider() {
	var provider object.Provider
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &provider)
	if err != nil {
		panic(err)
	}

	c.Data["json"] = wrapActionResponse(object.AddProvider(&provider))
	c.ServeJSON()
}

// @Title DeleteProvider
// @Tag Provider API
// @Description delete provider
// @Param   body    body   object.Provider  true        "The details of the provider"
// @Success 200 {object} controllers.Response The Response object
// @router /delete-provider [post]
func (c *ApiController) DeleteProvider() {
	var provider object.Provider
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &provider)
	if err != nil {
		panic(err)
	}

	c.Data["json"] = wrapActionResponse(object.DeleteProvider(&provider))
	c.ServeJSON()
}
