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
import {Card, Col, Row} from "antd";
import * as Setting from "../Setting";
import SingleCard from "./SingleCard";
import i18next from "i18next";

class HomePage extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      classes: props,
    };
  }

  getItems() {
    let items = [];
    if (Setting.isAdminUser(this.props.account)) {
      items = [
        {link: "/organizations", name: i18next.t("general:Organizations"), organizer: i18next.t("general:User containers")},
        {link: "/users", name: i18next.t("general:Users"), organizer: i18next.t("general:Users under all organizations")},
        {link: "/providers", name: i18next.t("general:Providers"), organizer: i18next.t("general:OAuth providers")},
        {link: "/applications", name: i18next.t("general:Applications"), organizer: i18next.t("general:Applications that require authentication")},
      ];
    } else {
      items = [
        {link: "/account", name: i18next.t("account:My Account"), organizer: i18next.t("account:Settings for your account")},
      ];
    }

    for (let i = 0; i < items.length; i ++) {
      let filename = items[i].link;
      if (filename === "/account") {
        filename = "/users";
      }
      items[i].logo = `https://static.bhojpur.net/image/${filename}.png`;
      items[i].createdTime = "";
    }

    return items
  }

  renderCards() {
    const items = this.getItems();

    if (Setting.isMobile()) {
      return (
        <Card bodyStyle={{padding: 0}}>
          {
            items.map(item => {
              return (
                <SingleCard logo={item.logo} link={item.link} title={item.name} desc={item.organizer} isSingle={items.length === 1} />
              )
            })
          }
        </Card>
      )
    } else {
      return (
        <div style={{marginRight:'15px',marginLeft:'15px'}}>
              <Row style={{marginLeft: "-20px", marginRight: "-20px", marginTop: "20px"}} gutter={24}>
                {
                  items.map(item => {
                    return (
                      <SingleCard logo={item.logo} link={item.link} title={item.name} desc={item.organizer} time={item.createdTime} isSingle={items.length === 1} key={item.name} />
                    )
                  })
                }
              </Row>
        </div>
      )
    }
  }

  render() {
    return (
      <div>
        <Row style={{width: "100%"}}>
          <Col span={24} style={{display: "flex", justifyContent:  "center"}} >
            {
              this.renderCards()
            }
          </Col>
        </Row>
      </div>
    )
  }
}

export default HomePage;