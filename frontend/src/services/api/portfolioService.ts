import { apiClient } from './client';

export interface Portfolio {
  id: string;
  name: string;
  title: string;
  bio: string;
  email: string;
  social_links?: {
    github?: string;
    linkedin?: string;
    phone?: string;
  };
  settings?: {
    theme?: string;
    language?: string;
  };
  created_at: string;
  updated_at: string;
}

export interface UpdatePortfolioDto {
  name?: string;
  title?: string;
  bio?: string;
  email?: string;
  social_links?: {
    github?: string;
    linkedin?: string;
    phone?: string;
  };
  settings?: {
    theme?: string;
    language?: string;
  };
}

export const portfolioService = {
  getPortfolio: async () => {
    const response = await apiClient.get('/api/v1/portfolio');
    return response.data;
  },

  updatePortfolio: async (portfolio: UpdatePortfolioDto) => {
    const response = await apiClient.put('/api/v1/admin/portfolio', portfolio);
    return response.data;
  },
};

