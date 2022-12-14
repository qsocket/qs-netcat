FROM golang:latest as builder
RUN go install github.com/qsocket/qs-netcat@latest
ENTRYPOINT ["qs-netcat"]