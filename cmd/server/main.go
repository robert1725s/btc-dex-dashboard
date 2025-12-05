// cmd/server/main.go
// このファイルはアプリケーションのエントリーポイント（開始点）です。

package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// Gin のデフォルトルーターを作成
	// Default() は Logger と Recovery ミドルウェアが含まれている
	r := gin.Default()

	// GET /api/health エンドポイントを定義
	// 第1引数: パス
	// 第2引数: ハンドラー関数（リクエストを処理する関数）
	r.GET("/api/health", func(c *gin.Context) {
		// c.JSON() で JSON レスポンスを返す
		// 第1引数: HTTP ステータスコード
		// 第2引数: レスポンスボディ（gin.H は map[string]any のエイリアス）
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	})

	// サーバーを起動（デフォルトは :8080）
	r.Run(":8080")
}
