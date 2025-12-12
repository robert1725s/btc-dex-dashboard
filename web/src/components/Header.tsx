import { useState, useEffect } from 'react';

export function Header() {
  const [time, setTime] = useState(new Date());

  useEffect(() => {
    const timer = setInterval(() => setTime(new Date()), 1000);
    return () => clearInterval(timer);
  }, []);

  const formatTime = (date: Date) => {
    return date.toLocaleTimeString('ja-JP', {
      hour12: false,
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
      timeZone: 'Asia/Tokyo',
    });
  };

  return (
    <header className="header">
      <div className="header-left">
        <h1 className="logo">BTC DEX ARBITRAGE DASHBOARD</h1>
        <span className="header-divider">|</span>
        <span className="header-time">JST {formatTime(time)}</span>
      </div>
      <div className="header-right">
        <a
          href="https://github.com/robert1725s/btc-dex-dashboard"
          target="_blank"
          rel="noopener noreferrer"
          className="github-link"
        >
          GitHub
        </a>
      </div>
    </header>
  );
}
