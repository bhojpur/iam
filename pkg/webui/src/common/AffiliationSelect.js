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
import {Cascader, Col, Input, Row, Select} from 'antd';
import i18next from "i18next";
import * as UserBackend from "../backend/UserBackend";
import * as Setting from "../Setting";

const { Option } = Select;

class AffiliationSelect extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      classes: props,
      addressOptions: [],
      affiliationOptions: [],
    };
  }

  componentWillMount() {
    this.getAddressOptions(this.props.application);
    this.getAffiliationOptions(this.props.application, this.props.user);
  }

  getAddressOptions(application) {
    if (application.affiliationUrl === "") {
      return;
    }

    const addressUrl = application.affiliationUrl.split("|")[0];
    UserBackend.getAddressOptions(addressUrl)
      .then((addressOptions) => {
        this.setState({
          addressOptions: addressOptions,
        });
      });
  }

  getAffiliationOptions(application, user) {
    if (application.affiliationUrl === "") {
      return;
    }

    const affiliationUrl = application.affiliationUrl.split("|")[1];
    const code = user.address[user.address.length - 1];
    UserBackend.getAffiliationOptions(affiliationUrl, code)
      .then((affiliationOptions) => {
        this.setState({
          affiliationOptions: affiliationOptions,
        });
      });
  }

  updateUserField(key, value) {
    this.props.onUpdateUserField(key, value);
  }

  render() {
    return (
      <React.Fragment>
        {
          this.props.application?.affiliationUrl === "" ? null : (
            <Row style={{marginTop: '20px'}} >
              <Col style={{marginTop: '5px'}} span={this.props.labelSpan}>
                {Setting.getLabel(i18next.t("user:Address"), i18next.t("user:Address - Tooltip"))} :
              </Col>
              <Col span={24 - this.props.labelSpan} >
                <Cascader style={{width: '100%', maxWidth: '400px'}} value={this.props.user.address} options={this.state.addressOptions} onChange={value => {
                  this.updateUserField('address', value);
                  this.updateUserField('affiliation', '');
                  this.updateUserField('score', 0);
                  this.getAffiliationOptions(this.props.application, this.props.user);
                }} placeholder={i18next.t("signup:Please input your address!")} />
              </Col>
            </Row>
          )
        }
        <Row style={{marginTop: '20px'}} >
          <Col style={{marginTop: '5px'}} span={this.props.labelSpan}>
            {Setting.getLabel(i18next.t("user:Affiliation"), i18next.t("user:Affiliation - Tooltip"))} :
          </Col>
          <Col span={22} >
            {
              this.props.application?.affiliationUrl === "" ? (
                <Input value={this.props.user.affiliation} onChange={e => {
                  this.updateUserField('affiliation', e.target.value);
                }} />
              ) : (
                <Select virtual={false} style={{width: '100%'}} value={this.props.user.affiliation} onChange={(value => {
                  const name = value;
                  const affiliationOption = Setting.getArrayItem(this.state.affiliationOptions, "name", name);
                  const id = affiliationOption.id;
                  this.updateUserField('affiliation', name);
                  this.updateUserField('score', id);
                })}>
                  {
                    <Option key={0} value={""}>(empty)</Option>
                  }
                  {
                    this.state.affiliationOptions.map((affiliationOption, index) => <Option key={affiliationOption.id} value={affiliationOption.name}>{affiliationOption.name}</Option>)
                  }
                </Select>
              )
            }
          </Col>
        </Row>
      </React.Fragment>
    )
  }
}

export default AffiliationSelect;