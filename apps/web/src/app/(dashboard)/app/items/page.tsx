'use client';

import { useState, useEffect } from 'react';
import { useItems, useDeleteItem } from '@/lib/api/hooks';
import { type Item } from '@/lib/api/types';
import { useItemsStore } from '@/lib/stores/items-store';
import { Button } from '@/components/ui/button';
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card';
import { Table, TableHeader, TableBody, TableHead, TableRow, TableCell } from '@/components/ui/table';
import { Plus, Trash2, RefreshCw, CheckCircle2, TrendingUp, TrendingDown, AlertTriangle as AlertTriangleIcon } from 'lucide-react';
import { toast } from 'sonner';
import { AddItemDialog } from '@/components/items/add-item-dialog';
import { PriceHistoryModal } from '@/components/items/price-history-modal';
import Link from 'next/link';

/**
 * Items Dashboard Page
 *
 * Features:
 * - Displays all tracked items in a responsive table
 * - Shows item details: Name, URL, Current Price, Target Price, Last Checked, In Stock status
 * - Price change indicators (up/down)
 * - Delete functionality with confirmation
 * - Add new items via modal dialog
 * - Real-time data fetching from API
 */
export default function ItemsDashboardPage() {
  const { data: itemsData, isLoading: isApiLoading, error: apiError, refetch } = useItems();
  const { mutate: deleteItem } = useDeleteItem();
  const [isAddDialogOpen, setIsAddDialogOpen] = useState(false);
  const [previousPrices, setPreviousPrices] = useState<Record<string, number | undefined>>({});

  // Use Zustand store as a fallback
  const {
    items: zustandItems,
    totalItems: zustandTotalItems,
    isLoading: isZustandLoading,
    error: zustandError,
    fetchItems: fetchZustandItems
  } = useItemsStore();

  const { removeItem: removeZustandItem } = useItemsStore();

  const handleDelete = (itemId: string) => {
    deleteItem(itemId, {
      onSuccess: () => {
        toast.success('Item deleted successfully');
        refetch();
        // Also update Zustand store
        removeZustandItem(itemId);
      },
      onError: (err) => {
        toast.error(`Failed to delete item: ${err.message}`);
      },
    });
  };

  const handleAddItemSuccess = () => {
    setIsAddDialogOpen(false);
    refetch();
    // Refresh Zustand store as well
    fetchZustandItems();
    toast.success('Item added successfully!');
  };

  // Determine which data source to use
  const apiItems = itemsData?.items || [];
  const apiTotalItems = itemsData?.total || 0;

  // Use API data if available, otherwise use Zustand data
  const items = apiItems.length > 0 ? apiItems : zustandItems;
  const totalItems = apiItems.length > 0 ? apiTotalItems : zustandTotalItems;
  const isLoading = isApiLoading || isZustandLoading;
  const error = apiError || zustandError;

  // Track price changes by comparing current prices with previous ones
  useEffect(() => {
    if (items.length > 0) {
      const newPreviousPrices: Record<string, number | undefined> = {};
      items.forEach((item) => {
        if (item.current_price !== null && item.current_price !== undefined) {
          // Only update if we have a previous price to compare with
          if (previousPrices[item.id] !== undefined && previousPrices[item.id] !== item.current_price) {
            newPreviousPrices[item.id] = previousPrices[item.id];
          } else {
            newPreviousPrices[item.id] = item.current_price;
          }
        }
      });
      setPreviousPrices(newPreviousPrices);
    }
  }, [items]);

  // Debug logging - moved here to maintain consistent hook order and ensure all variables are defined
  useEffect(() => {
    console.log('API Response:', itemsData);
    console.log('Zustand State:', { items: zustandItems, totalItems: zustandTotalItems });
    console.log('Final items to display:', items);
    console.log('Final total:', totalItems);
  }, [itemsData, zustandItems, totalItems]);

  if (isLoading) {
    return (
      <main className="flex min-h-screen items-center justify-center p-4">
        <Card className="w-full max-w-72">
          <CardHeader className="text-center">
            <CardTitle>Loading Items...</CardTitle>
          </CardHeader>
          <CardDescription className="text-center text-neutral-700 dark:text-neutral-300">
            Please wait while we fetch your tracked items
          </CardDescription>
        </Card>
      </main>
    );
  }

  if (error) {
    return (
      <main className="flex min-h-screen items-center justify-center p-4">
        <Card className="w-full max-w-72">
          <CardHeader className="text-center">
            <CardTitle className="flex items-center justify-center gap-2">
              <AlertTriangleIcon className="text-red-500" />
              <span>Error Loading Items</span>
            </CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-neutral-700 dark:text-neutral-300 mb-4">
              {error instanceof Error ? error.message : 'Unknown error'}
            </p>
            <Button
              onClick={() => {
                refetch();
                fetchZustandItems();
              }}
              variant="outline"
              className="w-full"
            >
              <RefreshCw className="mr-2 h-4 w-4" />
              Retry
            </Button>
          </CardContent>
        </Card>
      </main>
    );
  }

  return (
    <main className="container mx-auto py-8 px-4">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold">Your Tracked Items</h1>
        <Button
          onClick={() => setIsAddDialogOpen(true)}
          className="flex items-center gap-2"
        >
          <Plus className="h-4 w-4" />
          <span>Add Item</span>
        </Button>
      </div>

      {items.length === 0 ? (
        <Card className="mb-8">
          <CardHeader className="text-center">
            <CardTitle>No Items Found</CardTitle>
            <CardDescription>
              You haven't added any items to track yet.
            </CardDescription>
          </CardHeader>
          <CardContent className="text-center">
            <Button
              onClick={() => setIsAddDialogOpen(true)}
              className="flex items-center gap-2 mx-auto"
            >
              <Plus className="h-4 w-4" />
              <span>Add Your First Item</span>
            </Button>
          </CardContent>
        </Card>
      ) : (
        <Card>
          <CardHeader>
            <CardTitle>Items Dashboard</CardTitle>
            <CardDescription>
              View and manage all your price-tracked items
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="overflow-x-auto">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Name</TableHead>
                    <TableHead>URL</TableHead>
                    <TableHead>Current Price</TableHead>
                    <TableHead>Target Price</TableHead>
                    <TableHead>Last Checked</TableHead>
                    <TableHead>In Stock</TableHead>
                    <TableHead>Actions</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {items.map((item: Item) => (
                    <TableRow key={item.id}>
                      <TableCell className="font-medium">{item.name}</TableCell>
                      <TableCell className="truncate max-w-xs">
                        <Link
                          href={item.url}
                          target="_blank"
                          rel="noopener noreferrer"
                          className="text-blue-600 hover:underline truncate"
                        >
                          {item.url}
                        </Link>
                      </TableCell>
                      <TableCell>
                        <div className="flex items-center gap-2">
                          {item.current_price !== null && item.current_price !== undefined ? (
                            <>
                              <span>${item.current_price.toFixed(2)}</span>
                              {previousPrices[item.id] !== undefined && previousPrices[item.id] !== item.current_price && (
                                <span className={`text-xs flex items-center gap-1 ${
                                  item.current_price > previousPrices[item.id]! ? 'text-red-500' : 'text-green-500'
                                }`}>
                                  {item.current_price > previousPrices[item.id]! ? (
                                    <TrendingUp className="h-3 w-3" />
                                  ) : (
                                    <TrendingDown className="h-3 w-3" />
                                  )}
                                  {(((item.current_price - previousPrices[item.id]!) / previousPrices[item.id]!) * 100).toFixed(1)}%
                                </span>
                              )}
                            </>
                          ) : (
                            <span className="text-neutral-500">N/A</span>
                          )}
                        </div>
                      </TableCell>
                      <TableCell>
                        {item.target_price !== null && item.target_price !== undefined ? (
                          <span>${item.target_price.toFixed(2)}</span>
                        ) : (
                          <span className="text-neutral-500">Not set</span>
                        )}
                      </TableCell>
                      <TableCell>
                        {item.last_checked_at ? (
                          <span>{new Date(item.last_checked_at).toLocaleString()}</span>
                        ) : (
                          <span className="text-neutral-500">Never</span>
                        )}
                      </TableCell>
                      <TableCell>
                        {item.in_stock ? (
                          <div className="flex items-center gap-2">
                            <CheckCircle2 className="text-green-500" />
                            <span>In Stock</span>
                          </div>
                        ) : (
                          <div className="flex items-center gap-2">
                            <AlertTriangleIcon className="text-red-500" />
                            <span>Out of Stock</span>
                          </div>
                        )}
                      </TableCell>
                      <TableCell>
                        <div className="flex items-center gap-2">
                          <PriceHistoryModal
                            itemId={item.id}
                            itemName={item.name}
                            currentPrice={item.current_price || 0}
                            targetPrice={item.target_price}
                          />
                          <Button
                            onClick={() => handleDelete(item.id)}
                            variant="destructive"
                            className="flex items-center gap-1"
                          >
                            <Trash2 className="h-3 w-3" />
                            <span>Delete</span>
                          </Button>
                        </div>
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </div>
          </CardContent>
        </Card>
      )}

      <AddItemDialog
        open={isAddDialogOpen}
        onOpenChange={setIsAddDialogOpen}
        onSuccess={handleAddItemSuccess}
      />
    </main>
  );
}
