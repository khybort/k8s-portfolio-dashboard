import axios from 'axios';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';
const AUTH_URL = import.meta.env.VITE_AUTH_URL || 'http://localhost:8081';

export const apiClient = axios.create({
  baseURL: API_URL,
});

export const authClient = axios.create({
  baseURL: AUTH_URL,
});

// Request interceptor - Add token
apiClient.interceptors.request.use((config) => {
  const token = localStorage.getItem('access_token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Response interceptor - Handle token refresh
apiClient.interceptors.response.use(
  (response) => response,
  async (error) => {
    if (error.response?.status === 401) {
      const refreshToken = localStorage.getItem('refresh_token');
      if (refreshToken) {
        try {
          const response = await authClient.post('/api/v1/auth/refresh', {
            refresh_token: refreshToken,
          });
          localStorage.setItem('access_token', response.data.access_token);
          // Retry original request
          return apiClient.request(error.config);
        } catch (refreshError) {
          // Redirect to login
          localStorage.removeItem('access_token');
          localStorage.removeItem('refresh_token');
          window.location.href = '/admin/login';
        }
      }
    }
    return Promise.reject(error);
  }
);

