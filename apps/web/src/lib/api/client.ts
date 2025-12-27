// API Client using TanStack Query
import { QueryClient } from '@tanstack/react-query';

// Create a client
export const queryClient = new QueryClient();

// Helper function to get auth token from storage (if needed)
export function getAuthToken(): string | null {
  if (typeof window !== 'undefined') {
    return localStorage.getItem('authToken');
  }
  return null;
}

// Base API URL - should be configured from environment
export const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080';

// API Error type
export type ApiError = {
  message: string;
  statusCode: number;
  error?: string;
  details?: any;
};

// Generic API response type
export type ApiResponse<T> = {
  data: T;
  error?: ApiError;
  success: boolean;
};

// Fetch wrapper with error handling
export async function apiFetch<T>(
  endpoint: string,
  options?: RequestInit
): Promise<ApiResponse<T>> {
  try {
    // Get auth token if available
    const token = getAuthToken();
    
    const response = await fetch(`${API_BASE_URL}${endpoint}`, {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        ...(token && { 'Authorization': `Bearer ${token}` }),
        ...options?.headers,
      },
    });

    const data = await response.json();

    if (!response.ok) {
      return {
        data: null as unknown as T,
        error: {
          message: data.message || 'Request failed',
          statusCode: response.status,
          error: data.error,
          details: data.details,
        },
        success: false,
      };
    }

    return {
      data,
      success: true,
    };
  } catch (error) {
    return {
      data: null as unknown as T,
      error: {
        message: error instanceof Error ? error.message : 'Unknown error occurred',
        statusCode: 500,
        details: error,
      },
      success: false,
    };
  }
}