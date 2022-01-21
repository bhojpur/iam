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
import i18next from "i18next";

export function getGlobalUsers(page, pageSize, field = "", value = "", sortField = "", sortOrder = "") {
  return fetch(`${Setting.ServerUrl}/api/get-global-users?p=${page}&pageSize=${pageSize}&field=${field}&value=${value}&sortField=${sortField}&sortOrder=${sortOrder}`, {
    method: "GET",
    credentials: "include"
  }).then(res => res.json());
}

export function getUsers(owner, page = "", pageSize = "", field = "", value = "", sortField = "", sortOrder = "") {
  return fetch(`${Setting.ServerUrl}/api/get-users?owner=${owner}&p=${page}&pageSize=${pageSize}&field=${field}&value=${value}&sortField=${sortField}&sortOrder=${sortOrder}`, {
    method: "GET",
    credentials: "include"
  }).then(res => res.json());
}

export function getUser(owner, name) {
  return fetch(`${Setting.ServerUrl}/api/get-user?id=${owner}/${encodeURIComponent(name)}`, {
    method: "GET",
    credentials: "include"
  }).then(res => res.json());
}

export function updateUser(owner, name, user) {
  let newUser = Setting.deepCopy(user);
  return fetch(`${Setting.ServerUrl}/api/update-user?id=${owner}/${encodeURIComponent(name)}`, {
    method: 'POST',
    credentials: 'include',
    body: JSON.stringify(newUser),
  }).then(res => res.json());
}

export function addUser(user) {
  let newUser = Setting.deepCopy(user);
  return fetch(`${Setting.ServerUrl}/api/add-user`, {
    method: 'POST',
    credentials: 'include',
    body: JSON.stringify(newUser),
  }).then(res => res.json());
}

export function deleteUser(user) {
  let newUser = Setting.deepCopy(user);
  return fetch(`${Setting.ServerUrl}/api/delete-user`, {
    method: 'POST',
    credentials: 'include',
    body: JSON.stringify(newUser),
  }).then(res => res.json());
}

export function getAddressOptions(url) {
  return fetch(url, {
    method: "GET",
  }).then(res => res.json());
}

export function getAffiliationOptions(url, code) {
  return fetch(`${url}/${code}`, {
    method: "GET",
  }).then(res => res.json());
}

export function setPassword(userOwner, userName, oldPassword, newPassword) {
  let formData = new FormData();
  formData.append("userOwner", userOwner);
  formData.append("userName", userName);
  formData.append("oldPassword", oldPassword);
  formData.append("newPassword", newPassword);
  return fetch(`${Setting.ServerUrl}/api/set-password`, {
    method: "POST",
    credentials: "include",
    body: formData
  }).then(res => res.json());
}

export function sendCode(checkType, checkId, checkKey, dest, type, orgId, checkUser) {
  let formData = new FormData();
  formData.append("checkType", checkType);
  formData.append("checkId", checkId);
  formData.append("checkKey", checkKey);
  formData.append("dest", dest);
  formData.append("type", type);
  formData.append("organizationId", orgId);
  formData.append("checkUser", checkUser);
  return fetch(`${Setting.ServerUrl}/api/send-verification-code`, {
    method: "POST",
    credentials: "include",
    body: formData
  }).then(res => res.json()).then(res => {
    if (res.status === "ok") {
      Setting.showMessage("success", i18next.t("user:Code Sent"));
      return true;
    } else {
      Setting.showMessage("error", i18next.t("user:" + res.msg));
      return false;
    }
  });
}

export function resetEmailOrPhone(dest, type, code) {
  let formData = new FormData();
  formData.append("dest", dest);
  formData.append("type", type);
  formData.append("code", code);
  return fetch(`${Setting.ServerUrl}/api/reset-email-or-phone`, {
    method: "POST",
    credentials: "include",
    body: formData
  }).then(res => res.json());
}

export function getHumanCheck() {
  return fetch(`${Setting.ServerUrl}/api/get-human-check`, {
    method: "GET"
  }).then(res => res.json());
}