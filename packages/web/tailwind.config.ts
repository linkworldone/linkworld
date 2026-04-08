import type { Config } from "tailwindcss";

const config: Config = {
  content: ["./index.html", "./src/**/*.{js,ts,jsx,tsx}"],
  theme: {
    extend: {
      fontFamily: {
        orbitron: ["Orbitron", "sans-serif"],
      },
      colors: {
        brand: {
          orange: "#f97316",
          purple: "#a855f7",
        },
      },
    },
  },
  plugins: [],
};

export default config;
