package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `
		<!DOCTYPE html>
		<html lang="ja">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Render.com テストページ</title>
			<style>
				body {
					font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
					background-color: #f4f7f6;
					color: #333;
					display: flex;
					flex-direction: column;
					align-items: center;
					justify-content: center;
					height: 100vh;
					margin: 0;
				}
				.container {
					background: #fff;
					padding: 40px;
					border-radius: 8px;
					box-shadow: 0 4px 6px rgba(0,0,0,0.1);
					text-align: center;
					max-width: 600px;
					width: 90%%;
				}
				h1 {
					color: #2c3e50;
					margin-bottom: 20px;
				}
				p {
					line-height: 1.6;
					color: #555;
				}
				.status {
					margin-top: 30px;
					padding: 10px;
					background-color: #e8f8f5;
					border-left: 4px solid #1abc9c;
					color: #16a085;
					font-weight: bold;
				}
			</style>
		</head>
		<body>
			<div class="container">
				<h1>Hello from Go on Render!</h1>
				<p>これは <strong>Render.com</strong> で動作する Go 言語による Web アプリケーションのテストページです。</p>
				<p>環境変数 <code>PORT</code> によって動的にポートが割り当てられ、リッスンしています。</p>
				<div class="status">サーバーは正常に稼働しています 🎉</div>
			</div>
		</body>
		</html>
	`)
}

func main() {
	// Render.com は環境変数 PORT に動的にポート番号を割り当てます
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // ローカル開発用のデフォルトポート
	}

	http.HandleFunc("/", handler)

	log.Printf("サーバーをポート %s で起動しています...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
