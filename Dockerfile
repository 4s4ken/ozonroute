FROM alpine:3

WORKDIR /app
COPY . .

RUN go build -o notifier ./cmd/grpc/notifier

CMD ["./notifier"]
