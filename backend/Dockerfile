# ベースとなるGoの公式イメージを指定
FROM golang:1.22


# 作業ディレクトリを設定
WORKDIR /app

# Go Modulesを利用している場合、go.modとgo.sumを先にコピー
COPY go.mod go.sum ./

# 依存関係を取得
RUN go mod download

# 残りのアプリケーションコードをコピー
COPY . .

# アプリケーションをビルド
RUN go build -o app .

# コンテナ起動時に実行するコマンドを指定
CMD ["./app"]
