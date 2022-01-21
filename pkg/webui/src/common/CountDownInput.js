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

import {Button, Col, Input, Modal, Row} from "antd";
import React from "react";
import * as Setting from "../Setting";
import i18next from "i18next";
import * as UserBackend from "../backend/UserBackend";
import {SafetyOutlined} from "@ant-design/icons";
import * as Util from "../auth/Util";
import {isValidEmail, isValidPhone} from "../Setting";

const { Search } = Input;

export const CountDownInput = (props) => {
  const {disabled, textBefore, onChange, onButtonClickArgs} = props;
  const [visible, setVisible] = React.useState(false);
  const [key, setKey] = React.useState("");
  const [captchaImg, setCaptchaImg] = React.useState("");
  const [checkType, setCheckType] = React.useState("");
  const [checkId, setCheckId] = React.useState("");
  const [buttonLeftTime, setButtonLeftTime] = React.useState(0);
  const [buttonLoading, setButtonLoading] = React.useState(false);

  const handleCountDown = (leftTime = 60) => {
    let leftTimeSecond = leftTime
    setButtonLeftTime(leftTimeSecond)
    const countDown = () => {
      leftTimeSecond--;
      setButtonLeftTime(leftTimeSecond)
      if (leftTimeSecond === 0) {
        return;
      }
      setTimeout(countDown, 1000);
    }
    setTimeout(countDown, 1000);
  }

  const handleOk = () => {
    setVisible(false);
    if (isValidEmail(onButtonClickArgs[0])) {
        onButtonClickArgs[1] = "email";
    } else if (isValidPhone(onButtonClickArgs[0])) {
        onButtonClickArgs[1] = "phone";
    } else {
        Util.showMessage("error", i18next.t("login:Invalid Email or phone"))
        return;
    }
    setButtonLoading(true)
    UserBackend.sendCode(checkType, checkId, key, ...onButtonClickArgs).then(res => {
      setKey("");
      setButtonLoading(false)
      if (res) {
        handleCountDown(60);
      }
    })
  }

  const handleCancel = () => {
    setVisible(false);
    setKey("");
  }

  const loadHumanCheck = () => {
    UserBackend.getHumanCheck().then(res => {
      if (res.type === "none") {
        UserBackend.sendCode("none", "", "", ...onButtonClickArgs);
      } else if (res.type === "captcha") {
        setCheckId(res.captchaId);
        setCaptchaImg(res.captchaImage);
        setCheckType("captcha");
        setVisible(true);
      } else {
        Setting.showMessage("error", i18next.t("signup:Unknown Check Type"));
      }
    })
  }

  const renderCaptcha = () => {
    return (
      <Col>
        <Row
          style={{
            backgroundImage: `url('data:image/png;base64,${captchaImg}')`,
            backgroundRepeat: "no-repeat",
            height: "80px",
            width: "200px",
            borderRadius: "3px",
            border: "1px solid #ccc",
            marginBottom: 10
          }}
        />
        <Row>
          <Input autoFocus value={key} prefix={<SafetyOutlined/>} placeholder={i18next.t("general:Captcha")} onPressEnter={handleOk} onChange={e => setKey(e.target.value)}/>
        </Row>
      </Col>
    )
  }

  const renderCheck = () => {
    if (checkType === "captcha") return renderCaptcha();
    return null;
  }

  return (
    <React.Fragment>
      <Search
        addonBefore={textBefore}
        disabled={disabled}
        prefix={<SafetyOutlined/>}
        placeholder={i18next.t("code:Enter your code")}
        onChange={e => onChange(e.target.value)}
        enterButton={
          <Button style={{fontSize: 14}} type={"primary"} disabled={disabled || buttonLeftTime > 0} loading={buttonLoading}>
            {buttonLeftTime > 0 ? `${buttonLeftTime} s` : buttonLoading ? i18next.t("code:Sending Code") : i18next.t("code:Send Code")}
          </Button>
        }
        onSearch={loadHumanCheck}
      />
      <Modal
        closable={false}
        maskClosable={false}
        destroyOnClose={true}
        title={i18next.t("general:Captcha")}
        visible={visible}
        okText={i18next.t("user:OK")}
        cancelText={i18next.t("user:Cancel")}
        onOk={handleOk}
        onCancel={handleCancel}
        okButtonProps={{disabled: key.length !== 5}}
        width={248}
      >
        {
          renderCheck()
        }
      </Modal>
    </React.Fragment>
  );
}