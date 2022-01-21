package object

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

import (
	"bytes"
	"fmt"
	"io"

	"github.com/bhojpur/iam/pkg/proxy"
	websvr "github.com/bhojpur/web/pkg/engine"
)

var defaultStorageProvider *Provider = nil

func InitDefaultStorageProvider() {
	defaultStorageProviderStr := websvr.AppConfig.String("defaultStorageProvider")
	if defaultStorageProviderStr != "" {
		defaultStorageProvider = getProvider("admin", defaultStorageProviderStr)
	}
}

func downloadFile(url string) (*bytes.Buffer, error) {
	httpClient := proxy.GetHttpClient(url)

	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	fileBuffer := bytes.NewBuffer(nil)
	_, err = io.Copy(fileBuffer, resp.Body)
	if err != nil {
		return nil, err
	}

	return fileBuffer, nil
}

func getPermanentAvatarUrl(organization string, username string, url string) string {
	if defaultStorageProvider == nil {
		return ""
	}

	fullFilePath := fmt.Sprintf("/avatar/%s/%s.png", organization, username)
	uploadedFileUrl, _ := getUploadFileUrl(defaultStorageProvider, fullFilePath, false)

	fileBuffer, err := downloadFile(url)
	if err != nil {
		panic(err)
	}

	_, _, err = UploadFileSafe(defaultStorageProvider, fullFilePath, fileBuffer)
	if err != nil {
		panic(err)
	}

	return uploadedFileUrl
}
