FROM golang:1.17-alpine as builder
ENV GO111MODULE=on

WORKDIR /app
ADD . /app
RUN go mod download
RUN go build -o /space-booking-service ./cmd/bookingservice

FROM alpine:3.14
COPY --from=builder /space-booking-service /
EXPOSE 8080

CMD ["/space-booking-service"]