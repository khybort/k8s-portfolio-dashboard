import { authClient } from './client';

export interface LoginDto {
  email: string;
  password: string;
}

export interface RegisterDto {
  email: string;
  password: string;
  name: string;
}

export interface AuthResponse {
  access_token: string;
  refresh_token: string;
  token_type: string;
  expires_in: number;
  user: {
    id: string;
    email: string;
    name: string;
    role: string;
  };
}

export const authService = {
  login: async (credentials: LoginDto): Promise<AuthResponse> => {
    const response = await authClient.post('/api/v1/auth/login', credentials);
    return response.data;
  },

  register: async (data: RegisterDto): Promise<any> => {
    const response = await authClient.post('/api/v1/auth/register', data);
    return response.data;
  },

  refresh: async (refreshToken: string): Promise<{ access_token: string }> => {
    const response = await authClient.post('/api/v1/auth/refresh', {
      refresh_token: refreshToken,
    });
    return response.data;
  },

  verify: async (token: string): Promise<any> => {
    const response = await authClient.post('/api/v1/auth/verify', {
      token,
    });
    return response.data;
  },
};

