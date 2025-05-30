ARG BUILD_ENV="base"

# 本番と開発環境両方で使われるステージ
FROM golang:1.24.3 AS base
# タイムゾーンと証明書の設定
RUN apt-get update && apt-get install -y --no-install-recommends tzdata ca-certificates && \
  cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime

# 開発環境用のステージ
FROM base AS dev
WORKDIR /go/src/myblog
ENV TZ="Asia/Tokyo"
RUN apt-get update && apt-get install -y build-essential default-mysql-server jq
RUN go install github.com/air-verse/air@latest

# ビルド用のステージ
# BUILD_ENVが"dev"に設定されている場合のみベースイメージとして"dev"ステージを使う
# BUILD_ENVが指定されていない場合はベースイメージとして"base"ステージを使うため、"dev"ステージがスキップされる
FROM $BUILD_ENV AS builder
WORKDIR /go/src/myblog
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o ./cmd/api/bin/api ./cmd/api/main.go

# リリース用のステージ
FROM gcr.io/distroless/static AS api
COPY --from=builder /go/src/myblog/cmd/api/bin/api /
COPY --from=builder /etc/localtime /etc/localtime
COPY --from=builder /go/src/myblog/db/migrations /db/migrations
EXPOSE 8080
ENTRYPOINT ["/api"]
