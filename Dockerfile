FROM golang:1.15

WORKDIR /go/src/github.com/fin-assistant/
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /usr/local/bin/fin-assistant github.com/fin-assistant/cmd/assistant

###

FROM alpine:3.9

COPY --from=0 /usr/local/bin/fin-assistant /usr/local/bin/fin-assistant
RUN apk add --no-cache ca-certificates

ENTRYPOINT ["assistant"]
