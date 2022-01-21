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
import {Switch, Table} from 'antd';
import * as Setting from "./Setting";
import * as RecordBackend from "./backend/RecordBackend";
import i18next from "i18next";
import moment from "moment";
import BaseListPage from "./BaseListPage";

class RecordListPage extends BaseListPage {

  UNSAFE_componentWillMount() {
    this.state.pagination.pageSize = 20;
    const { pagination } = this.state;
    this.fetch({ pagination });
  }

  newRecord() {
    return {
      owner: "built-in",
      name: "1234",
      id : "1234",
      clientIp: "::1",
      timestamp: moment().format(),
      organization: "built-in",
      username: "admin",
      requestUri: "/api/get-account",
      action: "login",
      isTriggered: false,
    }
  }

  renderTable(records) {
    const columns = [
      {
        title: i18next.t("general:Name"),
        dataIndex: 'name',
        key: 'name',
        width: '320px',
        sorter: true,
        ...this.getColumnSearchProps('name'),
      },
      {
        title: i18next.t("general:ID"),
        dataIndex: 'id',
        key: 'id',
        width: '90px',
        sorter: true,
        ...this.getColumnSearchProps('id'),
      },
      {
        title: i18next.t("general:Client IP"),
        dataIndex: 'clientIp',
        key: 'clientIp',
        width: '150px',
        sorter: true,
        ...this.getColumnSearchProps('clientIp'),
        render: (text, record, index) => {
          return (
            <a target="_blank" rel="noreferrer" href={`https://db-ip.com/${text}`}>
              {text}
            </a>
          )
        }
      },
      {
        title: i18next.t("general:Timestamp"),
        dataIndex: 'createdTime',
        key: 'createdTime',
        width: '180px',
        sorter: true,
        render: (text, record, index) => {
          return Setting.getFormattedDate(text);
        }
      },
      {
        title: i18next.t("general:Organization"),
        dataIndex: 'organization',
        key: 'organization',
        width: '110px',
        sorter: true,
        ...this.getColumnSearchProps('organization'),
        render: (text, record, index) => {
          return (
            <Link to={`/organizations/${text}`}>
              {text}
            </Link>
          )
        }
      },
      {
        title: i18next.t("general:User"),
        dataIndex: 'user',
        key: 'user',
        width: '120px',
        sorter: true,
        ...this.getColumnSearchProps('user'),
        render: (text, record, index) => {
          return (
            <Link to={`/users/${record.organization}/${record.user}`}>
              {text}
            </Link>
          )
        }
      },
      {
        title: i18next.t("general:Method"),
        dataIndex: 'method',
        key: 'method',
        width: '110px',
        sorter: true,
        filterMultiple: false,
        filters: [
          {text: 'GET', value: 'GET'},
          {text: 'HEAD', value: 'HEAD'},
          {text: 'POST', value: 'POST'},
          {text: 'PUT', value: 'PUT'},
          {text: 'DELETE', value: 'DELETE'},
          {text: 'CONNECT', value: 'CONNECT'},
          {text: 'OPTIONS', value: 'OPTIONS'},
          {text: 'TRACE', value: 'TRACE'},
          {text: 'PATCH', value: 'PATCH'},
        ],
      },
      {
        title: i18next.t("general:Request URI"),
        dataIndex: 'requestUri',
        key: 'requestUri',
        // width: '300px',
        sorter: true,
        ...this.getColumnSearchProps('requestUri'),
      },
      {
        title: i18next.t("general:Action"),
        dataIndex: 'action',
        key: 'action',
        width: '200px',
        sorter: true,
        ...this.getColumnSearchProps('action'),
        fixed: (Setting.isMobile()) ? "false" : "right",
        render: (text, record, index) => {
          return text;
        }
      },
      {
        title: i18next.t("record:Is Triggered"),
        dataIndex: 'isTriggered',
        key: 'isTriggered',
        width: '140px',
        sorter: true,
        fixed: (Setting.isMobile()) ? "false" : "right",
        render: (text, record, index) => {
          if (!["signup", "login", "logout", "update-user"].includes(record.action)) {
            return null;
          }

          return (
            <Switch disabled checkedChildren="ON" unCheckedChildren="OFF" checked={text} />
          )
        }
      },
    ];

    const paginationProps = {
      total: this.state.pagination.total,
      pageSize: this.state.pagination.pageSize,
      showQuickJumper: true,
      showSizeChanger: true,
      showTotal: () => i18next.t("general:{total} in total").replace("{total}", this.state.pagination.total),
    };

    return (
      <div>
        <Table scroll={{x: 'max-content'}} columns={columns} dataSource={records} rowKey="id" size="middle" bordered pagination={paginationProps}
               title={() => (
                 <div>
                   {i18next.t("general:Records")}&nbsp;&nbsp;&nbsp;&nbsp;
                 </div>
               )}
               loading={this.state.loading}
               onChange={this.handleTableChange}
        />
      </div>
    );
  }

  fetch = (params = {}) => {
    let field = params.searchedColumn, value = params.searchText;
    let sortField = params.sortField, sortOrder = params.sortOrder;
    if (params.method !== undefined && params.method !== null) {
      field = "method";
      value = params.method;
    }
    this.setState({ loading: true });
    RecordBackend.getRecords(params.pagination.current, params.pagination.pageSize, field, value, sortField, sortOrder)
      .then((res) => {
        if (res.status === "ok") {
          this.setState({
            loading: false,
            data: res.data,
            pagination: {
              ...params.pagination,
              total: res.data2,
            },
            searchText: params.searchText,
            searchedColumn: params.searchedColumn,
          });
        }
      });
  };
}

export default RecordListPage;