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
import {Button, Card, Col, Input, InputNumber, Row, Select} from "antd";
import {EyeInvisibleOutlined, EyeTwoTone} from "@ant-design/icons";
import * as LddpBackend from "./backend/LdapBackend";
import * as OrganizationBackend from "./backend/OrganizationBackend";
import * as Setting from "./Setting";
import i18next from "i18next";

const {Option} = Select;

class LdapEditPage extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      ldapId: props.match.params.ldapId,
      ldap: null,
      organizations: [],
    };
  }

  UNSAFE_componentWillMount() {
    this.getLdap();
    this.getOrganizations();
  }

  getLdap() {
    LddpBackend.getLdap(this.state.ldapId)
      .then((res) => {
        if (res.status === "ok") {
          this.setState({
            ldap: res.data
          })
        } else {
          Setting.showMessage("error", res.msg);
        }
      })
  }

  getOrganizations() {
    OrganizationBackend.getOrganizations("admin")
      .then((res) => {
        this.setState({
          organizations: (res.msg === undefined) ? res : [],
        });
      });
  }

  updateLdapField(key, value) {
    this.setState((prevState) => {
      prevState.ldap[key] = value;
      return prevState;
    });
  }

  renderAutoSyncWarn() {
    if (this.state.ldap.autoSync > 0) {
      return (
        <span style={{
          color: "#faad14",
          marginLeft: "20px"
        }}>{i18next.t("ldap:The Auto Sync option will sync all users to specify organization")}</span>
      )
    }
  }

  renderLdap() {
    return (
      <Card size="small" title={
        <div>
          {i18next.t("ldap:Edit LDAP")}&nbsp;&nbsp;&nbsp;&nbsp;
          <Button type="primary" onClick={() => this.submitLdapEdit()}>{i18next.t("general:Save")}</Button>
        </div>
      } style={{marginLeft: "5px"}} type="inner">
        <Row style={{marginTop: "10px"}}>
          <Col style={{lineHeight: "32px", textAlign: "right", paddingRight: "25px"}} span={3}>
            {Setting.getLabel(i18next.t("general:Organization"), i18next.t("general:Organization - Tooltip"))} :
          </Col>
          <Col span={21}>
            <Select virtual={false} style={{width: "100%"}} disabled={!Setting.isAdminUser(this.props.account)}
                    value={this.state.ldap.owner} onChange={(value => {
              this.updateLdapField("owner", value);
            })}>
              {
                this.state.organizations.map((organization, index) => <Option key={index}
                                                                              value={organization.name}>{organization.name}</Option>)
              }
            </Select>
          </Col>
        </Row>
        <Row style={{marginTop: "20px"}}>
          <Col style={{lineHeight: "32px", textAlign: "right", paddingRight: "25px"}} span={3}>
            {Setting.getLabel(i18next.t("ldap:ID"), i18next.t("general:ID - Tooltip"))} :
          </Col>
          <Col span={21}>
            <Input value={this.state.ldap.id} disabled={true}/>
          </Col>
        </Row>
        <Row style={{marginTop: "20px"}}>
          <Col style={{lineHeight: "32px", textAlign: "right", paddingRight: "25px"}} span={3}>
            {Setting.getLabel(i18next.t("ldap:Server Name"), i18next.t("ldap:Server Name - Tooltip"))} :
          </Col>
          <Col span={21}>
            <Input value={this.state.ldap.serverName} onChange={e => {
              this.updateLdapField("serverName", e.target.value);
            }}/>
          </Col>
        </Row>
        <Row style={{marginTop: "20px"}}>
          <Col style={{lineHeight: "32px", textAlign: "right", paddingRight: "25px"}} span={3}>
            {Setting.getLabel(i18next.t("ldap:Server Host"), i18next.t("ldap:Server Host - Tooltip"))} :
          </Col>
          <Col span={21}>
            <Input value={this.state.ldap.host} onChange={e => {
              this.updateLdapField("host", e.target.value);
            }}/>
          </Col>
        </Row>
        <Row style={{marginTop: "20px"}}>
          <Col style={{lineHeight: "32px", textAlign: "right", paddingRight: "25px"}} span={3}>
            {Setting.getLabel(i18next.t("ldap:Server Port"), i18next.t("ldap:Server Port - Tooltip"))} :
          </Col>
          <Col span={21}>
            <InputNumber min={0} max={65535} formatter={value => value.replace(/\$\s?|(,*)/g, "")}
                         value={this.state.ldap.port} onChange={value => {
              this.updateLdapField("port", value);
            }}/>
          </Col>
        </Row>
        <Row style={{marginTop: "20px"}}>
          <Col style={{lineHeight: "32px", textAlign: "right", paddingRight: "25px"}} span={3}>
            {Setting.getLabel(i18next.t("ldap:Base DN"), i18next.t("ldap:Base DN - Tooltip"))} :
          </Col>
          <Col span={21}>
            <Input value={this.state.ldap.baseDn} onChange={e => {
              this.updateLdapField("baseDn", e.target.value);
            }}/>
          </Col>
        </Row>
        <Row style={{marginTop: "20px"}}>
          <Col style={{lineHeight: "32px", textAlign: "right", paddingRight: "25px"}} span={3}>
            {Setting.getLabel(i18next.t("ldap:Admin"), i18next.t("ldap:Admin - Tooltip"))} :
          </Col>
          <Col span={21}>
            <Input value={this.state.ldap.admin} onChange={e => {
              this.updateLdapField("admin", e.target.value);
            }}/>
          </Col>
        </Row>
        <Row style={{marginTop: "20px"}}>
          <Col style={{lineHeight: "32px", textAlign: "right", paddingRight: "25px"}} span={3}>
            {Setting.getLabel(i18next.t("ldap:Admin Password"), i18next.t("ldap:Admin Password - Tooltip"))} :
          </Col>
          <Col span={21}>
            <Input.Password
              iconRender={visible => (visible ? <EyeTwoTone/> : <EyeInvisibleOutlined/>)} value={this.state.ldap.passwd}
              onChange={e => {
                this.updateLdapField("passwd", e.target.value);
              }}
            />
          </Col>
        </Row>
        <Row style={{marginTop: "20px"}}>
          <Col style={{lineHeight: "32px", textAlign: "right", paddingRight: "25px"}} span={3}>
            {Setting.getLabel(i18next.t("ldap:Auto Sync"), i18next.t("ldap:Auto Sync - Tooltip"))} :
          </Col>
          <Col span={21}>
            <InputNumber min={0} formatter={value => value.replace(/\$\s?|(,*)/g, "")} disabled={false}
                         value={this.state.ldap.autoSync} onChange={value => {
              this.updateLdapField("autoSync", value);
            }}/><span>&nbsp;mins</span>
            {this.renderAutoSyncWarn()}
          </Col>
        </Row>
      </Card>
    )
  }

  submitLdapEdit() {
    LddpBackend.updateLdap(this.state.ldap)
      .then((res) => {
        if (res.status === "ok") {
          Setting.showMessage("success", `Update LDAP server success`);
          this.setState((prevState) => {
            prevState.ldap = res.data2;
          })
        } else {
          Setting.showMessage("error", res.msg);
        }
      })
      .catch(error => {
        Setting.showMessage("error", `Update LDAP server failed: ${error}`);
      });
  }

  render() {
    return (
      <div>
        <Row style={{width: "100%"}}>
          <Col span={1}>
          </Col>
          <Col span={22}>
            {
              this.state.ldap !== null ? this.renderLdap() : null
            }
          </Col>
          <Col span={1}>
          </Col>
        </Row>
        <Row style={{margin: 10}}>
          <Col span={2}>
          </Col>
          <Col span={18}>
            <Button type="primary" size="large"
                    onClick={() => this.submitLdapEdit()}>{i18next.t("general:Save")}</Button>
          </Col>
        </Row>
      </div>
    );
  }
}

export default LdapEditPage;