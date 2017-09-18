FROM centos:latest
LABEL MAINTAINER luopengift<870148195@qq.com>
ENV VERSION=0.0.2 GOROOT=/usr/local/go GOPATH=/data/golang PATH=/usr/local/go/bin:$PATH
WORKDIR $GOPATH/src/github.com/luopengift
RUN yum -y install wget git && wget https://storage.googleapis.com/golang/go1.9.linux-amd64.tar.gz && \
    tar -xvf go1.9.linux-amd64.tar.gz && rm -rf go1.9.linux-amd64.tar.gz && mv go /usr/local
RUN go get github.com/luopengift/transport && cd $GOPATH/src/github.com/luopengift/transport && \
    go get ./... && ./init.sh build cmd/main.go && mv transport /usr/local/bin && mkdir -p /usr/local/bin/var
COPY test/docker-test.json .
ENTRYPOINT ["transport"]
