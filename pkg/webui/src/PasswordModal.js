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
import * as UserBackend from "./backend/UserBackend";
import * as Setting from "./Setting";

export const PasswordModal = (props) => {
  const [visible, setVisible] = React.useState(false);
  const [confirmLoading, setConfirmLoading] = React.useState(false);
  const [oldPassword, setOldPassword] = React.useState("");
  const [newPassword, setNewPassword] = React.useState("");
  const [rePassword, setRePassword] = React.useState("");
  const {user} = props;
  const {account} = props;

  const showModal = () => {
    setVisible(true);
  };

  const handleCancel = () => {
    setVisible(false);
  };

  const handleOk = () => {
    if (newPassword === "" || rePassword === "") {
      Setting.showMessage("error", i18next.t("user:Empty input!"));
      return;
    }
    if (newPassword !== rePassword) {
      Setting.showMessage("error", i18next.t("user:Two passwords you typed do not match."));
      return;
    }
    setConfirmLoading(true);
    UserBackend.setPassword(user.owner, user.name, oldPassword, newPassword).then((res) => {
      setConfirmLoading(false);
      if (res.status === "ok") {
        Setting.showMessage("success", i18next.t("user:Password Set"));
        setVisible(false);
      }
      else Setting.showMessage("error", i18next.t(`user:${res.msg}`));
    })
  }

  let hasOldPassword = user.password !== "";

  return (
    <Row>
      <Button type="default" disabled={props.disabled} onClick={showModal}>
        { hasOldPassword ? i18next.t("user:Modify password...") : i18next.t("user:Set password...")}
      </Button>
      <Modal
        maskClosable={false}
        title={i18next.t("user:Password")}
        visible={visible}
        okText={i18next.t("user:Set Password")}
        cancelText={i18next.t("user:Cancel")}
        confirmLoading={confirmLoading}
        onCancel={handleCancel}
        onOk={handleOk}
        width={600}
      >
        <Col style={{margin: "0px auto 40px auto", width: 1000, height: 300}}>
          { (hasOldPassword && !Setting.isAdminUser(account)) ? (
            <Row style={{width: "100%", marginBottom: "20px"}}>
              <Input.Password addonBefore={i18next.t("user:Old Password")} placeholder={i18next.t("user:input password")} onChange={(e) => setOldPassword(e.target.value)}/>
            </Row>
          ) : null}
          <Row style={{width: "100%", marginBottom: "20px"}}>
            <Input.Password addonBefore={i18next.t("user:New Password")} placeholder={i18next.t("user:input password")} onChange={(e) => setNewPassword(e.target.value)}/>
          </Row>
          <Row style={{width: "100%", marginBottom: "20px"}}>
            <Input.Password addonBefore={i18next.t("user:Re-enter New")} placeholder={i18next.t("user:input password")} onChange={(e) => setRePassword(e.target.value)}/>
          </Row>
        </Col>
      </Modal>
    </Row>
  )
}

export default PasswordModal;