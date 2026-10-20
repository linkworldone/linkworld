// Empty default keeps production API requests same-origin; Vite proxies /api during local development.
export const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || "";
