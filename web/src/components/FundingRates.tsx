import { useMemo } from 'react';
import type { FundingRate } from '../types/spread';

interface Props {
  rates: FundingRate[];
}

const exchangeColors: Record<string, string> = {
  Hyperliquid: '#58a6ff',
  Lighter: '#3fb950',
  Aster: '#d29922',
};

export function FundingRates({ rates }: Props) {
  // Calculate best funding arbitrage opportunity
  const bestArb = useMemo(() => {
    if (rates.length < 2) return null;

    let bestLong = rates[0];
    let bestShort = rates[0];

    for (const rate of rates) {
      // For funding arb: Long where rate is lowest (you receive), Short where rate is highest (you receive)
      if (rate.rate < bestLong.rate) bestLong = rate;
      if (rate.rate > bestShort.rate) bestShort = rate;
    }

    // Funding rate diff (short receives when positive, long pays when positive)
    const rateDiff = bestShort.rate - bestLong.rate;
    if (rateDiff <= 0) return null;

    return {
      longExchange: bestLong.exchange,
      shortExchange: bestShort.exchange,
      rateDiff: rateDiff * 100, // Convert to percentage
      annualized: rateDiff * 100 * 3 * 365, // 8h funding * 3 * 365
    };
  }, [rates]);

  return (
    <div className="card funding-rates">
      <h2 className="card-title">Funding Rates</h2>
      <table className="funding-table">
        <thead>
          <tr>
            <th>DEX</th>
            <th>Symbol</th>
            <th>Current Rate</th>
          </tr>
        </thead>
        <tbody>
          {rates.map((rate) => {
            const isPositive = rate.rate >= 0;
            return (
              <tr key={rate.exchange}>
                <td>
                  <span
                    className="exchange-indicator"
                    style={{ backgroundColor: exchangeColors[rate.exchange] || '#8b949e' }}
                  />
                  {rate.exchange}
                </td>
                <td className="symbol-cell">BTC-PERP</td>
                <td>
                  <span className={`rate-badge ${isPositive ? 'positive' : 'negative'}`}>
                    {isPositive ? '+' : ''}{(rate.rate * 100).toFixed(4)}%
                  </span>
                </td>
              </tr>
            );
          })}
        </tbody>
      </table>

      {bestArb && (
        <div className="funding-arb-info">
          <div className="funding-arb-label">Best FR Arbitrage</div>
          <div className="funding-arb-strategy">
            <span className="arb-long">Long {bestArb.longExchange}</span>
            <span className="arb-sep">/</span>
            <span className="arb-short">Short {bestArb.shortExchange}</span>
          </div>
          <div className="funding-arb-rate">
            <span className="rate-value">+{bestArb.rateDiff.toFixed(4)}%</span>
            <span className="rate-period">/8h</span>
            <span className="rate-annual">(~{bestArb.annualized.toFixed(1)}% APR)</span>
          </div>
        </div>
      )}
    </div>
  );
}
