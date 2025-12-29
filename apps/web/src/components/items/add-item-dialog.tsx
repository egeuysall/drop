'use client';

import { useState } from 'react';
import { useCreateItem } from '@/lib/api/hooks';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from '@/components/ui/dialog';
import { Plus, Loader2, AlertTriangle } from 'lucide-react';

/**
 * Add Item Dialog Component
 *
 * Features:
 * - Modal dialog for adding new items
 * - Form fields: URL (required), Name (optional), Target Price (optional)
 * - URL validation
 * - Loading state during API calls
 * - Success/error feedback
 * - Form validation and error handling
 */
interface AddItemDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSuccess: () => void;
}

export function AddItemDialog({ open, onOpenChange, onSuccess }: AddItemDialogProps) {
  const [url, setUrl] = useState('');
  const [name, setName] = useState('');
  const [targetPrice, setTargetPrice] = useState('');
  const [error, setError] = useState<string | null>(null);
  const { mutate: createItem, isPending } = useCreateItem();

  const validateUrl = (url: string): boolean => {
    try {
      new URL(url);
      return true;
    } catch {
      return false;
    }
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    // Reset error
    setError(null);

    // Validate URL
    if (!url.trim()) {
      setError('URL is required');
      return;
    }

    if (!validateUrl(url)) {
      setError('Please enter a valid URL (include http:// or https://)');
      return;
    }

    // Prepare item data
    const itemData = {
      url: url.trim(),
      name: name.trim() || undefined,
      target_price: targetPrice ? parseFloat(targetPrice) : undefined,
    };

    // Create item
    createItem(itemData, {
      onSuccess: () => {
        // Reset form
        setUrl('');
        setName('');
        setTargetPrice('');
        onSuccess();
      },
      onError: (err) => {
        setError(err.message || 'Failed to add item');
      },
    });
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-106 bg-neutral-100 dark:bg-neutral-900 border border-neutral-300 dark:border-neutral-700 rounded-md p-6">
        <DialogHeader>
          <DialogTitle className="text-lg font-semibold">Add New Item</DialogTitle>
          <DialogDescription className="text-sm text-neutral-700 dark:text-neutral-300 mt-2">
            Enter the product URL and optional details to start tracking prices.
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit} className="space-y-4 py-4">
          <div className="space-y-2">
            <Label htmlFor="url" className="text-sm font-medium">
              Product URL <span className="text-red-500">*</span>
            </Label>
            <Input
              id="url"
              type="url"
              placeholder="https://example.com/product"
              value={url}
              onChange={(e) => setUrl(e.target.value)}
              required
              className="w-full"
            />
            <p className="text-xs text-neutral-500 dark:text-neutral-400">
              Enter the full URL of the product you want to track
            </p>
          </div>

          <div className="space-y-2">
            <Label htmlFor="name" className="text-sm font-medium">
              Product Name (Optional)
            </Label>
            <Input
              id="name"
              type="text"
              placeholder="Product name"
              value={name}
              onChange={(e) => setName(e.target.value)}
              className="w-full"
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="targetPrice" className="text-sm font-medium">
              Target Price (Optional)
            </Label>
            <Input
              id="targetPrice"
              type="number"
              placeholder="e.g., 99.99"
              value={targetPrice}
              onChange={(e) => setTargetPrice(e.target.value)}
              min="0"
              step="0.01"
              className="w-full"
            />
            <p className="text-xs text-neutral-500 dark:text-neutral-400">
              Get notified when price drops to this amount
            </p>
          </div>

          {error && (
            <div className="flex items-center gap-2 p-2 bg-red-100 dark:bg-red-900 border border-red-300 dark:border-red-700 rounded-sm">
              <AlertTriangle className="h-4 w-4 text-red-500" />
              <span className="text-sm text-red-600 dark:text-red-300">{error}</span>
            </div>
          )}

          <DialogFooter className="flex justify-end gap-2 pt-4">
            <Button
              type="button"
              variant="outline"
              onClick={() => onOpenChange(false)}
              disabled={isPending}
            >
              Cancel
            </Button>
            <Button
              type="submit"
              disabled={isPending}
              className="flex items-center gap-2"
            >
              {isPending ? (
                <>
                  <Loader2 className="h-4 w-4 animate-spin" />
                  <span>Adding...</span>
                </>
              ) : (
                <>
                  <Plus className="h-4 w-4" />
                  <span>Add Item</span>
                </>
              )}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
