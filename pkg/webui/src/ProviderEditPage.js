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
import {Button, Card, Col, Input, InputNumber, Row, Select, Switch} from 'antd';
import {LinkOutlined} from "@ant-design/icons";
import * as ProviderBackend from "./backend/ProviderBackend";
import * as Setting from "./Setting";
import i18next from "i18next";
import { authConfig } from "./auth/Auth";
import copy from 'copy-to-clipboard';

const { Option } = Select;
const { TextArea } = Input;

class ProviderEditPage extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      classes: props,
      providerName: props.match.params.providerName,
      provider: null,
    };
  }

  UNSAFE_componentWillMount() {
    this.getProvider();
  }

  getProvider() {
    ProviderBackend.getProvider("admin", this.state.providerName)
      .then((provider) => {
        this.setState({
          provider: provider,
        });
      });
  }

  parseProviderField(key, value) {
    if (["port"].includes(key)) {
      value = Setting.myParseInt(value);
    }
    return value;
  }

  updateProviderField(key, value) {
    value = this.parseProviderField(key, value);

    let provider = this.state.provider;
    provider[key] = value;
    this.setState({
      provider: provider,
    });
  }

  getClientIdLabel() {
    switch (this.state.provider.category) {
      case "Email":
        return Setting.getLabel(i18next.t("signup:Username"), i18next.t("signup:Username - Tooltip"));
      case "SMS":
        if (this.state.provider.type === "Volc Engine SMS")
          return Setting.getLabel(i18next.t("provider:Access key"), i18next.t("provider:Access key - Tooltip"));
      default:
        return Setting.getLabel(i18next.t("provider:Client ID"), i18next.t("provider:Client ID - Tooltip"));
    }
  }

  getClientSecretLabel() {
    switch (this.state.provider.category) {
      case "Email":
        return Setting.getLabel(i18next.t("login:Password"), i18next.t("login:Password - Tooltip"));
      case "SMS":
        if (this.state.provider.type === "Volc Engine SMS")
          return Setting.getLabel(i18next.t("provider:Secret access key"), i18next.t("provider:SecretAccessKey - Tooltip"));
      default:
        return Setting.getLabel(i18next.t("provider:Client secret"), i18next.t("provider:Client secret - Tooltip"));
    }
  }

  getAppIdRow() {
    let text, tooltip;
    if (this.state.provider.category === "SMS" && this.state.provider.type === "Tencent Cloud SMS") {
      text = "provider:App ID";
      tooltip = "provider:App ID - Tooltip";
    } else if (this.state.provider.category === "SMS" && this.state.provider.type === "Volc Engine SMS") {
      text = "provider:SMS account";
      tooltip = "provider:SMS account - Tooltip";
    } else {
      return null;
    }

    return <Row style={{marginTop: '20px'}} >
      <Col style={{marginTop: '5px'}} span={(Setting.isMobile()) ? 22 : 2}>
        {Setting.getLabel(i18next.t(text), i18next.t(tooltip))} :
      </Col>
      <Col span={22} >
        <Input value={this.state.provider.appId} onChange={e => {
          this.updateProviderField('appId', e.target.value);
        }} />
      </Col>
    </Row>;
  }

  loadSamlConfiguration() {
    var parser = new DOMParser();
    var xmlDoc = parser.parseFromString(this.state.provider.metadata, "text/xml");
    var cert = xmlDoc.getElementsByTagName("ds:X509Certificate")[0].childNodes[0].nodeValue;
    var endpoint = xmlDoc.getElementsByTagName("md:SingleSignOnService")[0].getAttribute("Location");
    var issuerUrl = xmlDoc.getElementsByTagName("md:EntityDescriptor")[0].getAttribute("entityID");
    this.updateProviderField("idP", cert);
    this.updateProviderField("endpoint", endpoint);
    this.updateProviderField("issuerUrl", issuerUrl);
  }

  renderProvider() {
    return (
      <Card size="small" title={
        <div>
          {i18next.t("provider:Edit Provider")}&nbsp;&nbsp;&nbsp;&nbsp;
          <Button onClick={() => this.submitProviderEdit(false)}>{i18next.t("general:Save")}</Button>
          <Button style={{marginLeft: '20px'}} type="primary" onClick={() => this.submitProviderEdit(true)}>{i18next.t("general:Save & Exit")}</Button>
        </div>
      } style={(Setting.isMobile())? {margin: '5px'}:{}} type="inner">
        <Row style={{marginTop: '10px'}} >
          <Col style={{marginTop: '5px'}} span={(Setting.isMobile()) ? 22 : 2}>
            {Setting.getLabel(i18next.t("general:Name"), i18next.t("general:Name - Tooltip"))} :
          </Col>
          <Col span={22} >
            <Input value={this.state.provider.name} onChange={e => {
              this.updateProviderField('name', e.target.value);
            }} />
          </Col>
        </Row>
        <Row style={{marginTop: '20px'}} >
          <Col style={{marginTop: '5px'}} span={(Setting.isMobile()) ? 22 : 2}>
            {Setting.getLabel(i18next.t("general:Display name"), i18next.t("general:Display name - Tooltip"))} :
          </Col>
          <Col span={22} >
            <Input value={this.state.provider.displayName} onChange={e => {
              this.updateProviderField('displayName', e.target.value);
            }} />
          </Col>
        </Row>
        <Row style={{marginTop: '20px'}} >
          <Col style={{marginTop: '5px'}} span={(Setting.isMobile()) ? 22 : 2}>
            {Setting.getLabel(i18next.t("provider:Category"), i18next.t("provider:Category - Tooltip"))} :
          </Col>
          <Col span={22} >
            <Select virtual={false} style={{width: '100%'}} value={this.state.provider.category} onChange={(value => {
              this.updateProviderField('category', value);
              if (value === "OAuth") {
                this.updateProviderField('type', 'GitHub');
              } else if (value === "Email") {
                this.updateProviderField('type', 'Default');
                this.updateProviderField('title', 'Bhojpur IAM Verification Code');
                this.updateProviderField('content', 'You have requested a verification code at Bhojpur IAM. Here is your code: %s, please enter in 5 minutes.');
              } else if (value === "SMS") {
                this.updateProviderField('type', 'Aliyun SMS');
              } else if (value === "Storage") {
                this.updateProviderField('type', 'Local File System');
                this.updateProviderField('domain', Setting.getFullServerUrl());
              } else if (value === "SAML") {
                this.updateProviderField('type', 'Aliyun IDaaS');
              }
            })}>
              {
                [
                  {id: 'OAuth', name: 'OAuth'},
                  {id: 'Email', name: 'Email'},
                  {id: 'SMS', name: 'SMS'},
                  {id: 'Storage', name: 'Storage'},
                  {id: 'SAML', name: 'SAML'},
                ].map((providerCategory, index) => <Option key={index} value={providerCategory.id}>{providerCategory.name}</Option>)
              }
            </Select>
          </Col>
        </Row>
        <Row style={{marginTop: '20px'}} >
          <Col style={{marginTop: '5px'}} span={(Setting.isMobile()) ? 22 : 2}>
            {Setting.getLabel(i18next.t("provider:Type"), i18next.t("provider:Type - Tooltip"))} :
          </Col>
          <Col span={22} >
            <Select virtual={false} style={{width: '100%'}} value={this.state.provider.type} onChange={(value => {
              this.updateProviderField('type', value);
              if (value === "Local File System") {
                this.updateProviderField('domain', Setting.getFullServerUrl());
              }
            })}>
              {
                Setting.getProviderTypeOptions(this.state.provider.category).map((providerType, index) => <Option key={index} value={providerType.id}>{providerType.name}</Option>)
              }
            </Select>
          </Col>
        </Row>
        {
          this.state.provider.type !== "WeCom" ? null : (
            <Row style={{marginTop: '20px'}} >
              <Col style={{marginTop: '5px'}} span={2}>
                {Setting.getLabel(i18next.t("provider:Method"), i18next.t("provider:Method - Tooltip"))} :
              </Col>
              <Col span={22} >
                <Select virtual={false} style={{width: '100%'}} value={this.state.provider.method} onChange={value => {
                  this.updateProviderField('method', value);
                }}>
                  {
                    [{name: "Normal"}, {name: "Silent"}].map((method, index) => <Option key={index} value={method.name}>{method.name}</Option>)
                  }
                </Select>
              </Col>
            </Row>
          )
        }
        <Row style={{marginTop: '20px'}} >
          <Col style={{marginTop: '5px'}} span={(Setting.isMobile()) ? 22 : 2}>
            {this.getClientIdLabel()}
          </Col>
          <Col span={22} >
            <Input value={this.state.provider.clientId} onChange={e => {
              this.updateProviderField('clientId', e.target.value);
            }} />
          </Col>
        </Row>
        <Row style={{marginTop: '20px'}} >
          <Col style={{marginTop: '5px'}} span={(Setting.isMobile()) ? 22 : 2}>
            {this.getClientSecretLabel()}
          </Col>
          <Col span={22} >
            <Input value={this.state.provider.clientSecret} onChange={e => {
              this.updateProviderField('clientSecret', e.target.value);
            }} />
          </Col>
        </Row>
        {
          this.state.provider.type !== "WeChat" ? null : (
            <React.Fragment>
              <Row style={{marginTop: '20px'}} >
                <Col style={{marginTop: '5px'}} span={(Setting.isMobile()) ? 22 : 2}>
                  {Setting.getLabel(i18next.t("provider:Client ID 2"), i18next.t("provider:Client ID 2 - Tooltip"))}
                </Col>
                <Col span={22} >
                  <Input value={this.state.provider.clientId2} onChange={e => {
                    this.updateProviderField('clientId2', e.target.value);
                  }} />
                </Col>
              </Row>
              <Row style={{marginTop: '20px'}} >
                <Col style={{marginTop: '5px'}} span={(Setting.isMobile()) ? 22 : 2}>
                  {Setting.getLabel(i18next.t("provider:Client secret 2"), i18next.t("provider:Client secret 2 - Tooltip"))}
                </Col>
                <Col span={22} >
                  <Input value={this.state.provider.clientSecret2} onChange={e => {
                    this.updateProviderField('clientSecret2', e.target.value);
                  }} />
                </Col>
              </Row>
            </React.Fragment>
          )
        }
        {this.state.provider.category === "Storage" ? (
          <div>
            <Row style={{marginTop: '20px'}} >
              <Col style={{marginTop: '5px'}} span={2}>
                {Setting.getLabel(i18next.t("provider:Endpoint"), i18next.t("provider:Region endpoint for Internet"))} :
              </Col>
              <Col span={22} >
                <Input value={this.state.provider.endpoint} onChange={e => {
                  this.updateProviderField('endpoint', e.target.value);
                }} />
              </Col>
            </Row>
            <Row style={{marginTop: '20px'}} >
              <Col style={{marginTop: '5px'}} span={2}>
                {Setting.getLabel(i18next.t("provider:Endpoint (Intranet)"), i18next.t("provider:Region endpoint for Intranet"))} :
              </Col>
              <Col span={22} >
                <Input value={this.state.provider.intranetEndpoint} onChange={e => {
                  this.updateProviderField('intranetEndpoint', e.target.value);
                }} />
              </Col>
            </Row>
            <Row style={{marginTop: '20px'}} >
              <Col style={{marginTop: '5px'}} span={2}>
                {Setting.getLabel(i18next.t("provider:Bucket"), i18next.t("provider:Bucket - Tooltip"))} :
              </Col>
              <Col span={22} >
                <Input value={this.state.provider.bucket} onChange={e => {
                  this.updateProviderField('bucket', e.target.value);
                }} />
              </Col>
            </Row>
            <Row style={{marginTop: '20px'}} >
              <Col style={{marginTop: '5px'}} span={2}>
                {Setting.getLabel(i18next.t("provider:Domain"), i18next.t("provider:Domain - Tooltip"))} :
              </Col>
              <Col span={22} >
                <Input value={this.state.provider.domain} onChange={e => {
                  this.updateProviderField('domain', e.target.value);
                }} />
              </Col>
            </Row>
            {this.state.provider.type === "AWS S3" || this.state.provider.type === "Tencent Cloud COS" ? (
              <Row style={{marginTop: '20px'}} >
                <Col style={{marginTop: '5px'}} span={2}>
                  {Setting.getLabel(i18next.t("provider:Region ID"), i18next.t("provider:Region ID - Tooltip"))} :
                </Col>
                <Col span={22} >
                  <Input value={this.state.provider.regionId} onChange={e => {
                    this.updateProviderField('regionId', e.target.value);
                  }} />
                </Col>
              </Row>
            ) : null}
          </div>
        ) : null}
        {
          this.state.provider.category === "Email" ? (
            <React.Fragment>
              <Row style={{marginTop: '20px'}} >
                <Col style={{marginTop: '5px'}} span={(Setting.isMobile()) ? 22 : 2}>
                  {Setting.getLabel(i18next.t("provider:Host"), i18next.t("provider:Host - Tooltip"))} :
                </Col>
                <Col span={22} >
                  <Input prefix={<LinkOutlined/>} value={this.state.provider.host} onChange={e => {
                    this.updateProviderField('host', e.target.value);
                  }} />
                </Col>
              </Row>
              <Row style={{marginTop: '20px'}} >
                <Col style={{marginTop: '5px'}} span={(Setting.isMobile()) ? 22 : 2}>
                  {Setting.getLabel(i18next.t("provider:Port"), i18next.t("provider:Port - Tooltip"))} :
                </Col>
                <Col span={22} >
                  <InputNumber value={this.state.provider.port} onChange={value => {
                    this.updateProviderField('port', value);
                  }} />
                </Col>
              </Row>
              <Row style={{marginTop: '20px'}} >
                <Col style={{marginTop: '5px'}} span={(Setting.isMobile()) ? 22 : 2}>
                  {Setting.getLabel(i18next.t("provider:Email Title"), i18next.t("provider:Email Title - Tooltip"))} :
                </Col>
                <Col span={22} >
                  <Input value={this.state.provider.title} onChange={e => {
                    this.updateProviderField('title', e.target.value);
                  }} />
                </Col>
              </Row>
              <Row style={{marginTop: '20px'}} >
                <Col style={{marginTop: '5px'}} span={(Setting.isMobile()) ? 22 : 2}>
                  {Setting.getLabel(i18next.t("provider:Email Content"), i18next.t("provider:Email Content - Tooltip"))} :
                </Col>
                <Col span={22} >
                  <TextArea autoSize={{minRows: 1, maxRows: 100}} value={this.state.provider.content} onChange={e => {
                    this.updateProviderField('content', e.target.value);
                  }} />
                </Col>
              </Row>
            </React.Fragment>
          ) : this.state.provider.category === "SMS" ? (
            <React.Fragment>
              <Row style={{marginTop: '20px'}} >
                <Col style={{marginTop: '5px'}} span={(Setting.isMobile()) ? 22 : 2}>
                  {Setting.getLabel(i18next.t("provider:Sign Name"), i18next.t("provider:Sign Name - Tooltip"))} :
                </Col>
                <Col span={22} >
                  <Input value={this.state.provider.signName} onChange={e => {
                    this.updateProviderField('signName', e.target.value);
                  }} />
                </Col>
              </Row>
              <Row style={{marginTop: '20px'}} >
                <Col style={{marginTop: '5px'}} span={(Setting.isMobile()) ? 22 : 2}>
                  {Setting.getLabel(i18next.t("provider:Template Code"), i18next.t("provider:Template Code - Tooltip"))} :
                </Col>
                <Col span={22} >
                  <Input value={this.state.provider.templateCode} onChange={e => {
                    this.updateProviderField('templateCode', e.target.value);
                  }} />
                </Col>
              </Row>
            </React.Fragment>
          ) : this.state.provider.category === "SAML" ? (
            <React.Fragment>
              <Row style={{marginTop: '20px'}} >
                <Col style={{marginTop: '5px'}} span={(Setting.isMobile()) ? 22 : 2}>
                  {Setting.getLabel(i18next.t("provider:Sign request"), i18next.t("provider:Sign request - Tooltip"))} :
                </Col>
                <Col span={22} >
                  <Switch checked={this.state.provider.enableSignAuthnRequest} onChange={checked => {
                    this.updateProviderField('enableSignAuthnRequest', checked);
                  }} />
                </Col>
              </Row>
              <Row style={{marginTop: '20px'}} >
                <Col style={{marginTop: '5px'}} span={(Setting.isMobile()) ? 22 : 2}>
                  {Setting.getLabel(i18next.t("provider:Metadata"), i18next.t("provider:Metadata - Tooltip"))} :
                </Col>
                <Col span={22}>
                  <TextArea rows={4} value={this.state.provider.metadata} onChange={e => {
                    this.updateProviderField('metadata', e.target.value);
                  }} />
                </Col>
              </Row>
              <Row style={{marginTop: '20px'}}>
                <Col style={{marginTop: '5px'}} span={2}></Col>
                <Col span={2}>
                  <Button type="primary" onClick={() => {
                      try {
                        this.loadSamlConfiguration();
                        Setting.showMessage("success", i18next.t("provider:Parse Metadata successfully"));
                      } catch (err) {
                        Setting.showMessage("error", i18next.t("provider:Can not parse Metadata"));
                      }
                    }}>
                    {i18next.t("provider:Parse")}
                  </Button>
                </Col>
              </Row>
              <Row style={{marginTop: '20px'}} >
                <Col style={{marginTop: '5px'}} span={(Setting.isMobile()) ? 22 : 2}>
                  {Setting.getLabel(i18next.t("provider:Endpoint"), i18next.t("provider:SAML 2.0 Endpoint (HTTP)"))} :
                </Col>
                <Col span={22} >
                  <Input value={this.state.provider.endpoint} onChange={e => {
                    this.updateProviderField('endpoint', e.target.value);
                  }} />
                </Col>
              </Row>
              <Row style={{marginTop: '20px'}} >
                <Col style={{marginTop: '5px'}} span={(Setting.isMobile()) ? 22 : 2}>
                  {Setting.getLabel(i18next.t("provider:IdP"), i18next.t("provider:IdP public key"))} :
                </Col>
                <Col span={22} >
                  <Input value={this.state.provider.idP} onChange={e => {
                    this.updateProviderField('idP', e.target.value);
                  }} />
                </Col>
              </Row>
              <Row style={{marginTop: '20px'}} >
                <Col style={{marginTop: '5px'}} span={(Setting.isMobile()) ? 22 : 2}>
                  {Setting.getLabel(i18next.t("provider:Issuer URL"), i18next.t("provider:Issuer URL - Tooltip"))} :
                </Col>
                <Col span={22} >
                  <Input value={this.state.provider.issuerUrl} onChange={e => {
                    this.updateProviderField('issuerUrl', e.target.value);
                  }} />
                </Col>
              </Row>
              <Row style={{marginTop: '20px'}} >
                <Col style={{marginTop: '5px'}} span={(Setting.isMobile()) ? 22 : 2}>
                  {Setting.getLabel(i18next.t("provider:SP ACS URL"), i18next.t("provider:SP ACS URL - Tooltip"))} :
                </Col>
                <Col span={21} >
                  <Input value={`${authConfig.serverUrl}/api/acs`} readOnly="readonly" />
                </Col>
                <Col span={1}>
                  <Button type="primary" onClick={() => {
                    copy(`${authConfig.serverUrl}/api/acs`);
                    Setting.showMessage("success", i18next.t("provider:Link copied to clipboard successfully"));
                  }}>
                    {i18next.t("provider:Copy")}
                  </Button>
                </Col>
              </Row>
              <Row style={{marginTop: '20px'}} >
                <Col style={{marginTop: '5px'}} span={(Setting.isMobile()) ? 22 : 2}>
                  {Setting.getLabel(i18next.t("provider:SP Entity ID"), i18next.t("provider:SP ACS URL - Tooltip"))} :
                </Col>
                <Col span={21} >
                  <Input value={`${authConfig.serverUrl}/api/acs`} readOnly="readonly" />
                </Col>
                <Col span={1}>
                  <Button type="primary" onClick={() => {
                    copy(`${authConfig.serverUrl}/api/acs`);
                    Setting.showMessage("success", i18next.t("provider:Link copied to clipboard successfully"));
                  }}>
                    {i18next.t("provider:Copy")}
                  </Button>
                </Col>
              </Row>
            </React.Fragment>
          ) : null
        }
        {this.getAppIdRow()}
        <Row style={{marginTop: '20px'}} >
          <Col style={{marginTop: '5px'}} span={(Setting.isMobile()) ? 22 : 2}>
            {Setting.getLabel(i18next.t("provider:Provider URL"), i18next.t("provider:Provider URL - Tooltip"))} :
          </Col>
          <Col span={22} >
            <Input prefix={<LinkOutlined/>} value={this.state.provider.providerUrl} onChange={e => {
              this.updateProviderField('providerUrl', e.target.value);
            }} />
          </Col>
        </Row>
      </Card>
    )
  }

  submitProviderEdit(willExist) {
    let provider = Setting.deepCopy(this.state.provider);
    ProviderBackend.updateProvider(this.state.provider.owner, this.state.providerName, provider)
      .then((res) => {
        if (res.msg === "") {
          Setting.showMessage("success", `Successfully saved`);
          this.setState({
            providerName: this.state.provider.name,
          });

          if (willExist) {
            this.props.history.push(`/providers`);
          } else {
            this.props.history.push(`/providers/${this.state.provider.name}`);
          }
        } else {
          Setting.showMessage("error", res.msg);
          this.updateProviderField('name', this.state.providerName);
        }
      })
      .catch(error => {
        Setting.showMessage("error", `Failed to connect to server: ${error}`);
      });
  }

  render() {
    return (
      <div>
        {
          this.state.provider !== null ? this.renderProvider() : null
        }
        <div style={{marginTop: '20px', marginLeft: '40px'}}>
          <Button size="large" onClick={() => this.submitProviderEdit(false)}>{i18next.t("general:Save")}</Button>
          <Button style={{marginLeft: '20px'}} type="primary" size="large" onClick={() => this.submitProviderEdit(true)}>{i18next.t("general:Save & Exit")}</Button>
        </div>
      </div>
    );
  }
}

export default ProviderEditPage;