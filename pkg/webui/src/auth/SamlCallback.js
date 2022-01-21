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
import {Spin} from "antd";
import {withRouter} from "react-router-dom";
import * as AuthBackend from "./AuthBackend";
import * as Util from "./Util";
import * as Setting from "../Setting";
import i18next from "i18next";
import {authConfig} from "./Auth";

class SamlCallback extends React.Component {
    constructor(props) {
      super(props);
      this.state = {
        classes: props,
        msg: null,
      };
    }

    getResponseType(redirectUri) {
        const authServerUrl = authConfig.serverUrl;
        // Bhojpur IAM's own login page, so "code" is not necessary
        if (redirectUri === "null") {
            return "login";
        }
        const realRedirectUrl = new URL(redirectUri).origin;
        // For Bhojpur IAM itself, we use "login" directly
        if (authServerUrl === realRedirectUrl) {
            return "login";
        } else {
            return "code";
        }
    }

    UNSAFE_componentWillMount() {
        const params = new URLSearchParams(this.props.location.search);
        let relayState = params.get('relayState')
        let samlResponse = params.get('samlResponse')
        const messages = atob(relayState).split('&');
        const clientId = messages[0];
        const applicationName = messages[1] === "null" ? "app-built-in" : messages[1];
        const providerName = messages[2];
        const redirectUri = messages[3];
        const responseType = this.getResponseType(redirectUri);

        const body = {
            type: responseType,
            application: applicationName,
            provider: providerName,
            state: applicationName,
            redirectUri: `${window.location.origin}/callback`,
            method: "signup",
            relayState: relayState,
            samlResponse: encodeURIComponent(samlResponse),
          };

        let param;
        if (clientId === null || clientId === "") {
            param = ""
        } else {
            param = `?clientId=${clientId}&responseType=${responseType}&redirectUri=${redirectUri}&scope=read&state=${applicationName}`
        }

        AuthBackend.loginWithSaml(body, param)
          .then((res) => {
            if (res.status === 'ok') {
                const responseType = this.getResponseType(redirectUri);
                if (responseType === "login") {
                    Util.showMessage("success", `Logged in successfully`);
                    Setting.goToLink("/");
                } else if (responseType === "code") {
                    const code = res.data;
                    Setting.goToLink(`${redirectUri}?code=${code}&state=${applicationName}`);
                }
              } else {
                this.setState({
                  msg: res.msg,
                });
              }
          });
    }

    render() {
        return (
          <div style={{textAlign: "center"}}>
            {
              (this.state.msg === null) ? (
                <Spin size="large" tip={i18next.t("login:Signing in...")} style={{paddingTop: "10%"}} />
              ) : (
                Util.renderMessageLarge(this, this.state.msg)
              )
            }
          </div>
        )
      }
}
export default withRouter(SamlCallback);