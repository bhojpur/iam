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

import {createButton} from "react-social-login-buttons";

function Icon({ width = 24, height = 24, color }) {
  return <svg xmlns="http://www.w3.org/2000/svg" height="48" width="32" viewBox="-18.15 -35.9725 157.3 215.835"><path fill="#faab07" d="M60.503 142.237c-12.533 0-24.038-4.195-31.445-10.46-3.762 1.124-8.574 2.932-11.61 5.175-2.6 1.918-2.275 3.874-1.807 4.663 2.056 3.47 35.273 2.216 44.862 1.136zm0 0c12.535 0 24.039-4.195 31.447-10.46 3.76 1.124 8.573 2.932 11.61 5.175 2.598 1.918 2.274 3.874 1.805 4.663-2.056 3.47-35.272 2.216-44.862 1.136zm0 0"/><path d="M60.576 67.119c20.698-.14 37.286-4.147 42.907-5.683 1.34-.367 2.056-1.024 2.056-1.024.005-.189.085-3.37.085-5.01C105.624 27.768 92.58.001 60.5 0 28.42.001 15.375 27.769 15.375 55.401c0 1.642.08 4.822.086 5.01 0 0 .583.615 1.65.913 5.19 1.444 22.09 5.65 43.312 5.795zm56.245 23.02c-1.283-4.129-3.034-8.944-4.808-13.568 0 0-1.02-.126-1.537.023-15.913 4.623-35.202 7.57-49.9 7.392h-.153c-14.616.175-33.774-2.737-49.634-7.315-.606-.175-1.802-.1-1.802-.1-1.774 4.624-3.525 9.44-4.808 13.568-6.119 19.69-4.136 27.838-2.627 28.02 3.239.392 12.606-14.821 12.606-14.821 0 15.459 13.957 39.195 45.918 39.413h.848c31.96-.218 45.917-23.954 45.917-39.413 0 0 9.368 15.213 12.607 14.822 1.508-.183 3.491-8.332-2.627-28.021"/><path fill="#fff" d="M49.085 40.824c-4.352.197-8.07-4.76-8.304-11.063-.236-6.305 3.098-11.576 7.45-11.773 4.347-.195 8.064 4.76 8.3 11.065.238 6.306-3.097 11.577-7.446 11.771m31.133-11.063c-.233 6.302-3.951 11.26-8.303 11.063-4.35-.195-7.684-5.465-7.446-11.77.236-6.305 3.952-11.26 8.3-11.066 4.352.197 7.686 5.468 7.449 11.773"/><path fill="#faab07" d="M87.952 49.725C86.79 47.15 75.077 44.28 60.578 44.28h-.156c-14.5 0-26.212 2.87-27.375 5.446a.863.863 0 00-.085.367.88.88 0 00.16.496c.98 1.427 13.985 8.487 27.3 8.487h.156c13.314 0 26.319-7.058 27.299-8.487a.873.873 0 00.16-.498.856.856 0 00-.085-.365"/><path d="M54.434 29.854c.199 2.49-1.167 4.702-3.046 4.943-1.883.242-3.568-1.58-3.768-4.07-.197-2.492 1.167-4.704 3.043-4.944 1.886-.244 3.574 1.58 3.771 4.07m11.956.833c.385-.689 3.004-4.312 8.427-2.993 1.425.347 2.084.857 2.223 1.057.205.296.262.718.053 1.286-.412 1.126-1.263 1.095-1.734.875-.305-.142-4.082-2.66-7.562 1.097-.24.257-.668.346-1.073.04-.407-.308-.574-.93-.334-1.362"/><path fill="#fff" d="M60.576 83.08h-.153c-9.996.12-22.116-1.204-33.854-3.518-1.004 5.818-1.61 13.132-1.09 21.853 1.316 22.043 14.407 35.9 34.614 36.1h.82c20.208-.2 33.298-14.057 34.616-36.1.52-8.723-.087-16.035-1.092-21.854-11.739 2.315-23.862 3.64-33.86 3.518"/><path fill="#eb1923" d="M32.102 81.235v21.693s9.937 2.004 19.893.616V83.535c-6.307-.357-13.109-1.152-19.893-2.3"/><path fill="#eb1923" d="M105.539 60.412s-19.33 6.102-44.963 6.275h-.153c-25.591-.172-44.896-6.255-44.962-6.275L8.987 76.57c16.193 4.882 36.261 8.028 51.436 7.845h.153c15.175.183 35.242-2.963 51.437-7.845zm0 0"/></svg>;
}

const config = {
  text: "Sign in with QQ",
  icon: Icon,
  iconFormat: name => `fa fa-${name}`,
  style: {background: "rgb(94,188,249)"},
  activeStyle: {background: "rgb(76,143,208)"},
};

const QqLoginButton = createButton(config);

export default QqLoginButton;