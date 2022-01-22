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
import {Button, Popconfirm, Switch, Table} from 'antd';
import moment from "moment";
import * as Setting from "./Setting";
import * as SyncerBackend from "./backend/SyncerBackend";
import i18next from "i18next";
import BaseListPage from "./BaseListPage";

class SyncerListPage extends BaseListPage {

  newSyncer() {
    const randomName = Setting.getRandomName();
    return {
      owner: "admin",
      name: `syncer_${randomName}`,
      createdTime: moment().format(),
      organization: "built-in",
      type: "Database",
      host: "localhost",
      port: 3306,
      user: "root",
      password: "123456",
      databaseType: "mysql",
      database: "dbName",
      table: "tableName",
      tableColumns: [],
      affiliationTable: "",
      avatarBaseUrl: "",
      syncInterval: 10,
      isEnabled: true,
    }
  }

  addSyncer() {
    const newSyncer = this.newSyncer();
    SyncerBackend.addSyncer(newSyncer)
      .then((res) => {
          Setting.showMessage("success", `Syncer added successfully`);
          this.props.history.push(`/syncers/${newSyncer.name}`);
        }
      )
      .catch(error => {
        Setting.showMessage("error", `Syncer failed to add: ${error}`);
      });
  }

  deleteSyncer(i) {
    SyncerBackend.deleteSyncer(this.state.data[i])
      .then((res) => {
          Setting.showMessage("success", `Syncer deleted successfully`);
          this.setState({
            data: Setting.deleteRow(this.state.data, i),
            pagination: {total: this.state.pagination.total - 1},
          });
        }
      )
      .catch(error => {
        Setting.showMessage("error", `Syncer failed to delete: ${error}`);
      });
  }

  renderTable(syncers) {
    const columns = [
      {
        title: i18next.t("general:Organization"),
        dataIndex: 'organization',
        key: 'organization',
        width: '120px',
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
        title: i18next.t("general:Name"),
        dataIndex: 'name',
        key: 'name',
        width: '150px',
        fixed: 'left',
        sorter: true,
        ...this.getColumnSearchProps('name'),
        render: (text, record, index) => {
          return (
            <Link to={`/syncers/${text}`}>
              {text}
            </Link>
          )
        }
      },
      {
        title: i18next.t("general:Created time"),
        dataIndex: 'createdTime',
        key: 'createdTime',
        width: '160px',
        sorter: true,
        render: (text, record, index) => {
          return Setting.getFormattedDate(text);
        }
      },
      {
        title: i18next.t("provider:Type"),
        dataIndex: 'type',
        key: 'type',
        width: '100px',
        sorter: true,
        filterMultiple: false,
        filters: [
          {text: 'Database', value: 'Database'},
          {text: 'LDAP', value: 'LDAP'},
        ],
      },
      {
        title: i18next.t("provider:Host"),
        dataIndex: 'host',
        key: 'host',
        width: '120px',
        sorter: true,
        ...this.getColumnSearchProps('host'),
      },
      {
        title: i18next.t("provider:Port"),
        dataIndex: 'port',
        key: 'port',
        width: '100px',
        sorter: true,
        ...this.getColumnSearchProps('port'),
      },
      {
        title: i18next.t("general:User"),
        dataIndex: 'user',
        key: 'user',
        width: '120px',
        sorter: true,
        ...this.getColumnSearchProps('user'),
      },
      {
        title: i18next.t("general:Password"),
        dataIndex: 'password',
        key: 'password',
        width: '120px',
        sorter: true,
        ...this.getColumnSearchProps('password'),
      },
      {
        title: i18next.t("syncer:Database type"),
        dataIndex: 'databaseType',
        key: 'databaseType',
        width: '120px',
        sorter: (a, b) => a.databaseType.localeCompare(b.databaseType),
      },
      {
        title: i18next.t("syncer:Database"),
        dataIndex: 'database',
        key: 'database',
        width: '120px',
        sorter: true,
      },
      {
        title: i18next.t("syncer:Table"),
        dataIndex: 'table',
        key: 'table',
        width: '120px',
        sorter: true,
      },
      {
        title: i18next.t("syncer:Sync interval"),
        dataIndex: 'syncInterval',
        key: 'syncInterval',
        width: '130px',
        sorter: true,
        ...this.getColumnSearchProps('syncInterval'),
      },
      {
        title: i18next.t("general:Is enabled"),
        dataIndex: 'isEnabled',
        key: 'isEnabled',
        width: '120px',
        sorter: true,
        render: (text, record, index) => {
          return (
            <Switch disabled checkedChildren="ON" unCheckedChildren="OFF" checked={text} />
          )
        }
      },
      {
        title: i18next.t("general:Action"),
        dataIndex: '',
        key: 'op',
        width: '170px',
        fixed: (Setting.isMobile()) ? "false" : "right",
        render: (text, record, index) => {
          return (
            <div>
              <Button style={{marginTop: '10px', marginBottom: '10px', marginRight: '10px'}} type="primary" onClick={() => this.props.history.push(`/syncers/${record.name}`)}>{i18next.t("general:Edit")}</Button>
              <Popconfirm
                title={`Sure to delete syncer: ${record.name} ?`}
                onConfirm={() => this.deleteSyncer(index)}
              >
                <Button style={{marginBottom: '10px'}} type="danger">{i18next.t("general:Delete")}</Button>
              </Popconfirm>
            </div>
          )
        }
      },
    ];

    const paginationProps = {
      total: this.state.pagination.total,
      showQuickJumper: true,
      showSizeChanger: true,
      showTotal: () => i18next.t("general:{total} in total").replace("{total}", this.state.pagination.total),
    };

    return (
      <div>
        <Table scroll={{x: 'max-content'}} columns={columns} dataSource={syncers} rowKey="name" size="middle" bordered pagination={paginationProps}
               title={() => (
                 <div>
                   {i18next.t("general:Syncers")}&nbsp;&nbsp;&nbsp;&nbsp;
                   <Button type="primary" size="small" onClick={this.addSyncer.bind(this)}>{i18next.t("general:Add")}</Button>
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
    if (params.type !== undefined && params.type !== null) {
      field = "type";
      value = params.type;
    }
    this.setState({ loading: true });
    SyncerBackend.getSyncers("admin", params.pagination.current, params.pagination.pageSize, field, value, sortField, sortOrder)
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

export default SyncerListPage;