import axios from "axios";
import { API_BASE_URL } from "../../config/api";

export const apiClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
  headers: { "Content-Type": "application/json" },
});

// 响应拦截器：直接返回 data
apiClient.interceptors.response.use(
  (res) => res.data,
  (err) => {
    const message = err.response?.data?.error || err.message || "Network error";
    return Promise.reject(new Error(message));
  }
);
