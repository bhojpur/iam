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
import {Button, Col, Popconfirm, Row, Table} from 'antd';
import * as Setting from "./Setting";
import i18next from "i18next";
import * as LdapBackend from "./backend/LdapBackend";
import {Link} from "react-router-dom";

class LdapTable extends React.Component {
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

  newLdap() {
    return {
      id: "",
      owner: this.props.organizationName,
      createdTime: "",
      serverName: "Bhojpur LDAP Server",
      host: "ldap.bhojpur.net",
      port: 389,
      admin: "cn=shashi.rai,dc=bhojpur,dc=net",
      passwd: "123",
      baseDn: "ou=People,dc=bhojpur,dc=net",
      autosync: 0,
      lastSync: ""
    }
  }

  addRow(table) {
    const newLdap = this.newLdap();
    LdapBackend.addLdap(newLdap)
      .then((res) => {
          if (res.status === "ok") {
            Setting.showMessage("success", `Add LDAP server success`);
            if (table === undefined) {
              table = [];
            }
            table = Setting.addRow(table, res.data2);
            this.updateTable(table);
          } else {
            Setting.showMessage("error", res.msg);
          }
        }
      )
      .catch(error => {
        Setting.showMessage("error", `Add LDAP server failed: ${error}`);
      });
  }

  deleteRow(table, i) {
    LdapBackend.deleteLdap(table[i])
      .then((res) => {
          if (res.status === "ok") {
            Setting.showMessage("success", `Delete LDAP server success`);
            table = Setting.deleteRow(table, i);
            this.updateTable(table);
          } else {
            Setting.showMessage("error", res.msg);
          }
        }
      )
      .catch(error => {
        Setting.showMessage("error", `Delete LDAP server failed: ${error}`);
      });
  }

  renderTable(table) {
    const columns = [
      {
        title: i18next.t("ldap:Server Name"),
        dataIndex: "serverName",
        key: "serverName",
        width: "160px",
        sorter: (a, b) => a.serverName.localeCompare(b.serverName),
        render: (text, record, index) => {
          return (
            <Link to={`/ldaps/${record.id}`}>
              {text}
            </Link>
          )
        }
      },
      {
        title: i18next.t("ldap:Server"),
        dataIndex: "host",
        key: "host",
        ellipsis: true,
        sorter: (a, b) => a.host.localeCompare(b.host),
        render: (text, record, index) => {
          return `${text}:${record.port}`
        }
      },
      {
        title: i18next.t("ldap:Base DN"),
        dataIndex: "baseDn",
        key: "baseDn",
        ellipsis: true,
        sorter: (a, b) => a.baseDn.localeCompare(b.baseDn),
      },
      {
        title: i18next.t("ldap:Auto Sync"),
        dataIndex: "autoSync",
        key: "autoSync",
        width: "120px",
        sorter: (a, b) => a.autoSync.localeCompare(b.autoSync),
        render: (text, record, index) => {
          return text === 0 ? (<span style={{color: "#faad14"}}>Disable</span>) : (
            <span style={{color: "#52c41a"}}>{text + " mins"}</span>)
        }
      },
      {
        title: i18next.t("ldap:Last Sync"),
        dataIndex: "lastSync",
        key: "lastSync",
        ellipsis: true,
        sorter: (a, b) => a.lastSync.localeCompare(b.lastSync),
        render: (text, record, index) => {
          return text
        }
      },
      {
        title: i18next.t("general:Action"),
        dataIndex: "",
        key: "op",
        width: "240px",
        render: (text, record, index) => {
          return (
            <div>
              <Button style={{marginTop: "10px", marginBottom: "10px", marginRight: "10px"}}
                      type="primary"
                      onClick={() => Setting.goToLink(`/ldap/sync/${record.id}`)}>{i18next.t("ldap:Sync")}</Button>
              <Button style={{marginTop: "10px", marginBottom: "10px", marginRight: "10px"}}
                      onClick={() => Setting.goToLink(`/ldap/${record.id}`)}>{i18next.t("general:Edit")}</Button>
              <Popconfirm
                title={`Sure to delete LDAP Config: ${record.serverName} ?`}
                onConfirm={() => this.deleteRow(table, index)}
              >
                <Button style={{marginBottom: "10px"}}
                        type="danger">{i18next.t("general:Delete")}</Button>
              </Popconfirm>
            </div>
          )
        }
      },
    ];

    return (
      <Table scroll={{x: 'max-content'}} rowKey="id" columns={columns} dataSource={table} size="middle" bordered pagination={false}
             title={() => (
               <div>
                 {this.props.title}&nbsp;&nbsp;&nbsp;&nbsp;
                 <Button style={{marginRight: "5px"}} type="primary" size="small"
                         onClick={() => this.addRow(table)}>{i18next.t("general:Add")}</Button>
               </div>
             )}
      />
    );
  }

  render() {
    return (
      <div>
        <Row style={{marginTop: '20px'}}>
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

export default LdapTable;