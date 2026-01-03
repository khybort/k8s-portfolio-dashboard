import { apiClient } from './client';

export interface Project {
  id: string;
  name: string;
  description: string;
  github_url: string;
  live_url?: string;
  technologies: string[];
  featured: boolean;
  created_at: string;
  updated_at: string;
}

export interface CreateProjectDto {
  name: string;
  description: string;
  github_url: string;
  live_url?: string;
  technologies: string[];
  featured?: boolean;
}

export const projectService = {
  getProjects: async (page: number = 1, limit: number = 10): Promise<{ data: Project[]; total: number }> => {
    const response = await apiClient.get(`/api/v1/projects?page=${page}&limit=${limit}`);
    return response.data;
  },

  getProject: async (id: string) => {
    const response = await apiClient.get(`/api/v1/projects/${id}`);
    return response.data;
  },

  createProject: async (project: CreateProjectDto) => {
    const response = await apiClient.post('/api/v1/admin/projects', project);
    return response.data;
  },

  updateProject: async (id: string, project: Partial<CreateProjectDto>) => {
    const response = await apiClient.put(`/api/v1/admin/projects/${id}`, project);
    return response.data;
  },

  deleteProject: async (id: string) => {
    const response = await apiClient.delete(`/api/v1/admin/projects/${id}`);
    return response.data;
  },
};

