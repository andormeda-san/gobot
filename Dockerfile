FROM golang:latest as builder
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

ADD . /go/src/gobot
RUN go get github.com/PuerkitoBio/goquery
RUN go get github.com/nlopes/slack
RUN go install gobot

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /go/bin/gobot .
CMD ["./gobot"]