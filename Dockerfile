FROM golang:1.7

RUN go get github.com/PuerkitoBio/goquery
RUN go get github.com/nlopes/slack
WORKDIR /go/src/
ADD . /go/src/bot
RUN go install bot

ENTRYPOINT /go/bin/bot