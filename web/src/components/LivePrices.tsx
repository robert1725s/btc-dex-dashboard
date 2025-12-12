import { useMemo } from 'react';
import type { PriceInfo } from '../types/spread';

interface Props {
  prices: PriceInfo[];
  avgPrice: number;
}

const exchangeColors: Record<string, string> = {
  hyperliquid: '#58a6ff',
  lighter: '#3fb950',
  aster: '#d29922',
};

export function LivePrices({ prices, avgPrice }: Props) {
  const calculateChange = (price: number) => {
    if (avgPrice === 0) return 0;
    return ((price - avgPrice) / avgPrice) * 100;
  };

  // Find min and max prices using mid_price from API
  const { minPrice, maxPrice } = useMemo(() => {
    if (prices.length === 0) return { minPrice: 0, maxPrice: 0 };
    const midPrices = prices.map((p) => p.mid_price);
    return {
      minPrice: Math.min(...midPrices),
      maxPrice: Math.max(...midPrices),
    };
  }, [prices]);

  const getPriceClass = (midPrice: number) => {
    if (midPrice === maxPrice) return 'price-highest';
    if (midPrice === minPrice) return 'price-lowest';
    return 'price-neutral';
  };

  return (
    <div className="card live-prices">
      <h2 className="card-title">Live Prices</h2>
      <table className="prices-table">
        <thead>
          <tr>
            <th>DEX</th>
            <th>Price</th>
            <th>vs Avg</th>
          </tr>
        </thead>
        <tbody>
          {prices.map((price) => {
            const change = calculateChange(price.mid_price);
            const isPositive = change >= 0;
            const colorKey = price.exchange_key.toLowerCase();
            const priceClass = getPriceClass(price.mid_price);

            return (
              <tr key={price.exchange_key}>
                <td>
                  <span
                    className="exchange-indicator"
                    style={{ backgroundColor: exchangeColors[colorKey] || '#8b949e' }}
                  />
                  {price.exchange_name}
                </td>
                <td className="price-cell">
                  <span className={priceClass}>
                    ${price.mid_price.toLocaleString(undefined, { minimumFractionDigits: 0, maximumFractionDigits: 0 })}
                  </span>
                </td>
                <td className={`change-cell ${isPositive ? 'positive' : 'negative'}`}>
                  {isPositive ? '+' : ''}{change.toFixed(2)}%
                </td>
              </tr>
            );
          })}
        </tbody>
      </table>
    </div>
  );
}
