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

import {Button, Col, Modal, Row, Input,} from "antd";
import i18next from "i18next";
import React from "react";
import * as Setting from "./Setting"
import * as UserBackend from "./backend/UserBackend"
import {CountDownInput} from "./common/CountDownInput";
import {MailOutlined, PhoneOutlined} from "@ant-design/icons";

export const ResetModal = (props) => {
  const [visible, setVisible] = React.useState(false);
  const [confirmLoading, setConfirmLoading] = React.useState(false);
  const [dest, setDest] = React.useState("");
  const [code, setCode] = React.useState("");
  const {buttonText, destType, org} = props;

  const showModal = () => {
    setVisible(true);
  };

  const handleCancel = () => {
    setVisible(false);
  };

  const handleOk = () => {
    if (dest === "") {
      Setting.showMessage("error", i18next.t("user:Empty " + destType));
      return;
    }
    if (code === "") {
      Setting.showMessage("error", i18next.t("code:Empty Code"));
      return;
    }
    setConfirmLoading(true);
    UserBackend.resetEmailOrPhone(dest, destType, code).then(res => {
      if (res.status === "ok") {
        Setting.showMessage("success", i18next.t("user:" + destType + " reset"));
        window.location.reload();
      } else {
        Setting.showMessage("error", i18next.t("user:" + res.msg));
        setConfirmLoading(false);
      }
    })
  }

  let placeHolder = "";
  if (destType === "email") placeHolder = i18next.t("user:Input your email");
  else if (destType === "phone") placeHolder = i18next.t("user:Input your phone number");

  return (
    <Row>
      <Button type="default" onClick={showModal}>
        {buttonText}
      </Button>
      <Modal
        maskClosable={false}
        title={buttonText}
        visible={visible}
        okText={buttonText}
        cancelText={i18next.t("user:Cancel")}
        confirmLoading={confirmLoading}
        onCancel={handleCancel}
        onOk={handleOk}
        width={600}
      >
        <Col style={{margin: "0px auto 40px auto", width: 1000, height: 300}}>
          <Row style={{width: "100%", marginBottom: "20px"}}>
            <Input
              addonBefore={destType === "email" ? i18next.t("user:New Email") : i18next.t("user:New phone")}
              prefix={destType === "email" ? <MailOutlined /> : <PhoneOutlined />}
              placeholder={placeHolder}
              onChange={e => setDest(e.target.value)}
            />
          </Row>
          <Row style={{width: "100%", marginBottom: "20px"}}>
            <CountDownInput
              textBefore={i18next.t("code:Code You Received")}
              onChange={setCode}
              onButtonClickArgs={[dest, destType, `${org?.owner}/${org?.name}`]}
            />
          </Row>
        </Col>
      </Modal>
    </Row>
  )
}

export default ResetModal;