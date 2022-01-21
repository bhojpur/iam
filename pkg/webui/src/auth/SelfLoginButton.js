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
import {createButton} from "react-social-login-buttons";

class SelfLoginButton extends React.Component {
    generateIcon() {
        const avatar = this.props.account.avatar;
        return () => {
            return <img width={36} height={36} src={avatar} alt="Sign in with Google"/>;
        };
    }

    render() {
        const config = {
            icon: this.generateIcon(),
            iconFormat: name => `fa fa-${name}`,
            style: {background: "#ffffff", color: "#000000"},
            activeStyle: {background: "#eff0ee"},
        };

        const SelfLoginButton = createButton(config);
        return <SelfLoginButton text={`${this.props.account.name} (${this.props.account.displayName})`} onClick={() => this.props.onClick()} align={"center"} />
    }
}

export default SelfLoginButton;