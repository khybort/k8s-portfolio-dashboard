import { apiClient } from './client';

export interface Article {
  id: string;
  title: string;
  slug: string;
  excerpt: string;
  content: string;
  author_id: string;
  published: boolean;
  published_at?: string;
  created_at: string;
  updated_at: string;
}

export interface CreateArticleDto {
  title: string;
  slug: string;
  excerpt?: string;
  content: string;
  published: boolean;
}

export const articleService = {
  getArticles: async (page = 1, limit = 10): Promise<{ data: Article[]; total: number; pagination: any }> => {
    const response = await apiClient.get(`/api/v1/articles?page=${page}&limit=${limit}`);
    return response.data;
  },

  getArticle: async (id: string): Promise<Article> => {
    const response = await apiClient.get(`/api/v1/articles/${id}`);
    return response.data;
  },

  getArticleBySlug: async (slug: string): Promise<Article> => {
    const response = await apiClient.get(`/api/v1/articles/slug/${slug}`);
    return response.data;
  },

  createArticle: async (article: CreateArticleDto): Promise<Article> => {
    const response = await apiClient.post('/api/v1/admin/articles', article);
    return response.data;
  },

  updateArticle: async (id: string, article: Partial<CreateArticleDto>): Promise<Article> => {
    const response = await apiClient.put(`/api/v1/admin/articles/${id}`, article);
    return response.data;
  },

  deleteArticle: async (id: string): Promise<void> => {
    await apiClient.delete(`/api/v1/admin/articles/${id}`);
  },
};

