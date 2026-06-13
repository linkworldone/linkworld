import { useTranslation } from "react-i18next";

/**
 * 语言切换按钮：并排 EN | 中，高亮当前语言。
 * 仅靠 localStorage("lw_lang") 记忆，默认英文，不做浏览器语言探测。
 */
export function LanguageToggle() {
  const { i18n } = useTranslation();
  const current = i18n.language?.startsWith("zh") ? "zh" : "en";

  const setLang = (lang: "en" | "zh") => {
    if (lang === current) return;
    i18n.changeLanguage(lang);
    localStorage.setItem("lw_lang", lang);
  };

  return (
    <div className="flex items-center text-xs text-text-secondary border border-border rounded-lg overflow-hidden">
      <button
        type="button"
        onClick={() => setLang("en")}
        className={`px-2 py-1 transition-colors ${
          current === "en" ? "text-brand-gold font-semibold" : "hover:text-text-primary"
        }`}
      >
        EN
      </button>
      <span className="text-border">|</span>
      <button
        type="button"
        onClick={() => setLang("zh")}
        className={`px-2 py-1 transition-colors ${
          current === "zh" ? "text-brand-gold font-semibold" : "hover:text-text-primary"
        }`}
      >
        中
      </button>
    </div>
  );
}
