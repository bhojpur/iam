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
import {Link} from "react-router-dom";
import {Button, Checkbox, Col, Form, Input, Result, Row, Spin} from "antd";
import {LockOutlined, UserOutlined} from "@ant-design/icons";
import * as AuthBackend from "./AuthBackend";
import * as ApplicationBackend from "../backend/ApplicationBackend";
import * as Provider from "./Provider";
import * as Util from "./Util";
import * as Setting from "../Setting";
import SelfLoginButton from "./SelfLoginButton";
import {GithubLoginButton, GoogleLoginButton} from "react-social-login-buttons";
import FacebookLoginButton from "./FacebookLoginButton";
import QqLoginButton from "./QqLoginButton";
import DingTalkLoginButton from "./DingTalkLoginButton";
import GiteeLoginButton from "./GiteeLoginButton";
import WechatLoginButton from "./WechatLoginButton";
import WeiboLoginButton from "./WeiboLoginButton";
import i18next from "i18next";
import LinkedInLoginButton from "./LinkedInLoginButton";
import WeComLoginButton from "./WeComLoginButton";
import LarkLoginButton from "./LarkLoginButton";
import GitLabLoginButton from "./GitLabLoginButton";
import AppleLoginButton from "./AppleLoginButton"
import AzureADLoginButton from "./AzureADLoginButton";
import SlackLoginButton from "./SlackLoginButton";
import CustomGithubCorner from "../CustomGithubCorner";
import {CountDownInput} from "../common/CountDownInput";

class LoginPage extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      classes: props,
      type: props.type,
      applicationName: props.applicationName !== undefined ? props.applicationName : (props.match === undefined ? null : props.match.params.applicationName),
      application: null,
      mode: props.mode !== undefined ? props.mode : (props.match === undefined ? null : props.match.params.mode), // "signup" or "signin"
      isCodeSignin: false,
      msg: null,
      username: null,
      validEmailOrPhone: false
    };
  }

  UNSAFE_componentWillMount() {
    if (this.state.type === "login") {
      this.getApplication();
    } else if (this.state.type === "code") {
      this.getApplicationLogin();
    } else {
      Util.showMessage("error", `Unknown authentication type: ${this.state.type}`);
    }
  }

  getApplicationLogin() {
    const oAuthParams = Util.getOAuthGetParameters();
    AuthBackend.getApplicationLogin(oAuthParams)
      .then((res) => {
        if (res.status === "ok") {
          this.setState({
            application: res.data,
          });
        } else {
          // Util.showMessage("error", res.msg);
          this.setState({
            application: res.data,
            msg: res.msg,
          });
        }
      });
  }

  getApplication() {
    if (this.state.applicationName === null) {
      return;
    }

    ApplicationBackend.getApplication("admin", this.state.applicationName)
      .then((application) => {
        this.setState({
          application: application,
        });
      });
  }

  getApplicationObj() {
    if (this.props.application !== undefined) {
      return this.props.application;
    } else {
      return this.state.application;
    }
  }

  onUpdateAccount(account) {
    this.props.onUpdateAccount(account);
  }

  onFinish(values) {
    const application = this.getApplicationObj();
    const ths = this;
    values["type"] = this.state.type;
    values["phonePrefix"] = this.getApplicationObj()?.organizationObj.phonePrefix;
    const oAuthParams = Util.getOAuthGetParameters();

    AuthBackend.login(values, oAuthParams)
      .then((res) => {
        if (res.status === 'ok') {
          const responseType = this.state.type;
          if (responseType === "login") {
            Util.showMessage("success", `Logged in successfully`);
            Setting.goToLink("/");
          } else if (responseType === "code") {
            const code = res.data;

            if (Setting.hasPromptPage(application)) {
              AuthBackend.getAccount("")
                .then((res) => {
                  let account = null;
                  if (res.status === "ok") {
                    account = res.data;
                    account.organization = res.data2;

                    this.onUpdateAccount(account);

                    if (Setting.isPromptAnswered(account, application)) {
                      Setting.goToLink(`${oAuthParams.redirectUri}?code=${code}&state=${oAuthParams.state}`);
                    } else {
                      Setting.goToLinkSoft(ths, `/prompt/${application.name}?redirectUri=${oAuthParams.redirectUri}&code=${code}&state=${oAuthParams.state}`);
                    }
                  } else {
                    Setting.showMessage("error", `Failed to sign in: ${res.msg}`);
                  }
                });
            } else {
              Setting.goToLink(`${oAuthParams.redirectUri}?code=${code}&state=${oAuthParams.state}`);
            }

            // Util.showMessage("success", `Authorization code: ${res.data}`);
          }
        } else {
          Util.showMessage("error", `Failed to log in: ${res.msg}`);
        }
      });
  };

  getSigninButton(type) {
    const text = i18next.t("login:Sign in with {type}").replace("{type}", type);
    if (type === "GitHub") {
      return <GithubLoginButton text={text} align={"center"} />
    } else if (type === "Google") {
      return <GoogleLoginButton text={text} align={"center"} />
    } else if (type === "QQ") {
      return <QqLoginButton text={text} align={"center"} />
    } else if (type === "Facebook") {
      return <FacebookLoginButton text={text} align={"center"} />
    } else if (type === "Weibo") {
      return <WeiboLoginButton text={text} align={"center"} />
    } else if (type === "Gitee") {
      return <GiteeLoginButton text={text} align={"center"} />
    } else if (type === "WeChat") {
      return <WechatLoginButton text={text} align={"center"} />
    } else if (type === "DingTalk") {
      return <DingTalkLoginButton text={text} align={"center"} />
    } else if (type === "LinkedIn"){
      return <LinkedInLoginButton text={text} align={"center"} />
    } else if (type === "WeCom") {
      return <WeComLoginButton text={text} align={"center"} />
    } else if (type === "Lark") {
      return <LarkLoginButton text={text} align={"center"} />
    } else if (type === "GitLab") {
      return <GitLabLoginButton text={text} align={"center"} />
    } else if (type === "Apple") {
      return <AppleLoginButton text={text} align={"center"} />
    } else if (type === "AzureAD") {
      return <AzureADLoginButton text={text} align={"center"} />
    } else if (type === "Slack") {
      return <SlackLoginButton text={text} align={"center"} />
    }

    return text;
  }

  getSamlUrl(provider) {
    const params = new URLSearchParams(this.props.location.search);
    let clientId = params.get("client_id");
    let application = params.get("state");
    let realRedirectUri = params.get("redirect_uri");
    let redirectUri = `${window.location.origin}/callback/saml`;
    let providerName = provider.name;
    let relayState = `${clientId}&${application}&${providerName}&${realRedirectUri}&${redirectUri}`;
    AuthBackend.getSamlLogin(`${provider.owner}/${providerName}`, btoa(relayState)).then((res) => {
      if (res.data2 === "POST") {
        document.write(res.data)
      } else {
        window.location.href = res.data
      }
    });
  }

  renderProviderLogo(provider, application, width, margin, size) {
    if (size === "small") {
      if (provider.category === "OAuth") {
        return (
          <a key={provider.displayName} href={Provider.getAuthUrl(application, provider, "signup")}>
            <img width={width} height={width} src={Provider.getProviderLogo(provider)} alt={provider.displayName} style={{margin: margin}} />
          </a>
        )
      } else if (provider.category === "SAML") {
        return (
          <a key={provider.displayName} onClick={this.getSamlUrl.bind(this, provider)}>
            <img width={width} height={width} src={Provider.getProviderLogo(provider)} alt={provider.displayName} style={{margin: margin}} />
          </a>
        )
      }
      
    } else {
      return (
        <div key={provider.displayName} style={{marginBottom: "10px"}}>
          <a href={Provider.getAuthUrl(application, provider, "signup")}>
            {
              this.getSigninButton(provider.type)
            }
          </a>
        </div>
      )
    }
  }

  isProviderVisible(providerItem) {
    if (this.state.mode === "signup") {
      return Setting.isProviderVisibleForSignUp(providerItem);
    } else {
      return Setting.isProviderVisibleForSignIn(providerItem);
    }
  }

  renderForm(application) {
    if (this.state.msg !== null) {
      return Util.renderMessage(this.state.msg)
    }

    if (this.state.mode === "signup" && !application.enableSignUp) {
      return (
        <Result
          status="error"
          title="Sign Up Error"
          subTitle={"The application does not allow to sign up new account"}
          extra={[
            <Link onClick={() => {
              Setting.goToLogin(this, application);
            }}>
              <Button type="primary" key="signin">
                Sign In
              </Button>
            </Link>
          ]}
        >
        </Result>
      )
    }

    if (application.enablePassword) {
      return (
        <Form
          name="normal_login"
          initialValues={{
            organization: application.organization,
            application: application.name,
            autoSignin: true,
          }}
          onFinish={(values) => {this.onFinish(values)}}
          style={{width: "300px"}}
          size="large"
        >
          <Form.Item
            style={{height: 0, visibility: "hidden"}}
            name="application"
            rules={[
              {
                required: true,
                message: 'Please input your application!',
              },
            ]}
          >
          </Form.Item>
          <Form.Item
            style={{height: 0, visibility: "hidden"}}
            name="organization"
            rules={[
              {
                required: true,
                message: 'Please input your organization!',
              },
            ]}
          >
          </Form.Item>
          <Form.Item
            name="username"
            rules={[
                {
                  required: true,
                  message: i18next.t("login:Please input your username, Email or phone!")
                },
                {
                  validator: (_, value) => {
                    if (this.state.isCodeSignin) {
                      if (this.state.email !== "" && !Setting.isValidEmail(this.state.username) && !Setting.isValidPhone(this.state.username)) {
                        this.setState({validEmailOrPhone: false});
                        return Promise.reject(i18next.t("login:The input is not valid Email or Phone!"));
                      }
                    }
                    this.setState({validEmailOrPhone: true});
                    return Promise.resolve();
                  }
                }
              ]}
          >
            <Input
              prefix={<UserOutlined className="site-form-item-icon" />}
              placeholder={ this.state.isCodeSignin ? i18next.t("login:Email or phone") : i18next.t("login:username, Email or phone") }
              disabled={!application.enablePassword}
              onChange={e => {
                this.setState({
                  username: e.target.value,
                });
              }}
            />
          </Form.Item>
          {
            this.state.isCodeSignin ? (
              <Form.Item
                name="code"
                rules={[{ required: true, message: i18next.t("login:Please input your code!") }]}
              >
                <CountDownInput
                  disabled={this.state.username?.length === 0 || !this.state.validEmailOrPhone}
                  onButtonClickArgs={[this.state.username, "", Setting.getApplicationOrgName(application), true]}
                />
              </Form.Item>
            ) : (
              <Form.Item
                name="password"
                rules={[{ required: true, message: i18next.t("login:Please input your password!") }]}
              >
                <Input
                  prefix={<LockOutlined className="site-form-item-icon" />}
                  type="password"
                  placeholder={i18next.t("login:Password")}
                  disabled={!application.enablePassword}
                />
              </Form.Item>
            )
          }
          <Form.Item>
            <Form.Item name="autoSignin" valuePropName="checked" noStyle>
              <Checkbox style={{float: "left"}} disabled={!application.enablePassword}>
                {i18next.t("login:Auto sign in")}
              </Checkbox>
            </Form.Item>
            <a style={{float: "right"}} onClick={() => {
              Setting.goToForget(this, application);
            }}>
              {i18next.t("login:Forgot password?")}
            </a>
          </Form.Item>
          <Form.Item>
            <Button
              type="primary"
              htmlType="submit"
              style={{width: "100%", marginBottom: '5px'}}
              disabled={!application.enablePassword}
            >
              {i18next.t("login:Sign In")}
            </Button>
            {
              !application.enableSignUp ? null : this.renderFooter(application)
            }
          </Form.Item>
          <Form.Item>
            {
              application.providers.filter(providerItem => this.isProviderVisible(providerItem)).map(providerItem => {
                return this.renderProviderLogo(providerItem.provider, application, 30, 5, "small");
              })
            }
          </Form.Item>
        </Form>
      );
    } else {
      return (
        <div style={{marginTop: "20px"}}>
          <div style={{fontSize: 16, textAlign: "left"}}>
            {i18next.t("login:To access")}&nbsp;
            <a target="_blank" rel="noreferrer" href={application.homepageUrl}>
              <span style={{fontWeight: "bold"}}>
                {application.displayName}
              </span>
            </a>
            :
          </div>
          <br/>
          {
            application.providers.filter(providerItem => this.isProviderVisible(providerItem)).map(providerItem => {
              return this.renderProviderLogo(providerItem.provider, application, 40, 10, "big");
            })
          }
          {
            !application.enableSignUp ? null : (
              <div>
                <br/>
                {
                  this.renderFooter(application)
                }
              </div>
            )
          }
        </div>
      )
    }
  }

  renderFooter(application) {
    if (this.state.mode === "signup") {
      return (
        <div style={{float: "right"}}>
          {i18next.t("signup:Have account?")}&nbsp;
          <Link onClick={() => {
            Setting.goToLogin(this, application);
          }}>
            {i18next.t("signup:sign in now")}
          </Link>
        </div>
      )
    } else {
      return (
        <React.Fragment>
          <span style={{float: "left"}}>
            {
              !application.enableCodeSignin ? null : (
                <a onClick={() => {
                  this.setState({
                    isCodeSignin: !this.state.isCodeSignin,
                  });
                }}>
                  {this.state.isCodeSignin ? i18next.t("login:Sign in with password") : i18next.t("login:Sign in with code")}
                </a>
              )
            }
          </span>
          <span style={{float: "right"}}>
            {i18next.t("login:No account?")}&nbsp;
            <a onClick={() => {
              Setting.goToSignup(this, application);
            }}>
              {i18next.t("login:sign up now")}
            </a>
          </span>
        </React.Fragment>
      )
    }
  }

  renderSignedInBox() {
    if (this.props.account === undefined || this.props.account === null) {
      return null;
    }
    let application = this.getApplicationObj()
    if (this.props.account.owner !== application.organization) {
      return null;
    }

    return (
      <div>
        <div style={{fontSize: 16, textAlign: "left"}}>
          {i18next.t("login:Continue with")}&nbsp;:
        </div>
        <br/>
        <SelfLoginButton account={this.props.account} onClick={() => {
          let values = {};
          values["application"] = this.state.application.name;
          this.onFinish(values);
        }} />
        <br/>
        <br/>
        <div style={{fontSize: 16, textAlign: "left"}}>
          {i18next.t("login:Or sign in with another account")}&nbsp;:
        </div>
      </div>
    )
  }

  render() {
    const application = this.getApplicationObj();
    if (application === null) {
      return Util.renderMessageLarge(this, this.state.msg);
    }

    if (application.signinHtml !== "") {
      return (
        <div dangerouslySetInnerHTML={{ __html: application.signinHtml}} />
      )
    }

    const visibleOAuthProviderItems = application.providers.filter(providerItem => this.isProviderVisible(providerItem));
    if (this.props.application === undefined && !application.enablePassword && visibleOAuthProviderItems.length === 1) {
      Setting.goToLink(Provider.getAuthUrl(application, visibleOAuthProviderItems[0].provider, "signup"));
      return (
        <div style={{textAlign: "center"}}>
          <Spin size="large" tip={i18next.t("login:Signing in...")} style={{paddingTop: "10%"}} />
        </div>
      )
    }

    return (
      <Row>
        <Col span={24} style={{display: "flex", justifyContent: "center"}}>
          <div style={{marginTop: "80px", marginBottom: "50px", textAlign: "center"}}>
            {
              Setting.renderHelmet(application)
            }
            <CustomGithubCorner />
            {
              Setting.renderLogo(application)
            }
            {/*{*/}
            {/*  this.state.clientId !== null ? "Redirect" : null*/}
            {/*}*/}
            {
              this.renderSignedInBox()
            }
            {
              this.renderForm(application)
            }
          </div>
        </Col>
      </Row>
    )
  }
}

export default LoginPage;