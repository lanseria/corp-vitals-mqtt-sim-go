# =========================================================================
# Stage 1: Builder
# =========================================================================
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
ENV GOPROXY=https://goproxy.cn,direct
RUN go mod download

COPY . .

ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} go build -a -o /app/server ./cmd/app

# =========================================================================
# Stage 2: Final
# =========================================================================
FROM scratch

WORKDIR /app
COPY --from=builder /app/server .

ENTRYPOINT ["/app/server"]
