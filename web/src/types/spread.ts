export interface PriceInfo {
  exchange_key: string;
  exchange_name: string;
  bid: number;
  ask: number;
  mid_price: number;
}

export interface ArbitrageInfo {
  buy_exchange: string;
  sell_exchange: string;
  buy_price: number;
  sell_price: number;
  spread_abs: number;
  spread_pct: number;
}

export interface HistoryPoint {
  timestamp: string;
  hyperliquid: number;
  lighter: number;
  aster: number;
  spread: number;
  spread_pct: number;
}

export interface MaxSpreadInfo {
  value: number;
  pct: number;
  timestamp: string;
  high_exchange: string;
  low_exchange: string;
  high_price: number;
  low_price: number;
}

export interface SpreadStats {
  max_spread: MaxSpreadInfo | null;
  avg_spread: number;
  avg_spread_pct: number;
  avg_price: number;
  period_minutes: number;
}

export interface SpreadResult {
  prices: PriceInfo[];
  buy_opportunity: ArbitrageInfo | null;
  sell_opportunity: ArbitrageInfo | null;
  history: HistoryPoint[];
  stats: SpreadStats;
}

export interface FundingRate {
  exchange: string;
  rate: number;
  nextUpdate: Date;
}
