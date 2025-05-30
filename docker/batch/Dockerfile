# build (for development)
FROM golang:1.24.3 AS builder
WORKDIR /go/src/myblog
ENV TZ="Asia/Tokyo"
RUN apt-get update && apt-get install -y --no-install-recommends tzdata && \
  cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN make build-batch

# release
FROM gcr.io/distroless/static AS batch
COPY --from=builder /go/src/myblog/cmd/batch/bin/batch /usr/local/bin/myblog-batch
COPY --from=builder /etc/localtime /etc/localtime
ENTRYPOINT ["myblog-batch"]
