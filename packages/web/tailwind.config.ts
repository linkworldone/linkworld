import type { Config } from "tailwindcss";

/*
 * 色值单一出口：所有 colors 引用 src/index.css :root 的 CSS 变量，禁止裸 HEX。
 * 换主题只改 :root。业务语义类（surface/brand/text/status）与 shadcn 原子类
 * （primary/card/muted/...）同喂一套变量。
 */
const config: Config = {
  darkMode: "class",
  content: ["./index.html", "./src/**/*.{js,ts,jsx,tsx}"],
  theme: {
    extend: {
      fontFamily: {
        sans: ["var(--font-sans)"],
        display: ["var(--font-display)"],
        data: ["var(--font-data)"],
      },
      colors: {
        // ── 业务语义 token（layout/shared/wallet/pages 消费） ──
        surface: {
          DEFAULT: "var(--surface-canvas)",
          card: "var(--surface-card)",
          "card-elevated": "var(--surface-card-elevated)",
          input: "var(--surface-input)",
          "card-line": "var(--surface-card-line)",
        },
        brand: {
          navy: "var(--brand-navy)",
          "navy-deep": "var(--brand-navy-deep)",
          royal: "var(--brand-royal)",
          "royal-bright": "var(--brand-royal-bright)",
          gold: "var(--brand-gold)",
          "gold-hover": "var(--brand-gold-hover)",
          "gold-press": "var(--brand-gold-press)",
          champagne: "var(--brand-champagne)",
        },
        status: {
          success: "var(--status-success)",
          warning: "var(--status-warning)",
          danger: "var(--status-danger)",
          info: "var(--status-info)",
        },
        text: {
          // 默认（body / 深底）三级
          primary: "var(--text-primary)",
          secondary: "var(--text-secondary)",
          muted: "var(--text-muted)",
          // 浅底（米白卡内）
          "on-light-primary": "var(--text-on-light-primary)",
          "on-light-secondary": "var(--text-on-light-secondary)",
          "on-light-muted": "var(--text-on-light-muted)",
          // 深底（navy 画布上）
          "on-dark-primary": "var(--text-on-dark-primary)",
          "on-dark-secondary": "var(--text-on-dark-secondary)",
          "on-dark-muted": "var(--text-on-dark-muted)",
          "on-dark-gold": "var(--text-on-dark-gold)",
        },
        // ── shadcn 语义 token（ui/* 原子组件消费，指向同一真源） ──
        background: "var(--background)",
        foreground: "var(--foreground)",
        card: {
          DEFAULT: "var(--card)",
          foreground: "var(--card-foreground)",
        },
        popover: {
          DEFAULT: "var(--popover)",
          foreground: "var(--popover-foreground)",
        },
        primary: {
          DEFAULT: "var(--primary)",
          foreground: "var(--primary-foreground)",
        },
        secondary: {
          DEFAULT: "var(--secondary)",
          foreground: "var(--secondary-foreground)",
        },
        muted: {
          DEFAULT: "var(--muted)",
          foreground: "var(--muted-foreground)",
        },
        accent: {
          DEFAULT: "var(--accent)",
          foreground: "var(--accent-foreground)",
        },
        destructive: "var(--destructive)",
        border: "var(--border)",
        input: "var(--input)",
        ring: "var(--ring)",
      },
      backgroundImage: {
        "gradient-hero": "var(--gradient-hero)",
        "canvas": "var(--bg-canvas)",
        "gold-line": "var(--gradient-gold-line)",
      },
      borderRadius: {
        lg: "var(--radius)",
        md: "calc(var(--radius) - 0.25rem)",
        sm: "calc(var(--radius) - 0.5rem)",
      },
      boxShadow: {
        card: "0 4px 16px rgba(10, 27, 51, 0.25)",
      },
      maxWidth: {
        mobile: "430px",
      },
    },
  },
  plugins: [],
};

export default config;
