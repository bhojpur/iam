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

import * as Setting from "../Setting";

export function getLdaps(owner) {
  return fetch(`${Setting.ServerUrl}/api/get-ldaps?owner=${owner}`, {
    method: "POST",
    credentials: "include",
  }).then(res => res.json());
}

export function getLdap(id) {
  return fetch(`${Setting.ServerUrl}/api/get-ldap?id=${id}`, {
    method: "POST",
    credentials: "include",
  }).then(res => res.json());
}

export function addLdap(body) {
  return fetch(`${Setting.ServerUrl}/api/add-ldap`, {
    method: "POST",
    credentials: "include",
    body: JSON.stringify(body),
  }).then(res => res.json());
}

export function deleteLdap(body) {
  return fetch(`${Setting.ServerUrl}/api/delete-ldap`, {
    method: "POST",
    credentials: "include",
    body: JSON.stringify(body),
  }).then(res => res.json());
}

export function updateLdap(body) {
  return fetch(`${Setting.ServerUrl}/api/update-ldap`, {
    method: "POST",
    credentials: "include",
    body: JSON.stringify(body),
  }).then(res => res.json());
}

export function getLdapUser(body) {
  return fetch(`${Setting.ServerUrl}/api/get-ldap-user`, {
    method: "POST",
    credentials: "include",
    body: JSON.stringify(body),
  }).then(res => res.json());
}

export function syncUsers(owner, ldapId, body) {
  return fetch(`${Setting.ServerUrl}/api/sync-ldap-users?owner=${owner}&ldapId=${ldapId}`, {
    method: "POST",
    credentials: "include",
    body: JSON.stringify(body),
  }).then(res => res.json());
}

export function checkLdapUsersExist(owner, body) {
  return fetch(`${Setting.ServerUrl}/api/check-ldap-users-exist?owner=${owner}`, {
    method: "POST",
    credentials: "include",
    body: JSON.stringify(body),
  }).then(res => res.json());
}