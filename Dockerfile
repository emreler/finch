FROM golang:1.7.1

WORKDIR /go/src/gitlab.com/emreler/finch

ADD . .

RUN go get

RUN go install

ENTRYPOINT finch --config=/etc/finch/config.json

EXPOSE 8081
