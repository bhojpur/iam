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
	"strings"

	"github.com/bhojpur/iam/pkg/object"
	"github.com/bhojpur/iam/pkg/utils"
	pagination "github.com/bhojpur/web/pkg/pagination"
)

// GetGlobalUsers
// @Title GetGlobalUsers
// @Tag User API
// @Description get global users
// @Success 200 {array} object.User The Response object
// @router /get-global-users [get]
func (c *ApiController) GetGlobalUsers() {
	webform, _ := c.Input()
	limit := webform.Get("pageSize")
	page := webform.Get("p")
	field := webform.Get("field")
	value := webform.Get("value")
	sortField := webform.Get("sortField")
	sortOrder := webform.Get("sortOrder")
	if limit == "" || page == "" {
		c.Data["json"] = object.GetMaskedUsers(object.GetGlobalUsers())
		c.ServeJSON()
	} else {
		limit := utils.ParseInt(limit)
		paginator := pagination.SetPaginator(c.Ctx, limit, int64(object.GetGlobalUserCount(field, value)))
		users := object.GetPaginationGlobalUsers(paginator.Offset(), limit, field, value, sortField, sortOrder)
		c.ResponseOk(users, paginator.Nums())
	}
}

// GetUsers
// @Title GetUsers
// @Tag User API
// @Description
// @Param   owner     query    string  true        "The owner of users"
// @Success 200 {array} object.User The Response object
// @router /get-users [get]
func (c *ApiController) GetUsers() {
	webform, _ := c.Input()
	owner := webform.Get("owner")
	limit := webform.Get("pageSize")
	page := webform.Get("p")
	field := webform.Get("field")
	value := webform.Get("value")
	sortField := webform.Get("sortField")
	sortOrder := webform.Get("sortOrder")
	if limit == "" || page == "" {
		c.Data["json"] = object.GetMaskedUsers(object.GetUsers(owner))
		c.ServeJSON()
	} else {
		limit := utils.ParseInt(limit)
		paginator := pagination.SetPaginator(c.Ctx, limit, int64(object.GetUserCount(owner, field, value)))
		users := object.GetPaginationUsers(owner, paginator.Offset(), limit, field, value, sortField, sortOrder)
		c.ResponseOk(users, paginator.Nums())
	}
}

// GetUser
// @Title GetUser
// @Tag User API
// @Description get user
// @Param   id     query    string  true        "The id of the user"
// @Success 200 {object} object.User The Response object
// @router /get-user [get]
func (c *ApiController) GetUser() {
	webform, _ := c.Input()
	id := webform.Get("id")
	owner := webform.Get("owner")
	email := webform.Get("email")

	var user *object.User
	if email == "" {
		user = object.GetUser(id)
	} else {
		user = object.GetUserByEmail(owner, email)
	}

	c.Data["json"] = object.GetMaskedUser(user)
	c.ServeJSON()
}

// UpdateUser
// @Title UpdateUser
// @Tag User API
// @Description update user
// @Param   id     query    string  true        "The id of the user"
// @Param   body    body   object.User  true        "The details of the user"
// @Success 200 {object} controllers.Response The Response object
// @router /update-user [post]
func (c *ApiController) UpdateUser() {
	webform, _ := c.Input()
	id := webform.Get("id")
	columnsStr := webform.Get("columns")

	var user object.User
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &user)
	if err != nil {
		panic(err)
	}

	if user.DisplayName == "" {
		c.ResponseError("Display name cannot be empty")
		return
	}

	columns := []string{}
	if columnsStr != "" {
		columns = strings.Split(columnsStr, ",")
	}

	isGlobalAdmin := c.IsGlobalAdmin()
	affected := object.UpdateUser(id, &user, columns, isGlobalAdmin)
	if affected {
		object.UpdateUserToOriginalDatabase(&user)
	}

	c.Data["json"] = wrapActionResponse(affected)
	c.ServeJSON()
}

// AddUser
// @Title AddUser
// @Tag User API
// @Description add user
// @Param   body    body   object.User  true        "The details of the user"
// @Success 200 {object} controllers.Response The Response object
// @router /add-user [post]
func (c *ApiController) AddUser() {
	var user object.User
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &user)
	if err != nil {
		panic(err)
	}

	c.Data["json"] = wrapActionResponse(object.AddUser(&user))
	c.ServeJSON()
}

// DeleteUser
// @Title DeleteUser
// @Tag User API
// @Description delete user
// @Param   body    body   object.User  true        "The details of the user"
// @Success 200 {object} controllers.Response The Response object
// @router /delete-user [post]
func (c *ApiController) DeleteUser() {
	var user object.User
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &user)
	if err != nil {
		panic(err)
	}

	c.Data["json"] = wrapActionResponse(object.DeleteUser(&user))
	c.ServeJSON()
}

// GetEmailAndPhone
// @Title GetEmailAndPhone
// @Tag User API
// @Description get email and phone by username
// @Param   username    formData   string  true        "The username of the user"
// @Param   organization    formData   string  true        "The organization of the user"
// @Success 200 {object} controllers.Response The Response object
// @router /get-email-and-phone [post]
func (c *ApiController) GetEmailAndPhone() {
	var form RequestForm
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &form)
	if err != nil {
		panic(err)
	}

	user := object.GetUserByFields(form.Organization, form.Username)
	if user == nil {
		c.ResponseError("No such user.")
		return
	}

	respUser := object.User{Email: user.Email, Phone: user.Phone, Name: user.Name}
	var contentType string
	switch form.Username {
	case user.Email:
		contentType = "email"
	case user.Phone:
		contentType = "phone"
	case user.Name:
		contentType = "username"
	}

	c.ResponseOk(respUser, contentType)
}

// SetPassword
// @Title SetPassword
// @Tag Account API
// @Description set password
// @Param   userOwner   formData    string  true        "The owner of the user"
// @Param   userName   formData    string  true        "The name of the user"
// @Param   oldPassword   formData    string  true        "The old password of the user"
// @Param   newPassword   formData    string  true        "The new password of the user"
// @Success 200 {object} controllers.Response The Response object
// @router /set-password [post]
func (c *ApiController) SetPassword() {
	userOwner := c.Ctx.Request.Form.Get("userOwner")
	userName := c.Ctx.Request.Form.Get("userName")
	oldPassword := c.Ctx.Request.Form.Get("oldPassword")
	newPassword := c.Ctx.Request.Form.Get("newPassword")

	requestUserId := c.GetSessionUsername()
	if requestUserId == "" {
		c.ResponseError("Please login first.")
		return
	}

	userId := fmt.Sprintf("%s/%s", userOwner, userName)
	targetUser := object.GetUser(userId)
	if targetUser == nil {
		c.ResponseError(fmt.Sprintf("The user: %s doesn't exist", userId))
		return
	}

	hasPermission := false
	if strings.HasPrefix(requestUserId, "app/") {
		hasPermission = true
	} else {
		requestUser := object.GetUser(requestUserId)
		if requestUser == nil {
			c.ResponseError("Session outdated. Please login again.")
			return
		}
		if requestUser.IsGlobalAdmin {
			hasPermission = true
		} else if requestUserId == userId {
			hasPermission = true
		} else if targetUser.Owner == requestUser.Owner && requestUser.IsAdmin {
			hasPermission = true
		}
	}
	if !hasPermission {
		c.ResponseError("You don't have the permission to do this.")
		return
	}

	if oldPassword != "" {
		msg := object.CheckPassword(targetUser, oldPassword)
		if msg != "" {
			c.ResponseError(msg)
			return
		}
	}

	if strings.Contains(newPassword, " ") {
		c.ResponseError("New password cannot contain blank space.")
		return
	}

	if len(newPassword) <= 5 {
		c.ResponseError("New password must have at least 6 characters")
		return
	}

	targetUser.Password = newPassword
	object.SetUserField(targetUser, "password", targetUser.Password)
	c.Data["json"] = Response{Status: "ok"}
	c.ServeJSON()
}

// @Title CheckUserPassword
// @router /check-user-password [post]
// @Tag User API
func (c *ApiController) CheckUserPassword() {
	var user object.User
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &user)
	if err != nil {
		panic(err)
	}

	_, msg := object.CheckUserPassword(user.Owner, user.Name, user.Password)
	if msg == "" {
		c.ResponseOk()
	} else {
		c.ResponseError(msg)
	}
}

// GetSortedUsers
// @Title GetSortedUsers
// @Tag User API
// @Description
// @Param   owner     query    string  true        "The owner of users"
// @Param   sorter     query    string  true        "The DB column name to sort by, e.g., created_time"
// @Param   limit     query    string  true        "The count of users to return, e.g., 25"
// @Success 200 {array} object.User The Response object
// @router /get-sorted-users [get]
func (c *ApiController) GetSortedUsers() {
	webform, _ := c.Input()
	owner := webform.Get("owner")
	sorter := webform.Get("sorter")
	limit := utils.ParseInt(webform.Get("limit"))

	c.Data["json"] = object.GetMaskedUsers(object.GetSortedUsers(owner, sorter, limit))
	c.ServeJSON()
}

// GetUserCount
// @Title GetUserCount
// @Tag User API
// @Description
// @Param   owner     query    string  true        "The owner of users"
// @Param   isOnline     query    string  true        "The filter for query, 1 for online, 0 for offline, empty string for all users"
// @Success 200 {int} int The count of filtered users for an organization
// @router /get-user-count [get]
func (c *ApiController) GetUserCount() {
	webform, _ := c.Input()
	owner := webform.Get("owner")
	isOnline := webform.Get("isOnline")

	count := 0
	if isOnline == "" {
		count = object.GetUserCount(owner, "", "")
	} else {
		count = object.GetOnlineUserCount(owner, utils.ParseInt(isOnline))
	}

	c.Data["json"] = count
	c.ServeJSON()
}
