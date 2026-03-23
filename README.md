# Render.com 用 Go Web アプリケーション

これは Render.com にデプロイするための最小構成の Go Web アプリケーションテストページです。

## Render.com への導入（デプロイ）手順

### 1. GitHub へのプッシュ
まず、このディレクトリ (`C:\works\testWebpage2\`) の内容を GitHub リポジトリにプッシュします。
1. GitHub で新しいリポジトリを作成します（例：`render-go-test`）。
2. コマンドプロンプト等で以下のコマンドを実行し、コードをプッシュします：
   ```bash
   cd C:\works\testWebpage2
   git init
   git add .
   git commit -m "Initial commit"
   git branch -M main
   git remote add origin https://github.com/あなたのユーザー名/リポジトリ名.git
   git push -u origin main
   ```

### 2. Render.com でのサービス作成
1. [Render.com](https://dashboard.render.com/) のダッシュボードにログインします。
2. 右上の **「New」** ボタンをクリックし、**「Web Service」** を選択します。
3. **「Build and deploy from a Git repository」** を選択し、「Next」をクリックします。
4. 先ほど作成した GitHub のリポジトリを検索して **「Connect」** をクリックします。

### 3. デプロイ設定
以下の設定項目を入力・確認します。
- **Name**: サービスの名前（例：`go-test-app`）
- **Region**: お好みのリージョン（例：`Singapore` や `Oregon`など）
- **Branch**: `main`
- **Root Directory**: （空白のままでOK）
- **Environment**: `Go` （通常は自動で認識されます）
- **Build Command**: `go build -o main .`
- **Start Command**: `./main`
- **Instance Type**: 料金プラン（まずは無料の `Free` プランを選択してください）

### 4. デプロイの実行
- ページ下部の **「Create Web Service」** をクリックします。
- Render.com 側で自動的にビルドとデプロイが進みます。
- デプロイが完了すると、左上に割り当てられた URL（例：`https://go-test-app-xxxx.onrender.com`）が表示されるので、クリックしてアクセス確認してください。

### 仕様のポイント
Render.com などのクラウド環境では、外部からアクセスするためのポート番号が動的に決定されます。
そのため、`main.go` の中で `os.Getenv("PORT")` を使用して環境変数 `PORT` からポート番号を取得する形にしています。設定されていない場合（ローカルでの実行時など）は `8080` を使うフォールバック処理を入れています。
