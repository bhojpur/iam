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
import { Select } from "antd";

const { Option } = Select;

class SelectRegionBox extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            classes: props,
            value: "",
        };
    }

    onChange(e) {
        this.props.onChange(e);
        this.setState({value: e})
    };

    render() {
        return (
          <Select virtual={false}
                  showSearch
                  optionFilterProp="label"
                  style={{width: '100%'}}
                  defaultValue={this.props.defaultValue || undefined}
                  placeholder="Please select country/region"
                  onChange={(value => {this.onChange(value);})}
                  filterOption={(input, option) =>
                      option.label.indexOf(input) >= 0
                  }
          >
            {
                Setting.CountryRegionData.map((item, index) => (
                    <Option key={index} value={item.code} label={item.code} >
                        <img src={`${Setting.StaticBaseUrl}/image/flags/3x2/${item.code}.svg`} alt={item.name} height={20} style={{marginRight: 10}}/>
                        {`${item.name} (${item.code})`}
                    </Option>
                ))
            }
          </Select>
        )
    };
}

export default SelectRegionBox;