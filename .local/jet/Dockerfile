FROM golang:1.23

RUN go install github.com/go-jet/jet/v2/cmd/jet@latest

COPY ./entrypoint.sh /entrypoint.sh

RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]