FROM golang:1.24.4-alpine AS builder

WORKDIR /build

COPY hosting-provisioning-service/go.mod hosting-provisioning-service/go.sum ./hosting-provisioning-service/
COPY hosting-events-contract/go.mod hosting-events-contract/go.sum ./hosting-events-contract/
COPY hosting-kit/go.mod hosting-kit/go.sum ./hosting-kit/

WORKDIR /build/hosting-provisioning-service
RUN go mod download

WORKDIR /build
COPY hosting-events-contract/ ./hosting-events-contract/
COPY hosting-kit/ ./hosting-kit/
COPY hosting-provisioning-service/ ./hosting-provisioning-service/

WORKDIR /build/hosting-provisioning-service
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /server cmd/server/main.go

FROM alpine:3.21

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

COPY --from=builder /server /server

USER appuser
EXPOSE 8080
CMD ["/server"]