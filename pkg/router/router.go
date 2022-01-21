package router

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
	websvr "github.com/bhojpur/web/pkg/engine"

	"github.com/bhojpur/iam/pkg/controllers"
)

func init() {
	initAPI()
}

func initAPI() {
	ns :=
		websvr.NewNamespace("/",
			websvr.NSNamespace("/api",
				websvr.NSInclude(
					&controllers.ApiController{},
				),
			),
			websvr.NSNamespace("",
				websvr.NSInclude(
					&controllers.RootController{},
				),
			),
		)
	websvr.AddNamespace(ns)

	websvr.Router("/api/signup", &controllers.ApiController{}, "POST:Signup")
	websvr.Router("/api/login", &controllers.ApiController{}, "POST:Login")
	websvr.Router("/api/get-app-login", &controllers.ApiController{}, "GET:GetApplicationLogin")
	websvr.Router("/api/logout", &controllers.ApiController{}, "POST:Logout")
	websvr.Router("/api/get-account", &controllers.ApiController{}, "GET:GetAccount")
	websvr.Router("/api/unlink", &controllers.ApiController{}, "POST:Unlink")
	websvr.Router("/api/get-saml-login", &controllers.ApiController{}, "GET:GetSamlLogin")
	websvr.Router("/api/acs", &controllers.ApiController{}, "POST:HandleSamlLogin")

	websvr.Router("/api/get-organizations", &controllers.ApiController{}, "GET:GetOrganizations")
	websvr.Router("/api/get-organization", &controllers.ApiController{}, "GET:GetOrganization")
	websvr.Router("/api/update-organization", &controllers.ApiController{}, "POST:UpdateOrganization")
	websvr.Router("/api/add-organization", &controllers.ApiController{}, "POST:AddOrganization")
	websvr.Router("/api/delete-organization", &controllers.ApiController{}, "POST:DeleteOrganization")

	websvr.Router("/api/get-global-users", &controllers.ApiController{}, "GET:GetGlobalUsers")
	websvr.Router("/api/get-users", &controllers.ApiController{}, "GET:GetUsers")
	websvr.Router("/api/get-sorted-users", &controllers.ApiController{}, "GET:GetSortedUsers")
	websvr.Router("/api/get-user-count", &controllers.ApiController{}, "GET:GetUserCount")
	websvr.Router("/api/get-user", &controllers.ApiController{}, "GET:GetUser")
	websvr.Router("/api/update-user", &controllers.ApiController{}, "POST:UpdateUser")
	websvr.Router("/api/add-user", &controllers.ApiController{}, "POST:AddUser")
	websvr.Router("/api/delete-user", &controllers.ApiController{}, "POST:DeleteUser")
	websvr.Router("/api/upload-users", &controllers.ApiController{}, "POST:UploadUsers")

	websvr.Router("/api/get-roles", &controllers.ApiController{}, "GET:GetRoles")
	websvr.Router("/api/get-role", &controllers.ApiController{}, "GET:GetRole")
	websvr.Router("/api/update-role", &controllers.ApiController{}, "POST:UpdateRole")
	websvr.Router("/api/add-role", &controllers.ApiController{}, "POST:AddRole")
	websvr.Router("/api/delete-role", &controllers.ApiController{}, "POST:DeleteRole")

	websvr.Router("/api/get-permissions", &controllers.ApiController{}, "GET:GetPermissions")
	websvr.Router("/api/get-permission", &controllers.ApiController{}, "GET:GetPermission")
	websvr.Router("/api/update-permission", &controllers.ApiController{}, "POST:UpdatePermission")
	websvr.Router("/api/add-permission", &controllers.ApiController{}, "POST:AddPermission")
	websvr.Router("/api/delete-permission", &controllers.ApiController{}, "POST:DeletePermission")

	websvr.Router("/api/set-password", &controllers.ApiController{}, "POST:SetPassword")
	websvr.Router("/api/check-user-password", &controllers.ApiController{}, "POST:CheckUserPassword")
	websvr.Router("/api/get-email-and-phone", &controllers.ApiController{}, "POST:GetEmailAndPhone")
	websvr.Router("/api/send-verification-code", &controllers.ApiController{}, "POST:SendVerificationCode")
	websvr.Router("/api/reset-email-or-phone", &controllers.ApiController{}, "POST:ResetEmailOrPhone")
	websvr.Router("/api/get-human-check", &controllers.ApiController{}, "GET:GetHumanCheck")

	websvr.Router("/api/get-ldap-user", &controllers.ApiController{}, "POST:GetLdapUser")
	websvr.Router("/api/get-ldaps", &controllers.ApiController{}, "POST:GetLdaps")
	websvr.Router("/api/get-ldap", &controllers.ApiController{}, "POST:GetLdap")
	websvr.Router("/api/add-ldap", &controllers.ApiController{}, "POST:AddLdap")
	websvr.Router("/api/update-ldap", &controllers.ApiController{}, "POST:UpdateLdap")
	websvr.Router("/api/delete-ldap", &controllers.ApiController{}, "POST:DeleteLdap")
	websvr.Router("/api/check-ldap-users-exist", &controllers.ApiController{}, "POST:CheckLdapUsersExist")
	websvr.Router("/api/sync-ldap-users", &controllers.ApiController{}, "POST:SyncLdapUsers")

	websvr.Router("/api/get-providers", &controllers.ApiController{}, "GET:GetProviders")
	websvr.Router("/api/get-provider", &controllers.ApiController{}, "GET:GetProvider")
	websvr.Router("/api/update-provider", &controllers.ApiController{}, "POST:UpdateProvider")
	websvr.Router("/api/add-provider", &controllers.ApiController{}, "POST:AddProvider")
	websvr.Router("/api/delete-provider", &controllers.ApiController{}, "POST:DeleteProvider")

	websvr.Router("/api/get-applications", &controllers.ApiController{}, "GET:GetApplications")
	websvr.Router("/api/get-application", &controllers.ApiController{}, "GET:GetApplication")
	websvr.Router("/api/get-user-application", &controllers.ApiController{}, "GET:GetUserApplication")
	websvr.Router("/api/update-application", &controllers.ApiController{}, "POST:UpdateApplication")
	websvr.Router("/api/add-application", &controllers.ApiController{}, "POST:AddApplication")
	websvr.Router("/api/delete-application", &controllers.ApiController{}, "POST:DeleteApplication")

	websvr.Router("/api/get-resources", &controllers.ApiController{}, "GET:GetResources")
	websvr.Router("/api/get-resource", &controllers.ApiController{}, "GET:GetResource")
	websvr.Router("/api/update-resource", &controllers.ApiController{}, "POST:UpdateResource")
	websvr.Router("/api/add-resource", &controllers.ApiController{}, "POST:AddResource")
	websvr.Router("/api/delete-resource", &controllers.ApiController{}, "POST:DeleteResource")
	websvr.Router("/api/upload-resource", &controllers.ApiController{}, "POST:UploadResource")

	websvr.Router("/api/get-tokens", &controllers.ApiController{}, "GET:GetTokens")
	websvr.Router("/api/get-token", &controllers.ApiController{}, "GET:GetToken")
	websvr.Router("/api/update-token", &controllers.ApiController{}, "POST:UpdateToken")
	websvr.Router("/api/add-token", &controllers.ApiController{}, "POST:AddToken")
	websvr.Router("/api/delete-token", &controllers.ApiController{}, "POST:DeleteToken")
	websvr.Router("/api/login/oauth/code", &controllers.ApiController{}, "POST:GetOAuthCode")
	websvr.Router("/api/login/oauth/access_token", &controllers.ApiController{}, "POST:GetOAuthToken")
	websvr.Router("/api/login/oauth/refresh_token", &controllers.ApiController{}, "POST:RefreshToken")

	websvr.Router("/api/get-records", &controllers.ApiController{}, "GET:GetRecords")
	websvr.Router("/api/get-records-filter", &controllers.ApiController{}, "POST:GetRecordsByFilter")

	websvr.Router("/api/get-webhooks", &controllers.ApiController{}, "GET:GetWebhooks")
	websvr.Router("/api/get-webhook", &controllers.ApiController{}, "GET:GetWebhook")
	websvr.Router("/api/update-webhook", &controllers.ApiController{}, "POST:UpdateWebhook")
	websvr.Router("/api/add-webhook", &controllers.ApiController{}, "POST:AddWebhook")
	websvr.Router("/api/delete-webhook", &controllers.ApiController{}, "POST:DeleteWebhook")

	websvr.Router("/api/get-syncers", &controllers.ApiController{}, "GET:GetSyncers")
	websvr.Router("/api/get-syncer", &controllers.ApiController{}, "GET:GetSyncer")
	websvr.Router("/api/update-syncer", &controllers.ApiController{}, "POST:UpdateSyncer")
	websvr.Router("/api/add-syncer", &controllers.ApiController{}, "POST:AddSyncer")
	websvr.Router("/api/delete-syncer", &controllers.ApiController{}, "POST:DeleteSyncer")

	websvr.Router("/api/get-certs", &controllers.ApiController{}, "GET:GetCerts")
	websvr.Router("/api/get-cert", &controllers.ApiController{}, "GET:GetCert")
	websvr.Router("/api/update-cert", &controllers.ApiController{}, "POST:UpdateCert")
	websvr.Router("/api/add-cert", &controllers.ApiController{}, "POST:AddCert")
	websvr.Router("/api/delete-cert", &controllers.ApiController{}, "POST:DeleteCert")

	websvr.Router("/api/send-email", &controllers.ApiController{}, "POST:SendEmail")
	websvr.Router("/api/send-sms", &controllers.ApiController{}, "POST:SendSms")

	websvr.Router("/.well-known/openid-configuration", &controllers.RootController{}, "GET:GetOidcDiscovery")
	websvr.Router("/api/certs", &controllers.RootController{}, "*:GetOidcCert")
}
