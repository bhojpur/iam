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
	"github.com/bhojpur/web/pkg/pagination"
)

// GetPayments
// @Title GetPayments
// @Tag Payment API
// @Description get payments
// @Param   owner     query    string  true        "The owner of payments"
// @Success 200 {array} object.Payment The Response object
// @router /get-payments [get]
func (c *ApiController) GetPayments() {
	webform, _ := c.Input()
	owner := webform.Get("owner")
	limit := webform.Get("pageSize")
	page := webform.Get("p")
	field := webform.Get("field")
	value := webform.Get("value")
	sortField := webform.Get("sortField")
	sortOrder := webform.Get("sortOrder")
	if limit == "" || page == "" {
		c.Data["json"] = object.GetPayments(owner)
		c.ServeJSON()
	} else {
		limit := utils.ParseInt(limit)
		paginator := pagination.SetPaginator(c.Ctx, limit, int64(object.GetPaymentCount(owner, field, value)))
		payments := object.GetPaginationPayments(owner, paginator.Offset(), limit, field, value, sortField, sortOrder)
		c.ResponseOk(payments, paginator.Nums())
	}
}

// @Title GetPayment
// @Tag Payment API
// @Description get payment
// @Param   id    query    string  true        "The id of the payment"
// @Success 200 {object} object.Payment The Response object
// @router /get-payment [get]
func (c *ApiController) GetPayment() {
	webform, _ := c.Input()
	id := webform.Get("id")

	c.Data["json"] = object.GetPayment(id)
	c.ServeJSON()
}

// @Title UpdatePayment
// @Tag Payment API
// @Description update payment
// @Param   id    query    string  true        "The id of the payment"
// @Param   body    body   object.Payment  true        "The details of the payment"
// @Success 200 {object} controllers.Response The Response object
// @router /update-payment [post]
func (c *ApiController) UpdatePayment() {
	webform, _ := c.Input()
	id := webform.Get("id")

	var payment object.Payment
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &payment)
	if err != nil {
		panic(err)
	}

	c.Data["json"] = wrapActionResponse(object.UpdatePayment(id, &payment))
	c.ServeJSON()
}

// @Title AddPayment
// @Tag Payment API
// @Description add payment
// @Param   body    body   object.Payment  true        "The details of the payment"
// @Success 200 {object} controllers.Response The Response object
// @router /add-payment [post]
func (c *ApiController) AddPayment() {
	var payment object.Payment
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &payment)
	if err != nil {
		panic(err)
	}

	c.Data["json"] = wrapActionResponse(object.AddPayment(&payment))
	c.ServeJSON()
}

// @Title DeletePayment
// @Tag Payment API
// @Description delete payment
// @Param   body    body   object.Payment  true        "The details of the payment"
// @Success 200 {object} controllers.Response The Response object
// @router /delete-payment [post]
func (c *ApiController) DeletePayment() {
	var payment object.Payment
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &payment)
	if err != nil {
		panic(err)
	}

	c.Data["json"] = wrapActionResponse(object.DeletePayment(&payment))
	c.ServeJSON()
}
