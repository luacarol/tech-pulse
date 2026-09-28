FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/ingestion ./cmd/ingestion

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
COPY --from=build /out/ingestion /usr/local/bin/ingestion
COPY configs /app/configs
WORKDIR /app
ENTRYPOINT ["ingestion"]
