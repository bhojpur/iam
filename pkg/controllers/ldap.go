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
)

type LdapServer struct {
	Host   string `json:"host"`
	Port   int    `json:"port"`
	Admin  string `json:"admin"`
	Passwd string `json:"passwd"`
	BaseDn string `json:"baseDn"`
}

type LdapResp struct {
	//Groups []LdapRespGroup `json:"groups"`
	Users []object.LdapRespUser `json:"users"`
}

//type LdapRespGroup struct {
//	GroupId   string
//	GroupName string
//}

type LdapSyncResp struct {
	Exist  []object.LdapRespUser `json:"exist"`
	Failed []object.LdapRespUser `json:"failed"`
}

// @Tag Account API
// @Title GetLdapser
// @router /get-ldap-user [post]
func (c *ApiController) GetLdapUser() {
	ldapServer := LdapServer{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &ldapServer)
	if err != nil || utils.IsStrsEmpty(ldapServer.Host, ldapServer.Admin, ldapServer.Passwd, ldapServer.BaseDn) {
		c.ResponseError("Missing parameter")
		return
	}

	var resp LdapResp

	conn, err := object.GetLdapConn(ldapServer.Host, ldapServer.Port, ldapServer.Admin, ldapServer.Passwd)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	//groupsMap, err := conn.GetLdapGroups(ldapServer.BaseDn)
	//if err != nil {
	//  c.ResponseError(err.Error())
	//	return
	//}

	//for _, group := range groupsMap {
	//	resp.Groups = append(resp.Groups, LdapRespGroup{
	//		GroupId:   group.GidNumber,
	//		GroupName: group.Cn,
	//	})
	//}

	users, err := conn.GetLdapUsers(ldapServer.BaseDn)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	for _, user := range users {
		resp.Users = append(resp.Users, object.LdapRespUser{
			UidNumber: user.UidNumber,
			Uid:       user.Uid,
			Cn:        user.Cn,
			GroupId:   user.GidNumber,
			//GroupName: groupsMap[user.GidNumber].Cn,
			Uuid:    user.Uuid,
			Email:   utils.GetMaxLenStr(user.Mail, user.Email, user.EmailAddress),
			Phone:   utils.GetMaxLenStr(user.TelephoneNumber, user.Mobile, user.MobileTelephoneNumber),
			Address: utils.GetMaxLenStr(user.RegisteredAddress, user.PostalAddress),
		})
	}

	c.Data["json"] = Response{Status: "ok", Data: resp}
	c.ServeJSON()
}

// @Tag Account API
// @Title GetLdaps
// @router /get-ldaps [post]
func (c *ApiController) GetLdaps() {
	webform, _ := c.Input()
	owner := webform.Get("owner")

	c.Data["json"] = Response{Status: "ok", Data: object.GetLdaps(owner)}
	c.ServeJSON()
}

// @Tag Account API
// @Title GetLdap
// @router /get-ldap [post]
func (c *ApiController) GetLdap() {
	webform, _ := c.Input()
	id := webform.Get("id")

	if utils.IsStrsEmpty(id) {
		c.ResponseError("Missing parameter")
		return
	}

	c.Data["json"] = Response{Status: "ok", Data: object.GetLdap(id)}
	c.ServeJSON()
}

// @Tag Account API
// @Title AddLdap
// @router /add-ldap [post]
func (c *ApiController) AddLdap() {
	var ldap object.Ldap
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &ldap)
	if err != nil {
		c.ResponseError("Missing parameter")
		return
	}

	if utils.IsStrsEmpty(ldap.Owner, ldap.ServerName, ldap.Host, ldap.Admin, ldap.Passwd, ldap.BaseDn) {
		c.ResponseError("Missing parameter")
		return
	}

	if object.CheckLdapExist(&ldap) {
		c.ResponseError("Ldap server exist")
		return
	}

	affected := object.AddLdap(&ldap)
	resp := wrapActionResponse(affected)
	if affected {
		resp.Data2 = ldap
	}
	if ldap.AutoSync != 0 {
		object.GetLdapAutoSynchronizer().StartAutoSync(ldap.Id)
	}

	c.Data["json"] = resp
	c.ServeJSON()
}

// @Tag Account API
// @Title UpdateLdap
// @router /update-ldap [post]
func (c *ApiController) UpdateLdap() {
	var ldap object.Ldap
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &ldap)
	if err != nil || utils.IsStrsEmpty(ldap.Owner, ldap.ServerName, ldap.Host, ldap.Admin, ldap.Passwd, ldap.BaseDn) {
		c.ResponseError("Missing parameter")
		return
	}

	affected := object.UpdateLdap(&ldap)
	resp := wrapActionResponse(affected)
	if affected {
		resp.Data2 = ldap
	}
	if ldap.AutoSync != 0 {
		object.GetLdapAutoSynchronizer().StartAutoSync(ldap.Id)
	}

	c.Data["json"] = resp
	c.ServeJSON()
}

// @Tag Account API
// @Title DeleteLdap
// @router /delete-ldap [post]
func (c *ApiController) DeleteLdap() {
	var ldap object.Ldap
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &ldap)
	if err != nil {
		panic(err)
	}

	object.GetLdapAutoSynchronizer().StopAutoSync(ldap.Id)
	c.Data["json"] = wrapActionResponse(object.DeleteLdap(&ldap))
	c.ServeJSON()
}

// @Tag Account API
// @Title SyncLdapUsers
// @router /sync-ldap-users [post]
func (c *ApiController) SyncLdapUsers() {
	webform, _ := c.Input()
	owner := webform.Get("owner")
	ldapId := webform.Get("ldapId")
	var users []object.LdapRespUser
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &users)
	if err != nil {
		panic(err)
	}

	object.UpdateLdapSyncTime(ldapId)

	exist, failed := object.SyncLdapUsers(owner, users)
	c.Data["json"] = &Response{Status: "ok", Data: &LdapSyncResp{
		Exist:  *exist,
		Failed: *failed,
	}}
	c.ServeJSON()
}

// @Tag Account API
// @Title CheckLdapUserExist
// @router /check-ldap-users-exist [post]
func (c *ApiController) CheckLdapUsersExist() {
	webform, _ := c.Input()
	owner := webform.Get("owner")
	var uuids []string
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &uuids)
	if err != nil {
		panic(err)
	}

	exist := object.CheckLdapUuidExist(owner, uuids)
	c.Data["json"] = &Response{Status: "ok", Data: exist}
	c.ServeJSON()
}
