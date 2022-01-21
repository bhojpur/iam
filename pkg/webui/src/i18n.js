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

import i18n from "i18next";
import zh from "./locales/zh/data.json";
import en from "./locales/en/data.json";
import fr from "./locales/fr/data.json";
import de from "./locales/de/data.json";
import ko from "./locales/ko/data.json";
import ru from "./locales/ru/data.json";
import ja from "./locales/ja/data.json";
import * as Conf from "./Conf";
import * as Setting from "./Setting";

const resources = {
  en: en,
  zh: zh,
  fr: fr,
  de: de,
  ko: ko,
  ru: ru,
  ja: ja,
};

function initLanguage() {
  let language = localStorage.getItem("language");
  if (language === undefined || language == null) {
    if (Conf.ForceLanguage !== "") {
      language = Conf.ForceLanguage;
    } else {
      let userLanguage;
      userLanguage = navigator.language;
      switch (userLanguage) {
        case "zh-CN":
          language = "zh";
          break;
        case "zh":
          language = "zh";
          break;
        case "en":
          language = "en";
          break;
        case "en-US":
          language = "en";
          break;
        case "fr":
          language = "fr";
          break;
        case "de":
          language = "de";
          break;
        case "ko":
          language = "ko";
          break;
        case "ru":
          language = "ru";
          break;
        case "ja":
          language = "ja";
          break;
        default:
          language = Conf.DefaultLanguage;
      }
    }
  }
  Setting.changeMomentLanguage(language);

  return language;
}

i18n.init({
  lng: initLanguage(),

  resources: resources,

  keySeparator: false,

  interpolation: {
    escapeValue: false,
  },
  //debug: true,
  saveMissing: true,
});

export default i18n;