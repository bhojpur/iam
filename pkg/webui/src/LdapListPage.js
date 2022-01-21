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
import {Button, Col, Popconfirm, Row, Table} from "antd";
import * as Setting from "./Setting";
import * as LdapBackend from "./backend/LdapBackend";
import i18next from "i18next";

class LdapListPage extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      ldaps: null
    };
  }

  UNSAFE_componentWillMount() {
    this.getLdaps()
  }

  getLdaps() {
    LdapBackend.getLdaps("")
      .then((res) => {
        let ldapsData = [];
        if (res.status === "ok") {
          ldapsData = res.data;
        } else {
          Setting.showMessage("error", res.msg);
        }
        this.setState((prevState) => {
          prevState.ldaps = ldapsData;
          return prevState;
        })
      });
  }

  deleteLdap(index) {

  }

  renderTable(ldaps) {
    const columns = [
      {
        title: i18next.t("ldap:Server Name"),
        dataIndex: "serverName",
        key: "serverName",
        width: "200px",
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
        title: i18next.t("general:Organization"),
        dataIndex: "owner",
        key: "owner",
        width: "140px",
        sorter: (a, b) => a.owner.localeCompare(b.owner),
        render: (text, record, index) => {
          return (
            <Link to={`/organizations/${text}`}>
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
        title: i18next.t("ldap:Admin"),
        dataIndex: "admin",
        key: "admin",
        ellipsis: true,
        sorter: (a, b) => a.admin.localeCompare(b.admin),
      },
      {
        title: i18next.t("ldap:Auto Sync"),
        dataIndex: "autoSync",
        key: "autoSync",
        width: "100px",
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
                onConfirm={() => this.deleteLdap(index)}
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
      <div>
        <Table columns={columns} dataSource={ldaps} rowKey="id" size="middle" bordered
               pagination={{pageSize: 100}}
               title={() => (
                 <div>
                   <span>{i18next.t("general:LDAPs")}</span>
                   <Button type="primary" size="small" style={{marginLeft: "10px"}}
                           onClick={() => {
                             this.addLdap()
                           }}>{i18next.t("general:Add")}</Button>
                 </div>
               )}
               loading={ldaps === null}
        />
      </div>
    );
  }

  render() {
    return (
      <div>
        <Row style={{width: "100%"}}>
          <Col span={1}>
          </Col>
          <Col span={22}>
            {
              this.renderTable(this.state.ldaps)
            }
          </Col>
          <Col span={1}>
          </Col>
        </Row>
      </div>
    );
  }
}

export default LdapListPage;