FROM golang:1.23-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -o /orders-api ./cmd/server

FROM alpine:3.22
RUN addgroup -S app && adduser -S app -G app
COPY --from=build /orders-api /app/orders-api
USER app
EXPOSE 8080
EXPOSE 9090
ENTRYPOINT ["/app/orders-api"]
