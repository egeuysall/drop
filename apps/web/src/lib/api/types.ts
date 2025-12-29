// Actual API response structure
export interface ItemsResponse {
  items: Item[];
  total: number;
}

// Item types - matches the actual API response structure
export interface Item {
  id: string;
  user_id: string;
  name: string;
  url: string;
  current_price: number;
  target_price?: number;
  in_stock: boolean;
  created_at: string;
  last_checked_at: string;
}

// Price history types
export interface PriceHistory {
  price: number;
  scraped_at: string;
}

// Price statistics types
export interface PriceStats {
  item_id: string;
  avg_price: number;
  min_price: number;
  max_price: number;
  history_count: number;
}

// Item creation/update DTO
export interface ItemCreateDTO {
  name?: string;
  url: string;
  target_price?: number;
}

// Item update DTO
export interface ItemUpdateDTO {
  name?: string;
  target_price?: number;
}

// Example API endpoint types
export interface ExampleData {
  id: string;
  name: string;
  description?: string;
}

// API endpoint paths
export const API_ENDPOINTS = {
  ITEMS: {
    LIST: '/v1/items',
    CREATE: '/v1/items',
    DETAIL: (id: string) => `/v1/items/${id}`,
    UPDATE: (id: string) => `/v1/items/${id}`,
    DELETE: (id: string) => `/v1/items/${id}`,
    UPDATE_PRICE: (id: string) => `/v1/items/${id}/price`,
    HISTORY: (id: string) => `/v1/items/${id}/history`,
    STATS: (id: string) => `/v1/items/${id}/stats`,
    CHECK: (id: string) => `/v1/items/${id}/check`,
  },
};
