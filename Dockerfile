FROM golang:1.7.1

WORKDIR /go/src/github.com/emreler/finch

RUN go get github.com/tools/godep

ADD Godeps Godeps

RUN godep restore

ADD . .

RUN go install

ENTRYPOINT finch --config=/etc/finch/config.json
