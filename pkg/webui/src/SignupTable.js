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
import {DownOutlined, DeleteOutlined, UpOutlined} from '@ant-design/icons';
import {Button, Col, Row, Select, Switch, Table, Tooltip} from 'antd';
import * as Setting from "./Setting";
import i18next from "i18next";

const { Option } = Select;

class SignupTable extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      classes: props,
    };
  }

  updateTable(table) {
    this.props.onUpdateTable(table);
  }

  updateField(table, index, key, value) {
    table[index][key] = value;
    this.updateTable(table);
  }

  addRow(table) {
    let row = {name: Setting.getNewRowNameForTable(table, "Please select a signup item"), visible: true, required: true, rule: "None"};
    if (table === undefined) {
      table = [];
    }
    table = Setting.addRow(table, row);
    this.updateTable(table);
  }

  deleteRow(table, i) {
    table = Setting.deleteRow(table, i);
    this.updateTable(table);
  }

  upRow(table, i) {
    table = Setting.swapRow(table, i - 1, i);
    this.updateTable(table);
  }

  downRow(table, i) {
    table = Setting.swapRow(table, i, i + 1);
    this.updateTable(table);
  }

  renderTable(table) {
    const columns = [
      {
        title: i18next.t("provider:Name"),
        dataIndex: 'name',
        key: 'name',
        render: (text, record, index) => {
          const items = [
            {id: 'Username', name: 'Username'},
            {id: 'ID', name: 'ID'},
            {id: 'Display name', name: 'Display name'},
            {id: 'Affiliation', name: 'Affiliation'},
            {id: 'Country/Region', name: 'Country/Region'},
            {id: 'ID card', name: 'ID card'},
            {id: 'Email', name: 'Email'},
            {id: 'Password', name: 'Password'},
            {id: 'Confirm password', name: 'Confirm password'},
            {id: 'Phone', name: 'Phone'},
            {id: 'Agreement', name: 'Agreement'},
          ];

          return (
            <Select virtual={false} style={{width: '100%'}}
                    value={text}
                    onChange={value => {
                      this.updateField(table, index, 'name', value);
                    }} >
              {
                Setting.getDeduplicatedArray(items, table, "name").map((item, index) => <Option key={index} value={item.name}>{item.name}</Option>)
              }
            </Select>
          )
        }
      },
      {
        title: i18next.t("provider:visible"),
        dataIndex: 'visible',
        key: 'visible',
        width: '120px',
        render: (text, record, index) => {
          if (record.name === "ID") {
            return null;
          }

          return (
            <Switch checked={text} onChange={checked => {
              this.updateField(table, index, 'visible', checked);
              if (!checked) {
                this.updateField(table, index, 'required', false);
              } else {
                this.updateField(table, index, 'required', true);
              }
            }} />
          )
        }
      },
      {
        title: i18next.t("provider:required"),
        dataIndex: 'required',
        key: 'required',
        width: '120px',
        render: (text, record, index) => {
          if (!record.visible) {
            return null;
          }

          return (
            <Switch checked={text} onChange={checked => {
              this.updateField(table, index, 'required', checked);
            }} />
          )
        }
      },
      {
        title: i18next.t("provider:prompted"),
        dataIndex: 'prompted',
        key: 'prompted',
        width: '120px',
        render: (text, record, index) => {
          if (record.name === "ID") {
            return null;
          }

          if (record.visible) {
            return null;
          }

          return (
            <Switch checked={text} onChange={checked => {
              this.updateField(table, index, 'prompted', checked);
            }} />
          )
        }
      },
      {
        title: i18next.t("provider:rule"),
        dataIndex: 'rule',
        key: 'rule',
        width: '120px',
        render: (text, record, index) => {
          let options = [];
          if (record.name === "ID") {
            options = [
              {id: 'Random', name: 'Random'},
              {id: 'Incremental', name: 'Incremental'},
            ];
          } if (record.name === "Display name") {
            options = [
              {id: 'None', name: 'None'},
              {id: 'Personal', name: 'Personal'},
            ];
          }

          if (options.length === 0) {
            return null;
          }

          return (
            <Select virtual={false} style={{width: '100%'}} value={text} onChange={(value => {
              this.updateField(table, index, 'rule', value);
            })}>
              {
                options.map((item, index) => <Option key={index} value={item.id}>{item.name}</Option>)
              }
            </Select>
          )
        }
      },
      {
        title: i18next.t("general:Action"),
        key: 'action',
        width: '100px',
        render: (text, record, index) => {
          return (
            <div>
              <Tooltip placement="bottomLeft" title={i18next.t("general:Up")}>
                <Button style={{marginRight: "5px"}} disabled={index === 0} icon={<UpOutlined />} size="small" onClick={() => this.upRow(table, index)} />
              </Tooltip>
              <Tooltip placement="topLeft" title={i18next.t("general:Down")}>
                <Button style={{marginRight: "5px"}} disabled={index === table.length - 1} icon={<DownOutlined />} size="small" onClick={() => this.downRow(table, index)} />
              </Tooltip>
              <Tooltip placement="topLeft" title={i18next.t("general:Delete")}>
                <Button icon={<DeleteOutlined />} size="small" onClick={() => this.deleteRow(table, index)} />
              </Tooltip>
            </div>
          );
        }
      },
    ];

    return (
      <Table scroll={{x: 'max-content'}} rowKey="name" columns={columns} dataSource={table} size="middle" bordered pagination={false}
             title={() => (
               <div>
                 {this.props.title}&nbsp;&nbsp;&nbsp;&nbsp;
                 <Button style={{marginRight: "5px"}} type="primary" size="small" onClick={() => this.addRow(table)}>{i18next.t("general:Add")}</Button>
               </div>
             )}
      />
    );
  }

  render() {
    return (
      <div>
        <Row style={{marginTop: '20px'}} >
          <Col span={24}>
            {
              this.renderTable(this.props.table)
            }
          </Col>
        </Row>
      </div>
    )
  }
}

export default SignupTable;