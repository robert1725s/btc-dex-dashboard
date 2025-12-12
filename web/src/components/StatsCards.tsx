import type { ArbitrageInfo, SpreadStats } from '../types/spread';

interface Props {
  stats: SpreadStats;
  topArb: ArbitrageInfo | null;
}

export function StatsCards({ stats, topArb }: Props) {
  return (
    <div className="stats-cards">
      <div className="stat-card">
        <div className="stat-label">Max Spread</div>
        <div className="stat-content">
          <span className="stat-main">${stats.max_spread?.value.toFixed(0) ?? '0'}</span>
          <span className="stat-sub">({stats.max_spread?.pct.toFixed(3) ?? '0.000'}%)</span>
        </div>
      </div>

      <div className="stat-card">
        <div className="stat-label">{stats.period_minutes}m Avg Spread</div>
        <div className="stat-content">
          <span className="stat-main">${stats.avg_spread.toFixed(0)}</span>
          <span className="stat-sub">({stats.avg_spread_pct.toFixed(3)}%)</span>
        </div>
      </div>

      <div className="stat-card">
        <div className="stat-label">Top Arb Opp</div>
        <div className="stat-content">
          {topArb ? (
            <div className="arb-info">
              <span className="arb-long">Long {topArb.buy_exchange}</span>
              <span className="arb-sep">/</span>
              <span className="arb-short">Short {topArb.sell_exchange}</span>
            </div>
          ) : (
            <span className="stat-none">No opportunity</span>
          )}
        </div>
      </div>
    </div>
  );
}
