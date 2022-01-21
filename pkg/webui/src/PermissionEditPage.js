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
import {Button, Card, Col, Input, Row, Select, Switch} from 'antd';
import * as PermissionBackend from "./backend/PermissionBackend";
import * as OrganizationBackend from "./backend/OrganizationBackend";
import * as UserBackend from "./backend/UserBackend";
import * as Setting from "./Setting";
import i18next from "i18next";
import * as RoleBackend from "./backend/RoleBackend";

const { Option } = Select;

class PermissionEditPage extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      classes: props,
      organizationName: props.organizationName !== undefined ? props.organizationName : props.match.params.organizationName,
      permissionName: props.match.params.permissionName,
      permission: null,
      organizations: [],
      users: [],
      roles: [],
    };
  }

  UNSAFE_componentWillMount() {
    this.getPermission();
    this.getOrganizations();
  }

  getPermission() {
    PermissionBackend.getPermission(this.state.organizationName, this.state.permissionName)
      .then((permission) => {
        this.setState({
          permission: permission,
        });

        this.getUsers(permission.owner);
        this.getRoles(permission.owner);
      });
  }

  getOrganizations() {
    OrganizationBackend.getOrganizations("admin")
      .then((res) => {
        this.setState({
          organizations: (res.msg === undefined) ? res : [],
        });
      });
  }

  getUsers(organizationName) {
    UserBackend.getUsers(organizationName)
      .then((res) => {
        this.setState({
          users: res,
        });
      });
  }

  getRoles(organizationName) {
    RoleBackend.getRoles(organizationName)
      .then((res) => {
        this.setState({
          roles: res,
        });
      });
  }

  parsePermissionField(key, value) {
    if ([""].includes(key)) {
      value = Setting.myParseInt(value);
    }
    return value;
  }

  updatePermissionField(key, value) {
    value = this.parsePermissionField(key, value);

    let permission = this.state.permission;
    permission[key] = value;
    this.setState({
      permission: permission,
    });
  }

  renderPermission() {
    return (
      <Card size="small" title={
        <div>
          {i18next.t("permission:Edit Permission")}&nbsp;&nbsp;&nbsp;&nbsp;
          <Button onClick={() => this.submitPermissionEdit(false)}>{i18next.t("general:Save")}</Button>
          <Button style={{marginLeft: '20px'}} type="primary" onClick={() => this.submitPermissionEdit(true)}>{i18next.t("general:Save & Exit")}</Button>
        </div>
      } style={(Setting.isMobile())? {margin: '5px'}:{}} type="inner">
        <Row style={{marginTop: '10px'}} >
          <Col style={{marginTop: '5px'}} span={(Setting.isMobile()) ? 22 : 2}>
            {Setting.getLabel(i18next.t("general:Organization"), i18next.t("general:Organization - Tooltip"))} :
          </Col>
          <Col span={22} >
            <Select virtual={false} style={{width: '100%'}} value={this.state.permission.owner} onChange={(owner => {
              this.updatePermissionField('owner', owner);

              this.getUsers(owner);
              this.getRoles(owner);
            })}>
              {
                this.state.organizations.map((organization, index) => <Option key={index} value={organization.name}>{organization.name}</Option>)
              }
            </Select>
          </Col>
        </Row>
        <Row style={{marginTop: '20px'}} >
          <Col style={{marginTop: '5px'}} span={(Setting.isMobile()) ? 22 : 2}>
            {Setting.getLabel(i18next.t("general:Name"), i18next.t("general:Name - Tooltip"))} :
          </Col>
          <Col span={22} >
            <Input value={this.state.permission.name} onChange={e => {
              this.updatePermissionField('name', e.target.value);
            }} />
          </Col>
        </Row>
        <Row style={{marginTop: '20px'}} >
          <Col style={{marginTop: '5px'}} span={(Setting.isMobile()) ? 22 : 2}>
            {Setting.getLabel(i18next.t("general:Display name"), i18next.t("general:Display name - Tooltip"))} :
          </Col>
          <Col span={22} >
            <Input value={this.state.permission.displayName} onChange={e => {
              this.updatePermissionField('displayName', e.target.value);
            }} />
          </Col>
        </Row>
        <Row style={{marginTop: '20px'}} >
          <Col style={{marginTop: '5px'}} span={(Setting.isMobile()) ? 22 : 2}>
            {Setting.getLabel(i18next.t("role:Sub users"), i18next.t("role:Sub users - Tooltip"))} :
          </Col>
          <Col span={22} >
            <Select virtual={false} mode="tags" style={{width: '100%'}} value={this.state.permission.users} onChange={(value => {this.updatePermissionField('users', value);})}>
              {
                this.state.users.map((user, index) => <Option key={index} value={`${user.owner}/${user.name}`}>{`${user.owner}/${user.name}`}</Option>)
              }
            </Select>
          </Col>
        </Row>
        <Row style={{marginTop: '20px'}} >
          <Col style={{marginTop: '5px'}} span={(Setting.isMobile()) ? 22 : 2}>
            {Setting.getLabel(i18next.t("role:Sub roles"), i18next.t("role:Sub roles - Tooltip"))} :
          </Col>
          <Col span={22} >
            <Select virtual={false} mode="tags" style={{width: '100%'}} value={this.state.permission.roles} onChange={(value => {this.updatePermissionField('roles', value);})}>
              {
                this.state.roles.filter(roles => (roles.owner !== this.state.roles.owner || roles.name !== this.state.roles.name)).map((permission, index) => <Option key={index} value={`${permission.owner}/${permission.name}`}>{`${permission.owner}/${permission.name}`}</Option>)
              }
            </Select>
          </Col>
        </Row>
        <Row style={{marginTop: '20px'}} >
          <Col style={{marginTop: '5px'}} span={(Setting.isMobile()) ? 22 : 2}>
            {Setting.getLabel(i18next.t("permission:Resource type"), i18next.t("permission:Resource type - Tooltip"))} :
          </Col>
          <Col span={22} >
            <Select virtual={false} style={{width: '100%'}} value={this.state.permission.resourceType} onChange={(value => {
              this.updatePermissionField('resourceType', value);
            })}>
              {
                [
                  {id: 'Application', name: 'Application'},
                ].map((item, index) => <Option key={index} value={item.id}>{item.name}</Option>)
              }
            </Select>
          </Col>
        </Row>
        <Row style={{marginTop: '20px'}} >
          <Col style={{marginTop: '5px'}} span={(Setting.isMobile()) ? 22 : 2}>
            {Setting.getLabel(i18next.t("permission:Actions"), i18next.t("permission:Actions - Tooltip"))} :
          </Col>
          <Col span={22} >
            <Select virtual={false} mode="tags" style={{width: '100%'}} value={this.state.permission.actions} onChange={(value => {
              this.updatePermissionField('actions', value);
            })}>
              {
                [
                  {id: 'Read', name: 'Read'},
                  {id: 'Write', name: 'Write'},
                  {id: 'Admin', name: 'Admin'},
                ].map((item, index) => <Option key={index} value={item.id}>{item.name}</Option>)
              }
            </Select>
          </Col>
        </Row>
        <Row style={{marginTop: '20px'}} >
          <Col style={{marginTop: '5px'}} span={(Setting.isMobile()) ? 22 : 2}>
            {Setting.getLabel(i18next.t("permission:Effect"), i18next.t("permission:Effect - Tooltip"))} :
          </Col>
          <Col span={22} >
            <Select virtual={false} style={{width: '100%'}} value={this.state.permission.effect} onChange={(value => {
              this.updatePermissionField('effect', value);
            })}>
              {
                [
                  {id: 'Allow', name: 'Allow'},
                  {id: 'Deny', name: 'Deny'},
                ].map((item, index) => <Option key={index} value={item.id}>{item.name}</Option>)
              }
            </Select>
          </Col>
        </Row>
        <Row style={{marginTop: '20px'}} >
          <Col style={{marginTop: '5px'}} span={(Setting.isMobile()) ? 19 : 2}>
            {Setting.getLabel(i18next.t("general:Is enabled"), i18next.t("general:Is enabled - Tooltip"))} :
          </Col>
          <Col span={1} >
            <Switch checked={this.state.permission.isEnabled} onChange={checked => {
              this.updatePermissionField('isEnabled', checked);
            }} />
          </Col>
        </Row>
      </Card>
    )
  }

  submitPermissionEdit(willExist) {
    let permission = Setting.deepCopy(this.state.permission);
    PermissionBackend.updatePermission(this.state.organizationName, this.state.permissionName, permission)
      .then((res) => {
        if (res.msg === "") {
          Setting.showMessage("success", `Successfully saved`);
          this.setState({
            permissionName: this.state.permission.name,
          });

          if (willExist) {
            this.props.history.push(`/permissions`);
          } else {
            this.props.history.push(`/permissions/${this.state.permission.owner}/${this.state.permission.name}`);
          }
        } else {
          Setting.showMessage("error", res.msg);
          this.updatePermissionField('name', this.state.permissionName);
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
          this.state.permission !== null ? this.renderPermission() : null
        }
        <div style={{marginTop: '20px', marginLeft: '40px'}}>
          <Button size="large" onClick={() => this.submitPermissionEdit(false)}>{i18next.t("general:Save")}</Button>
          <Button style={{marginLeft: '20px'}} type="primary" size="large" onClick={() => this.submitPermissionEdit(true)}>{i18next.t("general:Save & Exit")}</Button>
        </div>
      </div>
    );
  }
}

export default PermissionEditPage;