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

import React from "react";
import {Alert, Button, message, Result} from "antd";

export function showMessage(type, text) {
  if (type === "success") {
    message.success(text);
  } else if (type === "error") {
    message.error(text);
  }
}

export function renderMessage(msg) {
  if (msg !== null) {
    return (
      <div style={{display: "inline"}}>
        <Alert
          message="Login Error"
          showIcon
          description={msg}
          type="error"
          action={
            <Button size="small" danger>
              Detail
            </Button>
          }
        />
      </div>
    )
  } else {
    return null;
  }
}

export function renderMessageLarge(ths, msg) {
  if (msg !== null) {
    return (
      <div style={{display: "inline"}}>
        <Result
          status="error"
          title="There was a problem signing you in.."
          subTitle={msg}
          extra={[
            <Button type="primary" key="back" onClick={() => {
              window.history.go(-2);
            }}>
              Back
            </Button>,
            // <Button key="home" onClick={() => Setting.goToLinkSoft(ths, "/")}>
            //   Home
            // </Button>,
            // <Button type="primary" key="signup" onClick={() => Setting.goToLinkSoft(ths, "/signup")}>
            //   Sign Up
            // </Button>,
          ]}
        >
        </Result>
      </div>
    )
  } else {
    return null;
  }
}

export function getOAuthGetParameters(params) {
  const queries = (params !== undefined) ? params : new URLSearchParams(window.location.search);
  const clientId = queries.get("client_id");
  const responseType = queries.get("response_type");
  const redirectUri = queries.get("redirect_uri");
  const scope = queries.get("scope");
  const state = queries.get("state");
  const nonce = queries.get("nonce")
  const challengeMethod = queries.get("code_challenge_method")
  const codeChallenge = queries.get("code_challenge")
  
  if (clientId === undefined || clientId === null) {
    // login
    return null;
  } else {
    // code
    return {
      clientId: clientId,
      responseType: responseType,
      redirectUri: redirectUri,
      scope: scope,
      state: state,
      nonce: nonce,
      challengeMethod: challengeMethod,
      codeChallenge: codeChallenge,
    };
  }
}

export function getQueryParamsToState(applicationName, providerName, method) {
  let query = window.location.search;
  query = `${query}&application=${applicationName}&provider=${providerName}&method=${method}`;
  if (method === "link") {
    query = `${query}&from=${window.location.pathname}`;
  }
  return btoa(query);
}

export function stateToGetQueryParams(state) {
  return atob(state);
}