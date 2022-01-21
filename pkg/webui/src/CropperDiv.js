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

import React, {useState} from "react";
import Cropper from "react-cropper";
import "cropperjs/dist/cropper.css";
import * as Setting from "./Setting";
import {Button, Row, Col, Modal} from 'antd';
import i18next from "i18next";
import * as ResourceBackend from "./backend/ResourceBackend";

export const CropperDiv = (props) => {
    const [image, setImage] = useState("");
    const [cropper, setCropper] = useState();
    const [visible, setVisible] = React.useState(false);
    const [confirmLoading, setConfirmLoading] = React.useState(false);
    const {title} = props;
    const {user} = props;
    const {account} = props;
    const {buttonText} = props;
    let uploadButton;

    const onChange = (e) => {
        e.preventDefault();
        let files;
        if (e.dataTransfer) {
            files = e.dataTransfer.files;
        } else if (e.target) {
            files = e.target.files;
        }
        const reader = new FileReader();
        reader.onload = () => {
            setImage(reader.result);
        };
        if (!(files[0] instanceof Blob)) {
            return;
        }
        reader.readAsDataURL(files[0]);
    };

    const uploadAvatar = () => {
        cropper.getCroppedCanvas().toBlob(blob => {
            if (blob === null) {
                Setting.showMessage("error", "You must select a picture first!");
                return false;
            }
            // Setting.showMessage("success", "uploading...");
            const extension = image.substring(image.indexOf('/') + 1, image.indexOf(';base64'));
            const fullFilePath = `avatar/${user.owner}/${user.name}.${extension}`;
            ResourceBackend.uploadResource(user.owner, user.name, "avatar", "CropperDiv", fullFilePath, blob)
              .then((res) => {
                  if (res.status === "ok") {
                      window.location.href = "/account";
                  } else {
                      Setting.showMessage("error", res.msg);
                  }
              });
            return true;
        });
    }

    const showModal = () => {
        setVisible(true);
    };

    const handleOk = () => {
        setConfirmLoading(true);
        if (!uploadAvatar()) {
            setConfirmLoading(false);
        }
    };

    const handleCancel = () => {
        console.log('Clicked cancel button');
        setVisible(false);
    };

    const selectFile = () => {
        uploadButton.click();
    }

    return (
        <div>
            <Button type="default" onClick={showModal}>
                {buttonText}
            </Button>
            <Modal
                maskClosable={false}
                title={title}
                visible={visible}
                okText={i18next.t("user:Upload a photo")}
                confirmLoading={confirmLoading}
                onCancel={handleCancel}
                width={600}
                footer={
                    [<Button block type="primary" onClick={handleOk}>{i18next.t("user:Set new profile picture")}</Button>]
                }
            >
                <Col style={{margin: "0px auto 40px auto", width: 1000, height: 300}}>
                    <Row style={{width: "100%", marginBottom: "20px"}}>
                        <input style={{display: "none"}} ref={input => uploadButton = input} type="file" accept="image/*" onChange={onChange}/>
                        <Button block onClick={selectFile}>{i18next.t("user:Select a photo...")}</Button>
                    </Row>
                    <Cropper
                      style={{height: "100%"}}
                      initialAspectRatio={1}
                      preview=".img-preview"
                      src={image}
                      viewMode={1}
                      guides={true}
                      minCropBoxHeight={10}
                      minCropBoxWidth={10}
                      background={false}
                      responsive={true}
                      autoCropArea={1}
                      checkOrientation={false}
                      onInitialized={(instance) => {
                          setCropper(instance);
                      }}
                    />
                </Col>
            </Modal>
        </div>
    )
}

export default CropperDiv;