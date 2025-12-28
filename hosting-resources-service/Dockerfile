FROM golang:1.25.5-alpine AS builder
WORKDIR /build

COPY hosting-resources-service/go.mod hosting-resources-service/go.sum ./hosting-resources-service/
COPY hosting-contracts/go.mod hosting-contracts/go.sum ./hosting-contracts/
COPY hosting-kit/go.mod hosting-kit/go.sum ./hosting-kit/

WORKDIR /build/hosting-resources-service
RUN go mod download

WORKDIR /build
COPY hosting-contracts/ ./hosting-contracts/
COPY hosting-kit/ ./hosting-kit/
COPY hosting-resources-service/ ./hosting-resources-service/

WORKDIR /build/hosting-resources-service
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /server cmd/server/main.go
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /migrator cmd/migrator/main.go

FROM alpine:3.21

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

COPY --from=builder /server /server
COPY --from=builder /migrator /migrator

USER appuser

EXPOSE 8080

CMD ["/server"]