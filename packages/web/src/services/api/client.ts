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
    const error = new Error(message) as Error & { status?: number };
    // 附带 HTTP 状态码（如 429 限流 / 400 校验失败），调用方据此区分提示文案；
    // 纯附加属性，不改变既有 message 行为，兼容所有现有 catch(e){e.message}。
    if (err.response?.status) error.status = err.response.status;
    return Promise.reject(error);
  }
);
