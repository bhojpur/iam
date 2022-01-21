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

export function getSyncers(owner, page = "", pageSize = "", field = "", value = "", sortField = "", sortOrder = "") {
  return fetch(`${Setting.ServerUrl}/api/get-syncers?owner=${owner}&p=${page}&pageSize=${pageSize}&field=${field}&value=${value}&sortField=${sortField}&sortOrder=${sortOrder}`, {
    method: "GET",
    credentials: "include"
  }).then(res => res.json());
}

export function getSyncer(owner, name) {
  return fetch(`${Setting.ServerUrl}/api/get-syncer?id=${owner}/${encodeURIComponent(name)}`, {
    method: "GET",
    credentials: "include"
  }).then(res => res.json());
}

export function updateSyncer(owner, name, syncer) {
  let newSyncer = Setting.deepCopy(syncer);
  return fetch(`${Setting.ServerUrl}/api/update-syncer?id=${owner}/${encodeURIComponent(name)}`, {
    method: 'POST',
    credentials: 'include',
    body: JSON.stringify(newSyncer),
  }).then(res => res.json());
}

export function addSyncer(syncer) {
  let newSyncer = Setting.deepCopy(syncer);
  return fetch(`${Setting.ServerUrl}/api/add-syncer`, {
    method: 'POST',
    credentials: 'include',
    body: JSON.stringify(newSyncer),
  }).then(res => res.json());
}

export function deleteSyncer(syncer) {
  let newSyncer = Setting.deepCopy(syncer);
  return fetch(`${Setting.ServerUrl}/api/delete-syncer`, {
    method: 'POST',
    credentials: 'include',
    body: JSON.stringify(newSyncer),
  }).then(res => res.json());
}