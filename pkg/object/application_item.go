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

func (application *Application) GetProviderByCategory(category string) *Provider {
	providers := GetProviders(application.Owner)
	m := map[string]*Provider{}
	for _, provider := range providers {
		if provider.Category != category {
			continue
		}

		m[provider.Name] = provider
	}

	for _, providerItem := range application.Providers {
		if provider, ok := m[providerItem.Name]; ok {
			return provider
		}
	}

	return nil
}

func (application *Application) GetEmailProvider() *Provider {
	return application.GetProviderByCategory("Email")
}

func (application *Application) GetSmsProvider() *Provider {
	return application.GetProviderByCategory("SMS")
}

func (application *Application) GetStorageProvider() *Provider {
	return application.GetProviderByCategory("Storage")
}

func (application *Application) getSignupItem(itemName string) *SignupItem {
	for _, signupItem := range application.SignupItems {
		if signupItem.Name == itemName {
			return signupItem
		}
	}
	return nil
}

func (application *Application) IsSignupItemVisible(itemName string) bool {
	signupItem := application.getSignupItem(itemName)
	if signupItem == nil {
		return false
	}

	return signupItem.Visible
}

func (application *Application) IsSignupItemRequired(itemName string) bool {
	signupItem := application.getSignupItem(itemName)
	if signupItem == nil {
		return false
	}

	return signupItem.Required
}

func (application *Application) GetSignupItemRule(itemName string) string {
	signupItem := application.getSignupItem(itemName)
	if signupItem == nil {
		return ""
	}

	return signupItem.Rule
}

func (application *Application) getAllPromptedProviderItems() []*ProviderItem {
	res := []*ProviderItem{}
	for _, providerItem := range application.Providers {
		if providerItem.isProviderPrompted() {
			res = append(res, providerItem)
		}
	}
	return res
}

func (application *Application) isAffiliationPrompted() bool {
	signupItem := application.getSignupItem("Affiliation")
	if signupItem == nil {
		return false
	}

	return signupItem.Prompted
}

func (application *Application) HasPromptPage() bool {
	providerItems := application.getAllPromptedProviderItems()
	if len(providerItems) != 0 {
		return true
	}

	return application.isAffiliationPrompted()
}
