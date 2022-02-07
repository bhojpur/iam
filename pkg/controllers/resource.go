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
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"path/filepath"

	"github.com/bhojpur/iam/pkg/object"
	"github.com/bhojpur/iam/pkg/utils"
	pagination "github.com/bhojpur/web/pkg/pagination"
)

// @router /get-resources [get]
// @Tag Resource API
// @Title GetResources
func (c *ApiController) GetResources() {
	webform, _ := c.Input()
	owner := webform.Get("owner")
	user := webform.Get("user")
	limit := webform.Get("pageSize")
	page := webform.Get("p")
	field := webform.Get("field")
	value := webform.Get("value")
	sortField := webform.Get("sortField")
	sortOrder := webform.Get("sortOrder")
	if limit == "" || page == "" {
		c.Data["json"] = object.GetResources(owner, user)
		c.ServeJSON()
	} else {
		limit := utils.ParseInt(limit)
		paginator := pagination.SetPaginator(c.Ctx, limit, int64(object.GetResourceCount(owner, user, field, value)))
		resources := object.GetPaginationResources(owner, user, paginator.Offset(), limit, field, value, sortField, sortOrder)
		c.ResponseOk(resources, paginator.Nums())
	}
}

// @Tag Resource API
// @Title GetResource
// @router /get-resource [get]
func (c *ApiController) GetResource() {
	webform, _ := c.Input()
	id := webform.Get("id")

	c.Data["json"] = object.GetResource(id)
	c.ServeJSON()
}

// @Tag Resource API
// @Title UpdateResource
// @router /update-resource [post]
func (c *ApiController) UpdateResource() {
	webform, _ := c.Input()
	id := webform.Get("id")

	var resource object.Resource
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &resource)
	if err != nil {
		panic(err)
	}

	c.Data["json"] = wrapActionResponse(object.UpdateResource(id, &resource))
	c.ServeJSON()
}

// @Tag Resource API
// @Title AddResource
// @router /add-resource [post]
func (c *ApiController) AddResource() {
	var resource object.Resource
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &resource)
	if err != nil {
		panic(err)
	}

	c.Data["json"] = wrapActionResponse(object.AddResource(&resource))
	c.ServeJSON()
}

// @Tag Resource API
// @Title DeleteResource
// @router /delete-resource [post]
func (c *ApiController) DeleteResource() {
	var resource object.Resource
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &resource)
	if err != nil {
		panic(err)
	}

	provider, _, ok := c.GetProviderFromContext("Storage")
	if !ok {
		return
	}

	err = object.DeleteFile(provider, resource.Name)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	c.Data["json"] = wrapActionResponse(object.DeleteResource(&resource))
	c.ServeJSON()
}

// @Tag Resource API
// @Title UploadResource
// @router /upload-resource [post]
func (c *ApiController) UploadResource() {
	webform, _ := c.Input()
	owner := webform.Get("owner")
	username := webform.Get("user")
	application := webform.Get("application")
	tag := webform.Get("tag")
	parent := webform.Get("parent")
	fullFilePath := webform.Get("fullFilePath")
	createdTime := webform.Get("createdTime")
	description := webform.Get("description")

	file, header, err := c.GetFile("file")
	if err != nil {
		c.ResponseError(err.Error())
		return
	}
	defer file.Close()

	if username == "" || fullFilePath == "" {
		c.ResponseError(fmt.Sprintf("username or fullFilePath is empty: username = %s, fullFilePath = %s", username, fullFilePath))
		return
	}

	filename := filepath.Base(fullFilePath)
	fileBuffer := bytes.NewBuffer(nil)
	if _, err = io.Copy(fileBuffer, file); err != nil {
		c.ResponseError(err.Error())
		return
	}

	provider, user, ok := c.GetProviderFromContext("Storage")
	if !ok {
		return
	}

	fileType := "unknown"
	contentType := header.Header.Get("Content-Type")
	fileType, _ = utils.GetOwnerAndNameFromId(contentType)

	if fileType != "image" && fileType != "video" {
		ext := filepath.Ext(filename)
		mimeType := mime.TypeByExtension(ext)
		fileType, _ = utils.GetOwnerAndNameFromId(mimeType)
	}

	fileUrl, objectKey, err := object.UploadFileSafe(provider, fullFilePath, fileBuffer)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	if createdTime == "" {
		createdTime = utils.GetCurrentTime()
	}
	fileFormat := filepath.Ext(fullFilePath)
	fileSize := int(header.Size)
	resource := &object.Resource{
		Owner:       owner,
		Name:        objectKey,
		CreatedTime: createdTime,
		User:        username,
		Provider:    provider.Name,
		Application: application,
		Tag:         tag,
		Parent:      parent,
		FileName:    filename,
		FileType:    fileType,
		FileFormat:  fileFormat,
		FileSize:    fileSize,
		Url:         fileUrl,
		Description: description,
	}
	object.AddOrUpdateResource(resource)

	switch tag {
	case "avatar":
		if user == nil {
			user = object.GetUserNoCheck(username)
			if user == nil {
				c.ResponseError("user is nil for tag: \"avatar\"")
				return
			}
		}

		user.Avatar = fileUrl
		object.UpdateUser(user.GetId(), user, []string{"avatar"}, false)
	case "termsOfUse":
		applicationId := fmt.Sprintf("admin/%s", parent)
		app := object.GetApplication(applicationId)
		app.TermsOfUse = fileUrl
		object.UpdateApplication(applicationId, app)
	}

	c.ResponseOk(fileUrl, objectKey)
}
