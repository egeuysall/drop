'use client';

import { useQuery, useMutation } from '@tanstack/react-query';
import type { UseQueryOptions, UseMutationOptions } from '@tanstack/react-query';
import { apiFetch } from './client';
import { createClient } from '../supabase/client';
import type { ItemsResponse, Item, ItemCreateDTO, ItemUpdateDTO, PriceHistory, PriceStats } from './types';
import { API_ENDPOINTS } from './types';

// Local User type for authentication
export interface User {
  id: string;
  email: string;
  name?: string;
  createdAt: string;
  updatedAt: string;
}

// Generic API hooks that can be used with any endpoint
// These hooks provide a foundation for API communication without requiring specific endpoints to exist

/**
 * Generic GET request hook
 * @template T - Response data type
 * @param endpoint - API endpoint URL
 * @param queryKey - Unique key for caching the query
 * @param options - Additional React Query options
 */
export function useGet<T>(endpoint: string, queryKey: string[], options?: UseQueryOptions<T>) {
  return useQuery<T, Error>({
    queryKey,
    queryFn: async () => {
      const response = await apiFetch<T>(endpoint, {
        method: 'GET',
      });

      if (!response.success || !response.data) {
        throw new Error(response.error?.message || 'Request failed');
      }

      return response.data;
    },
    ...options,
  });
}

/**
 * Generic POST request hook
 * @template T - Response data type
 * @template V - Request body type
 * @param endpoint - API endpoint URL
 * @param options - Additional React Query options
 */
export function usePost<T, V>(endpoint: string, options?: UseMutationOptions<T, Error, V>) {
  return useMutation<T, Error, V>({
    mutationFn: async data => {
      const response = await apiFetch<T>(endpoint, {
        method: 'POST',
        body: JSON.stringify(data),
      });

      if (!response.success || !response.data) {
        throw new Error(response.error?.message || 'Request failed');
      }

      return response.data;
    },
    ...options,
  });
}

/**
 * Generic PUT request hook
 * @template T - Response data type
 * @template V - Request body type
 * @param endpoint - API endpoint URL
 * @param options - Additional React Query options
 */
export function usePut<T, V>(endpoint: string, options?: UseMutationOptions<T, Error, V>) {
  return useMutation<T, Error, V>({
    mutationFn: async data => {
      const response = await apiFetch<T>(endpoint, {
        method: 'PUT',
        body: JSON.stringify(data),
      });

      if (!response.success || !response.data) {
        throw new Error(response.error?.message || 'Request failed');
      }

      return response.data;
    },
    ...options,
  });
}

/**
 * Generic DELETE request hook
 * @template T - Response data type
 * @param endpoint - API endpoint URL
 * @param options - Additional React Query options
 */
export function useDelete<T>(endpoint: string, options?: UseMutationOptions<T, Error, void>) {
  return useMutation<T, Error, void>({
    mutationFn: async () => {
      const response = await apiFetch<T>(endpoint, {
        method: 'DELETE',
      });

      if (!response.success || !response.data) {
        throw new Error(response.error?.message || 'Request failed');
      }

      return response.data;
    },
    ...options,
  });
}

// Health check hook - uses the existing /health endpoint
export function useHealthCheck() {
  return useQuery<string, Error>({
    queryKey: ['health'],
    queryFn: async () => {
      const response = await apiFetch<string>('/health', {
        method: 'GET',
      });

      if (!response.success || !response.data) {
        throw new Error(response.error?.message || 'Health check failed');
      }

      return response.data;
    },
    retry: 3,
    retryDelay: 1000,
  });
}

// GitHub OAuth hook using Supabase
export function useGitHubAuth() {
  const loginWithGitHub = async () => {
    const supabase = createClient();

    const { data, error } = await supabase.auth.signInWithOAuth({
      provider: 'github',
      options: {
        redirectTo: `${window.location.origin}/callback`,
      },
    });

    if (error) {
      throw new Error(error.message);
    }

    return data;
  };

  return { loginWithGitHub };
}

// User session hook
export function useUserSession() {
  return useQuery<User | null, Error>({
    queryKey: ['user-session'],
    queryFn: async () => {
      const supabase = createClient();
      const {
        data: { user },
        error,
      } = await supabase.auth.getUser();

      if (error) {
        throw new Error(error.message);
      }

      if (!user) {
        return null;
      }

      // Map Supabase user to our User type
      return {
        id: user.id,
        email: user.email || '',
        name: user.user_metadata?.full_name || user.user_metadata?.name || '',
        createdAt: user.created_at
          ? new Date(user.created_at).toISOString()
          : new Date().toISOString(),
        updatedAt: user.updated_at
          ? new Date(user.updated_at).toISOString()
          : new Date().toISOString(),
      };
    },
  });
}

// Sign out hook
export function useSignOut() {
  const signOut = async () => {
    const supabase = createClient();
    const { error } = await supabase.auth.signOut();

    if (error) {
      throw new Error(error.message);
    }
  };

  return useMutation<void, Error, void>({
    mutationFn: signOut,
  });
}

// Item-specific hooks
export function useItems() {
  return useQuery<ItemsResponse, Error>({
    queryKey: ['items'],
    queryFn: async () => {
      const response = await apiFetch<ItemsResponse>(API_ENDPOINTS.ITEMS.LIST, {
        method: 'GET',
      });

      if (!response.success || !response.data) {
        throw new Error(response.error?.message || 'Failed to fetch items');
      }

      return response.data;
    },
  });
}

export function useItem(itemId: string) {
  return useQuery<Item, Error>({
    queryKey: ['items', itemId],
    queryFn: async () => {
      const response = await apiFetch<Item>(API_ENDPOINTS.ITEMS.DETAIL(itemId), {
        method: 'GET',
      });

      if (!response.success || !response.data) {
        throw new Error(response.error?.message || 'Failed to fetch item');
      }

      return response.data;
    },
    enabled: !!itemId,
  });
}

export function useCreateItem() {
  return useMutation<Item, Error, ItemCreateDTO>({
    mutationFn: async (itemData) => {
      const response = await apiFetch<Item>(API_ENDPOINTS.ITEMS.CREATE, {
        method: 'POST',
        body: JSON.stringify(itemData),
      });

      if (!response.success || !response.data) {
        throw new Error(response.error?.message || 'Failed to create item');
      }

      return response.data;
    },
  });
}

export function useUpdateItem() {
  return useMutation<Item, Error, { id: string; data: ItemUpdateDTO }>({
    mutationFn: async ({ id, data }) => {
      const response = await apiFetch<Item>(API_ENDPOINTS.ITEMS.UPDATE(id), {
        method: 'PUT',
        body: JSON.stringify(data),
      });

      if (!response.success || !response.data) {
        throw new Error(response.error?.message || 'Failed to update item');
      }

      return response.data;
    },
  });
}

export function useDeleteItem() {
  return useMutation<void, Error, string>({
    mutationFn: async (itemId) => {
      const response = await apiFetch<void>(API_ENDPOINTS.ITEMS.DELETE(itemId), {
        method: 'DELETE',
      });

      if (!response.success) {
        throw new Error(response.error?.message || 'Failed to delete item');
      }
    },
  });
}

export function useItemPriceHistory(itemId: string) {
  return useQuery<PriceHistory[], Error>({
    queryKey: ['items', itemId, 'history'],
    queryFn: async () => {
      const response = await apiFetch<PriceHistory[]>(API_ENDPOINTS.ITEMS.HISTORY(itemId), {
        method: 'GET',
      });

      if (!response.success || !response.data) {
        throw new Error(response.error?.message || 'Failed to fetch price history');
      }

      return response.data;
    },
    enabled: !!itemId,
  });
}

export function useItemStats(itemId: string) {
  return useQuery<PriceStats, Error>({
    queryKey: ['items', itemId, 'stats'],
    queryFn: async () => {
      const response = await apiFetch<PriceStats>(API_ENDPOINTS.ITEMS.STATS(itemId), {
        method: 'GET',
      });

      if (!response.success || !response.data) {
        throw new Error(response.error?.message || 'Failed to fetch item stats');
      }

      return response.data;
    },
    enabled: !!itemId,
  });
}

export function useCheckPriceDrop() {
  return useMutation<void, Error, string>({
    mutationFn: async (itemId) => {
      const response = await apiFetch<void>(API_ENDPOINTS.ITEMS.CHECK(itemId), {
        method: 'GET',
      });

      if (!response.success) {
        throw new Error(response.error?.message || 'Failed to check price drop');
      }
    },
  });
}
