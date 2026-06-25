import axios from "axios";

/**
 * Pre-configured Axios instance for the Jewellery Billing API.
 *
 * - Base URL is proxied by Vite in dev (/api → localhost:8080/api)
 * - Automatically attaches the JWT Bearer token from localStorage
 * - Intercepts 401 responses and redirects to /login
 */
const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || "/api",
  headers: {
    "Content-Type": "application/json",
  },
});

// ── Request interceptor: attach JWT ────────────────────────────────
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem("access_token");
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// ── Response interceptor: handle 401 ───────────────────────────────
api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;

    // If we get a 401 and haven't tried refreshing yet, attempt a refresh.
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;

      const refreshToken = localStorage.getItem("refresh_token");
      if (refreshToken) {
        try {
          const { data } = await axios.post(`${import.meta.env.VITE_API_URL || "/api"}/auth/refresh`, {
            refresh_token: refreshToken,
          });

          const newAccessToken = data.data.access_token;
          const newRefreshToken = data.data.refresh_token;

          localStorage.setItem("access_token", newAccessToken);
          localStorage.setItem("refresh_token", newRefreshToken);

          originalRequest.headers.Authorization = `Bearer ${newAccessToken}`;
          return api(originalRequest);
        } catch {
          // Refresh failed — clear tokens and redirect to login.
          localStorage.removeItem("access_token");
          localStorage.removeItem("refresh_token");
          localStorage.removeItem("user");
          localStorage.removeItem("organization");
          // Use hash-based path so HashRouter handles it correctly
          window.location.href = "/#/login";
        }
      } else {
        // No refresh token — redirect to login.
        localStorage.removeItem("access_token");
        localStorage.removeItem("user");
        localStorage.removeItem("organization");
        // BrowserRouter handles /login correctly via vercel.json rewrite
        window.location.href = "/#/login";
      }
    }

    return Promise.reject(error);
  }
);

export default api;
