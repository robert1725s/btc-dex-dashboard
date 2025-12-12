# BTC DEX Dashboard

Go + React を使った DEX（分散型取引所）間の BTC 価格スプレッド監視ダッシュボード。

## 概要

複数の DEX（Hyperliquid、Lighter、Aster）から BTC 価格をリアルタイムで取得し、スプレッド（価格差）とアービトラージ機会を可視化するダッシュボードです。

## 目的

このプロジェクトは **Go と React/TypeScript の学習** を目的として作成しました。

- Go: Gin フレームワーク、GORM、クリーンアーキテクチャ
- React: TypeScript、Recharts、Hooks

## 機能

- 3つの DEX からリアルタイム価格取得（2秒間隔）
- スプレッド計算とアービトラージ機会の検出
- 15分間の価格履歴チャート
- 統計情報（最大スプレッド、平均スプレッド）

## 技術スタック

### Backend
- Go 1.21+
- Gin (Web フレームワーク)
- GORM (ORM)
- SQLite
- Viper (設定管理)

### Frontend
- React 19
- TypeScript 5.9
- Vite 7
- Recharts

## アーキテクチャ

```
├── cmd/server/          # エントリーポイント
├── internal/
│   ├── api/             # HTTP ハンドラー・ミドルウェア
│   ├── config/          # 設定管理
│   ├── domain/model/    # ドメインモデル
│   ├── infrastructure/  # DB・外部 API クライアント
│   ├── job/             # 定期実行ジョブ
│   ├── repository/      # データアクセス層
│   └── service/         # ビジネスロジック
└── web/                 # React フロントエンド
```

## セットアップ

### Backend

```bash
go mod download
go run cmd/server/main.go
```

### Frontend

```bash
cd web
npm install
npm run dev
```

## 使用方法

1. Backend サーバーを起動（http://localhost:8080）
2. Frontend 開発サーバーを起動（http://localhost:5173）
3. ブラウザで http://localhost:5173 にアクセス

## API エンドポイント

| エンドポイント | 説明 |
|--------------|------|
| GET /api/health | ヘルスチェック |
| GET /api/spread | スプレッド・価格情報 |
| GET /api/exchanges | 取引所一覧 |
| GET /api/funding-rates | ファンディングレート |

