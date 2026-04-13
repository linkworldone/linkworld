import type { Config } from "tailwindcss";

const config: Config = {
  darkMode: "class",
  content: ["./index.html", "./src/**/*.{js,ts,jsx,tsx}"],
  theme: {
    extend: {
      fontFamily: {
        orbitron: ["Orbitron", "sans-serif"],
      },
      colors: {
        surface: {
          DEFAULT: "#0a0a14",
          card: "#0f0f1a",
          secondary: "#1a1a2e",
          gradient: {
            from: "#1a1a3e",
            to: "#0f1a2e",
          },
        },
        brand: {
          blue: "#3b82f6",
          purple: "#8b5cf6",
          cyan: "#06b6d4",
        },
        status: {
          success: "#22c55e",
          warning: "#f59e0b",
          danger: "#ef4444",
        },
        text: {
          primary: "#e0e0e0",
          secondary: "#888888",
          muted: "#555555",
        },
        border: {
          DEFAULT: "#1a1a2e",
        },
      },
      maxWidth: {
        mobile: "430px",
      },
    },
  },
  plugins: [],
};

export default config;
