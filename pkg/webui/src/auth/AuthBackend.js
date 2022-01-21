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

import {authConfig} from "./Auth";

export function getAccount(query) {
  return fetch(`${authConfig.serverUrl}/api/get-account${query}`, {
    method: 'GET',
    credentials: 'include'
  }).then(res => res.json());
}

export function signup(values) {
  return fetch(`${authConfig.serverUrl}/api/signup`, {
    method: 'POST',
    credentials: "include",
    body: JSON.stringify(values),
  }).then(res => res.json());
}

export function getEmailAndPhone(values) {
  return fetch(`${authConfig.serverUrl}/api/get-email-and-phone`, {
    method: "POST",
    credentials: "include",
    body: JSON.stringify(values),
  }).then((res) => res.json());
}

function oAuthParamsToQuery(oAuthParams) {
  // login
  if (oAuthParams === null) {
    return "";
  }

  // code
  return `?clientId=${oAuthParams.clientId}&responseType=${oAuthParams.responseType}&redirectUri=${oAuthParams.redirectUri}&scope=${oAuthParams.scope}&state=${oAuthParams.state}&nonce=${oAuthParams.nonce}&code_challenge_method=${oAuthParams.challengeMethod}&code_challenge=${oAuthParams.codeChallenge}`;
}

export function getApplicationLogin(oAuthParams) {
  return fetch(`${authConfig.serverUrl}/api/get-app-login${oAuthParamsToQuery(oAuthParams)}`, {
    method: 'GET',
    credentials: 'include',
  }).then(res => res.json());
}

export function login(values, oAuthParams) {
  return fetch(`${authConfig.serverUrl}/api/login${oAuthParamsToQuery(oAuthParams)}`, {
    method: 'POST',
    credentials: "include",
    body: JSON.stringify(values),
  }).then(res => res.json());
}

export function logout() {
  return fetch(`${authConfig.serverUrl}/api/logout`, {
    method: 'POST',
    credentials: "include",
  }).then(res => res.json());
}

export function unlink(values) {
  return fetch(`${authConfig.serverUrl}/api/unlink`, {
    method: 'POST',
    credentials: "include",
    body: JSON.stringify(values),
  }).then(res => res.json());
}

export function getSamlLogin(providerId, relayState) {
  return fetch(`${authConfig.serverUrl}/api/get-saml-login?id=${providerId}&relayState=${relayState}`, {
    method: 'GET',
    credentials: 'include',
  }).then(res => res.json());
}

export function loginWithSaml(values, param) {
  return fetch(`${authConfig.serverUrl}/api/login${param}`, {
    method: 'POST',
    credentials: "include",
    body: JSON.stringify(values),
  }).then(res => res.json());
}