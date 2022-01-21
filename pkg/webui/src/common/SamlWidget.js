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
import {Col, Row} from "antd";
import * as Setting from "../Setting";

class SamlWidget extends React.Component {
	constructor(props) {
		super(props);
		this.state = {
			classes: props,
			addressOptions: [],
			affiliationOptions: [],
		};
	}

	renderIdp(user, application, providerItem) {
		const provider = providerItem.provider;
		const name = user.name;

		return (
			<Row key={provider.name} style={{marginTop: '20px'}}>
				<Col style={{marginTop: '5px'}} span={this.props.labelSpan}>
					{
						Setting.getProviderLogo(provider)
					}
					<span style={{marginLeft: '5px'}}>
					{
						`${provider.type}:`
					}
					</span>
				</Col>
				<Col span={24 - this.props.labelSpan} style={{marginTop: '5px'}}>
					<span style={{
						width: this.props.labelSpan === 3 ? '300px' : '130px',
						display: (Setting.isMobile()) ? 'inline' : "inline-block"}}>{name}</span>
				</Col>
			</Row>
		)
	}

	render() {
		return this.renderIdp(this.props.user, this.props.application, this.props.providerItem)
	}
}

export default SamlWidget;