import { useMemo } from 'react';
import {
  ComposedChart,
  Area,
  XAxis,
  YAxis,
  Tooltip,
  ResponsiveContainer,
  ReferenceLine,
} from 'recharts';
import type { HistoryPoint, SpreadStats } from '../types/spread';

interface Props {
  history: HistoryPoint[];
  stats: SpreadStats;
}

const formatTime = (timestamp: string) => {
  return new Date(timestamp).toLocaleTimeString('ja-JP', {
    hour: '2-digit',
    minute: '2-digit',
    hour12: false,
    timeZone: 'Asia/Tokyo',
  });
};

const formatTimeWithSeconds = (timestamp: string) => {
  return new Date(timestamp).toLocaleTimeString('ja-JP', {
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    hour12: false,
    timeZone: 'Asia/Tokyo',
  });
};

const formatPrice = (value: number) => `$${value.toLocaleString()}`;

interface CustomTooltipProps {
  active?: boolean;
  payload?: Array<{
    name: string;
    value: number;
    color: string;
    payload?: { timestamp?: string };
  }>;
}

const CustomTooltip = ({ active, payload }: CustomTooltipProps) => {
  if (!active || !payload || payload.length === 0) return null;

  const timestamp = payload[0]?.payload?.timestamp;
  const timeWithSeconds = timestamp ? formatTimeWithSeconds(timestamp) : '';

  return (
    <div className="chart-tooltip">
      <div className="tooltip-time">{timeWithSeconds}</div>
      {payload.map((entry) => (
        <div key={entry.name} className="tooltip-row">
          <span className="tooltip-dot" style={{ backgroundColor: entry.color }} />
          <span className="tooltip-label">{entry.name}:</span>
          <span className="tooltip-value">{formatPrice(entry.value)}</span>
        </div>
      ))}
    </div>
  );
};

export function PriceChart({ history, stats }: Props) {
  // Format data for chart
  const { formattedData, tickIndices } = useMemo(() => {
    const formatted = history.map((item, index) => ({
      ...item,
      time: formatTime(item.timestamp),
      index,
    }));

    // Find indices where 10-minute boundary is crossed (for 10-minute interval ticks)
    const indices: number[] = [];
    let lastTenMinute = -1;

    formatted.forEach((item, index) => {
      const minute = new Date(item.timestamp).getMinutes();
      const tenMinute = Math.floor(minute / 10);
      if (tenMinute !== lastTenMinute) {
        indices.push(index);
        lastTenMinute = tenMinute;
      }
    });

    return { formattedData: formatted, tickIndices: indices };
  }, [history]);

  const allPrices = history.flatMap((d) => [d.hyperliquid, d.lighter, d.aster]).filter((p) => p > 0);
  const minPrice = allPrices.length > 0 ? Math.min(...allPrices) * 0.9995 : 0;
  const maxPrice = allPrices.length > 0 ? Math.max(...allPrices) * 1.0005 : 100000;

  const maxSpread = stats.max_spread;

  return (
    <div className="card price-chart">
      <div className="chart-header">
        <h2 className="card-title">BTC Price History (Last {stats.period_minutes}m)</h2>
        <div className="chart-legend">
          <div className="legend-item">
            <span className="legend-dot hyperliquid" />
            <span>Hyperliquid</span>
          </div>
          <div className="legend-item">
            <span className="legend-dot lighter" />
            <span>Lighter</span>
          </div>
          <div className="legend-item">
            <span className="legend-dot aster" />
            <span>Aster</span>
          </div>
        </div>
      </div>

      {maxSpread && (
        <div className="max-spread-info">
          <span className="max-spread-label">Spread</span>
          <span className="max-spread-value">
            ${maxSpread.value.toFixed(2)} ({maxSpread.pct.toFixed(3)}%)
          </span>
          <span className="max-spread-time">@ {formatTimeWithSeconds(maxSpread.timestamp)}</span>
          <span className="max-spread-detail">
            | Hi: {maxSpread.high_exchange} ${maxSpread.high_price.toLocaleString()}
            {' '}/ Lo: {maxSpread.low_exchange} ${maxSpread.low_price.toLocaleString()}
          </span>
        </div>
      )}

      <div className="chart-container">
        <ResponsiveContainer width="100%" height="100%">
          <ComposedChart data={formattedData} margin={{ top: 10, right: 10, left: 0, bottom: 0 }}>
            <defs>
              <linearGradient id="gradientHyperliquid" x1="0" y1="0" x2="0" y2="1">
                <stop offset="5%" stopColor="#58a6ff" stopOpacity={0.3} />
                <stop offset="95%" stopColor="#58a6ff" stopOpacity={0} />
              </linearGradient>
              <linearGradient id="gradientLighter" x1="0" y1="0" x2="0" y2="1">
                <stop offset="5%" stopColor="#3fb950" stopOpacity={0.3} />
                <stop offset="95%" stopColor="#3fb950" stopOpacity={0} />
              </linearGradient>
              <linearGradient id="gradientAster" x1="0" y1="0" x2="0" y2="1">
                <stop offset="5%" stopColor="#d29922" stopOpacity={0.3} />
                <stop offset="95%" stopColor="#d29922" stopOpacity={0} />
              </linearGradient>
            </defs>
            <XAxis
              dataKey="index"
              axisLine={false}
              tickLine={false}
              tick={{ fill: '#8b949e', fontSize: 12 }}
              tickFormatter={(index: number) => {
                const item = formattedData[index];
                if (!item || !tickIndices.includes(index)) return '';
                return item.time;
              }}
              interval={0}
              minTickGap={50}
            />
            <YAxis
              domain={[minPrice, maxPrice]}
              axisLine={false}
              tickLine={false}
              tick={{ fill: '#8b949e', fontSize: 12 }}
              tickFormatter={(value) => `$${(value / 1000).toFixed(0)}k`}
              width={55}
            />
            <Tooltip content={<CustomTooltip />} />
            <ReferenceLine y={minPrice + (maxPrice - minPrice) / 2} stroke="#30363d" strokeDasharray="3 3" />

            <Area
              type="monotone"
              dataKey="hyperliquid"
              name="Hyperliquid"
              stroke="#58a6ff"
              strokeWidth={2}
              fill="url(#gradientHyperliquid)"
            />
            <Area
              type="monotone"
              dataKey="lighter"
              name="Lighter"
              stroke="#3fb950"
              strokeWidth={2}
              fill="url(#gradientLighter)"
            />
            <Area
              type="monotone"
              dataKey="aster"
              name="Aster"
              stroke="#d29922"
              strokeWidth={2}
              fill="url(#gradientAster)"
            />
          </ComposedChart>
        </ResponsiveContainer>
      </div>
    </div>
  );
}
