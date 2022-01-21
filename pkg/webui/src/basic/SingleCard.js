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
import {Card, Col} from "antd";
import * as Setting from "../Setting";
import {withRouter} from "react-router-dom";

const { Meta } = Card;

class SingleCard extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      classes: props,
    };
  }

  renderCardMobile(logo, link, title, desc, time, isSingle) {
    const gridStyle = {
      width: '100vw',
      textAlign: 'center',
      cursor: 'pointer',
    };

    return (
      <Card.Grid style={gridStyle} onClick={() => Setting.goToLinkSoft(this, link)}>
        <img src={logo} alt="logo" height={60} style={{marginBottom: '20px'}}/>
        <Meta
          title={title}
          description={desc}
        />
      </Card.Grid>
    )
  }

  renderCard(logo, link, title, desc, time, isSingle) {
    return (
      <Col style={{paddingLeft: "20px", paddingRight: "20px", paddingBottom: "20px", marginBottom: "20px"}} span={6}>
        <Card
          hoverable
          cover={
            <img alt="logo" src={logo} width={"100%"} height={"100%"} />
          }
          onClick={() => Setting.goToLinkSoft(this, link)}
          style={isSingle ? {width: "320px"} : null}
        >
          <Meta title={title} description={desc} />
          <br/>
          <br/>
          <Meta title={""} description={Setting.getFormattedDateShort(time)} />
        </Card>
      </Col>
    )
  }

  render() {
    if (Setting.isMobile()) {
      return this.renderCardMobile(this.props.logo, this.props.link, this.props.title, this.props.desc, this.props.time, this.props.isSingle);
    } else {
      return this.renderCard(this.props.logo, this.props.link, this.props.title, this.props.desc, this.props.time, this.props.isSingle);
    }
  }
}

export default withRouter(SingleCard);