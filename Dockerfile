FROM golang:1.24-alpine AS builder

RUN apk add --no-cache ca-certificates git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath \
    -ldflags="-s -w -X main.version=${VERSION:-0.0.1}" \
    -o app cmd/api/main.go


FROM gcr.io/distroless/base-debian12

WORKDIR /app

COPY --from=builder /app/app /app/app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

USER nonroot:nonroot

EXPOSE 8080

ENTRYPOINT ["/app/app"]