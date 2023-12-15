# ステップ1: ビルド用のベースイメージを選択
FROM golang:1.21.1 as builder

# 作業ディレクトリを設定
WORKDIR /app


# 依存関係の管理ファイルをコピー
COPY go.mod ./
COPY go.sum ./

# 依存関係をダウンロード
RUN go mod download

# ソースコードをコピー
COPY . .

RUN go run github.com/steebchen/prisma-client-go db push
# アプリケーションをビルド
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server .

# ステップ2: 実行用の軽量なベースイメージを選択
FROM alpine:latest  

# セキュリティアップデート
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# ビルドしたバイナリをステップ1からコピー
COPY --from=builder /app/server .

# サーバーの実行
CMD ["./server"]