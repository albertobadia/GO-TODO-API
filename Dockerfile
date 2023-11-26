FROM golang as builder

WORKDIR /

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest

WORKDIR /

COPY --from=builder /main .

CMD ["/main"]
