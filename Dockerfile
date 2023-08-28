FROM golang:1.20 as builder
WORKDIR /app
COPY . ./
RUN make build

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/training-store-backend ./
copy .env ./
EXPOSE 9090
CMD ["./training-store-backend"]