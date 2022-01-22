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

import React from 'react';
import {Link} from "react-router-dom";
import {Form, Input, Checkbox, Button, Row, Col, Result, Modal} from 'antd';
import * as Setting from "../Setting";
import * as AuthBackend from "./AuthBackend";
import i18next from "i18next";
import * as Util from "./Util";
import {authConfig} from "./Auth";
import * as ApplicationBackend from "../backend/ApplicationBackend";
import {CountDownInput} from "../common/CountDownInput";
import SelectRegionBox from "../SelectRegionBox";
import CustomGithubCorner from "../CustomGithubCorner";

const formItemLayout = {
  labelCol: {
    xs: {
      span: 24,
    },
    sm: {
      span: 8,
    },
  },
  wrapperCol: {
    xs: {
      span: 24,
    },
    sm: {
      span: 18,
    },
  },
};

const tailFormItemLayout = {
  wrapperCol: {
    xs: {
      span: 24,
      offset: 0,
    },
    sm: {
      span: 16,
      offset: 8,
    },
  },
};

class SignupPage extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      classes: props,
      applicationName: props.match?.params.applicationName !== undefined ? props.match.params.applicationName : authConfig.appName,
      application: null,
      email: "",
      phone: "",
      emailCode: "",
      phoneCode: "",
      validEmail: false,
      validPhone: false,
      region: "",
      isTermsOfUseVisible: false,
      termsOfUseContent: "",
    };

    this.form = React.createRef();
  }

  UNSAFE_componentWillMount() {
    if (this.state.applicationName !== undefined) {
      this.getApplication();
    } else {
      Util.showMessage("error", `Unknown application name: ${this.state.applicationName}`);
    }
  }

  getApplication() {
    if (this.state.applicationName === undefined) {
      return;
    }

    ApplicationBackend.getApplication("admin", this.state.applicationName)
      .then((application) => {
        this.setState({
          application: application,
        });
        this.getTermsofuseContent(application.termsOfUse);
      });
  }

  getResultPath(application) {
    if (authConfig.appName === application.name) {
      return "/result";
    } else {
      if (Setting.hasPromptPage(application)) {
        return `/prompt/${application.name}`;
      } else {
        return `/result/${application.name}`;
      }
    }
  }

  getApplicationObj() {
    if (this.props.application !== undefined) {
      return this.props.application;
    } else {
      return this.state.application;
    }
  }

  getTermsofuseContent(url) {
    fetch(url, {
      method: "GET",
    }).then(r => {
      r.text().then(res => {
        this.setState({termsOfUseContent: res})
      })
    })
  }

  onUpdateAccount(account) {
    this.props.onUpdateAccount(account);
  }

  onFinish(values) {
    const application = this.getApplicationObj();
    values.phonePrefix = application.organizationObj.phonePrefix;
    AuthBackend.signup(values)
      .then((res) => {
        if (res.status === 'ok') {
          if (Setting.hasPromptPage(application)) {
            AuthBackend.getAccount("")
              .then((res) => {
                let account = null;
                if (res.status === "ok") {
                  account = res.data;
                  account.organization = res.data2;

                  this.onUpdateAccount(account);
                  Setting.goToLinkSoft(this, this.getResultPath(application));
                } else {
                  Setting.showMessage("error", `Failed to sign in: ${res.msg}`);
                }
              });
          } else {
            Setting.goToLinkSoft(this, this.getResultPath(application));
          }
        } else {
          Setting.showMessage("error", i18next.t(`signup:${res.msg}`));
        }
      });
  }

  onFinishFailed(values, errorFields, outOfDate) {
    this.form.current.scrollToField(errorFields[0].name);
  }

  renderFormItem(application, signupItem) {
    if (!signupItem.visible) {
      return null;
    }

    const required = signupItem.required;

    if (signupItem.name === "Username") {
      return (
        <Form.Item
          name="username"
          key="username"
          label={i18next.t("signup:Username")}
          rules={[
            {
              required: required,
              message: i18next.t("forget:Please input your username!"),
              whitespace: true,
            },
          ]}
        >
          <Input />
        </Form.Item>
      )
    } else if (signupItem.name === "Display name") {
      return (
        <Form.Item
          name="name"
          key="name"
          label={signupItem.rule === "Personal" ? i18next.t("general:Personal name") : i18next.t("general:Display name")}
          rules={[
            {
              required: required,
              message: signupItem.rule === "Personal" ? i18next.t("signup:Please input your personal name!") : i18next.t("signup:Please input your display name!"),
              whitespace: true,
            },
          ]}
        >
          <Input />
        </Form.Item>
      )
    } else if (signupItem.name === "Affiliation") {
      return (
        <Form.Item
          name="affiliation"
          key="affiliation"
          label={i18next.t("user:Affiliation")}
          rules={[
            {
              required: required,
              message: i18next.t("signup:Please input your affiliation!"),
              whitespace: true,
            },
          ]}
        >
          <Input />
        </Form.Item>
      )
    } else if (signupItem.name === "ID card") {
      return (
        <Form.Item
          name="idCard"
          key="idCard"
          label={i18next.t("user:ID card")}
          rules={[
            {
              required: required,
              message: i18next.t("signup:Please input your ID card number!"),
              whitespace: true,
            },
            {
              required: required,
              pattern: new RegExp(/^[1-9]\d{5}(18|19|20)\d{2}((0[1-9])|(10|11|12))(([0-2][1-9])|10|20|30|31)\d{3}[0-9X]$/, "g"),
              message: i18next.t("signup:Please input the correct ID card number!"),
            },
          ]}
        >
          <Input />
        </Form.Item>
      )
    } else if (signupItem.name === "Country/Region") {
      return (
        <Form.Item
            name="country_region"
            key="region"
            label={i18next.t("user:Country/Region")}
            rules={[
                {
                    required: required,
                    message: i18next.t("signup:Please select your country/region!"),
                },
            ]}
        >
          <SelectRegionBox onChange={(value) => {this.setState({region: value})}} />
        </Form.Item>
      )
    } else if (signupItem.name === "Email") {
      return (
        <React.Fragment>
          <Form.Item
            name="email"
            key="email"
            label={i18next.t("general:Email")}
            rules={[
              {
                required: required,
                message: i18next.t("signup:Please input your Email!"),
              },
              {
                validator: (_, value) => {
                  if (this.state.email !== "" && !Setting.isValidEmail(this.state.email)) {
                    this.setState({validEmail: false});
                    return Promise.reject(i18next.t("signup:The input is not valid Email!"));
                  }

                  this.setState({validEmail: true});
                  return Promise.resolve();
                }
              }
            ]}
          >
            <Input onChange={e => this.setState({email: e.target.value})} />
          </Form.Item>
          <Form.Item
            name="emailCode"
            key="emailCode"
            label={i18next.t("code:Email code")}
            rules={[{
              required: required,
              message: i18next.t("code:Please input your verification code!"),
            }]}
          >
            <CountDownInput
              disabled={!this.state.validEmail}
              onButtonClickArgs={[this.state.email, "email", Setting.getApplicationOrgName(application)]}
            />
          </Form.Item>
        </React.Fragment>
      )
    } else if (signupItem.name === "Phone") {
      return (
        <React.Fragment>
          <Form.Item
            name="phone"
            key="phone"
            label={i18next.t("general:Phone")}
            rules={[
              {
                required: required,
                message: i18next.t("signup:Please input your phone number!"),
              },
              {
                validator: (_, value) =>{
                  if (this.state.phone !== "" && !Setting.isValidPhone(this.state.phone)) {
                    this.setState({validPhone: false});
                    return Promise.reject(i18next.t("signup:The input is not valid Phone!"));
                  }

                  this.setState({validPhone: true});
                  return Promise.resolve();
                }
              }
            ]}
          >
            <Input
              style={{
                width: '100%',
              }}
              addonBefore={`+${this.state.application?.organizationObj.phonePrefix}`}
              onChange={e => this.setState({phone: e.target.value})}
            />
          </Form.Item>
          <Form.Item
            name="phoneCode"
            key="phoneCode"
            label={i18next.t("code:Phone code")}
            rules={[
              {
                required: required,
                message: i18next.t("code:Please input your phone verification code!"),
              },
            ]}
          >
            <CountDownInput
              disabled={!this.state.validPhone}
              onButtonClickArgs={[this.state.phone, "phone", Setting.getApplicationOrgName(application)]}
            />
          </Form.Item>
        </React.Fragment>
      )
    } else if (signupItem.name === "Password") {
      return (
        <Form.Item
          name="password"
          key="password"
          label={i18next.t("general:Password")}
          rules={[
            {
              required: required,
              min: 6,
              message: i18next.t("login:Please input your password, at least 6 characters!"),
            },
          ]}
          hasFeedback
        >
          <Input.Password />
        </Form.Item>
      )
    } else if (signupItem.name === "Confirm password") {
      return (
        <Form.Item
          name="confirm"
          key="confirm"
          label={i18next.t("signup:Confirm")}
          dependencies={['password']}
          hasFeedback
          rules={[
            {
              required: required,
              message: i18next.t("signup:Please confirm your password!"),
            },
            ({ getFieldValue }) => ({
              validator(rule, value) {
                if (!value || getFieldValue('password') === value) {
                  return Promise.resolve();
                }

                return Promise.reject(i18next.t("signup:Your confirmed password is inconsistent with the password!"));
              },
            }),
          ]}
        >
          <Input.Password />
        </Form.Item>
      )
    } else if (signupItem.name === "Agreement") {
      return (
        <Form.Item
          name="agreement"
          key="agreement"
          valuePropName="checked"
          rules={[
            {
              required: required,
              message: i18next.t("signup:Please accept the agreement!"),
            },
          ]}
          {...tailFormItemLayout}
        >
          <Checkbox>
            {i18next.t("signup:Accept")}&nbsp;
            <Link onClick={() => {
              this.setState({
                isTermsOfUseVisible: true,
              });
            }}>
              {i18next.t("signup:Terms of Use")}
            </Link>
          </Checkbox>
        </Form.Item>
      )
    }
  }

  renderModal() {
    return (
      <Modal
        title={i18next.t("signup:Terms of Use")}
        visible={this.state.isTermsOfUseVisible}
        width={"55vw"}
        closable={false}
        okText={i18next.t("signup:Accept")}
        cancelText={i18next.t("signup:Decline")}
        onOk={() => {
          this.form.current.setFieldsValue({agreement: true})
          this.setState({
            isTermsOfUseVisible: false,
          });
        }}
        onCancel={() => {
          this.form.current.setFieldsValue({agreement: false})
          this.setState({
            isTermsOfUseVisible: false,
          });
          this.props.history.goBack();
        }}
      >
        <iframe title={"terms"} style={{border: 0, width: "100%", height: "60vh"}} srcDoc={this.state.termsOfUseContent}/>
      </Modal>
    )
  }

  renderForm(application) {
    if (!application.enableSignUp) {
      return (
        <Result
          status="error"
          title="Sign Up Error"
          subTitle={"The application does not allow to sign up new account"}
          extra={[
            <Button type="primary" key="signin" onClick={() => {
              Setting.goToLogin(this, application);
            }}>
              Sign In
            </Button>
          ]}
        >
        </Result>
      )
    }
    return (
      <Form
        {...formItemLayout}
        ref={this.form}
        name="signup"
        onFinish={(values) => this.onFinish(values)}
        onFinishFailed={(errorInfo) => this.onFinishFailed(errorInfo.values, errorInfo.errorFields, errorInfo.outOfDate)}
        initialValues={{
          application: application.name,
          organization: application.organization,
        }}
        style={{width: !Setting.isMobile() ? "400px" : "250px"}}
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
        {
          application.signupItems?.map(signupItem => this.renderFormItem(application, signupItem))
        }
        <Form.Item {...tailFormItemLayout}>
          <Button type="primary" htmlType="submit">
            {i18next.t("account:Sign Up")}
          </Button>
          &nbsp;&nbsp;{i18next.t("signup:Have account?")}&nbsp;
          <a onClick={() => {
            Setting.goToLogin(this, application);
          }}>
            {i18next.t("signup:sign in now")}
          </a>
        </Form.Item>
      </Form>
    )
  }

  render() {
    const application = this.getApplicationObj();
    if (application === null) {
      return null;
    }

    if (application.signupHtml !== "") {
      return (
        <div dangerouslySetInnerHTML={{ __html: application.signupHtml}} />
      )
    }

    return (
      <div>
        <CustomGithubCorner />
        &nbsp;
        <Row>
          <Col span={24} style={{display: "flex", justifyContent:  "center"}} >
            <div style={{marginTop: "10px", textAlign: "center"}}>
              {
                Setting.renderHelmet(application)
              }
              {
                Setting.renderLogo(application)
              }
              {
                this.renderForm(application)
              }
            </div>
          </Col>
        </Row>
        {
          this.renderModal()
        }
      </div>
    )
  }
}

export default SignupPage;