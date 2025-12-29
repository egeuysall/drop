import { create } from 'zustand';
import { type Item } from '../api/types';

interface ItemsState {
  items: Item[];
  totalItems: number;
  isLoading: boolean;
  error: Error | null;
  fetchItems: () => Promise<void>;
  addItem: (item: Item) => void;
  removeItem: (itemId: string) => void;
}

export const useItemsStore = create<ItemsState>((set) => ({
  items: [],
  totalItems: 0,
  isLoading: false,
  error: null,

  fetchItems: async () => {
    try {
      set({ isLoading: true, error: null });

      const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080';
      
      // Get auth token from Supabase
      let authToken = null;
      try {
        const { createClient } = await import('../supabase/client');
        const supabase = createClient();
        const { data: { session } } = await supabase.auth.getSession();
        if (session) {
          authToken = session.access_token;
        }
      } catch (err) {
        console.error('Error getting auth token:', err);
      }

      const response = await fetch(`${API_BASE_URL}/v1/items`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
          ...(authToken && { 'Authorization': `Bearer ${authToken}` }),
        },
        credentials: 'include',
      });

      if (!response.ok) {
        throw new Error(`Failed to fetch items: ${response.status} ${response.statusText}`);
      }

      const data = await response.json();
      
      // Handle the actual API response structure: {data: {items: [], total: number}}
      const apiData = data.data;
      const items = apiData?.items || [];
      const totalItems = apiData?.total || 0;

      set({ items, totalItems, isLoading: false });
      
      console.log('Fetched items via Zustand:', { items, totalItems });
      
    } catch (error) {
      console.error('Error fetching items:', error);
      set({ 
        error: error instanceof Error ? error : new Error('Unknown error'),
        isLoading: false 
      });
    }
  },

  addItem: (item) => set((state) => ({
    items: [...state.items, item],
    totalItems: state.totalItems + 1
  })),

  removeItem: (itemId) => set((state) => ({
    items: state.items.filter(item => item.id !== itemId),
    totalItems: state.totalItems - 1
  })),
}));

// Initialize the store by fetching items
useItemsStore.getState().fetchItems();