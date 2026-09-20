ARG GO_VERSION=1.27.3
FROM golang:${GO_VERSION}-bookworm AS builder

WORKDIR /src
ENV GOTOOLCHAIN=local

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags='-s -w' -o /out/developers-cv-api ./cmd/api

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app
COPY --from=builder /out/developers-cv-api /app/developers-cv-api

USER nonroot:nonroot
EXPOSE 8080
STOPSIGNAL SIGTERM

ENTRYPOINT ["/app/developers-cv-api"]
