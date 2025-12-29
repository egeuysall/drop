'use client';

import { useState } from 'react';
import { useItemPriceHistory } from '@/lib/api/hooks';
import { Button } from '@/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog';
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from 'recharts';
import { Loader2, TrendingUp, TrendingDown, AlertTriangle } from 'lucide-react';


/**
 * Price History Modal Component
 *
 * Features:
 * - Shows price history chart for an item
 * - Loading states
 * - Responsive chart using Recharts
 * - Price trend analysis
 */
interface PriceHistoryModalProps {
  itemId: string;
  itemName: string;
  currentPrice: number;
  targetPrice?: number;
}

export function PriceHistoryModal({
  itemId,
  itemName,
  currentPrice,
  targetPrice,
}: PriceHistoryModalProps) {
  const [isOpen, setIsOpen] = useState(false);

  const { data: historyData, isLoading, error } = useItemPriceHistory(itemId);

  // Debug: Check the structure of historyData
  console.log('History data:', historyData);

  // Handle different response structures
  let historyArray: any[] = [];

  // Check if historyData is a direct array
  if (Array.isArray(historyData)) {
    historyArray = historyData;
  }
  // Check if historyData has a data property that's an array
  else if (historyData && typeof historyData === 'object' && 'data' in historyData && Array.isArray((historyData as {data?: any}).data)) {
    historyArray = (historyData as {data: any[]}).data;
  }
  // Log warning for unexpected structures
  else if (historyData !== undefined && historyData !== null) {
    console.warn('Unexpected history data structure:', historyData);
  }

  // Filter data based on time range
  const filteredData = historyArray.length > 0 ? historyArray
    .filter(entry => entry && entry.price !== undefined && entry.scraped_at)
    .map((entry, index) => ({
      date: new Date(entry.scraped_at).toLocaleDateString(),
      price: typeof entry.price === 'number' ? entry.price : parseFloat(entry.price),
      index: index + 1,
    })) : [];

  // Calculate price trends
  const firstPrice = filteredData.length > 0 ? filteredData[0]?.price : currentPrice;
  const lastPrice = filteredData.length > 0 ? filteredData[filteredData.length - 1]?.price : currentPrice;
  const priceChange = lastPrice - firstPrice;
  const priceChangePercent = firstPrice > 0 ? (priceChange / firstPrice) * 100 : 0;

  // Debug: Log calculated values
  console.log('Price trends:', { firstPrice, lastPrice, priceChange, priceChangePercent });

  // Determine if price is at or below target
  const isAtOrBelowTarget = targetPrice ? currentPrice <= targetPrice : false;

  return (
    <Dialog open={isOpen} onOpenChange={setIsOpen}>
      <DialogTrigger asChild>
        <Button variant="outline" className="h-8 gap-1 flex items-center">
          <TrendingUp className="h-3.5 w-3.5" />
          <span className="sr-only sm:not-sr-only sm:whitespace-nowrap">
            Price History
          </span>
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-200 max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            Price History: {itemName}
            {isAtOrBelowTarget && (
              <span className="text-green-500 text-sm font-medium">
                ✓ At or below target
              </span>
            )}
          </DialogTitle>
        </DialogHeader>

        {/* Price Summary */}
        <div className="flex flex-wrap gap-4 mb-4 p-4 rounded-lg">
          <div className="flex items-center gap-2">
            <span className="text-sm text-muted-foreground">Current Price:</span>
            <span className="font-semibold">${currentPrice.toFixed(2)}</span>
          </div>
          {targetPrice && (
            <div className="flex items-center gap-2">
              <span className="text-sm text-muted-foreground">Target Price:</span>
              <span className={`font-semibold ${isAtOrBelowTarget ? 'text-green-500' : 'text-red-500'}`}>
                ${targetPrice.toFixed(2)}
              </span>
              {isAtOrBelowTarget && (
                <span className="text-green-500 text-xs">✓ At target</span>
              )}
            </div>
          )}
          {priceChange !== 0 && (
            <div className="flex items-center gap-2">
              <span className="text-sm text-muted-foreground">Change:</span>
              <span className={`font-semibold ${priceChange >= 0 ? 'text-green-500' : 'text-red-500'}`}>
                {priceChange >= 0 ? '+' : ''}{priceChange.toFixed(2)} ({priceChange >= 0 ? '+' : ''}{priceChangePercent.toFixed(1)}%)
              </span>
              {priceChange >= 0 ? (
                <TrendingUp className="h-4 w-4 text-green-500" />
              ) : (
                <TrendingDown className="h-4 w-4 text-red-500" />
              )}
            </div>
          )}
        </div>



        {/* Chart Loading/Error States */}
        {isLoading ? (
          <div className="flex flex-col items-center justify-center py-12 gap-4">
            <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
            <p className="text-muted-foreground">Loading price history...</p>
          </div>
        ) : error ? (
          <div className="flex flex-col items-center justify-center py-12 gap-4">
            <AlertTriangle className="h-8 w-8 text-red-500" />
            <p className="text-red-500">Failed to load price history</p>
            <p className="text-sm text-muted-foreground">
              {error instanceof Error ? error.message : 'Unknown error'}
            </p>
          </div>
        ) : filteredData.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-12 gap-4">
            <TrendingUp className="h-8 w-8 text-muted-foreground" />
            <p className="text-muted-foreground">No price history available</p>
            <p className="text-sm text-muted-foreground text-center">
              Price history will appear after the item has been tracked for some time
            </p>
          </div>
        ) : filteredData.length > 0 ? (
          <div className="h-100 w-full">
            <ResponsiveContainer width="100%" height="100%">
              <LineChart
                data={filteredData}
                margin={{ top: 5, right: 30, left: 20, bottom: 5 }}
              >
                <CartesianGrid strokeDasharray="3 3" stroke="#f0f0f0" />
                <XAxis
                  dataKey="date"
                  tick={{ fontSize: 12, fill: '#6b7280' }}
                  tickLine={false}
                  axisLine={false}
                />
                <YAxis
                  tickFormatter={(value) => `$${value.toFixed(2)}`}
                />
                <Tooltip
                  formatter={(value: number | undefined) => value !== undefined ? [`$${value.toFixed(2)}`, 'Price'] : ['N/A', 'Price']}
                />
                <Legend verticalAlign="top" height={36} />
                <Line
                  type="monotone"
                  dataKey="price"
                  name="Price"
                  stroke="#3b82f6"
                  strokeWidth={2}
                  dot={{ r: 4, fill: '#3b82f6' }}
                  activeDot={{ r: 6, fill: '#3b82f6' }}
                />
                {targetPrice && (
                  <Line
                    type="monotone"
                    dataKey={() => targetPrice}
                    name="Target Price"
                    stroke="#10b981"
                    strokeWidth={1}
                    strokeDasharray="5 5"
                    dot={false}
                    isAnimationActive={false}
                  />
                )}
              </LineChart>
            </ResponsiveContainer>
          </div>
        ) : (
          <div className="flex flex-col items-center justify-center py-12 gap-4">
            <TrendingUp className="h-8 w-8 text-muted-foreground" />
            <p className="text-muted-foreground">No valid price data available</p>
            <p className="text-sm text-muted-foreground text-center">
              Price history will appear after the item has valid price entries
            </p>
          </div>
        )}

        {/* Price Insights */}
        {filteredData.length > 0 && (
          <div className="mt-6 p-4 rounded-lg bg-background">
            <h3 className="font-semibold mb-3">Price Insights</h3>
            <div className="grid grid-cols-2 gap-4 text-sm">
              <div className="flex justify-between">
                <span className="text-muted-foreground">Highest Price:</span>
                <span className="font-medium">
                  ${Math.max(...filteredData.map(d => d.price)).toFixed(2)}
                </span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted-foreground">Lowest Price:</span>
                <span className="font-medium">
                  ${Math.min(...filteredData.map(d => d.price)).toFixed(2)}
                </span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted-foreground">Average Price:</span>
                <span className="font-medium">
                  ${(filteredData.reduce((sum, d) => sum + d.price, 0) / filteredData.length).toFixed(2)}
                </span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted-foreground">Price Points:</span>
                <span className="font-medium">{filteredData.length}</span>
              </div>
            </div>
          </div>
        )}
      </DialogContent>
    </Dialog>
  );
}
