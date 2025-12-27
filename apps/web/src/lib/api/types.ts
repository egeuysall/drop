// API Types and Interfaces

// User types
export interface User {
  id: string;
  email: string;
  name?: string;
  createdAt: string;
  updatedAt: string;
}

// Auth types
export interface AuthResponse {
  token: string;
  user: User;
}

// Pagination types
export interface Pagination {
  page: number;
  limit: number;
  total: number;
  totalPages: number;
}

// Generic paginated response
export interface PaginatedResponse<T> {
  data: T[];
  pagination: Pagination;
}

// Example API endpoint types
export interface ExampleData {
  id: string;
  name: string;
  description?: string;
}

// API endpoint paths
export const API_ENDPOINTS = {
  AUTH: {
    LOGIN: '/api/auth/login',
    REGISTER: '/api/auth/register',
    ME: '/api/auth/me',
  },
  USERS: {
    LIST: '/api/users',
    DETAIL: (id: string) => `/api/users/${id}`,
  },
  EXAMPLES: {
    LIST: '/api/examples',
    DETAIL: (id: string) => `/api/examples/${id}`,
  },
};