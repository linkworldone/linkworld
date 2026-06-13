// i18n 初始化（react-i18next）。
// 默认英文；用户手动切换后记忆到 localStorage(key=lw_lang)，无浏览器语言探测。
// resources 采用单一 translation namespace，key 按页面/组件域分组（common / landing / ...）。
// 在 main.tsx 顶部 `import "./i18n"` 以在 App 渲染前完成初始化。
// JSON 为同步 import，故关闭 Suspense 避免加载闪烁。
import i18n from "i18next";
import { initReactI18next } from "react-i18next";
import en from "./locales/en.json";
import zh from "./locales/zh.json";

const STORAGE_KEY = "lw_lang";
const stored = localStorage.getItem(STORAGE_KEY);
const initialLang = stored === "en" || stored === "zh" ? stored : "en";

i18n.use(initReactI18next).init({
  resources: {
    en: { translation: en },
    zh: { translation: zh },
  },
  lng: initialLang,
  fallbackLng: "en",
  interpolation: {
    escapeValue: false,
  },
  react: {
    useSuspense: false,
  },
});

export default i18n;
