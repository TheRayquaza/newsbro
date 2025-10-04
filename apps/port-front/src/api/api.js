import { ENV } from "../env";

class ApiService {
  constructor() {
    this.cache = new Map();
    this.cacheTimeout = 5 * 60 * 1000;
  }

  getAuthHeaders() {
    const token = localStorage.getItem("auth_token");
    return {
      "Content-Type": "application/json",
      ...(token && { Authorization: `Bearer ${token}` }),
    };
  }

  async request(endpoint, options = {}) {
    const url = `${ENV.API_BASE_URL}${endpoint}`;
    const cacheKey = `${options.method || "GET"}_${url}_${JSON.stringify(options.body || {})}`;

    // Check cache for GET requests
    if (!options.method || options.method === "GET") {
      const cached = this.cache.get(cacheKey);
      if (cached && Date.now() - cached.timestamp < this.cacheTimeout) {
        return cached.data;
      }
    }

    const response = await fetch(url, {
      ...options,
      headers: {
        ...this.getAuthHeaders(),
        ...options.headers,
      },
    });

    const data = await response.json();

    if (!response.ok) {
      throw new Error(data.message || data.error || "Request failed");
    }

    // Cache successful GET requests
    if (!options.method || options.method === "GET") {
      this.cache.set(cacheKey, { data, timestamp: Date.now() });
    }

    return data;
  }

  clearCache() {
    this.cache.clear();
  }

  async login(email, password) {
    return this.request("/api/v1/auth/login", {
      method: "POST",
      body: JSON.stringify({ email, password }),
    });
  }

  async register(userData) {
    return this.request("/api/v1/auth/register", {
      method: "POST",
      body: JSON.stringify(userData),
    });
  }

  async logout() {
    return this.request("/api/v1/auth/logout", { method: "POST" });
  }

  forgeOAuthUrl() {
    return `${ENV.API_BASE_URL}/api/v1/auth/oauth/login`;
  }

  // Add more API methods as needed
}

const api = new ApiService();
export default api;
