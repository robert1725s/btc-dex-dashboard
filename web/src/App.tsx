import { useState, useEffect, useCallback } from 'react';
import './App.css';
import { Header } from './components/Header';
import { LivePrices } from './components/LivePrices';
import { StatsCards } from './components/StatsCards';
import { PriceChart } from './components/PriceChart';
import { FundingRates } from './components/FundingRates';
import type { SpreadResult, FundingRate } from './types/spread';

const API_URL = 'http://localhost:8080/api/spread';

function App() {
  const [data, setData] = useState<SpreadResult | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  const fetchData = useCallback(async () => {
    try {
      const response = await fetch(API_URL);
      if (!response.ok) {
        throw new Error(`HTTP error: ${response.status}`);
      }
      const result: SpreadResult = await response.json();
      setData(result);
      setError(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Unknown error');
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchData();
    const timer = setInterval(fetchData, 3000);
    return () => clearInterval(timer);
  }, [fetchData]);

  // Mock funding rates
  const fundingRates: FundingRate[] = [
    { exchange: 'Hyperliquid', rate: 0.0001, nextUpdate: new Date(Date.now() + 3600000) },
    { exchange: 'Lighter', rate: 0.00012, nextUpdate: new Date(Date.now() + 3600000) },
    { exchange: 'Aster', rate: 0.00009, nextUpdate: new Date(Date.now() + 3600000) },
  ];

  if (loading) {
    return (
      <div className="loading-screen">
        <div className="loading-spinner" />
        <div className="loading-text">Loading dashboard...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="error-screen">
        <div className="error-icon">!</div>
        <div className="error-message">Error: {error}</div>
        <button className="retry-button" onClick={fetchData}>
          Retry
        </button>
      </div>
    );
  }

  if (!data) {
    return (
      <div className="error-screen">
        <div className="error-message">No data available</div>
        <button className="retry-button" onClick={fetchData}>
          Retry
        </button>
      </div>
    );
  }

  return (
    <div className="app">
      <Header />
      <main className="dashboard">
        <section className="top-section">
          <LivePrices prices={data.prices} avgPrice={data.stats.avg_price} />
          <StatsCards
            stats={data.stats}
            topArb={data.buy_opportunity}
          />
        </section>
        <section className="bottom-section">
          <PriceChart history={data.history} stats={data.stats} />
          <FundingRates rates={fundingRates} />
        </section>
      </main>
    </div>
  );
}

export default App;
