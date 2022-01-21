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
import * as Setting from "./Setting";
import { Menu, Dropdown} from "antd";
import { createFromIconfontCN } from '@ant-design/icons';
import './App.less';

const IconFont = createFromIconfontCN({
  scriptUrl: '//at.alicdn.com/t/font_2680620_ffij16fkwdg.js',
});

class SelectLanguageBox extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      classes: props,
    };
  }

  render() {
    const menu = (
      <Menu onClick={(e) => {
        Setting.changeLanguage(e.key);
      }}>
        <Menu.Item key="en" icon={<IconFont type="icon-en" />}>English</Menu.Item>
        <Menu.Item key="zh" icon={<IconFont type="icon-zh" />}>简体中文</Menu.Item>
        <Menu.Item key="fr" icon={<IconFont type="icon-fr" />}>Français</Menu.Item>
        <Menu.Item key="de" icon={<IconFont type="icon-de" />}>Deutsch</Menu.Item>
        <Menu.Item key="ja" icon={<IconFont type="icon-ja" />}>日本語</Menu.Item>
        <Menu.Item key="ko" icon={<IconFont type="icon-ko" />}>한국어</Menu.Item>
        <Menu.Item key="ru" icon={<IconFont type="icon-ru" />}>Русский</Menu.Item>
      </Menu>
    );

    return (
      <Dropdown overlay={menu} >
        <div className="language_box" />
      </Dropdown>
    );
  }
}

export default SelectLanguageBox;