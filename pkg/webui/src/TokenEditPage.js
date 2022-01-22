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
import {Button, Card, Col, Input, Row} from 'antd';
import * as TokenBackend from "./backend/TokenBackend";
import * as Setting from "./Setting";
import i18next from "i18next";

class TokenEditPage extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      classes: props,
      tokenName: props.match.params.tokenName,
      token: null,
    };
  }

  UNSAFE_componentWillMount() {
    this.getToken();
  }

  getToken() {
    TokenBackend.getToken("admin", this.state.tokenName)
      .then((token) => {
        this.setState({
          token: token,
        });
      });
  }

  parseTokenField(key, value) {
    // if ([].includes(key)) {
    //   value = Setting.myParseInt(value);
    // }
    return value;
  }

  updateTokenField(key, value) {
    value = this.parseTokenField(key, value);

    let token = this.state.token;
    token[key] = value;
    this.setState({
      token: token,
    });
  }

  renderToken() {
    return (
      <Card size="small" title={
        <div>
          {i18next.t("token:Edit Token")}&nbsp;&nbsp;&nbsp;&nbsp;
          <Button onClick={() => this.submitTokenEdit(false)}>{i18next.t("general:Save")}</Button>
          <Button style={{marginLeft: '20px'}} type="primary" onClick={() => this.submitTokenEdit(true)}>{i18next.t("general:Save & Exit")}</Button>
        </div>
      } style={(Setting.isMobile())? {margin: '5px'}:{}} type="inner">
        <Row style={{marginTop: '10px'}} >
          <Col style={{marginTop: '5px'}} span={(Setting.isMobile()) ? 22 : 2}>
            {i18next.t("general:Name")}:
          </Col>
          <Col span={22} >
            <Input value={this.state.token.name} onChange={e => {
              this.updateTokenField('name', e.target.value);
            }} />
          </Col>
        </Row>
        <Row style={{marginTop: '20px'}} >
          <Col style={{marginTop: '5px'}} span={(Setting.isMobile()) ? 22 : 2}>
            {i18next.t("general:Application")}:
          </Col>
          <Col span={22} >
            <Input value={this.state.token.application} onChange={e => {
              this.updateTokenField('application', e.target.value);
            }} />
          </Col>
        </Row>
        <Row style={{marginTop: '20px'}} >
          <Col style={{marginTop: '5px'}} span={(Setting.isMobile()) ? 22 : 2}>
            {i18next.t("general:Organization")}:
          </Col>
          <Col span={22} >
            <Input value={this.state.token.organization} onChange={e => {
              this.updateTokenField('organization', e.target.value);
            }} />
          </Col>
        </Row>
        <Row style={{marginTop: '20px'}} >
          <Col style={{marginTop: '5px'}} span={(Setting.isMobile()) ? 22 : 2}>
            {i18next.t("general:User")}:
          </Col>
          <Col span={22} >
            <Input value={this.state.token.user} onChange={e => {
              this.updateTokenField('user', e.target.value);
            }} />
          </Col>
        </Row>
        <Row style={{marginTop: '20px'}} >
          <Col style={{marginTop: '5px'}} span={(Setting.isMobile()) ? 22 : 2}>
            {i18next.t("token:Authorization code")}:
          </Col>
          <Col span={22} >
            <Input value={this.state.token.code} onChange={e => {
              this.updateTokenField('code', e.target.value);
            }} />
          </Col>
        </Row>
        <Row style={{marginTop: '20px'}} >
          <Col style={{marginTop: '5px'}} span={(Setting.isMobile()) ? 22 : 2}>
            {i18next.t("token:Access token")}:
          </Col>
          <Col span={22} >
            <Input value={this.state.token.accessToken} onChange={e => {
              this.updateTokenField('accessToken', e.target.value);
            }} />
          </Col>
        </Row>
        <Row style={{marginTop: '20px'}} >
          <Col style={{marginTop: '5px'}} span={(Setting.isMobile()) ? 22 : 2}>
            {i18next.t("token:Expires in")}:
          </Col>
          <Col span={22} >
            <Input value={this.state.token.expiresIn} onChange={e => {
              this.updateTokenField('expiresIn', e.target.value);
            }} />
          </Col>
        </Row>
        <Row style={{marginTop: '20px'}} >
          <Col style={{marginTop: '5px'}} span={(Setting.isMobile()) ? 22 : 2}>
            {i18next.t("token:Scope")}:
          </Col>
          <Col span={22} >
            <Input value={this.state.token.scope} onChange={e => {
              this.updateTokenField('scope', e.target.value);
            }} />
          </Col>
        </Row>
        <Row style={{marginTop: '20px'}} >
          <Col style={{marginTop: '5px'}} span={(Setting.isMobile()) ? 22 : 2}>
            {i18next.t("token:Token type")}:
          </Col>
          <Col span={22} >
            <Input value={this.state.token.tokenType} onChange={e => {
              this.updateTokenField('tokenType', e.target.value);
            }} />
          </Col>
        </Row>
      </Card>
    )
  }

  submitTokenEdit(willExist) {
    let token = Setting.deepCopy(this.state.token);
    TokenBackend.updateToken(this.state.token.owner, this.state.tokenName, token)
      .then((res) => {
        if (res.msg === "") {
          Setting.showMessage("success", `Successfully saved`);
          this.setState({
            tokenName: this.state.token.name,
          });

          if (willExist) {
            this.props.history.push(`/tokens`);
          } else {
            this.props.history.push(`/tokens/${this.state.token.name}`);
          }
        } else {
          Setting.showMessage("error", res.msg);
          this.updateTokenField('name', this.state.tokenName);
        }
      })
      .catch(error => {
        Setting.showMessage("error", `failed to connect to server: ${error}`);
      });
  }

  render() {
    return (
      <div>
      {
        this.state.token !== null ? this.renderToken() : null
      }
      <div style={{marginTop: '20px', marginLeft: '40px'}}>
        <Button size="large" onClick={() => this.submitTokenEdit(false)}>{i18next.t("general:Save")}</Button>
        <Button style={{marginLeft: '20px'}} type="primary" size="large" onClick={() => this.submitTokenEdit(true)}>{i18next.t("general:Save & Exit")}</Button>
      </div>
    </div>
    );
  }
}

export default TokenEditPage;