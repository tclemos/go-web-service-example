FROM golang:1.13 as builder

RUN mkdir /build
WORKDIR /build

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

ADD . .

RUN CGO_ENABLED=0 GOOS=linux go build -o thing .

FROM alpine:latest

COPY --from=builder /build/thing /app/
RUN chmod +x /app/thing

EXPOSE 8080

CMD /app/thing