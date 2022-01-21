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
import {Button, Card, Col, Input, InputNumber, Row, Select} from 'antd';
import * as CertBackend from "./backend/CertBackend";
import * as Setting from "./Setting";
import i18next from "i18next";
import copy from "copy-to-clipboard";
import FileSaver from "file-saver";

const { Option } = Select;
const { TextArea } = Input;

class CertEditPage extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      classes: props,
      certName: props.match.params.certName,
      cert: null,
    };
  }

  UNSAFE_componentWillMount() {
    this.getCert();
  }

  getCert() {
    CertBackend.getCert("admin", this.state.certName)
      .then((cert) => {
        this.setState({
          cert: cert,
        });
      });
  }

  parseCertField(key, value) {
    if (["port"].includes(key)) {
      value = Setting.myParseInt(value);
    }
    return value;
  }

  updateCertField(key, value) {
    value = this.parseCertField(key, value);

    let cert = this.state.cert;
    cert[key] = value;
    this.setState({
      cert: cert,
    });
  }

  renderCert() {
    return (
      <Card size="small" title={
        <div>
          {i18next.t("cert:Edit Cert")}&nbsp;&nbsp;&nbsp;&nbsp;
          <Button onClick={() => this.submitCertEdit(false)}>{i18next.t("general:Save")}</Button>
          <Button style={{marginLeft: '20px'}} type="primary" onClick={() => this.submitCertEdit(true)}>{i18next.t("general:Save & Exit")}</Button>
        </div>
      } style={(Setting.isMobile())? {margin: '5px'}:{}} type="inner">
        <Row style={{marginTop: '10px'}} >
          <Col style={{marginTop: '5px'}} span={(Setting.isMobile()) ? 22 : 2}>
            {Setting.getLabel(i18next.t("general:Name"), i18next.t("general:Name - Tooltip"))} :
          </Col>
          <Col span={22} >
            <Input value={this.state.cert.name} onChange={e => {
              this.updateCertField('name', e.target.value);
            }} />
          </Col>
        </Row>
        <Row style={{marginTop: '20px'}} >
          <Col style={{marginTop: '5px'}} span={(Setting.isMobile()) ? 22 : 2}>
            {Setting.getLabel(i18next.t("general:Display name"), i18next.t("general:Display name - Tooltip"))} :
          </Col>
          <Col span={22} >
            <Input value={this.state.cert.displayName} onChange={e => {
              this.updateCertField('displayName', e.target.value);
            }} />
          </Col>
        </Row>
        <Row style={{marginTop: '20px'}} >
          <Col style={{marginTop: '5px'}} span={(Setting.isMobile()) ? 22 : 2}>
            {Setting.getLabel(i18next.t("cert:Scope"), i18next.t("cert:Scope - Tooltip"))} :
          </Col>
          <Col span={22} >
            <Select virtual={false} style={{width: '100%'}} value={this.state.cert.scope} onChange={(value => {
              this.updateCertField('scope', value);
            })}>
              {
                [
                  {id: 'JWT', name: 'JWT'},
                ].map((item, index) => <Option key={index} value={item.id}>{item.name}</Option>)
              }
            </Select>
          </Col>
        </Row>
        <Row style={{marginTop: '20px'}} >
          <Col style={{marginTop: '5px'}} span={(Setting.isMobile()) ? 22 : 2}>
            {Setting.getLabel(i18next.t("cert:Type"), i18next.t("cert:Type - Tooltip"))} :
          </Col>
          <Col span={22} >
            <Select virtual={false} style={{width: '100%'}} value={this.state.cert.type} onChange={(value => {
              this.updateCertField('type', value);
            })}>
              {
                [
                  {id: 'x509', name: 'x509'},
                ].map((item, index) => <Option key={index} value={item.id}>{item.name}</Option>)
              }
            </Select>
          </Col>
        </Row>
        <Row style={{marginTop: '20px'}} >
          <Col style={{marginTop: '5px'}} span={(Setting.isMobile()) ? 22 : 2}>
            {Setting.getLabel(i18next.t("cert:Crypto algorithm"), i18next.t("cert:Crypto algorithm - Tooltip"))} :
          </Col>
          <Col span={22} >
            <Select virtual={false} style={{width: '100%'}} value={this.state.cert.cryptoAlgorithm} onChange={(value => {
              this.updateCertField('cryptoAlgorithm', value);
            })}>
              {
                [
                  {id: 'RSA', name: 'RSA'},
                ].map((item, index) => <Option key={index} value={item.id}>{item.name}</Option>)
              }
            </Select>
          </Col>
        </Row>
        <Row style={{marginTop: '20px'}} >
          <Col style={{marginTop: '5px'}} span={(Setting.isMobile()) ? 22 : 2}>
            {Setting.getLabel(i18next.t("cert:Bit size"), i18next.t("cert:Bit size - Tooltip"))} :
          </Col>
          <Col span={22} >
            <InputNumber value={this.state.cert.bitSize} onChange={value => {
              this.updateCertField('bitSize', value);
            }} />
          </Col>
        </Row>
        <Row style={{marginTop: '20px'}} >
          <Col style={{marginTop: '5px'}} span={(Setting.isMobile()) ? 22 : 2}>
            {Setting.getLabel(i18next.t("cert:Expire in years"), i18next.t("cert:Expire in years - Tooltip"))} :
          </Col>
          <Col span={22} >
            <InputNumber value={this.state.cert.expireInYears} onChange={value => {
              this.updateCertField('expireInYears', value);
            }} />
          </Col>
        </Row>
        <Row style={{marginTop: '20px'}} >
          <Col style={{marginTop: '5px'}} span={(Setting.isMobile()) ? 22 : 2}>
            {Setting.getLabel(i18next.t("cert:Public key"), i18next.t("cert:Public key - Tooltip"))} :
          </Col>
          <Col span={9} >
            <Button style={{marginRight: '10px', marginBottom: '10px'}} onClick={() => {
              copy(this.state.cert.publicKey);
              Setting.showMessage("success", i18next.t("cert:Public key copied to clipboard successfully"));
            }}
            >
              {i18next.t("cert:Copy public key")}
            </Button>
            <Button type="primary" onClick={() => {
              const blob = new Blob([this.state.cert.publicKey], {type: "text/plain;charset=utf-8"});
              FileSaver.saveAs(blob, "token_jwt_key.pem");
            }}
            >
              {i18next.t("cert:Download public key")}
            </Button>
            <TextArea autoSize={{minRows: 30, maxRows: 30}} value={this.state.cert.publicKey} onChange={e => {
              this.updateCertField('publicKey', e.target.value);
            }} />
          </Col>
          <Col span={1} />
          <Col style={{marginTop: '5px'}} span={(Setting.isMobile()) ? 22 : 2}>
            {Setting.getLabel(i18next.t("cert:Private key"), i18next.t("cert:Private key - Tooltip"))} :
          </Col>
          <Col span={9} >
            <Button style={{marginRight: '10px', marginBottom: '10px'}} onClick={() => {
              copy(this.state.cert.privateKey);
              Setting.showMessage("success", i18next.t("cert:Private key copied to clipboard successfully"));
            }}
            >
              {i18next.t("cert:Copy private key")}
            </Button>
            <Button type="primary" onClick={() => {
              const blob = new Blob([this.state.cert.privateKey], {type: "text/plain;charset=utf-8"});
              FileSaver.saveAs(blob, "token_jwt_key.key");
            }}
            >
              {i18next.t("cert:Download private key")}
            </Button>
            <TextArea autoSize={{minRows: 30, maxRows: 30}} value={this.state.cert.privateKey} onChange={e => {
              this.updateCertField('privateKey', e.target.value);
            }} />
          </Col>
        </Row>
      </Card>
    )
  }

  submitCertEdit(willExist) {
    let cert = Setting.deepCopy(this.state.cert);
    CertBackend.updateCert(this.state.cert.owner, this.state.certName, cert)
      .then((res) => {
        if (res.msg === "") {
          Setting.showMessage("success", `Successfully saved`);
          this.setState({
            certName: this.state.cert.name,
          });

          if (willExist) {
            this.props.history.push(`/certs`);
          } else {
            this.props.history.push(`/certs/${this.state.cert.name}`);
          }
        } else {
          Setting.showMessage("error", res.msg);
          this.updateCertField('name', this.state.certName);
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
          this.state.cert !== null ? this.renderCert() : null
        }
        <div style={{marginTop: '20px', marginLeft: '40px'}}>
          <Button size="large" onClick={() => this.submitCertEdit(false)}>{i18next.t("general:Save")}</Button>
          <Button style={{marginLeft: '20px'}} type="primary" size="large" onClick={() => this.submitCertEdit(true)}>{i18next.t("general:Save & Exit")}</Button>
        </div>
      </div>
    );
  }
}

export default CertEditPage;