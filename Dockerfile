FROM golang as builder

WORKDIR /

COPY go.mod .

COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest

WORKDIR /app

COPY swagger-ui ./swagger-ui
COPY --from=builder /main .

CMD ["/app/main"]
