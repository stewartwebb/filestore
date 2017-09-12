FROM golang:1.8

RUN apt-get update \
    && apt-get install -y --no-install-recommends clamdscan \
    && rm -rf /var/lib/apt/lists/*

RUN mkdir -p /go/src/github.com/stewartwebb/filestore

ADD docker/app/startup.sh /startup.sh

RUN chmod +x /startup.sh

WORKDIR /go/src/github.com/stewartwebb/filestore/src

RUN go get github.com/githubnemo/CompileDaemon

ENTRYPOINT ["/startup.sh"]
