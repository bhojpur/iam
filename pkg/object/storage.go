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
	"strings"

	storage "github.com/bhojpur/drive/pkg/storage"
	"github.com/bhojpur/iam/pkg/utils"
	websvr "github.com/bhojpur/web/pkg/engine"
)

var isCloudIntranet bool

func init() {
	var err error
	isCloudIntranet, err = websvr.AppConfig.Bool("isCloudIntranet")
	if err != nil {
		//panic(err)
	}
}

func getProviderEndpoint(provider *Provider) string {
	endpoint := provider.Endpoint
	if provider.IntranetEndpoint != "" && isCloudIntranet {
		endpoint = provider.IntranetEndpoint
	}
	return endpoint
}

func getUploadFileUrl(provider *Provider, fullFilePath string, hasTimestamp bool) (string, string) {
	objectKey := utils.UrlJoin(utils.GetUrlPath(provider.Domain), fullFilePath)

	host := ""
	if provider.Type != "Local File System" {
		// provider.Domain = "https://iam.bhojpur.net/sso/"
		host = utils.GetUrlHost(provider.Domain)
		if !strings.HasPrefix(host, "http://") && !strings.HasPrefix(host, "https://") {
			host = fmt.Sprintf("https://%s", host)
		}
	} else {
		// provider.Domain = "http://localhost:8000" or "https://iam.bhojpur.net"
		host = utils.UrlJoin(provider.Domain, "/files")
	}

	fileUrl := utils.UrlJoin(host, objectKey)
	if hasTimestamp {
		fileUrl = fmt.Sprintf("%s?t=%s", utils.UrlJoin(host, objectKey), utils.GetCurrentUnixTime())
	}

	return fileUrl, objectKey
}

func uploadFile(provider *Provider, fullFilePath string, fileBuffer *bytes.Buffer) (string, string, error) {
	endpoint := getProviderEndpoint(provider)
	storageProvider := storage.GetStorageProvider(provider.Type, provider.ClientId, provider.ClientSecret, provider.RegionId, provider.Bucket, endpoint)
	if storageProvider == nil {
		return "", "", fmt.Errorf("the provider type: %s is not supported", provider.Type)
	}

	if provider.Domain == "" {
		provider.Domain = storageProvider.GetEndpoint()
		UpdateProvider(provider.GetId(), provider)
	}

	fileUrl, objectKey := getUploadFileUrl(provider, fullFilePath, true)

	_, err := storageProvider.Put(objectKey, fileBuffer)
	if err != nil {
		return "", "", err
	}

	return fileUrl, objectKey, nil
}

func UploadFileSafe(provider *Provider, fullFilePath string, fileBuffer *bytes.Buffer) (string, string, error) {
	var fileUrl string
	var objectKey string
	var err error
	times := 0
	for {
		fileUrl, objectKey, err = uploadFile(provider, fullFilePath, fileBuffer)
		if err != nil {
			times += 1
			if times >= 5 {
				return "", "", err
			}
		} else {
			break
		}
	}
	return fileUrl, objectKey, nil
}

func DeleteFile(provider *Provider, objectKey string) error {
	endpoint := getProviderEndpoint(provider)
	storageProvider := storage.GetStorageProvider(provider.Type, provider.ClientId, provider.ClientSecret, provider.RegionId, provider.Bucket, endpoint)
	if storageProvider == nil {
		return fmt.Errorf("the provider type: %s is not supported", provider.Type)
	}

	if provider.Domain == "" {
		provider.Domain = storageProvider.GetEndpoint()
		UpdateProvider(provider.GetId(), provider)
	}

	return storageProvider.Delete(objectKey)
}
