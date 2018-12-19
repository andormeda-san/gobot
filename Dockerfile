FROM golang:1.8-alpine
ADD . /go/src/gobot
RUN go install gobot

FROM alpine:latest
COPY --from=0 /go/bin/gobot .
CMD ["./gobot"]