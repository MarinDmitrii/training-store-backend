FROM golang:1.20 as builder
WORKDIR /app
COPY . ./
RUN make build

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/training-store-backend ./
COPY .env ./
COPY web/stripe/ ./web/stripe
EXPOSE 9090
CMD ["./training-store-backend"]