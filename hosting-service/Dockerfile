FROM golang:1.24.4-alpine AS builder
WORKDIR /build

COPY hosting-service/go.mod hosting-service/go.sum ./hosting-service/
COPY hosting-contracts/go.mod hosting-contracts/go.sum ./hosting-contracts/
COPY hosting-events-contract/go.mod hosting-events-contract/go.sum ./hosting-events-contract/
COPY hosting-kit/go.mod hosting-kit/go.sum ./hosting-kit/

WORKDIR /build/hosting-service
RUN go mod download

WORKDIR /build
COPY hosting-contracts/ ./hosting-contracts/
COPY hosting-events-contract/ ./hosting-events-contract/
COPY hosting-kit/ ./hosting-kit/
COPY hosting-service/ ./hosting-service/

WORKDIR /build/hosting-service
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /server cmd/server/main.go
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /migrator cmd/migrator/main.go

FROM alpine:3.21

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

COPY --from=builder /server /server
COPY --from=builder /migrator /migrator

USER appuser

EXPOSE 8080

CMD ["/server"]