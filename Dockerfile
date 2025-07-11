FROM alpine:3

WORKDIR /app
COPY . .

RUN go build -o notifier ./cmd/notifier

CMD ["./notifier"]
