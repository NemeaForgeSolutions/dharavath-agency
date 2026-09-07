FROM node:26-alpine AS css-builder

WORKDIR /app

COPY package.json package-lock.json ./

RUN npm ci

COPY web/ ./web/
COPY internal/ ./internal/

RUN npm run build:css

FROM golang:1.24-alpine AS go-builder

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY go.mod ./
RUN go mod download

COPY . .

COPY --from=css-builder /app/web/static/css/styles.css ./web/static/css/styles.css

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w -extldflags '-static'" \
    -o /app/bin/server ./cmd/server

FROM alpine:3.21 AS runtime

RUN apk add --no-cache ca-certificates tzdata

RUN addgroup -g 10001 -S appgroup && \
    adduser -u 10001 -S appuser -G appgroup

WORKDIR /app

COPY --from=go-builder /app/bin/server /app/server

USER appuser:appgroup

ENV HOST=0.0.0.0 \
    PORT=8080 \
    APP_ENV=production

EXPOSE 8080

ENTRYPOINT ["/app/server"]
